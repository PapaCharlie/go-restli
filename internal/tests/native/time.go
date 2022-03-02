package native

import (
	"time"

	"github.com/PapaCharlie/go-restli/fnv1a"
)

//go:restli_typeref
// {
//   "type": "int64",
//   "ref": "extras.Time",
//   "nativePackage": "time",
//   "nativeIdentifier": "Time",
//   "nonReceiverFuncsPackage": "github.com/PapaCharlie/go-restli/internal/tests/native"
// }

func UnmarshalRestLiTime(v int64) (time.Time, error) {
	return time.UnixMilli(v), nil
}

func MarshalRestLiTime(t time.Time) (int64, error) {
	return t.UnixMilli(), nil
}

func EqualsTime(left, right time.Time) (b bool) {
	return left.Equal(right)
}

func ComputeHashTime(t time.Time) fnv1a.Hash {
	return fnv1a.HashInt64(t.UnixMilli())
}

//go:restli_typeref {
//   "type": "int32",
//   "ref": "extras.Temperature",
//   "nativePackage": "github.com/PapaCharlie/go-restli/internal/tests/native",
//   "nativeIdentifier": "Temp"
// }

type Temp int32

func UnmarshalRestLiTemperature(v int32) (Temp, error) {
	return Temp(v), nil
}

func (t Temp) MarshalRestLi() (int32, error) {
	return int32(t), nil
}

func (t Temp) Equals(other Temp) bool {
	return t == other
}

func (t Temp) ComputeHash() fnv1a.Hash {
	return fnv1a.HashInt32(int32(t))
}
