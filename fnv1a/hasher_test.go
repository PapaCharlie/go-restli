package fnv1a

import "testing"

func TestAddMap(t *testing.T) {
	h := NewHash()
	AddMap(h, map[string]int32{
		"foo": 1,
		"bar": 2,
	}, Hash.AddInt32)
}
