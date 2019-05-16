package protocol

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	RestLiProtocolVersion = "2.0.0"

	RestLiHeader_Method          = "X-RestLi-Method"
	RestLiHeader_ProtocolVersion = "X-RestLi-Protocol-Version"
	RestLiHeader_ErrorResponse   = "X-RestLi-Error-Response"
)

type RestLiMethod string

const (
	MethodGet           = RestLiMethod("Get")
	MethodCreate        = RestLiMethod("Create")
	MethodDelete        = RestLiMethod("Delete")
	MethodUpdate        = RestLiMethod("Update")
	MethodPartialUpdate = RestLiMethod("PartialUpdate")

	MethodBatchGet           = RestLiMethod("BatchGet")
	MethodBatchCreate        = RestLiMethod("BatchCreate")
	MethodBatchDelete        = RestLiMethod("BatchDelete")
	MethodBatchUpdate        = RestLiMethod("BatchUpdate")
	MethodBatchPartialUpdate = RestLiMethod("BatchPartialUpdate")

	MethodGetAll = RestLiMethod("GetAll")

	NoMethod = RestLiMethod("")
)

var RestLiMethodNameMapping = map[string]RestLiMethod{
	"get":            MethodGet,
	"create":         MethodCreate,
	"delete":         MethodDelete,
	"update":         MethodUpdate,
	"partial_update": MethodPartialUpdate,

	"batch_get":            MethodBatchGet,
	"batch_create":         MethodBatchCreate,
	"batch_delete":         MethodBatchDelete,
	"batch_update":         MethodBatchUpdate,
	"batch_partial_update": MethodBatchPartialUpdate,

	"get_all": MethodGetAll,
}

var emptyBuffer = &bytes.Buffer{}

type RestLiError struct {
	Status         int
	Message        string
	ExceptionClass string
	StackTrace     string

	FullResponse         []byte `json:"-"`
	DeserializationError error  `json:"-"`
}

func (r *RestLiError) Error() string {
	return fmt.Sprintf("RestLiError(status: %d, exceptionClass: %s, message: %s)", r.Status, r.ExceptionClass, r.Message)
}

func IsErrorResponse(res *http.Response) error {
	var err error
	var body []byte

	if res.Header.Get(RestLiHeader_ErrorResponse) == "true" {
		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		restLiError := &RestLiError{
			FullResponse: body,
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

func (s *SimpleHostnameSupplier) GetHostnameForQuery(string) (*url.URL, error) {
	return s.Hostname, nil
}

type HostnameResolver interface {
	GetHostnameForQuery(query string) (*url.URL, error)
}

type RestLiClient struct {
	*http.Client
	HostnameResolver
}

// Assumes a leading slash
func getFirstPathSegment(path string) string {
	idx := strings.Index(path[1:], "/")
	if idx > 0 {
		return path[:idx+1]
	} else {
		return path
	}
}

func (c *RestLiClient) FormatQueryUrl(rawQuery string) (*url.URL, error) {
	hostUrl, err := c.GetHostnameForQuery(rawQuery)
	if err != nil {
		return nil, err
	}

	rawQuery = "/" + strings.TrimPrefix(rawQuery, "/")
	query, err := url.Parse(rawQuery)
	if err != nil {
		return nil, err
	}

	hostPath := hostUrl.EscapedPath()
	if hostPath == "" || hostPath == "/" {
		return hostUrl.ResolveReference(query), nil
	}
	// The restli spec allows for at most one context path segment. If not, it becomes impossible to know when the
	// context ends and the query begins
	firstHostSegment := getFirstPathSegment(hostPath)
	firstQuerySegment := getFirstPathSegment(query.EscapedPath())
	if firstHostSegment == firstQuerySegment {
		return hostUrl.ResolveReference(query), nil
	} else {
		return hostUrl.Parse(firstHostSegment + query.RequestURI())
	}
}

func SetJsonContentTypeHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}

func SetRestLiHeaders(req *http.Request, method RestLiMethod) {
	req.Header.Set(RestLiHeader_ProtocolVersion, RestLiProtocolVersion)
	if method != NoMethod {
		req.Header.Set(RestLiHeader_Method, string(method))
	}
}

func (c *RestLiClient) GetRequest(url *url.URL, method RestLiMethod) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url.String(), emptyBuffer)
	if err != nil {
		return nil, err
	}

	SetRestLiHeaders(req, method)

	return req, nil
}

func (c *RestLiClient) JsonPostRequest(url *url.URL, method RestLiMethod, contents interface{}) (*http.Request, error) {
	buf, err := json.Marshal(contents)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	SetRestLiHeaders(req, method)
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

func (c *RestLiClient) Do(req *http.Request) (res *http.Response, err error) {
	res, err = c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	err = IsErrorResponse(res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *RestLiClient) DoAndDecode(req *http.Request, v interface{}) error {
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = res.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}

func (c *RestLiClient) DoAndIgnore(req *http.Request) error {
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	err = res.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
