package restlidecoding

type Decodable interface {
	RestliDecode(Decoder) error
}

type Decoder interface {
	ReadObject(objectReader func(field string) error) error
	ReadMap(mapReader func(key string) error) error
	ReadArray(arrayReader func(index int) error) error

	Int32() (int32, error)
	Int64() (int64, error)
	Float32() (float32, error)
	Float64() (float64, error)
	Bool() (bool, error)
	String() (string, error)
	Bytes() ([]byte, error)
	Decodable(Decodable) error

	SubDecoder() Decoder
}
