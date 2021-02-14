package restlicodec

import (
	"io"
)

// Marshaler is the interface that should be implemented by objects that can be serialized to JSON and ROR2
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
	// Size returns the size of the data that was written out. Note that calling Size after Finalize or ReadCloser will
	// return 0
	Size() int
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
	// WriteRawBytes appends the given bytes to the underlying buffer, without validating the input. Use at your own
	// risk!
	WriteRawBytes([]byte)
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
	// IsKeyExcluded checks whether or not the given key or field name at the current scope should be included in
	// serialization. Exclusion behavior is already built into the writer but this is intended for Marhsalers that use
	// custom serialization logic on top of the Writer
	IsKeyExcluded(key string) bool
	// SetScope returns a copy of the current writer with the given scope. For internal use only by Marshalers that use
	// custom serialization logic on top of the Writer. Designed to work around the backing PathSpec for field
	// exclusion.
	SetScope(...string) Writer
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

	// The following are exposed directly by jwriter.Writer
	RawByte(byte)
	Raw([]byte, error)
	RawString(string)
	BuildBytes(...[]byte) ([]byte, error)
	ReadCloser() (io.ReadCloser, error)
	Size() int
}

type genericWriter struct {
	excludedFields PathSpec
	scope          []string
	rawWriter
}

func newGenericWriter(raw rawWriter, excludedFields PathSpec) *genericWriter {
	return &genericWriter{rawWriter: raw, excludedFields: excludedFields}
}

func (e *genericWriter) WriteRawBytes(data []byte) {
	e.rawWriter.Raw(data, nil)
}

func (e *genericWriter) WriteMap(mapWriter MapWriter) (err error) {
	e.rawWriter.writeMapStart()

	first := true
	sub := e.subWriter("")

	err = mapWriter(func(key string) Writer {
		var writer Writer
		if !e.IsKeyExcluded(key) {
			if first {
				first = false
			} else {
				e.rawWriter.writeEntryDelimiter()
			}
			e.rawWriter.writeKey(key)
			e.rawWriter.writeKeyDelimiter()
			writer = sub
			sub.scope[len(sub.scope)-1] = key
		} else {
			writer = noopWriter
		}
		return writer
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
	sub := e.subWriter(WildCard)
	err = arrayWriter(func() Writer {
		if first {
			first = false
		} else {
			e.rawWriter.writeArrayItemDelimiter()
		}
		return sub
	})
	if err != nil {
		return err
	}

	e.rawWriter.writeArrayEnd()
	return nil
}

func (e *genericWriter) IsKeyExcluded(key string) bool {
	e.scope = append(e.scope, key)
	excluded := e.excludedFields.Matches(e.scope)
	e.scope = e.scope[:len(e.scope)-1]
	return excluded
}

func (e *genericWriter) SetScope(scope ...string) Writer {
	var out genericWriter
	out = *e
	out.scope = scope
	return &out
}

func (e *genericWriter) Finalize() string {
	data, _ := e.rawWriter.BuildBytes()
	return string(data)
}

func (e *genericWriter) ReadCloser() io.ReadCloser {
	rc, _ := e.rawWriter.ReadCloser()
	return rc
}

func (e *genericWriter) subWriter(key string) *genericWriter {
	var out genericWriter
	out = *e
	out.scope = copyAndAppend(e.scope, key)
	return &out
}

func copyAndAppend(a []string, v string) (out []string) {
	out = append(out, a...)
	out = append(out, v)
	return out
}
