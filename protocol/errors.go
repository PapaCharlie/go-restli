package protocol

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
	"github.com/PapaCharlie/go-restli/protocol/stdtypes"
)

// RestLiError is returned by the Do* methods when the X-RestLi-Error-Response header is set to true.
type RestLiError struct {
	stdtypes.ErrorResponse
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
		io.WriteString(s, r.Error())
		if r.StackTrace != nil {
			io.WriteString(s, "\n"+*r.StackTrace)
		}
	case 's':
		io.WriteString(s, r.Error())
	}
}

func (r *RestLiError) Error() string {
	b := strings.Builder{}
	b.WriteString("RestLiError(status: ")

	if r.Status != nil {
		b.WriteString(strconv.Itoa(int(*r.Status)))
	} else {
		b.WriteString("UNKNOWN")
	}

	if r.ExceptionClass != nil {
		b.WriteString(", exceptionClass: ")
		b.WriteString(*r.ExceptionClass)
	}

	if r.Message != nil {
		b.WriteString(", message: ")
		b.WriteString(*r.Message)
	}

	b.WriteString(")")

	return b.String()
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
		restLiError.DeserializationError = restLiError.UnmarshalRestLi(restlicodec.NewJsonReader(body))
		if restLiError.Status == nil {
			restLiError.Status = Int32Pointer(int32(res.StatusCode))
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

// UnsupportedRestLiProtocolVersion is returned when the server returns a version other than the requested one, which is
// always RestLiProtocolVersion.
type UnsupportedRestLiProtocolVersion struct {
	ReturnedVersion string
}

func (u *UnsupportedRestLiProtocolVersion) Error() string {
	return fmt.Sprintf("go-restli: Unsupported rest.li protocol version: %s (the only supported version is %s)",
		u.ReturnedVersion, RestLiProtocolVersion)
}

// CreateResponseHasNoEntityHeaderError is used specifically when a Create request succeeds but the resource
// implementation does not set the X-RestLi-Id header. This error is recoverable and can be ignored if the response id
// is not required
type CreateResponseHasNoEntityHeaderError struct {
	Response *http.Response
}

func (c CreateResponseHasNoEntityHeaderError) Error() string {
	return "go-restli: response from CREATE request did not specify a " + RestLiHeader_ID + " header"
}

// IllegalPartialUpdateError is returned by PartialUpdateFieldChecker a partial update struct defines an illegal
// operation, such as deleting and setting the same field.
type IllegalPartialUpdateError struct {
	Message    string
	RecordType string
	Field      string
}

func (c *IllegalPartialUpdateError) Error() string {
	return fmt.Sprintf("go-restli: %s field %q of %q", c.Message, c.Field, c.RecordType)
}

// BatchRequestResponseError is returned by all the batch methods, and represents the keys on which the operation
// failed.
type BatchRequestResponseError[K comparable] map[K]*stdtypes.ErrorResponse

func (b BatchRequestResponseError[K]) Error() string {
	prettyErrors := make(map[string]string, len(b))
	for k, v := range b {
		w := restlicodec.NewCompactJsonWriter()
		_ = v.MarshalRestLi(w)
		prettyErrors[fmt.Sprintf("%+v", k)] = w.Finalize()
	}
	return fmt.Sprintf("go-restli: Not all batch operations successful: %+v", prettyErrors)
}

var NilQueryParams = fmt.Errorf("go-restli: Query params cannot be nil")
