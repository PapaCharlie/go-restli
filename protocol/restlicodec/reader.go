package restlicodec

import (
	"fmt"
)

// Unmarshaler is the interface that should be implemented by objects that can be deserialized from JSON and ROR2
type Unmarshaler interface {
	UnmarshalRestLi(Reader) error
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

type Reader interface {
	fmt.Stringer
	PrimitiveReader
	// ReadMap tells the Reader that it should expect a map/object as its next input. If it is not (e.g. it is an array
	// or a primitive) it will return an error.
	// Note that not using the inner Reader passed to the MapReader may result in undefined behavior.
	ReadMap(mapReader MapReader) error
	// ReadRecord tells the Reader that it should expect an object as its next input and calls recordReader for each
	// field of the object. If the next input is not an object, it will return an error.
	// Note that not using the inner Reader passed to the MapReader may result in undefined behavior.
	ReadRecord(requiredFields RequiredFields, recordReader MapReader) error
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

	AtInputStart() bool
	RecordMissingRequiredFields(missingRequiredFields map[string]struct{})
	CheckMissingFields() error
}

func readRecord(reader rawReader, requiredFields RequiredFields, mapReader MapReader) (err error) {
	atInputStart := reader.AtInputStart()
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

	reader.RecordMissingRequiredFields(requiredFieldsRemaining)

	if atInputStart {
		return reader.CheckMissingFields()
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

type RequiredFields []string

func (rf RequiredFields) toMap() map[string]struct{} {
	fields := make(map[string]struct{}, len(rf))
	for _, f := range rf {
		fields[f] = struct{}{}
	}
	return fields
}
