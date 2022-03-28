package restli

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PapaCharlie/go-restli/restlicodec"
)

const (
	ProtocolVersion = "2.0.0"

	IDHeader              = "X-RestLi-Id"
	MethodHeader          = "X-RestLi-Method"
	ProtocolVersionHeader = "X-RestLi-Protocol-Version"
	ErrorResponseHeader   = "X-RestLi-Error-Response"
)

type Method int

// Disabled until https://github.com/golang/go/issues/45218 is resolved: go:generate stringer -type=Method -trimprefix Method_
const (
	Method_Unknown = Method(iota)

	Method_get
	Method_create
	Method_delete
	Method_update
	Method_partial_update

	Method_batch_get
	Method_batch_create
	Method_batch_delete
	Method_batch_update
	Method_batch_partial_update

	Method_get_all

	Method_action
	Method_finder
)

var MethodNameMapping = func() map[string]Method {
	mapping := make(map[string]Method)
	for m := Method_get; m <= Method_finder; m++ {
		mapping[m.String()] = m
	}
	return mapping
}()

type Client struct {
	*http.Client
	HostnameResolver
	// Whether missing fields in a restli response should cause a MissingRequiredFields error to be returned. Note that
	// even if the error is returned, the response will still be fully deserialized.
	StrictResponseDeserialization bool
}

func (c *Client) FormatQueryUrl(rp ResourcePath, query QueryParamsEncoder) (*url.URL, error) {
	path, err := rp.ResourcePath()
	if err != nil {
		return nil, err
	}

	if query != nil {
		var params string
		params, err = query.EncodeQueryParams()
		if err != nil {
			return nil, err
		}
		path += "?" + params
	}

	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	root := rp.RootResource()
	hostUrl, err := c.ResolveHostnameAndContextForQuery(root, u)
	if err != nil {
		return nil, err
	}

	resolvedPath := "/" + strings.TrimSuffix(strings.TrimPrefix(hostUrl.EscapedPath(), "/"), "/")

	if resolvedPath == "/" {
		return hostUrl.ResolveReference(u), nil
	}

	if idx := strings.Index(resolvedPath, "/"+root); idx >= 0 &&
		(len(resolvedPath) == idx+len(root)+1 || resolvedPath[idx+len(root)+1] == '/') {
		resolvedPath = resolvedPath[:idx]
	}

	return hostUrl.Parse(resolvedPath + u.RequestURI())
}

func SetJsonAcceptHeader(req *http.Request) {
	req.Header.Set("Accept", "application/json")
}

func SetJsonContentTypeHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}

func SetRestLiHeaders(req *http.Request, method Method) {
	req.Header.Set(ProtocolVersionHeader, ProtocolVersion)
	req.Header.Set(MethodHeader, method.String())
}

type contextKey int

const (
	extraRequestHeadersKey contextKey = iota
	responseHeadersCaptorKey
	methodCtxKey
	resourcePathSegmentsCtxKey
	entitySegmentsCtxKey
	finderNameCtxKey
	actionNameCtxKey
)

// ExtraRequestHeaders returns a new context with the given headers. If a context with headers is passed into any Client
// request, the headers will be added to the request before being sent. Note that these headers will not override any
// existing headers such as Content-Type or MethodHeader. Only new headers will be added to the request.
func ExtraRequestHeaders(ctx context.Context, headers http.Header) context.Context {
	return context.WithValue(ctx, extraRequestHeadersKey, headers)
}

// AddResponseHeadersCaptor returns a new context and a http.Header. If the returned context is passed into any Client
// request, the returned http.Header will be populated with all the headers returned by the server.
func AddResponseHeadersCaptor(ctx context.Context) (context.Context, http.Header) {
	headers := http.Header{}
	ctx = context.WithValue(ctx, responseHeadersCaptorKey, headers)
	return ctx, headers
}

func newRequest(
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
	httpMethod string,
	method Method,
	body io.Reader,
) (*http.Request, error) {
	u, err := c.FormatQueryUrl(rp, query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, httpMethod, u.String(), body)
	if err != nil {
		return nil, err
	}

	SetRestLiHeaders(req, method)
	return req, nil
}

// NewGetRequest creates a GET http.Request and sets the expected rest.li headers
func NewGetRequest(
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
	method Method,
) (*http.Request, error) {
	req, err := newRequest(c, ctx, rp, query, http.MethodGet, method, http.NoBody)
	if err != nil {
		return nil, err
	}

	SetJsonAcceptHeader(req)

	return req, nil
}

// NewDeleteRequest creates a DELETE http.Request and sets the expected rest.li headers
func NewDeleteRequest(
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
	method Method,
) (*http.Request, error) {
	req, err := newRequest(c, ctx, rp, query, http.MethodDelete, method, http.NoBody)
	if err != nil {
		return nil, err
	}

	SetJsonAcceptHeader(req)

	return req, nil
}

func NewCreateRequest(
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
	method Method,
	create restlicodec.Marshaler,
	readOnlyFields restlicodec.PathSpec,
) (*http.Request, error) {
	return NewJsonRequest(c, ctx, rp, query, http.MethodPost, method, create, readOnlyFields)
}

// NewJsonRequest creates an http.Request with the given HTTP method and rest.li method, and populates the body of the
// request with the given restlicodec.Marshaler contents (see RawJsonRequest)
func NewJsonRequest(
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
	httpMethod string,
	restLiMethod Method,
	contents restlicodec.Marshaler,
	excludedFields restlicodec.PathSpec,
) (*http.Request, error) {
	writer := restlicodec.NewCompactJsonWriterWithExcludedFields(excludedFields)
	err := contents.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}

	size := writer.Size()
	req, err := newRequest(c, ctx, rp, query, httpMethod, restLiMethod, writer.ReadCloser())
	if err != nil {
		return nil, err
	}

	SetJsonAcceptHeader(req)
	SetJsonContentTypeHeader(req)

	req.ContentLength = int64(size)
	return req, nil
}

// Do is a very thin shim between the standard http.Client.Do. All it does it parse the response into a Error if
// the RestLi error header is set. A non-nil Response with a non-nil error will only occur if http.Client.Do returns
// such values (see the corresponding documentation). Otherwise, the response will only be non-nil if the error is nil.
// All (and only) network-related errors will be of type *url.Error. Other types of errors such as parse errors will use
// different error types.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	res, err := c.Client.Do(req)
	if err != nil {
		return res, err
	}

	err = IsErrorResponse(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DoAndUnmarshal calls Do and attempts to unmarshal the response into the given value. The response body will always be
// read to EOF and closed, to ensure the connection can be reused.
func DoAndUnmarshal[V any](
	c *Client,
	req *http.Request,
	unmarshaler restlicodec.GenericUnmarshaler[V],
) (v V, res *http.Response, err error) {
	data, res, err := c.do(req)
	if err != nil {
		return v, res, err
	}

	r, err := restlicodec.NewJsonReader(data)
	if err != nil {
		return v, res, err
	}
	v, err = unmarshaler(r)
	if _, mfe := err.(*restlicodec.MissingRequiredFieldsError); mfe && !c.StrictResponseDeserialization {
		err = nil
	}
	return v, res, err
}

// DoAndIgnore calls Do and drops the response's body. The response body will always be read to EOF and closed, to
// ensure the connection can be reused.
func DoAndIgnore(c *Client, req *http.Request) (*http.Response, error) {
	_, res, err := c.do(req)
	return res, err
}

func (c *Client) do(req *http.Request) ([]byte, *http.Response, error) {
	if extraHeaders, ok := req.Context().Value(extraRequestHeadersKey).(http.Header); ok {
		for k, v := range extraHeaders {
			if _, ok = req.Header[k]; !ok {
				req.Header[k] = v
			}
		}
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, res, err
	}

	if v := res.Header.Get(ProtocolVersionHeader); v != ProtocolVersion {
		return nil, nil, &UnsupportedRestLiProtocolVersion{v}
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, &url.Error{
			Op:  "ReadResponse",
			URL: req.URL.String(),
			Err: err,
		}
	}

	err = res.Body.Close()
	if err != nil {
		return nil, nil, &url.Error{
			Op:  "CloseResponse",
			URL: req.URL.String(),
			Err: err,
		}
	}

	if resHeaders, ok := req.Context().Value(responseHeadersCaptorKey).(http.Header); ok {
		for k, v := range res.Header {
			resHeaders[k] = v
		}
	}

	return data, res, nil
}
