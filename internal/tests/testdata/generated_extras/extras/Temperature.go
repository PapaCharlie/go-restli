package extras

import (
	"github.com/PapaCharlie/go-restli/v2/fnv1a"
)

type Temperature int32

func MarshalTemperature(t Temperature) (int32, error) {
	return int32(t), nil
}

func UnmarshalTemperature(t int32) (Temperature, error) {
	return Temperature(t), nil
}

func EqualsTemperature(t1, t2 Temperature) bool {
	return t1 == t2
}

func ComputeHashTemperature(t Temperature) fnv1a.Hash {
	return fnv1a.HashInt32(int32(t))
}

func (t Temperature) Pointer() *Temperature {
	return &t
}
