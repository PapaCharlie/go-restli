package restlicodec

import (
	"io"
	"log"
	"reflect"
	"sort"
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
		log.Panicf("Unknown primitive type: %s", reflect.TypeOf(v))
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
	empty := true
	sub := e.subWriter("")

	err = mapWriter(func(key string) Writer {
		var writer Writer
		if !e.IsKeyExcluded(key) {
			if empty {
				e.rawWriter.writeMapStart()
				empty = false
			} else {
				e.rawWriter.writeEntryDelimiter()
			}
			e.rawWriter.writeKey(key)
			e.rawWriter.writeKeyDelimiter()
			writer = sub
			sub.scope[len(sub.scope)-1] = key
		} else {
			writer = NoopWriter
		}
		return writer
	})
	if err != nil {
		return err
	}

	if empty {
		e.rawWriter.writeEmptyMap()
	} else {
		e.rawWriter.writeMapEnd()
	}
	return nil
}

func (e *genericWriter) WriteArray(arrayWriter ArrayWriter) (err error) {
	empty := true
	sub := e.subWriter(WildCard)
	err = arrayWriter(func() Writer {
		if empty {
			e.rawWriter.writeArrayStart()
			empty = false
		} else {
			e.rawWriter.writeArrayItemDelimiter()
		}
		return sub
	})
	if err != nil {
		return err
	}

	if empty {
		e.rawWriter.writeEmptyArray()
	} else {
		e.rawWriter.writeArrayEnd()
	}
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
	out = make([]string, 0, len(out)+1)
	out = append(out, a...)
	out = append(out, v)
	return out
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
		if len(entries) == 0 {
			return nil
		}
		sortedEntries := make([]struct {
			key   string
			value V
		}, len(entries))
		i := 0
		for k, v := range entries {
			sortedEntries[i].key, err = keyMarshaler(k)
			if err != nil {
				return err
			}
			sortedEntries[i].value = v
			i++
		}
		sort.Slice(sortedEntries, func(i, j int) bool {
			return sortedEntries[i].key < sortedEntries[j].key
		})

		for _, e := range sortedEntries {
			err = valueMarshaler(e.value, keyWriter(e.key))
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
