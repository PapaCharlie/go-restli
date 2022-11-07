package restli

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/PapaCharlie/go-restli/v2/restli/batchkeyset"
	"github.com/PapaCharlie/go-restli/v2/restlicodec"
	"github.com/PapaCharlie/go-restli/v2/restlidata/generated/com/linkedin/restli/common"
)

type handler func(
	ctx *RequestContext,
	segments []restlicodec.Reader,
	body []byte,
) (responseBody restlicodec.Marshaler, err error)

type rootNode struct {
	*pathNode
	prefix  string
	filters []Filter
}

type pathNode struct {
	ResourcePathSegment
	rootNode *rootNode
	methods  map[Method]handler
	finders  map[string]handler
	actions  map[string]handler
	subNodes map[string]*pathNode
}

func copyMap[K comparable, V any](m map[K]V) map[K]V {
	mCopy := make(map[K]V, len(m))
	for k, v := range m {
		mCopy[k] = v
	}
	return mCopy
}

func copyCloneableMap[K comparable, V interface{ clone() V }](m map[K]V) map[K]V {
	mCopy := make(map[K]V, len(m))
	for k, v := range m {
		mCopy[k] = v.clone()
	}
	return mCopy
}

func (p *pathNode) clone() *pathNode {
	return &pathNode{
		ResourcePathSegment: p.ResourcePathSegment,
		rootNode:            p.rootNode,
		methods:             copyMap(p.methods),
		finders:             copyMap(p.finders),
		actions:             copyMap(p.actions),
		subNodes:            copyCloneableMap(p.subNodes),
	}
}

func (r *rootNode) Handler() http.Handler {
	deepCopy := &rootNode{
		prefix:  r.prefix,
		filters: append([]Filter(nil), r.filters...),
	}
	p := new(pathNode)
	*p = *r.pathNode
	p.rootNode = deepCopy
	deepCopy.pathNode = p.clone()
	return deepCopy
}

func (r *rootNode) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.RawPath
	if path == "" {
		path = req.URL.Path
	}

	if !strings.HasPrefix(path, r.prefix) {
		http.NotFound(res, req)
		return
	}

	path = strings.TrimPrefix(path, r.prefix)

	segments := strings.Split(path, "/")
	if len(segments) == 0 || r.subNodes[segments[0]] == nil {
		http.NotFound(res, req)
		return
	}

	err := DecodeTunnelledQuery(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	sub := r.subNodes[segments[0]]
	ctx := &RequestContext{
		Request:         req,
		ResponseHeaders: res.Header(),
		ResponseStatus:  http.StatusOK,
	}

	res.Header().Set(ProtocolVersionHeader, ProtocolVersion)

	responseBody, err := sub.receive(ctx, nil, nil, segments)
	if err == nil {
		for i := len(r.filters) - 1; i >= 0; i-- {
			err = r.filters[i].PostRequest(ctx.Request.Context(), res.Header())
			if err != nil {
				break
			}
		}
	}

	if errRes, ok := err.(*common.ErrorResponse); ok {
		res.Header().Set(ErrorResponseHeader, "true")
		responseBody = errRes
		if errRes.Status != nil {
			ctx.ResponseStatus = int(*errRes.Status)
		} else {
			ctx.ResponseStatus = http.StatusInternalServerError
		}
		if errRes.Message == nil {
			errRes.Message = StringPointer(http.StatusText(int(*errRes.Status)))
		}
	} else if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if responseBody != nil {
		w := restlicodec.NewCompactJsonWriter()
		err = responseBody.MarshalRestLi(w)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		} else {
			data := []byte(w.Finalize())
			res.Header().Set("Content-Type", "application/json")
			res.Header().Set("Content-Length", strconv.Itoa(len(data)))

			res.WriteHeader(ctx.ResponseStatus)
			// not much we can do about this error
			_, _ = res.Write(data)
		}
	} else {
		res.WriteHeader(ctx.ResponseStatus)
	}
}

func (p *pathNode) receive(
	ctx *RequestContext,
	pathSegments []ResourcePathSegment,
	entitySegments []validatedRor2String,
	remainingSegments []string,
) (responseBody restlicodec.Marshaler, err error) {
	pathSegments = append(pathSegments, p.ResourcePathSegment)
	hasEntity := false
	if len(remainingSegments) >= 1 {
		if p.isCollection && len(remainingSegments) > 1 {
			err = restlicodec.ValidateRor2Input(remainingSegments[1])
			if err != nil {
				return newErrorResponsef(err, http.StatusNotFound, "Invalid path segment %q: %s", remainingSegments[1])
			}
			hasEntity = true
			entitySegments = append(entitySegments, validatedRor2String(remainingSegments[1]))

			remainingSegments = remainingSegments[2:]
		} else {
			remainingSegments = remainingSegments[1:]
		}
		if len(remainingSegments) != 0 {
			subResource := remainingSegments[0]
			if sub, ok := p.subNodes[subResource]; ok {
				return sub.receive(ctx, pathSegments, entitySegments, remainingSegments)
			} else {
				return newErrorResponsef(nil, http.StatusNotFound, "Unknown sub resource: %q", subResource)
			}
		}
	}

	restLiMethod := MethodNameMapping[ctx.Request.Header.Get(MethodHeader)]
	httpMethod := ctx.Request.Method
	params, err := restlicodec.ParseQueryParams(ctx.Request.URL.RawQuery)
	if err != nil {
		return nil, err
	}

	var finder string
	if q, ok := params["q"]; ok {
		finder = q.String()
	}

	var action string
	if a, ok := params["action"]; ok {
		action = a.String()
	}

	if p.isCollection {
		// For whatever reason rest.li makes the method header optional, because it supposedly can be inferred from the
		// method itself and query params. The documentation (https://linkedin.github.io/rest.li/spec/protocol#message-headers)
		// implies this procedure is unambiguous so the following assumes "q", "action" and "ids" are reserved query
		// parameter names for their corresponding HTTP method to make the routing simpler. In practice, this isn't
		// actually true since the Java implementation lets you define query parameters named "q" or "action" for
		// methods like GET, but until this bites, this is how this logic will be implemented.
		if restLiMethod == Method_Unknown {
			hasIds := params[batchkeyset.EntityIDsField] != nil

			switch httpMethod {
			case http.MethodGet:
				switch {
				// Only the GET method can specify an entity, so skip checking the query params
				case hasEntity:
					restLiMethod = Method_get
				case finder != "":
					restLiMethod = Method_finder
				case hasIds:
					restLiMethod = Method_batch_get
				default:
					restLiMethod = Method_get_all
				}
			case http.MethodPost:
				return newErrorResponsef(nil, http.StatusBadRequest,
					"Header %q is required for POST requests", MethodHeader)
			case http.MethodDelete:
				if hasIds {
					restLiMethod = Method_batch_delete
				} else {
					restLiMethod = Method_delete
				}
			case http.MethodPut:
				if hasIds {
					restLiMethod = Method_batch_update
				} else {
					restLiMethod = Method_update
				}
			}
		}

		switch restLiMethod {
		case Method_get, Method_delete, Method_update, Method_partial_update:
			if !hasEntity {
				return newErrorResponsef(nil, http.StatusBadRequest, "No entity provided for %q method", restLiMethod)
			}
		case Method_finder, Method_create, Method_batch_get, Method_batch_create, Method_batch_delete,
			Method_batch_update, Method_batch_partial_update, Method_get_all:
			if hasEntity {
				return newErrorResponsef(nil, http.StatusBadRequest, "Cannot provide an entity for %q", restLiMethod)
			}
		}
	} else {
		switch httpMethod {
		case http.MethodGet:
			restLiMethod = Method_get
		case http.MethodPut:
			restLiMethod = Method_update
		case http.MethodDelete:
			restLiMethod = Method_delete
		case http.MethodPost:
			if action != "" {
				restLiMethod = Method_action
			} else {
				restLiMethod = Method_partial_update
			}
		}
		if hasEntity {
			return newErrorResponsef(nil, http.StatusBadRequest, "Cannot provide an entity for simple resources")
		}
	}

	newCtx := context.WithValue(ctx.Request.Context(), methodCtxKey, restLiMethod)
	newCtx = context.WithValue(newCtx, resourcePathSegmentsCtxKey, pathSegments)
	newCtx = context.WithValue(newCtx, entitySegmentsCtxKey, entitySegments)

	var h handler
	if restLiMethod == Method_finder {
		h = p.finders[finder]
		if h == nil {
			return newErrorResponsef(nil, http.StatusBadRequest, "Finder %q not defined on %q", finder, p.name)
		}
		newCtx = context.WithValue(newCtx, finderNameCtxKey, finder)
	} else if restLiMethod == Method_action {
		h = p.actions[action]
		if h == nil {
			return newErrorResponsef(nil, http.StatusBadRequest, "Action %q not defined on %q", action, p.name)
		}
		newCtx = context.WithValue(newCtx, actionNameCtxKey, action)
	} else {
		h = p.methods[restLiMethod]
		if h == nil {
			return newErrorResponsef(nil, http.StatusBadRequest, "%q not defined on %q", restLiMethod, p.name)
		}
	}

	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			log.Printf("Failed to handle %q (%s)\n%s", ctx.Request.URL, r, stack)
			err = &common.ErrorResponse{
				Status:     Int32Pointer(http.StatusInternalServerError),
				Message:    StringPointer(fmt.Sprint(r)),
				StackTrace: &stack,
			}
		}
	}()
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, err
	}

	ctx.Request = ctx.Request.WithContext(newCtx)
	for _, f := range p.rootNode.filters {
		newCtx, err = f.PreRequest(ctx.Request)
		if err != nil {
			return nil, err
		}
		if newCtx != nil {
			ctx.Request = ctx.Request.WithContext(newCtx)
		}
	}

	return h(ctx, segmentReaders(entitySegments), body)
}

type validatedRor2String string

func segmentReaders(segments []validatedRor2String) (readers []restlicodec.Reader) {
	readers = make([]restlicodec.Reader, len(segments))
	for i, s := range segments {
		readers[i], _ = restlicodec.NewRor2Reader(string(s))
	}
	return readers
}

func GetMethodFromContext(ctx context.Context) Method {
	return ctx.Value(methodCtxKey).(Method)
}

func GetResourcePathSegmentsFromContext(ctx context.Context) []ResourcePathSegment {
	return ctx.Value(resourcePathSegmentsCtxKey).([]ResourcePathSegment)
}

func GetEntitySegmentsFromContext(ctx context.Context) []restlicodec.Reader {
	return segmentReaders(ctx.Value(entitySegmentsCtxKey).([]validatedRor2String))
}

func GetFinderNameFromContext(ctx context.Context) string {
	return ctx.Value(finderNameCtxKey).(string)
}

func GetActionNameFromContext(ctx context.Context) string {
	return ctx.Value(actionNameCtxKey).(string)
}

func newErrorResponsef(cause error, status int, format string, a ...any) (restlicodec.Marshaler, error) {
	if e, ok := cause.(*common.ErrorResponse); ok {
		return nil, e
	}
	if cause != nil {
		a = append(a, cause)
	}
	return nil, &common.ErrorResponse{
		Status:  Int32Pointer(int32(status)),
		Message: StringPointerf(format, a...),
	}
}

type Server interface {
	// AddToMux adds a http.Handler for each root resource registered against this Server. Note that resources
	// registered to this Server after AddToMux is called will not be reflected.
	AddToMux(mux *http.ServeMux)
	// Handler returns a raw http.Handler backed by a copy of this sever. Note that resources registered to this Server
	// after Handler will not be reflected. This is not meant to be used in conjunction with a http.ServeMux, but
	// instead with methods like http.ListenAndServe or http.ListenAndServeTLS
	Handler() http.Handler

	subNode(segments []ResourcePathSegment) *pathNode
}

func (p *pathNode) newSubNode(segment ResourcePathSegment) *pathNode {
	return &pathNode{
		ResourcePathSegment: segment,
		rootNode:            p.rootNode,
		methods:             map[Method]handler{},
		finders:             map[string]handler{},
		actions:             map[string]handler{},
		subNodes:            map[string]*pathNode{},
	}
}

// Filter implementations can enrich the request as it comes in by adding values to the context. They also receive a
// callback to add any response headers, if necessary. Filters' PreRequest methods will ve called in the order in which
// they were passed to NewServer, and PostRequest methods in the inverse order.
type Filter interface {
	// PreRequest is called after the request is parsed and the corresponding method is found. It is not called on any
	// invalid requests. The request's context will have corresponding values for the method, resource segments, entity
	// segments, finder name (only set if the method is Method_finder) and action name (only set if the method is
	// Method_action). Use the corresponding FromContext methods to get the values. If the returned context is non-nil,
	// it will replace the context passed to the actual resource implementation.
	PreRequest(req *http.Request) (context.Context, error)
	// PostRequest is called with the original request context and the response header map right before the response
	// header is written.
	PostRequest(ctx context.Context, responseHeaders http.Header) error
}

func NewPrefixedServer(prefix string, filters ...Filter) Server {
	if prefix == "" {
		prefix = "/"
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	r := &rootNode{
		pathNode: &pathNode{
			ResourcePathSegment: ResourcePathSegment{},
			methods:             map[Method]handler{},
			finders:             map[string]handler{},
			actions:             map[string]handler{},
			subNodes:            map[string]*pathNode{},
		},
		prefix:  "/",
		filters: filters,
	}
	r.rootNode = r
	return r
}

func NewServer(filters ...Filter) Server {
	return NewPrefixedServer("/", filters...)
}

type ResourcePathSegment struct {
	name         string
	isCollection bool
}

func NewResourcePathSegment(name string, isCollection bool) ResourcePathSegment {
	return ResourcePathSegment{
		name:         name,
		isCollection: isCollection,
	}
}

func (p *pathNode) subNode(segments []ResourcePathSegment) *pathNode {
	subNode := p
	for _, s := range segments {
		if subNode.subNodes == nil {
			subNode.subNodes = make(map[string]*pathNode)
		}
		if subNode.subNodes[s.name] == nil {
			subNode.subNodes[s.name] = subNode.newSubNode(s)
		}
		subNode = subNode.subNodes[s.name]
		if subNode.isCollection != s.isCollection {
			log.Panicf("go-restli: Inconsistent isCollection value for %v", segments)
		}
	}
	return subNode
}

func (r *rootNode) AddToMux(mux *http.ServeMux) {
	h := r.Handler()
	for rootResource := range r.subNodes {
		mux.Handle(r.prefix+rootResource, h)
	}
}

func registerMethod[RP ResourcePathUnmarshaler[RP], QP restlicodec.QueryParamsDecoder[QP]](
	s Server,
	segments []ResourcePathSegment,
	method Method,
	h func(
		ctx *RequestContext,
		rp RP,
		qp QP,
		body []byte,
	) (responseBody restlicodec.Marshaler, err error),
) {
	p := s.subNode(segments)
	if _, ok := p.methods[method]; ok {
		log.Panicf("go-restli: Cannot register method %q twice for %v", method, segments)
	}
	p.methods[method] = func(
		ctx *RequestContext,
		segments []restlicodec.Reader,
		body []byte,
	) (responseBody restlicodec.Marshaler, err error) {
		rp, err := UnmarshalResourcePath[RP](segments)
		if err != nil {
			return newErrorResponsef(err, http.StatusBadRequest, "Invalid path for %q: %s", method)
		}

		var queryParams QP
		if !common.IsEmptyRecord(queryParams) {
			queryParams, err = restlicodec.UnmarshalQueryParamsDecoder[QP](ctx.Request.URL.RawQuery)
			if err != nil {
				return newErrorResponsef(err, http.StatusBadRequest, "Invalid query params for %q: %s", method)
			}
		}

		responseBody, err = h(ctx, rp, queryParams, body)
		if _, ok := err.(*common.ErrorResponse); err != nil && !ok {
			return newErrorResponsef(err, http.StatusInternalServerError, "%q failed: %s", method)
		} else {
			return responseBody, err
		}
	}
}

func registerMethodWithNoBody[RP ResourcePathUnmarshaler[RP], QP restlicodec.QueryParamsDecoder[QP]](
	s Server,
	segments []ResourcePathSegment,
	method Method,
	h func(
		ctx *RequestContext,
		rp RP,
		qp QP,
	) (responseBody restlicodec.Marshaler, err error),
) {
	registerMethod(s, segments, method,
		func(ctx *RequestContext, rp RP, qp QP, body []byte) (responseBody restlicodec.Marshaler, err error) {
			if len(body) != 0 {
				return newErrorResponsef(nil, http.StatusBadRequest, "%q does not take a body", method)
			}

			return h(ctx, rp, qp)
		})
}

func registerMethodWithBody[RP ResourcePathUnmarshaler[RP], QP restlicodec.QueryParamsDecoder[QP], V any](
	s Server,
	segments []ResourcePathSegment,
	method Method,
	excludedFields restlicodec.PathSpec,
	leadingScopeToIgnore int,
	unmarshaler restlicodec.GenericUnmarshaler[V],
	h func(
		ctx *RequestContext,
		rp RP,
		v V,
		qp QP,
	) (responseBody restlicodec.Marshaler, err error),
) {
	registerMethod(s, segments, method,
		func(ctx *RequestContext, rp RP, qp QP, body []byte) (responseBody restlicodec.Marshaler, err error) {
			var v V
			r, err := restlicodec.NewJsonReaderWithExcludedFields(body, excludedFields, leadingScopeToIgnore)
			if err == nil {
				v, err = unmarshaler(r)
			}
			if err != nil {
				return newErrorResponsef(err, http.StatusBadRequest, "Invalid request body for %q: %s", method)
			}

			return h(ctx, rp, v, qp)
		})
}

type RequestContext struct {
	Request         *http.Request
	ResponseHeaders http.Header
	ResponseStatus  int
}

func (c *RequestContext) RequestPath() string {
	path := c.Request.URL.RawPath
	if path == "" {
		path = c.Request.URL.Path
	}
	return path
}

func SetLocation[K any](ctx *RequestContext, c *common.CreatedEntity[K]) error {
	c.Location = new(string)
	w := restlicodec.NewRor2HeaderWriter()
	err := restlicodec.MarshalRestLi[K](c.Id, w)
	if err != nil {
		return err
	}
	*c.Location = strings.TrimSuffix(ctx.RequestPath(), "/") + "/" + w.Finalize()
	return nil
}
