package restlicodec

import "io"

type Marshaler interface {
	MarshalRestLi(Writer) error
}

type Closer interface {
	Finalize() string
	ReadCloser() io.ReadCloser
}

type PrimitiveWriter interface {
	WriteInt32(v int32)
	WriteInt64(v int64)
	WriteFloat32(v float32)
	WriteFloat64(v float64)
	WriteBool(v bool)
	WriteString(v string)
	WriteBytes(v []byte)
}

type WriteCloser interface {
	Writer
	Closer
}

type Writer interface {
	PrimitiveWriter
	WriteMap(mapWriter func(keyWriter func(key string) Writer) error) error
	WriteArray(arrayWriter func(itemWriter func() Writer) error) error
}

type rawWriter interface {
	PrimitiveWriter

	writeMapStart()
	writeKey(key string)
	writeKeyDelimiter()
	writeEntryDelimiter()
	writeMapEnd()

	writeArrayStart()
	writeArrayItemDelimiter()
	writeArrayEnd()

	RawByte(byte)
	RawString(string)

	BuildBytes(reuse ...[]byte) ([]byte, error)
	ReadCloser() (io.ReadCloser, error)
}

type genericWriter struct {
	rawWriter
}

func newGenericWriter(raw rawWriter) *genericWriter {
	return &genericWriter{rawWriter: raw}
}

func (e *genericWriter) WriteMap(mapEncoder func(keyWriter func(key string) Writer) error) (err error) {
	e.rawWriter.writeMapStart()

	first := true
	err = mapEncoder(func(key string) Writer {
		if first {
			first = false
		} else {
			e.rawWriter.writeEntryDelimiter()
		}
		e.rawWriter.writeKey(key)
		e.rawWriter.writeKeyDelimiter()
		return e
	})
	if err != nil {
		return err
	}

	e.rawWriter.writeMapEnd()
	return nil
}

func (e *genericWriter) WriteArray(arrayWriter func(itemWriter func() Writer) error) (err error) {
	e.rawWriter.writeArrayStart()

	first := true
	err = arrayWriter(func() Writer {
		if first {
			first = false
		} else {
			e.rawWriter.writeArrayItemDelimiter()
		}
		return e
	})
	if err != nil {
		return err
	}

	e.rawWriter.writeArrayEnd()
	return nil
}

func (e *genericWriter) Finalize() string {
	data, _ := e.rawWriter.BuildBytes()
	return string(data)
}

func (e *genericWriter) ReadCloser() io.ReadCloser {
	rc, _ := e.rawWriter.ReadCloser()
	return rc
}
