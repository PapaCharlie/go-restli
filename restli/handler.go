package restli

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/PapaCharlie/go-restli/restli/batchkeyset"
	"github.com/PapaCharlie/go-restli/restlicodec"
	"github.com/PapaCharlie/go-restli/restlidata"
)

type handler func(
	ctx *RequestContext,
	segments []restlicodec.Reader,
	body []byte,
) (responseBody restlicodec.Marshaler, err error)

type pathNode struct {
	ResourcePathSegment
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
		methods:             copyMap(p.methods),
		finders:             copyMap(p.finders),
		actions:             copyMap(p.actions),
		subNodes:            copyCloneableMap(p.subNodes),
	}
}

func handle(rootNode *pathNode, res http.ResponseWriter, req *http.Request) {
	path := req.URL.RawPath
	if path == "" {
		path = req.URL.Path
	}
	path = strings.Trim(path, "/")

	segments := strings.Split(path, "/")
	if len(segments) == 0 || rootNode.subNodes[segments[0]] == nil {
		http.NotFound(res, req)
		return
	}

	sub := rootNode.subNodes[segments[0]]
	ctx := &RequestContext{
		Request:         req,
		ResponseHeaders: res.Header(),
		ResponseStatus:  http.StatusOK,
	}
	responseBody, err := sub.receive(ctx, nil, segments)
	if errRes, ok := err.(*restlidata.ErrorResponse); ok {
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

	res.Header().Set(ProtocolVersionHeader, ProtocolVersion)

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
	keySegments []restlicodec.Reader,
	remainingSegments []string,
) (responseBody restlicodec.Marshaler, err error) {
	hasEntity := false
	if len(remainingSegments) >= 1 {
		if p.isCollection && len(remainingSegments) > 1 {
			var r restlicodec.Reader
			r, err = restlicodec.NewRor2Reader(remainingSegments[1])
			if err != nil {
				return newErrorResponsef(err, http.StatusNotFound, "Invalid path segment %q: %s", remainingSegments[1])
			}
			hasEntity = true
			keySegments = append(keySegments, r)

			remainingSegments = remainingSegments[2:]
		} else {
			remainingSegments = remainingSegments[1:]
		}
		if len(remainingSegments) != 0 {
			subResource := remainingSegments[0]
			if sub, ok := p.subNodes[subResource]; ok {
				return sub.receive(ctx, keySegments, remainingSegments)
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
		case Method_finder, Method_create, Method_batch_get, Method_batch_create, Method_batch_delete, Method_batch_update,
			Method_batch_partial_update, Method_get_all:
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

	var h handler
	if restLiMethod == Method_finder {
		h = p.finders[finder]
		if h == nil {
			return newErrorResponsef(nil, http.StatusBadRequest, "Finder %q not defined on %q", finder, p.name)
		}
	} else if restLiMethod == Method_action {
		h = p.actions[action]
		if h == nil {
			return newErrorResponsef(nil, http.StatusBadRequest, "Action %q not defined on %q", action, p.name)
		}
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
			err = &restlidata.ErrorResponse{
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

	return h(ctx, keySegments, body)
}

func newErrorResponsef(cause error, status int, format string, a ...any) (restlicodec.Marshaler, error) {
	if e, ok := cause.(*restlidata.ErrorResponse); ok {
		return nil, e
	}
	if cause != nil {
		a = append(a, cause)
	}
	return nil, &restlidata.ErrorResponse{
		Status:  Int32Pointer(int32(status)),
		Message: StringPointer(fmt.Sprintf(format, a...)),
	}
}

type Server interface {
	AddToMux(mux Mux)

	subNode(segments []ResourcePathSegment) *pathNode
}

func newPathNode(segment ResourcePathSegment) *pathNode {
	return &pathNode{
		ResourcePathSegment: segment,
		methods:             map[Method]handler{},
		finders:             map[string]handler{},
		actions:             map[string]handler{},
		subNodes:            map[string]*pathNode{},
	}
}

func NewServer() Server {
	return newPathNode(ResourcePathSegment{})
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
			subNode.subNodes[s.name] = newPathNode(s)
		}
		subNode = subNode.subNodes[s.name]
		if subNode.isCollection != s.isCollection {
			log.Panicf("go-restli: Inconsistent isCollection value for %v", segments)
		}
	}
	return subNode
}

type Mux interface {
	Handle(string, http.Handler)
}

func (p *pathNode) AddToMux(mux Mux) {
	pCopy := p.clone()
	h := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handle(pCopy, res, req)
	})
	for r := range pCopy.subNodes {
		mux.Handle(r, h)
	}
}

// Note: this implementation of http.Handler is _not_ threadsafe. The fact that pathNode implements this method is never
// exposed outside this package and therefore should not be used. It is intended exclusively for testing purposes.
func (p *pathNode) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	handle(p, res, req)
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
		if !restlidata.IsEmptyRecord(queryParams) {
			queryParams, err = restlicodec.UnmarshalQueryParamsDecoder[QP](ctx.Request.URL.RawQuery)
			if err != nil {
				return newErrorResponsef(err, http.StatusBadRequest, "Invalid query params for %q: %s", method)
			}
		}

		responseBody, err = h(ctx, rp, queryParams, body)
		if _, ok := err.(*restlidata.ErrorResponse); err != nil && !ok {
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

func SetLocation[K any](ctx *RequestContext, c *restlidata.CreatedEntity[K]) error {
	c.Location = new(string)
	w := restlicodec.NewRor2PathWriter()
	err := restlicodec.MarshalRestLi[K](c.Id, w)
	if err != nil {
		return err
	}
	*c.Location = strings.TrimSuffix(ctx.Request.RequestURI, "/") + "/" + w.Finalize()
	return nil
}
