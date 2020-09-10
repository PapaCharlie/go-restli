package protocol

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

const (
	RestLiProtocolVersion = "2.0.0"

	RestLiHeader_ID              = "X-RestLi-Id"
	RestLiHeader_Method          = "X-RestLi-Method"
	RestLiHeader_ProtocolVersion = "X-RestLi-Protocol-Version"
	RestLiHeader_ErrorResponse   = "X-RestLi-Error-Response"
)

type RestLiMethod int

//go:generate stringer -type=RestLiMethod -trimprefix Method_
const (
	Method_Unknown = RestLiMethod(iota)

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

var RestLiMethodNameMapping = func() map[string]RestLiMethod {
	mapping := make(map[string]RestLiMethod)
	for m := Method_get; m <= Method_finder; m++ {
		mapping[m.String()] = m
	}
	return mapping
}()

type RestLiClient struct {
	*http.Client
	HostnameResolver
}

func (c *RestLiClient) FormatQueryUrl(resourceBasename, rawQuery string) (*url.URL, error) {
	rawQuery = "/" + strings.TrimPrefix(rawQuery, "/")
	query, err := url.Parse(rawQuery)
	if err != nil {
		return nil, err
	}

	hostUrl, err := c.ResolveHostnameAndContextForQuery(resourceBasename, query)
	if err != nil {
		return nil, err
	}

	resolvedPath := "/" + strings.TrimSuffix(strings.TrimPrefix(hostUrl.EscapedPath(), "/"), "/")

	if resolvedPath == "/" {
		return hostUrl.ResolveReference(query), nil
	}

	if idx := strings.Index(resolvedPath, resourceBasename); idx >= 0 {
		resolvedPath = resolvedPath[:idx-1]
	}

	return hostUrl.Parse(resolvedPath + query.RequestURI())
}

func SetJsonAcceptHeader(req *http.Request) {
	req.Header.Set("Accept", "application/json")
}

func SetJsonContentTypeHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}

func SetRestLiHeaders(req *http.Request, method RestLiMethod) {
	req.Header.Set(RestLiHeader_ProtocolVersion, RestLiProtocolVersion)
	req.Header.Set(RestLiHeader_Method, method.String())
}

// GetRequest creates a GET http.Request and sets the expected rest.li headers
func GetRequest(ctx context.Context, url *url.URL, method RestLiMethod) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), http.NoBody)
	if err != nil {
		return nil, err
	}

	SetRestLiHeaders(req, method)
	SetJsonAcceptHeader(req)

	return req, nil
}

// DeleteRequest creates a DELETE http.Request and sets the expected rest.li headers
func DeleteRequest(ctx context.Context, url *url.URL, method RestLiMethod) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url.String(), http.NoBody)
	if err != nil {
		return nil, err
	}

	SetRestLiHeaders(req, method)
	SetJsonAcceptHeader(req)

	return req, nil
}

// JsonRequest creates an http.Request with the given HTTP method and rest.li method, and populates the body of the
// request with the given restlicodec.Marshaler contents (see RawJsonRequest)
func JsonRequest(
	ctx context.Context,
	url *url.URL,
	httpMethod string,
	restLiMethod RestLiMethod,
	contents restlicodec.Marshaler,
	excludedFields restlicodec.PathSpec,
) (*http.Request, error) {
	writer := restlicodec.NewCompactJsonWriterWithExcludedFields(excludedFields)
	err := contents.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}

	return RawJsonRequest(ctx, url, httpMethod, restLiMethod, strings.NewReader(writer.Finalize()))
}

// JsonRequest creates an http.Request with the given HTTP method and rest.li method, and populates the body of the
// request with the given reader
func RawJsonRequest(ctx context.Context, url *url.URL, httpMethod string, restLiMethod RestLiMethod, contents io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, httpMethod, url.String(), contents)
	if err != nil {
		return nil, err
	}

	SetRestLiHeaders(req, restLiMethod)
	SetJsonAcceptHeader(req)
	SetJsonContentTypeHeader(req)

	return req, nil
}

// Do is a very thin shim between the standard http.Client.Do. All it does it parse the response into a RestLiError if
// the RestLi error header is set. A non-nil Response with a non-nil error will only occur if http.Client.Do returns
// such values (see the corresponding documentation). Otherwise, the response will only be non-nil if the error is nil.
func (c *RestLiClient) Do(req *http.Request) (*http.Response, error) {
	res, err := c.Client.Do(req)
	if err != nil {
		return res, err
	}

	err = IsErrorResponse(req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DoAndDecode calls Do and attempts to unmarshal the response into the given value. The response body will always be
// read to EOF and closed, to ensure the connection can be reused.
func (c *RestLiClient) DoAndDecode(req *http.Request, v restlicodec.Unmarshaler) (*http.Response, error) {
	return c.doAndConsumeBody(req, func(body []byte) error {
		return v.UnmarshalRestLi(restlicodec.NewJsonReader(body))
	})
}

// DoAndDecode calls Do and drops the response's body. The response body will always be read to EOF and closed, to
// ensure the connection can be reused.
func (c *RestLiClient) DoAndIgnore(req *http.Request) (*http.Response, error) {
	return c.doAndConsumeBody(req, func([]byte) error {
		return nil
	})
}

func (c *RestLiClient) doAndConsumeBody(req *http.Request, bodyConsumer func(body []byte) error) (*http.Response, error) {
	res, err := c.Do(req)
	if err != nil {
		return res, err
	}

	if v := res.Header.Get(RestLiHeader_ProtocolVersion); v != RestLiProtocolVersion {
		return nil, fmt.Errorf("go-restli: Unsupported rest.li protocol version: %s", v)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	err = bodyConsumer(data)
	if err != nil {
		return nil, err
	}

	return res, nil
}
