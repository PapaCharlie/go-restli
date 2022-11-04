package restli

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PapaCharlie/go-restli/v2/restlidata/generated/com/linkedin/restli/common"
)

// Error is returned by the Do* methods when the X-RestLi-Error-Response header is set to true.
type Error struct {
	common.ErrorResponse
	// Will be non-nil if an error occurred when attempting to deserialize the actual JSON response fields (i.e. Status,
	// Message, ExceptionClass and StackTrace)
	DeserializationError error `json:"-"`
	// The raw response that this error was parsed from. Note that to ensure that the connection can be reused, the Body
	// of the response is fully read into ResponseBody then closed
	Response     *http.Response `json:"-"`
	ResponseBody []byte         `json:"-"`
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
	return fmt.Sprintf("unexpected response code from %s: %s", u.Response.Request.URL, u.Response.Status)
}

// IsErrorResponse checks the contents of the given http.Response and if the X-RestLi-Error-Response is set to `true`,
// parses the body of the response into a Error. If the header is not set, but the status code isn't a 2xx code,
// an UnexpectedStatusCodeError will be returned instead. Note that an UnexpectedStatusCodeError contains the
// http.Request and http.Response that resulted in this error, therefore an expected non-2xx can always be manually
// handled/recovered (e.g. a 3xx code redirecting to the HTTPS endpoint).
func IsErrorResponse(res *http.Response) error {
	var err error
	var body []byte

	if strings.ToLower(res.Header.Get(ErrorResponseHeader)) == "true" {
		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		restLiError := &Error{
			ResponseBody: body,
			Response:     res,
		}
		restLiError.DeserializationError = restLiError.UnmarshalJSON(body)
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
// always ProtocolVersion.
type UnsupportedRestLiProtocolVersion struct {
	ReturnedVersion string
}

func (u *UnsupportedRestLiProtocolVersion) Error() string {
	return fmt.Sprintf("go-restli: Unsupported rest.li protocol version: %s (the only supported version is %s)",
		u.ReturnedVersion, ProtocolVersion)
}

// CreateResponseHasNoEntityHeaderError is used specifically when a Create request succeeds but the resource
// implementation does not set the X-RestLi-Id header. This error is recoverable and can be ignored if the response id
// is not required
type CreateResponseHasNoEntityHeaderError struct {
	Response *http.Response
}

func (c CreateResponseHasNoEntityHeaderError) Error() string {
	return "go-restli: response from CREATE request did not specify a " + IDHeader + " header"
}

var NilQueryParams = fmt.Errorf("go-restli: Query params cannot be nil")

type IllegalEnumConstant struct {
	Enum     string
	Constant int
}

func (i *IllegalEnumConstant) Error() string {
	return fmt.Sprintf("go-restli: Illegal constant for %q enum: %d", i.Enum, i.Constant)
}

type UnknownEnumValue struct {
	Enum  string
	Value string
}

func (u *UnknownEnumValue) Error() string {
	return fmt.Sprintf("go-restli: Unknown enum value for %q: %q", u.Enum, u.Value)
}
