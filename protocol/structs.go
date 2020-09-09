package protocol

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

type partialUpdateRequest struct {
	Patch restlicodec.Marshaler
}

func (p *partialUpdateRequest) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(fieldNameWriter func(fieldName string) restlicodec.Writer) error {
		return p.Patch.MarshalRestLi(fieldNameWriter("patch").WithoutLastKey())
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
	Errors  []byte
}

func (b *batchGetRequestResponse) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	resultsReceived := false
	err = reader.ReadMap(func(reader restlicodec.Reader, field string) (err error) {
		switch field {
		case "results":
			resultsReceived = true
			err = reader.ReadMap(b.Results)
		case "errors":
			b.Errors, err = reader.Raw()
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
	if len(b.Errors) > 0 {
		return restlicodec.NewJsonReader(b.Errors).ReadMap(func(restlicodec.Reader, string) error {
			// if the MapReader is getting called then we received a non-empty "errors" map, so immediately return the
			// raw errors map as the error itself
			return fmt.Errorf(string(b.Errors))
		})
	}
	return nil
}
