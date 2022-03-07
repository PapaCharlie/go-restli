package testsuite

import (
	"time"

	"github.com/PapaCharlie/go-fnv1a"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

// go-restli:typeref {
//   "type": "int64",
//   "ref": "extras.Time",
//   "package": "time",
//   "name": "Time",
//   "nonReceiverFuncsPackage": "github.com/PapaCharlie/go-restli/internal/tests/native/testsuite"
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

// go-restli:generated {
//   "name": "Temperature",
//   "namespace": "extras",
//   "package": "github.com/PapaCharlie/go-restli/internal/tests/native/testsuite"
// }

type Temperature int32

func UnmarshalRestLiTemperature(reader restlicodec.Reader) (Temperature, error) {
	v, err := reader.ReadInt32()
	return Temperature(v), err
}

func (t *Temperature) UnmarshalRestLi(reader restlicodec.Reader) error {
	v, err := reader.ReadInt32()
	if err != nil {
		return err
	}
	*t = Temperature(v)
	return nil
}

func (t Temperature) MarshalRestLi(writer restlicodec.Writer) error {
	writer.WriteInt32(int32(t))
	return nil
}

func (t Temperature) Equals(other Temperature) bool {
	return t == other
}

func (t Temperature) ComputeHash() fnv1a.Hash {
	return fnv1a.HashInt32(int32(t))
}
