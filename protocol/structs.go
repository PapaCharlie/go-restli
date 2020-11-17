package protocol

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/fnv1a"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

type RestLiObject interface {
	restlicodec.Marshaler
	restlicodec.Unmarshaler
	ComputeHash() fnv1a.Hash
	EqualsInterface(interface{}) bool
}

type partialUpdateRequest struct {
	Patch restlicodec.Marshaler
}

func (p *partialUpdateRequest) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(fieldNameWriter func(fieldName string) restlicodec.Writer) error {
		return p.Patch.MarshalRestLi(fieldNameWriter("patch").SetScope())
	})
}

var EmptyRecord emptyRecord

type emptyRecord struct{}

func (e emptyRecord) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadMap(func(restlicodec.Reader, string) error { return reader.Skip() })
}

func (e emptyRecord) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(func(string) restlicodec.Writer) error { return nil })
}

type batchGetRequestResponse struct {
	Results restlicodec.MapReader
	Errors  *BatchMethodError
}

func (b *batchGetRequestResponse) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	resultsReceived := false
	err = reader.ReadMap(func(reader restlicodec.Reader, field string) (err error) {
		switch field {
		case "results":
			resultsReceived = true
			err = reader.ReadMap(b.Results)
		case "errors":
			b.Errors.Errors, err = reader.ReadRawBytes()
		case "statuses":
			b.Errors.Statuses, err = reader.ReadRawBytes()
		default:
			err = reader.Skip()
		}
		return err
	})
	if err != nil {
		return err
	}
	if !resultsReceived {
		return fmt.Errorf("no results received")
	}
	if len(b.Errors.Errors) > 0 {
		return restlicodec.NewJsonReader(b.Errors.Errors).ReadMap(func(restlicodec.Reader, string) error {
			// if the MapReader is getting called then we received a non-empty and non-null "errors" map, so immediately
			// return the raw errors map as the error itself
			return b.Errors
		})
	}
	return nil
}
