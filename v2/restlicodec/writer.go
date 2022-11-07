package restlicodec

import (
	"io"
	"sort"

	"github.com/mailru/easyjson/jwriter"
)

// Marshaler is the interface that should be implemented by objects that can be serialized to JSON and ROR2
type Marshaler interface {
	MarshalRestLi(Writer) error
}

// MarshalRestLi calls the corresponding PrimitiveWriter method if T is a Primitive (or an int), and directly calls
// Marshaler.MarshalRestLi if T is a Marshaler. Otherwise, this function panics.
func MarshalRestLi[T any](t T, writer Writer) (err error) {
	switch v := any(t).(type) {
	case Marshaler:
		return v.MarshalRestLi(writer)
	case int:
		writer.WriteInt(v)
	case int32:
		writer.WriteInt32(v)
	case int64:
		writer.WriteInt64(v)
	case float32:
		writer.WriteFloat32(v)
	case float64:
		writer.WriteFloat64(v)
	case bool:
		writer.WriteBool(v)
	case string:
		writer.WriteString(v)
	case []byte:
		writer.WriteBytes(v)
	default:
		return loadAdapter[T]().marshaler(t, writer)
	}
	return nil
}

// PrimitiveWriter provides the set of functions needed to write the supported rest.li primitives to the backing buffer,
// according to the rest.li serialization spec: https://linkedin.github.io/rest.li/how_data_is_serialized_for_transport.
type PrimitiveWriter interface {
	WriteInt(v int)
	WriteInt32(v int32)
	WriteInt64(v int64)
	WriteFloat32(v float32)
	WriteFloat64(v float64)
	WriteBool(v bool)
	WriteString(v string)
	WriteBytes(v []byte)
}

type (
	// The MarshalerFunc type is an adapter to allow the use of ordinary functions as marshalers, useful for inlining
	// marshalers instead of defining new types
	MarshalerFunc                   func(Writer) error
	PrimitiveMarshaler[T Primitive] func(writer Writer, t T)
	GenericMarshaler[T any]         func(t T, writer Writer) error
	ArrayWriter                     func(itemWriter func() Writer) (err error)
	MapWriter                       func(keyWriter func(key string) Writer) (err error)
)

func (m MarshalerFunc) MarshalRestLi(writer Writer) error {
	return m(writer)
}

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
	// Note that not using the inner Writer returned by the keyWriter will result in undefined behavior. Additionally,
	// reusing a Writer returned by a previous call to keyWriter will also result in undefined behavior.
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
	// Note that not using the inner Writer returned by the itemWriter may result in undefined behavior. Additionally,
	// reusing a Writer returned by a previous call to itemWriter will also result in undefined behavior.
	WriteArray(arrayWriter ArrayWriter) error
	// IsKeyExcluded checks whether the given key or field name at the current scope should be included in
	// serialization. Exclusion behavior is already built into the writer but this is intended for Marshalers that use
	// custom serialization logic on top of the Writer
	IsKeyExcluded(key string) bool
	// SetScope returns a copy of the current writer with the given scope. For internal use only by Marshalers that use
	// custom serialization logic on top of the Writer. Designed to work around the backing PathSpec for field
	// exclusion.
	SetScope(...string) Writer
	// Finalize returns the created object as a string and releases the pooled underlying buffer. Subsequent calls will
	// return the empty string
	Finalize() string
}

type rawWriter interface {
	PrimitiveWriter

	writeMapStart()
	writeKey(key string)
	writeKeyDelimiter()
	writeEntryDelimiter()
	writeMapEnd()
	writeEmptyMap()

	writeArrayStart()
	writeArrayItemDelimiter()
	writeArrayEnd()
	writeEmptyArray()

	getWriter() *jwriter.Writer
	setWriter(w *jwriter.Writer)

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

func (gw *genericWriter) Write(p []byte) (n int, err error) {
	gw.getWriter().Buffer.AppendBytes(p)
	return len(p), nil
}

func (gw *genericWriter) WriteRawBytes(data []byte) {
	gw.Raw(data, nil)
}

func (gw *genericWriter) enterScope(scope string) {
	gw.scope = append(gw.scope, scope)
}

func (gw *genericWriter) exitScope() {
	gw.scope = gw.scope[:len(gw.scope)-1]
}

func (gw *genericWriter) WriteMap(mapWriter MapWriter) (err error) {
	type entry struct {
		key    string
		writer *jwriter.Writer
	}
	var entries []entry

	buf := gw.getWriter()

	started := false
	err = mapWriter(func(key string) Writer {
		if started {
			// If started is true then a scope was entered and needs to be cleared
			gw.exitScope()
		} else {
			started = true
		}
		gw.enterScope(key)
		if gw.excludedFields.Matches(gw.scope) {
			return NoopWriter
		}

		if len(entries) == 0 {
			gw.writeMapStart()
		}

		e := entry{
			key:    key,
			writer: new(jwriter.Writer),
		}
		entries = append(entries, e)
		gw.setWriter(e.writer)
		return gw
	})
	if err != nil {
		return err
	}
	gw.setWriter(buf)

	if started {
		// If started is true then a scope was entered and needs to be cleared
		gw.exitScope()
	}

	if len(entries) == 0 {
		gw.writeEmptyMap()
		return nil
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].key < entries[j].key })
	for i, e := range entries {
		if i != 0 {
			gw.writeEntryDelimiter()
		}
		gw.writeKey(e.key)
		gw.writeKeyDelimiter()
		_, _ = e.writer.DumpTo(gw)
	}
	gw.writeMapEnd()
	return nil
}

func (gw *genericWriter) WriteArray(arrayWriter ArrayWriter) (err error) {
	empty := true
	gw.enterScope(WildCard)
	err = arrayWriter(func() Writer {
		if empty {
			gw.writeArrayStart()
			empty = false
		} else {
			gw.writeArrayItemDelimiter()
		}
		return gw
	})
	if err != nil {
		return err
	}
	gw.exitScope()

	if empty {
		gw.writeEmptyArray()
	} else {
		gw.writeArrayEnd()
	}
	return nil
}

func (gw *genericWriter) IsKeyExcluded(key string) bool {
	gw.enterScope(key)
	defer gw.exitScope()
	return gw.excludedFields.Matches(gw.scope)
}

func (gw *genericWriter) SetScope(scope ...string) Writer {
	var out genericWriter
	out = *gw
	out.scope = scope
	return &out
}

func (gw *genericWriter) Finalize() string {
	data, _ := gw.BuildBytes()
	return string(data)
}

func (gw *genericWriter) ReadCloser() io.ReadCloser {
	rc, _ := gw.rawWriter.ReadCloser()
	return rc
}

type ComparablePrimitive interface {
	int32 | int64 | float32 | float64 | bool | string
}

type Primitive interface {
	ComparablePrimitive | []byte
}

func WriteArray[T any](writer Writer, array []T, marshaler GenericMarshaler[T]) (err error) {
	return writer.WriteArray(func(itemWriter func() Writer) (err error) {
		for _, v := range array {
			err = marshaler(v, itemWriter())
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func WriteMap[V any](writer Writer, entries map[string]V, marshaler GenericMarshaler[V]) (err error) {
	return WriteGenericMap(writer, entries, func(s string) (string, error) { return s, nil }, marshaler)
}

func WriteGenericMap[K comparable, V any](
	writer Writer,
	entries map[K]V,
	keyMarshaler func(K) (string, error),
	valueMarshaler GenericMarshaler[V],
) (err error) {
	return writer.WriteMap(func(keyWriter func(key string) Writer) (err error) {
		var key string
		for k, v := range entries {
			key, err = keyMarshaler(k)
			if err != nil {
				return err
			}
			err = valueMarshaler(v, keyWriter(key))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func WriteInt32(v int32, w Writer) error {
	w.WriteInt32(v)
	return nil
}

func WriteInt64(v int64, w Writer) error {
	w.WriteInt64(v)
	return nil
}

func WriteFloat32(v float32, w Writer) error {
	w.WriteFloat32(v)
	return nil
}

func WriteFloat64(v float64, w Writer) error {
	w.WriteFloat64(v)
	return nil
}

func WriteBool(v bool, w Writer) error {
	w.WriteBool(v)
	return nil
}

func WriteString(v string, w Writer) error {
	w.WriteString(v)
	return nil
}

func WriteBytes(v []byte, w Writer) error {
	w.WriteBytes(v)
	return nil
}
