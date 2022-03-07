/*
DO NOT EDIT

Code automatically generated by github.com/PapaCharlie/go-restli
Source file: https://github.com/PapaCharlie/go-restli/blob/master/internal/codegen/resources/errorresponse.go
*/

package stdtypes

import (
	gofnv1a "github.com/PapaCharlie/go-fnv1a"
	equals "github.com/PapaCharlie/go-restli/protocol/equals"
	restlicodec "github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

type ErrorResponse struct {
	// The HTTP status code.
	Status *int32
	// A human-readable explanation of the error.
	Message *string
	// The FQCN of the exception thrown by the server.
	ExceptionClass *string
	// The full stack trace of the exception thrown by the server.
	StackTrace *string
}

func (e *ErrorResponse) Equals(other *ErrorResponse) bool {
	if e == other {
		return true
	}
	if e == nil || other == nil {
		return false
	}

	return equals.ComparablePointer(e.Status, other.Status) &&
		equals.ComparablePointer(e.Message, other.Message) &&
		equals.ComparablePointer(e.ExceptionClass, other.ExceptionClass) &&
		equals.ComparablePointer(e.StackTrace, other.StackTrace)
}

func (e *ErrorResponse) ComputeHash() gofnv1a.Hash {
	if e == nil {
		return gofnv1a.ZeroHash()
	}
	hash := gofnv1a.NewHash()

	if e.Status != nil {
		hash.AddInt32(*e.Status)
	}

	if e.Message != nil {
		hash.AddString(*e.Message)
	}

	if e.ExceptionClass != nil {
		hash.AddString(*e.ExceptionClass)
	}

	if e.StackTrace != nil {
		hash.AddString(*e.StackTrace)
	}

	return hash
}

func (e *ErrorResponse) MarshalRestLi(writer restlicodec.Writer) (err error) {
	return writer.WriteMap(func(keyWriter func(string) restlicodec.Writer) (err error) {
		if e.ExceptionClass != nil {
			keyWriter("exceptionClass").WriteString(*e.ExceptionClass)
		}
		if e.Message != nil {
			keyWriter("message").WriteString(*e.Message)
		}
		if e.StackTrace != nil {
			keyWriter("stackTrace").WriteString(*e.StackTrace)
		}
		if e.Status != nil {
			keyWriter("status").WriteInt32(*e.Status)
		}
		return nil
	})
}

func (e *ErrorResponse) MarshalJSON() (data []byte, err error) {
	writer := restlicodec.NewCompactJsonWriter()
	err = e.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}
	return []byte(writer.Finalize()), nil
}

func UnmarshalRestLiErrorResponse(reader restlicodec.Reader) (e *ErrorResponse, err error) {
	e = new(ErrorResponse)
	err = e.UnmarshalRestLi(reader)
	return e, err
}

func (e *ErrorResponse) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	err = reader.ReadRecord(nil, func(reader restlicodec.Reader, field string) (err error) {
		switch field {
		case "status":
			e.Status = new(int32)
			*e.Status, err = reader.ReadInt32()
		case "message":
			e.Message = new(string)
			*e.Message, err = reader.ReadString()
		case "exceptionClass":
			e.ExceptionClass = new(string)
			*e.ExceptionClass, err = reader.ReadString()
		case "stackTrace":
			e.StackTrace = new(string)
			*e.StackTrace, err = reader.ReadString()
		default:
			err = reader.Skip()
		}
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *ErrorResponse) UnmarshalJSON(data []byte) error {
	reader := restlicodec.NewJsonReader(data)
	return e.UnmarshalRestLi(reader)
}
