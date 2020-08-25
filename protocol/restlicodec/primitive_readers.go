package restlicodec

type primitiveReader struct {
	read func(Reader) (err error)
}

func (p *primitiveReader) UnmarshalRestLi(reader Reader) error {
	return p.read(reader)
}

func NewInt32PrimitiveUnmarshaler(v *int32) Unmarshaler {
	return &primitiveReader{read: func(reader Reader) (err error) {
		*v, err = reader.ReadInt32()
		return err
	}}
}

func NewInt64PrimitiveUnmarshaler(v *int64) Unmarshaler {
	return &primitiveReader{read: func(reader Reader) (err error) {
		*v, err = reader.ReadInt64()
		return err
	}}
}

func NewFloat32PrimitiveUnmarshaler(v *float32) Unmarshaler {
	return &primitiveReader{read: func(reader Reader) (err error) {
		*v, err = reader.ReadFloat32()
		return err
	}}
}

func NewFloat64PrimitiveUnmarshaler(v *float64) Unmarshaler {
	return &primitiveReader{read: func(reader Reader) (err error) {
		*v, err = reader.ReadFloat64()
		return err
	}}
}

func NewBoolPrimitiveUnmarshaler(v *bool) Unmarshaler {
	return &primitiveReader{read: func(reader Reader) (err error) {
		*v, err = reader.ReadBool()
		return err
	}}
}

func NewStringPrimitiveUnmarshaler(v *string) Unmarshaler {
	return &primitiveReader{read: func(reader Reader) (err error) {
		*v, err = reader.ReadString()
		return err
	}}
}

func NewBytesPrimitiveUnmarshaler(v *[]byte) Unmarshaler {
	return &primitiveReader{read: func(reader Reader) (err error) {
		*v, err = reader.ReadBytes()
		return err
	}}
}
