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
	RestliContentType     = "application/json"
	RestliProtocolVersion = "2.0.0"
)

const (
	RestliHeader_Method          = "X-RestLi-Method"
	RestliHeader_ProtocolVersion = "X-RestLi-Protocol-Version"
	RestliHeader_ErrorResponse   = "X-RestLi-Error-Response"
	RestliHeader_Id              = "X-RestLi-Id"
)

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

	if res.Header.Get(RestliHeader_ErrorResponse) == "true" {
		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		restliError := &RestLiError{
			FullResponse: body,
		}
		if err := json.NewDecoder(bytes.NewReader(body)).Decode(restliError); err != nil {
			restliError.DeserializationError = err
		}
		return restliError
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

func (c *RestLiClient) GetRequest(url *url.URL, method string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url.String(), emptyBuffer)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", RestliContentType)
	req.Header.Set(RestliHeader_ProtocolVersion, RestliProtocolVersion)
	if method != "" {
		req.Header.Set(RestliHeader_Method, method)
	}

	return req, nil
}

func (c *RestLiClient) PostRequest(url *url.URL, method string, contents interface{}) (*http.Request, error) {
	buf, err := json.Marshal(contents)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", RestliContentType)
	req.Header.Set(RestliHeader_ProtocolVersion, RestliProtocolVersion)
	if method != "" {
		req.Header.Set(RestliHeader_Method, method)
	}

	return req, nil
}

func (c *RestLiClient) Do(req *http.Request) (res *http.Response, err error) {
	res, err = c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	err = IsErrorResponse(res)
	return res, err
}
