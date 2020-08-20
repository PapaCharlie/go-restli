package restlicodec

type primitiveReader struct {
	read func(Reader) (err error)
}

func (p *primitiveReader) UnmarshalRestLi(reader Reader) error {
	return p.read(reader)
}

// NewInt32PrimitiveUnmarshaler wraps the given int32 pointer in an Unmarshaler which will set the pointer's value to
// the value read from the Reader given to UnmarshalRestLi.
// Note that if the given pointer is nil, it will cause a panic.
func NewInt32PrimitiveUnmarshaler(v *int32) Unmarshaler {
	return &primitiveReader{read: func(reader Reader) (err error) {
		*v, err = reader.ReadInt32()
		return err
	}}
}

// NewInt64PrimitiveUnmarshaler wraps the given int64 pointer in an Unmarshaler which will set the pointer's value to
// the value read from the Reader given to UnmarshalRestLi.
// Note that if the given pointer is nil, it will cause a panic.
func NewInt64PrimitiveUnmarshaler(v *int64) Unmarshaler {
	return &primitiveReader{read: func(reader Reader) (err error) {
		*v, err = reader.ReadInt64()
		return err
	}}
}

// NewFloat32PrimitiveUnmarshaler wraps the given float32 pointer in an Unmarshaler which will set the pointer's value
// to the value read from the Reader given to UnmarshalRestLi.
// Note that if the given pointer is nil, it will cause a panic.
func NewFloat32PrimitiveUnmarshaler(v *float32) Unmarshaler {
	return &primitiveReader{read: func(reader Reader) (err error) {
		*v, err = reader.ReadFloat32()
		return err
	}}
}

// NewFloat64PrimitiveUnmarshaler wraps the given float64 pointer in an Unmarshaler which will set the pointer's value
// to the value read from the Reader given to UnmarshalRestLi.
// Note that if the given pointer is nil, it will cause a panic.
func NewFloat64PrimitiveUnmarshaler(v *float64) Unmarshaler {
	return &primitiveReader{read: func(reader Reader) (err error) {
		*v, err = reader.ReadFloat64()
		return err
	}}
}

// NewBoolPrimitiveUnmarshaler wraps the given bool pointer in an Unmarshaler which will set the pointer's value to the
// value read from the Reader given to UnmarshalRestLi.
// Note that if the given pointer is nil, it will cause a panic.
func NewBoolPrimitiveUnmarshaler(v *bool) Unmarshaler {
	return &primitiveReader{read: func(reader Reader) (err error) {
		*v, err = reader.ReadBool()
		return err
	}}
}

// NewStringPrimitiveUnmarshaler wraps the given string pointer in an Unmarshaler which will set the pointer's value to
// the value read from the Reader given to UnmarshalRestLi.
// Note that if the given pointer is nil, it will cause a panic.
func NewStringPrimitiveUnmarshaler(v *string) Unmarshaler {
	return &primitiveReader{read: func(reader Reader) (err error) {
		*v, err = reader.ReadString()
		return err
	}}
}

// NewBytesPrimitiveUnmarshaler wraps the given byte slice pointer in an Unmarshaler which will set the pointer's value
// to the value read from the Reader given to UnmarshalRestLi.
// Note that if the given pointer is nil, it will cause a panic.
func NewBytesPrimitiveUnmarshaler(v *[]byte) Unmarshaler {
	return &primitiveReader{read: func(reader Reader) (err error) {
		*v, err = reader.ReadBytes()
		return err
	}}
}
