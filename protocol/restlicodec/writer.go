package restlicodec

import "io"

// Unmarshaler is the interface that should be implemented by objects that can be serialized to JSON and ROR2
type Marshaler interface {
	MarshalRestLi(Writer) error
}

// Closer provides methods for the underlying Writer to release its buffer and return the constructed objects to the
// caller. Note that once either Finalize or ReadCloser is called, subsequent calls to either will return the empty
// string or the empty reader respectively.
type Closer interface {
	// Finalize returns the created object as a string and releases the pooled underlying buffer. Subsequent calls will
	// return the empty string
	Finalize() string
	// ReadCloser returns an io.ReadCloser that will read the pooled underlying buffer, releasing each byte back into
	// the pool as they are read. Close will release any remaining unread bytes to the pool.
	ReadCloser() io.ReadCloser
}

// PrimitiveWriter provides the set of functions needed to write the supported rest.li primitives to the backing buffer,
// according to the rest.li serialization spec: https://linkedin.github.io/rest.li/how_data_is_serialized_for_transport.
type PrimitiveWriter interface {
	WriteInt32(v int32)
	WriteInt64(v int64)
	WriteFloat32(v float32)
	WriteFloat64(v float64)
	WriteBool(v bool)
	WriteString(v string)
	WriteBytes(v []byte)
}

// WriteCloser groups the Writer and Closer functions.
type WriteCloser interface {
	Writer
	Closer
}

type (
	ArrayWriter func(itemWriter func() Writer) error
	MapWriter   func(keyWriter func(key string) Writer) error
)

// Writer is the interface implemented by all serialization mechanisms supported by rest.li. See the New*Writer
// functions provided in package for all the supported serialization mechanisms.
type Writer interface {
	PrimitiveWriter
	// WriteMap writes the map keys/object fields written by the given lambda between object delimiters. The lambda
	// takes a function that is used to write the key/field name into the object and returns a nested Writer. This
	// Writer should be used to write inner fields. Take the following JSON object:
	//   {
	//     "foo": "bar",
	//     "baz": 42
	//   }
	// This would be written as follows using a Writer:
	//		err := writer.WriteMap(func(keyWriter func(string) restlicodec.Writer) error {
	//			keyWriter("foo").WriteString("bar")
	//			keyWriter("baz").WriteInt32(42)
	//			return nil
	//		}
	// Note that not using the inner Writer returned by the keyWriter may result in undefined behavior.
	WriteMap(mapWriter MapWriter) error
	// WriteArray writes the array items written by the given lambda between array delimiters. The lambda
	// takes a function that is used to signal that a new item is starting and returns a nested Writer. This
	// Writer should be used to write inner fields. Take the following JSON object:
	//   [
	//     "foo",
	//     "bar"
	//   ]
	// This would be written as follows using a Writer:
	//		err := writer.WriteArray(func(itemWriter func() restlicodec.Writer) error {
	//			itemWriter().WriteString("foo")
	//			itemWriter().WriteString("bar")
	//			return nil
	//		}
	// Note that not using the inner Writer returned by the itemWriter may result in undefined behavior.
	WriteArray(arrayWriter ArrayWriter) error
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

func (e *genericWriter) WriteMap(mapWriter MapWriter) (err error) {
	e.rawWriter.writeMapStart()

	first := true
	err = mapWriter(func(key string) Writer {
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

func (e *genericWriter) WriteArray(arrayWriter ArrayWriter) (err error) {
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
