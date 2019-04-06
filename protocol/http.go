package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
		return fmt.Errorf("%s", string(body))
	}

	return nil
}

type RestLiClient struct {
	*http.Client
	Hostname string
}

func (c *RestLiClient) FormatUrl(url string, segments ...string) string {
	a := make([]interface{}, len(segments))
	for i, s := range segments {
		a[i] = s
	}
	return c.Hostname + fmt.Sprintf(url, a...)
}

func (c *RestLiClient) GetRequest(url string, method string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, &bytes.Buffer{})
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

func (c *RestLiClient) PostRequest(url string, method string, contents interface{}) (*http.Request, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(contents)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, buf)
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
		return
	}

	err = IsErrorResponse(res)
	return
}
