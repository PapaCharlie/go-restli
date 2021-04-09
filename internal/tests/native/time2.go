package native

import (
	"time"

	"github.com/PapaCharlie/go-fnv1a"
)

func NewTime2(ts int64) time.Time {
	return time.Unix(0, 0).Add(time.Millisecond * time.Duration(ts))
}

func MarshalTime2(ts time.Time) (int64, error) {
	return MarshalTime(ts)
}

func UnmarshalTime2(ts int64) (time.Time, error) {
	return UnmarshalTime(ts)
}

func EqualsTime2(left time.Time, right time.Time) bool {
	return EqualsTime(left, right)
}

func ComputeHashTime2(ts time.Time) fnv1a.Hash {
	return ComputeHashTime(ts)
}

func ZeroValueTime2() (t time.Time) {
	return t
}
