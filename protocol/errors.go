package protocol

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// RestLiError is returned by the Do* methods when the X-RestLi-Error-Response header is set to true.
type RestLiError struct {
	Message        string `json:"message"`
	ExceptionClass string `json:"exceptionClass"`
	StackTrace     string `json:"stackTrace"`

	// Will be non-nil if an error occurred when attempting to deserialize the actual JSON response fields (i.e. Status,
	// Message, ExceptionClass and StackTrace)
	DeserializationError error `json:"-"`
	// The raw response that this error was parsed from. Note that to ensure that the connection can be reused, the Body
	// of the response is fully read into ResponseBody then closed
	Response     *http.Response `json:"-"`
	ResponseBody []byte         `json:"-"`
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
	return fmt.Sprintf("RestLiError(status: %d, exceptionClass: %s, message: %s)",
		r.Response.StatusCode, r.ExceptionClass, r.Message)
}

// UnexpectedStatusCodeError is returned by the Do* methods when the target rest.li service responded with non-2xx code
// but did not set the expected X-RestLi-Error-Response header.
type UnexpectedStatusCodeError struct {
	// The raw response that of the failed call. Note that to ensure that the connection can be reused, the Body
	// of the response is fully read into ResponseBody then closed
	Response     *http.Response
	ResponseBody []byte
}

func (u *UnexpectedStatusCodeError) Error() string {
	return fmt.Sprintf("unexpected response code from %s: %s", u.Response.Request.RequestURI, u.Response.Status)
}

// IsErrorResponse checks the contents of the given http.Response and if the X-RestLi-Error-Response is set to `true`,
// parses the body of the response into a RestLiError. If the header is not set, but the status code isn't a 2xx code,
// an UnexpectedStatusCodeError will be returned instead. Note that an UnexpectedStatusCodeError contains the
// http.Request and http.Response that resulted in this error, therefore an expected non-2xx can always be manually
// handled/recovered (e.g. a 3xx code redirecting to the HTTPS endpoint).
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
			ResponseBody: body,
			Response:     res,
		}
		if deserializationError := json.Unmarshal(body, restLiError); deserializationError != nil {
			restLiError.DeserializationError = deserializationError
		}
		return restLiError
	}

	if res.StatusCode/100 != 2 {
		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return &UnexpectedStatusCodeError{
			Response:     res,
			ResponseBody: body,
		}
	}

	return nil
}

// BatchMethodError is returned by BATCH_* requests when the "errors" field in the response is present and not null. It
// contains the raw bytes of said field as its contents are untyped.
type BatchMethodError struct {
	// The contents of the "errors" fields in the batch request's response
	// TODO: This is a map of the key that caused the error to a message describing the error. Ideally instead of having
	//  raw bytes, each batch method would generate an error handler that correctly deserializes this response
	Errors []byte
	// The contents of the "statuses" fields in the batch request's response
	// TODO: See above, this should also be deserialized into an actual type
	Statuses []byte
	// The request that resulted in this error
	Request *http.Request
	// The raw response that this error was parsed from
	Response *http.Response
}

func (b *BatchMethodError) Error() string {
	return string(b.Errors)
}

type UnsupportedRestLiProtocolVersion struct {
	ReturnedVersion string
}

func (u *UnsupportedRestLiProtocolVersion) Error() string {
	return fmt.Sprintf("go-restli: Unsupported rest.li protocol version: %s", u.ReturnedVersion)
}
