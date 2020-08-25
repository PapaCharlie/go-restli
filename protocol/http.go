package protocol

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	// "github.com/pkg/errors"
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

var emptyBuffer = &bytes.Buffer{}

type RestLiError struct {
	Message        string
	ExceptionClass string
	StackTrace     string

	Status               int         `json:"-"`
	FullResponse         []byte      `json:"-"`
	ResponseHeaders      http.Header `json:"-"`
	DeserializationError error       `json:"-"`
}

func (r *RestLiError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		io.WriteString(s, r.Error()+"\n")
		io.WriteString(s, r.StackTrace)
	case 's':
		io.WriteString(s, r.Error()+"\n")
	}
}

func (r *RestLiError) Error() string {
	return fmt.Sprintf("RestLiError(status: %d, exceptionClass: %s, message: %s)", r.Status, r.ExceptionClass, r.Message)
}

func IsErrorResponse(res *http.Response) error {
	var err error
	var body []byte

	if strings.ToLower(res.Header.Get(RestLiHeader_ErrorResponse)) == "true" {
		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		restLiError := &RestLiError{
			Status:          res.StatusCode,
			FullResponse:    body,
			ResponseHeaders: res.Header,
		}
		if deserializationError := json.Unmarshal(body, restLiError); deserializationError != nil {
			restLiError.DeserializationError = deserializationError
		}
		return restLiError
	}

	if res.StatusCode >= 500 {
		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}

	return nil
}

type SimpleHostnameSupplier struct {
	Hostname *url.URL
}

func (s *SimpleHostnameSupplier) ResolveHostnameAndContextForQuery(string, *url.URL) (*url.URL, error) {
	return s.Hostname, nil
}

type HostnameResolver interface {
	// ResolveHostnameAndContextForQuery takes in the name of the service for which to resolve the hostname, along with
	// the URL for the query that is about to be sent. The service name is often the top-level parent resource's name,
	// but can be any unique identifier for a D2 endpoint. Some HostnameResolver implementations will choose to ignore
	// this parameter and resolve hostnames using a different strategy. By default, the generated code will always pass
	// in the top-level parent resource's name.
	ResolveHostnameAndContextForQuery(serviceName string, query *url.URL) (*url.URL, error)
}

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

func (c *RestLiClient) GetRequest(ctx context.Context, url *url.URL, method RestLiMethod) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), emptyBuffer)
	if err != nil {
		return nil, err
	}

	SetRestLiHeaders(req, method)
	SetJsonAcceptHeader(req)

	return req, nil
}

func (c *RestLiClient) DeleteRequest(ctx context.Context, url *url.URL, method RestLiMethod) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url.String(), emptyBuffer)
	if err != nil {
		return nil, err
	}

	SetRestLiHeaders(req, method)
	SetJsonAcceptHeader(req)

	return req, nil
}

func (c *RestLiClient) JsonPutRequest(ctx context.Context, url *url.URL, restLiMethod RestLiMethod, contents interface{}) (*http.Request, error) {
	return jsonRequest(ctx, url, http.MethodPut, restLiMethod, contents)
}

func (c *RestLiClient) JsonPostRequest(ctx context.Context, url *url.URL, restLiMethod RestLiMethod, contents interface{}) (*http.Request, error) {
	return jsonRequest(ctx, url, http.MethodPost, restLiMethod, contents)
}

func jsonRequest(ctx context.Context, url *url.URL, httpMethod string, restLiMethod RestLiMethod, contents interface{}) (*http.Request, error) {
	buf, err := json.Marshal(contents)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, httpMethod, url.String(), bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	SetRestLiHeaders(req, restLiMethod)
	SetJsonAcceptHeader(req)
	SetJsonContentTypeHeader(req)

	return req, nil
}

func (c *RestLiClient) RawPostRequest(url *url.URL, method RestLiMethod, contents []byte) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(contents))
	if err != nil {
		return nil, err
	}

	SetRestLiHeaders(req, method)

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

	err = IsErrorResponse(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DoAndDecode calls Do and attempts to unmarshal the response into the given value. The response body will always be
// read to EOF and closed, to ensure the connection can be reused.
func (c *RestLiClient) DoAndDecode(req *http.Request, v interface{}) (res *http.Response, err error) {
	return c.doAndConsumeBody(req, func(body []byte) error {
		return json.Unmarshal(body, v)
	})
}

// DoAndDecode calls Do and drops the response's body. The response body will always be read to EOF and closed, to
// ensure the connection can be reused.
func (c *RestLiClient) DoAndIgnore(req *http.Request) (res *http.Response, err error) {
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
