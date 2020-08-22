package restliencoding

import (
	"io"
)

type Encodable interface {
	RestLiEncode(*Encoder) error
}

//go:generate go run ./internal
type Encoder struct {
	encoder encoder
}

func (e *Encoder) WriteObjectStart() {
	e.encoder.WriteObjectStart()
}

func (e *Encoder) WriteFieldDelimiter() {
	e.encoder.WriteFieldDelimiter()
}

func (e *Encoder) WriteObjectEnd() {
	e.encoder.WriteObjectEnd()
}

func (e *Encoder) WriteFieldNameAndDelimiter(name string) {
	e.encoder.WriteFieldName(name)
	e.encoder.WriteFieldNameDelimiter()
}

func (e *Encoder) Map(mapEncoder func(keyWriter func(key string)) error) (err error) {
	e.encoder.WriteMapStart()

	first := true
	err = mapEncoder(func(key string) {
		if first {
			first = false
		} else {
			e.encoder.WriteMapEntryDelimiter()
		}
		e.encoder.WriteMapKey(key)
		e.encoder.WriteMapKeyDelimiter()
	})
	if err != nil {
		return err
	}

	e.encoder.WriteMapEnd()
	return nil
}

func (e *Encoder) Array(arrayEncoder func(indexWriter func(index int)) error) (err error) {
	e.encoder.WriteArrayStart()
	err = arrayEncoder(func(index int) {
		if index > 0 {
			e.encoder.WriteArrayItemDelimiter()
		}
	})
	if err != nil {
		return err
	}
	e.encoder.WriteArrayEnd()
	return nil
}

func (e *Encoder) Int32(v int32)     { e.encoder.Int32(v) }
func (e *Encoder) Int64(v int64)     { e.encoder.Int64(v) }
func (e *Encoder) Float32(v float32) { e.encoder.Float32(v) }
func (e *Encoder) Float64(v float64) { e.encoder.Float64(v) }
func (e *Encoder) Bool(v bool)       { e.encoder.Bool(v) }
func (e *Encoder) String(v string)   { e.encoder.String(v) }
func (e *Encoder) Bytes(v []byte)    { e.encoder.Bytes(v) }
func (e *Encoder) Encodable(v Encodable) error {
	return v.RestLiEncode(&Encoder{encoder: e.encoder.SubEncoder()})
}

func (e *Encoder) Finalize() string {
	return e.encoder.Finalize()
}

func (e *Encoder) ReadCloser() io.ReadCloser {
	return e.encoder.ReadCloser()
}

type encoder interface {
	WriteObjectStart()
	WriteFieldName(name string)
	WriteFieldNameDelimiter()
	WriteFieldDelimiter()
	WriteObjectEnd()

	WriteMapStart()
	WriteMapKey(key string)
	WriteMapKeyDelimiter()
	WriteMapEntryDelimiter()
	WriteMapEnd()

	WriteArrayStart()
	WriteArrayItemDelimiter()
	WriteArrayEnd()

	Int32(v int32)
	Int64(v int64)
	Float32(v float32)
	Float64(v float64)
	Bool(v bool)
	String(v string)
	Bytes(v []byte)

	SubEncoder() encoder

	Finalize() string
	ReadCloser() io.ReadCloser
}
