package native

import (
	"time"

	"github.com/PapaCharlie/go-restli/fnv1a"
)

func toInt(ts time.Time) int64 {
	return int64(time.Duration(ts.UnixNano()) / time.Millisecond)
}

func NewTime(ts int64) time.Time {
	return time.Unix(0, 0).Add(time.Millisecond * time.Duration(ts))
}

func MarshalTime(ts time.Time) (int64, error) {
	return toInt(ts), nil
}

func UnmarshalTime(ts int64) (time.Time, error) {
	return NewTime(ts), nil
}

func EqualsTime(left time.Time, right time.Time) bool {
	return left.Equal(right)
}

func ComputeHashTime(ts time.Time) fnv1a.Hash {
	h := fnv1a.NewHash()
	h.AddInt64(toInt(ts))
	return h
}

func ZeroValueTime() (t time.Time) {
	return t
}
