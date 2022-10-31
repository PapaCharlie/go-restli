package restlicodec

import (
	"errors"
	"fmt"
)

// Unmarshaler is the interface that should be implemented by objects that can be deserialized from JSON and ROR2
type Unmarshaler interface {
	UnmarshalRestLi(Reader) error
}

// PointerUnmarshaler represents an interface implemented by records and other objects that use pointer receivers for
// all their methods (unlike enums that use direct receivers).
type PointerUnmarshaler[T any] interface {
	Unmarshaler
	NewInstance() T
}

// UnmarshalRestLi calls the corresponding PrimitiveReader method if T is a Primitive (or an int). If T implements
// Unmarshaler, it is expected to implement PointerUnmarshaler as well and its NewInstance method will be called, then
// Unmarshaler.UnmarshalRestLi is called on the new pointer. If *T implements Unmarshaler then
// Unmarshaler.UnmarshalRestLi is called directly on a pointer to a 0-value of T. Otherwise, this function panics.
func UnmarshalRestLi[T any](reader Reader) (t T, err error) {
	if pt, ok := any(t).(PointerUnmarshaler[T]); ok {
		t = pt.NewInstance()
		return t, any(t).(Unmarshaler).UnmarshalRestLi(reader)
	}
	if u, ok := any(&t).(Unmarshaler); ok {
		return t, u.UnmarshalRestLi(reader)
	}
	cast := func(v any, err error) (T, error) {
		return v.(T), err
	}
	switch any(t).(type) {
	case int:
		return cast(reader.ReadInt())
	case int32:
		return cast(reader.ReadInt32())
	case int64:
		return cast(reader.ReadInt64())
	case float32:
		return cast(reader.ReadFloat32())
	case float64:
		return cast(reader.ReadFloat64())
	case bool:
		return cast(reader.ReadBool())
	case string:
		return cast(reader.ReadString())
	case []byte:
		return cast(reader.ReadBytes())
	default:
		return loadAdapter[T]().unmarshaler(reader)
	}
}

// The UnmarshalerFunc type is an adapter to allow the use of ordinary functions as unmarshalers, useful for inlining
// marshalers instead of defining new types
type UnmarshalerFunc func(Reader) error

func (u UnmarshalerFunc) UnmarshalRestLi(reader Reader) error {
	return u(reader)
}

// PrimitiveReader describes the set of functions that read the supported rest.li primitives from the input. Note that
// if the reader's next input is not a primitive (i.e. it is an object/map or an array), each of these methods will
// return errors. The encoding spec can be found here:
// https://linkedin.github.io/rest.li/how_data_is_serialized_for_transport
type PrimitiveReader interface {
	ReadInt() (int, error)
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)
	ReadFloat32() (float32, error)
	ReadFloat64() (float64, error)
	ReadBool() (bool, error)
	ReadString() (string, error)
	ReadBytes() ([]byte, error)
}

type (
	GenericUnmarshaler[T any] func(reader Reader) (T, error)
	MapReader                 func(reader Reader, field string) (err error)
	ArrayReader               func(reader Reader) (err error)
)

// NoSuchFieldErr should be returned to signal that
var NoSuchFieldErr = errors.New("go-restli: No such field")

type Reader interface {
	fmt.Stringer
	PrimitiveReader
	KeyChecker
	// ReadMap tells the Reader that it should expect a map/object as its next input. If it is not (e.g. it is an array
	// or a primitive) it will return an error.
	// Note that not using the inner Reader passed to the MapReader may result in undefined behavior.
	ReadMap(mapReader MapReader) error
	// ReadRecord tells the Reader that it should expect an object as its next input and calls recordReader for each
	// field of the object. If the next input is not an object, it will return an error.
	// Note that not using the inner Reader passed to the MapReader may result in undefined behavior.
	ReadRecord(requiredFields *RequiredFields, recordReader MapReader) error
	// ReadArray tells the reader that it should expect an array as its next input. If it is not, it will return an
	// error
	// Note that not using the inner Reader passed to the ArrayReader may result in undefined behavior.
	ReadArray(arrayReader ArrayReader) error
	// ReadInterface reads an interface{} analogous to the 'encoding/json' package. It is a best-effort attempt to
	// deserialize the underlying data into map[string]interface{}, []interface{} or raw primitive types accordingly.
	// Note that for ROR2, because all primitives are encoded as strings, it is impossible to tell what the field's type
	// is intended to be without its schema. Therefore all primitive values are interpreted as strings
	ReadInterface() (interface{}, error)
	// ReadRawBytes returns the next primitive/array/map as a raw, unvalidated byte slice.
	ReadRawBytes() ([]byte, error)

	// Skip skips the next primitive/array/map completely.
	Skip() error
}

type rawReader interface {
	ReadMap(mapReader MapReader) error

	atInputStart() bool
	recordMissingRequiredFields(missingRequiredFields map[string]struct{})
	checkMissingFields() error
}

func readRecord(reader rawReader, requiredFields *RequiredFields, mapReader MapReader) (err error) {
	atInputStart := reader.atInputStart()
	requiredFieldsRemaining := requiredFields.toMap()

	err = reader.ReadMap(func(reader Reader, field string) (err error) {
		err = mapReader(reader, field)
		if err != nil {
			return err
		}

		delete(requiredFieldsRemaining, field)
		return nil
	})
	if err != nil {
		return err
	}

	reader.recordMissingRequiredFields(requiredFieldsRemaining)

	if atInputStart {
		return reader.checkMissingFields()
	} else {
		return nil
	}
}

func ReadMap[V any](reader Reader, unmarshaler GenericUnmarshaler[V]) (result map[string]V, err error) {
	result = make(map[string]V)
	err = reader.ReadMap(func(reader Reader, field string) (err error) {
		result[field], err = unmarshaler(reader)
		return err
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func ReadArray[V any](reader Reader, unmarshaler GenericUnmarshaler[V]) (result []V, err error) {
	err = reader.ReadArray(func(reader Reader) (err error) {
		item, err := unmarshaler(reader)
		if err != nil {
			return err
		}
		result = append(result, item)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

type DeserializationError struct {
	Scope string
	Err   error
}

func (d *DeserializationError) Error() string {
	return fmt.Sprintf("go-restli: Failed to deserialize %q (%+v)", d.Scope, d.Err)
}

type RequiredFields struct {
	fields []string
}

func NewRequiredFields(included ...*RequiredFields) (rf *RequiredFields) {
	rf = new(RequiredFields)
	for _, i := range included {
		rf.fields = append(rf.fields, i.fields...)
	}
	return rf
}

func (rf *RequiredFields) Add(fields ...string) *RequiredFields {
	rf.fields = append(rf.fields, fields...)
	return rf
}

func (rf *RequiredFields) toMap() map[string]struct{} {
	if rf == nil || len(rf.fields) == 0 {
		return nil
	}
	fields := make(map[string]struct{}, len(rf.fields))
	for _, f := range rf.fields {
		fields[f] = struct{}{}
	}
	return fields
}

func readBytes(s string, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}

	if s == "" {
		return nil, nil
	} else {
		return []byte(s), nil
	}
}
