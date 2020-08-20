package restlicodec

// Unmarshaler is the interface that should be implemented by objects that can be deserialized from JSON and ROR2
type Unmarshaler interface {
	UnmarshalRestLi(Reader) error
}

// PrimitiveReader describes the set of functions that read the supported rest.li primitives from the input. Note that
// if the reader's next input is not a primitive (i.e. it is an object/map or an array), each of these methods will
// return errors. The encoding spec can be found here:
// https://linkedin.github.io/rest.li/how_data_is_serialized_for_transport
type PrimitiveReader interface {
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)
	ReadFloat32() (float32, error)
	ReadFloat64() (float64, error)
	ReadBool() (bool, error)
	ReadString() (string, error)
	ReadBytes() ([]byte, error)
}

type Reader interface {
	PrimitiveReader
	// ReadMap tells the Reader that it should expect a map/object as its next input. If it is not (e.g. it is an array
	// or a primitive) it will return an error
	ReadMap(mapReader func(field string) error) error
	// ReadArray tells the reader that it should expect an array as its next input. If it is not, it will return an error
	ReadArray(arrayReader func() error) error

	Skip() error
}
