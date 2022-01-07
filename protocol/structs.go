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

var EmptyRecord emptyRecord

type emptyRecord struct{}

func (e emptyRecord) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadMap(func(restlicodec.Reader, string) error { return reader.Skip() })
}

func (e emptyRecord) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(func(string) restlicodec.Writer) error { return nil })
}

type BatchEntities map[string]restlicodec.Marshaler

func (b BatchEntities) Add(key restlicodec.Marshaler, value restlicodec.Marshaler) (err error) {
	writer := restlicodec.NewRor2HeaderWriter()
	err = key.MarshalRestLi(writer)
	if err != nil {
		return err
	}

	b[writer.Finalize()] = value
	return nil
}

func (b BatchEntities) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(keyWriter func(key string) restlicodec.Writer) error {
		return keyWriter("entities").WriteMap(func(keyWriter func(key string) restlicodec.Writer) error {
			for k, v := range b {
				err := keyWriter(k).WriteMap(func(keyWriter func(key string) restlicodec.Writer) error {
					return v.MarshalRestLi(keyWriter("patch"))
				})
				if err != nil {
					return err
				}
			}

			return nil
		})
	})
}

type BatchResultsReader func(keyReader restlicodec.Reader, valueReader restlicodec.Reader) (err error)

type batchRequestResponse struct {
	Results BatchResultsReader
	Errors  *BatchMethodError
}

func (b *batchRequestResponse) UnmarshalRestLi(reader restlicodec.Reader) (err error) {
	resultsReceived := false
	err = reader.ReadMap(func(reader restlicodec.Reader, field string) (err error) {
		switch field {
		case "results":
			resultsReceived = true
			err = reader.ReadMap(func(valueReader restlicodec.Reader, rawKey string) (err error) {
				keyReader, err := restlicodec.NewRor2Reader(rawKey)
				if err != nil {
					return err
				}

				return b.Results(keyReader, valueReader)
			})
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

type BatchEntityUpdateResponse struct {
	Status int
}

func (b *BatchEntityUpdateResponse) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadMap(func(reader restlicodec.Reader, field string) (err error) {
		switch field {
		case "status":
			b.Status, err = reader.ReadInt()
		default:
			err = reader.Skip()
		}
		return err
	})
}

func NewPagingContext(start int32, count int32) PagingContext {
	return PagingContext{
		Start: &start,
		Count: &count,
	}
}
