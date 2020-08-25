package restlicodec

type Unmarshaler interface {
	UnmarshalRestLi(Reader) error
}

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
	ReadMap(mapReader func(field string) error) error
	ReadArray(arrayReader func() error) error

	Skip() error
}
