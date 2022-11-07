package restli

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PapaCharlie/go-restli/v2/restlicodec"
)

const (
	ProtocolVersion = "2.0.0"

	IDHeader              = "X-RestLi-Id"
	MethodHeader          = "X-RestLi-Method"
	ProtocolVersionHeader = "X-RestLi-Protocol-Version"
	ErrorResponseHeader   = "X-RestLi-Error-Response"
	MethodOverrideHeader  = "X-HTTP-Method-Override"

	ContentTypeHeader          = "Content-Type"
	MultipartMixedContentType  = "multipart/mixed"
	MultipartBoundary          = "boundary"
	ApplicationJsonContentType = "application/json"
	FormUrlEncodedContentType  = "application/x-www-form-urlencoded"
)

type Method int

//go:generate stringer -type=Method -trimprefix Method_
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
	HostnameResolver HostnameResolver
	// Whether missing fields in a restli response should cause a MissingRequiredFields error to be returned. Note that
	// even if the error is returned, the response will still be fully deserialized.
	StrictResponseDeserialization bool
	// When greater than 0, this enables request tunnelling. When a request's query is longer than this value, the
	// request will instead be sent via POST, with the query encoded as a form query and the MethodOverrideHeader set to
	// the original HTTP method.
	QueryTunnellingThreshold int
}

func (c *Client) formatQueryUrl(rp ResourcePath, query QueryParamsEncoder) (*url.URL, error) {
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
	hostUrl, err := c.HostnameResolver.ResolveHostnameAndContextForQuery(root, u)
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

// ExtraRequestHeaders returns a context.Context to be passed into any generated client methods. Upon request creation,
// the given function will be executed, and the headers will be added to the request before being sent. Note that these
// headers will not override any existing headers such as Content-Type. Only new headers will be added to the request.
func ExtraRequestHeaders(ctx context.Context, f func() (http.Header, error)) context.Context {
	return context.WithValue(ctx, extraRequestHeadersKey, f)
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
	contents restlicodec.Marshaler,
	excludedFields restlicodec.PathSpec,
) (req *http.Request, err error) {
	u, err := c.formatQueryUrl(rp, query)
	if err != nil {
		return nil, err
	}

	var body []byte
	if contents != nil {
		writer := restlicodec.NewCompactJsonWriterWithExcludedFields(excludedFields)
		err = contents.MarshalRestLi(writer)
		if err != nil {
			return nil, err
		}
		body = []byte(writer.Finalize())
	}

	headers := http.Header{}
	if body != nil {
		headers.Set(ContentTypeHeader, ApplicationJsonContentType)
	}

	if c.QueryTunnellingThreshold > 0 && len(u.RawQuery) > c.QueryTunnellingThreshold {
		var tunnelHeaders http.Header
		body, tunnelHeaders = EncodeTunnelledQuery(httpMethod, u.RawQuery, body)
		for k := range tunnelHeaders {
			headers.Set(k, tunnelHeaders.Get(k))
		}
		httpMethod = http.MethodPost
		u.RawQuery = ""
	}

	req, err = http.NewRequestWithContext(ctx, httpMethod, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set(ProtocolVersionHeader, ProtocolVersion)
	req.Header.Set(MethodHeader, method.String())
	req.Header.Set("Accept", ApplicationJsonContentType)
	for k, v := range headers {
		req.Header[k] = v
	}

	if extraHeaders, ok := req.Context().Value(extraRequestHeadersKey).(func() (http.Header, error)); ok {
		extras, err := extraHeaders()
		if err != nil {
			return nil, err
		}

		for k, v := range extras {
			if _, ok = req.Header[k]; !ok {
				req.Header[k] = v
			}
		}
	}

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
	return newRequest(c, ctx, rp, query, http.MethodGet, method, nil, nil)
}

// NewDeleteRequest creates a DELETE http.Request and sets the expected rest.li headers
func NewDeleteRequest(
	c *Client,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
	method Method,
) (*http.Request, error) {
	return newRequest(c, ctx, rp, query, http.MethodDelete, method, nil, nil)
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
	if contents == nil {
		return nil, fmt.Errorf("go-restli: Must provide non-nil contents")
	}
	return newRequest(c, ctx, rp, query, httpMethod, restLiMethod, contents, excludedFields)
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
