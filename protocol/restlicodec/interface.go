package restlicodec

import "io"

type Encodable interface {
	RestLiEncode(*Encoder) error
}

type (
	MapEncoder   func(key string, value Encodable) error
	ArrayEncoder func(index int, item Encodable) error
)

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

func (e *Encoder) writeField(name string) {
	e.encoder.WriteFieldName(name)
	e.encoder.WriteFieldNameDelimiter()
}

func (e *Encoder) encode(v Encodable) error {
	return v.RestLiEncode(&Encoder{encoder: e.encoder.SubEncoder()})
}

func (e *Encoder) Field(fieldName string, fieldValue Encodable) error {
	e.writeField(fieldName)
	return e.encode(fieldValue)
}

func (e *Encoder) MapField(fieldName string, encoderFunc func(MapEncoder) error) error {
	e.writeField(fieldName)
	e.encoder.WriteMapStart()
	first := true
	err := encoderFunc(func(key string, value Encodable) error {
		if first {
			first = false
		} else {
			e.encoder.WriteMapEntryDelimiter()
		}
		e.encoder.WriteMapKey(key)
		e.encoder.WriteMapKeyDelimiter()
		return e.encode(value)
	})
	if err != nil {
		return err
	}
	e.encoder.WriteMapEntryDelimiter()
	return nil
}

func (e *Encoder) ArrayField(fieldName string, encoderFunc func(ArrayEncoder) error) error {
	e.writeField(fieldName)
	e.encoder.WriteArrayStart()
	err := encoderFunc(func(i int, v Encodable) error {
		if i > 0 {
			e.encoder.WriteArrayItemDelimiter()
		}
		return e.encode(v)
	})
	if err != nil {
		return err
	}
	e.encoder.WriteArrayEnd()
	return nil
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
