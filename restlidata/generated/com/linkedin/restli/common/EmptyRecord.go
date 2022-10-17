package common

import (
	"reflect"

	"github.com/PapaCharlie/go-restli/restlicodec"
)

type EmptyRecord struct{}

func (e EmptyRecord) NewInstance() EmptyRecord {
	return e
}

func (e EmptyRecord) DecodeQueryParams(restlicodec.QueryParamsReader) error {
	return nil
}

func (e EmptyRecord) UnmarshalRestLi(reader restlicodec.Reader) error {
	return reader.ReadMap(func(restlicodec.Reader, string) error { return reader.Skip() })
}

func (e EmptyRecord) MarshalRestLi(writer restlicodec.Writer) error {
	return writer.WriteMap(func(func(string) restlicodec.Writer) error { return nil })
}

func IsEmptyRecord[T any](t T) bool {
	return reflect.TypeOf(t) == reflect.TypeOf(EmptyRecord{})
}
