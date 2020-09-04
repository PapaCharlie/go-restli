package fnv1a

import (
	"math"
)

const (
	zeroHash   = Hash(2166136261)
	multiplier = 16777619
	mask       = Hash(0xFF)
)

func NewHash() Hash {
	return zeroHash
}

// Hash provides the set of functions needed to compute a running fnv1a of the supported rest.li primitives
type Hash uint32

func (h *Hash) addUint32(v uint32) {
	hash := *h
	hash ^= Hash(v) & mask
	hash *= multiplier
	hash ^= Hash(v>>8) & mask
	hash *= multiplier
	hash ^= Hash(v>>16) & mask
	hash *= multiplier
	hash ^= Hash(v>>24) & mask
	hash *= multiplier
	*h = hash
}

func (h *Hash) addUint64(v uint64) {
	hash := *h
	hash ^= Hash(v) & mask
	hash *= multiplier
	hash ^= Hash(v>>8) & mask
	hash *= multiplier
	hash ^= Hash(v>>16) & mask
	hash *= multiplier
	hash ^= Hash(v>>24) & mask
	hash *= multiplier
	hash ^= Hash(v>>32) & mask
	hash *= multiplier
	hash ^= Hash(v>>40) & mask
	hash *= multiplier
	hash ^= Hash(v>>48) & mask
	hash *= multiplier
	hash ^= Hash(v>>56) & mask
	hash *= multiplier
	*h = hash
}

func (h *Hash) AddInt32(v int32) {
	h.addUint32(uint32(v))
}

func (h *Hash) AddInt64(v int64) {
	h.addUint64(uint64(v))
}

func (h *Hash) AddFloat32(v float32) {
	h.addUint32(math.Float32bits(v))
}

func (h *Hash) AddFloat64(v float64) {
	h.addUint64(math.Float64bits(v))
}

func (h *Hash) AddBool(v bool) {
	hash := *h
	var b Hash
	if v {
		b = 1
	} else {
		b = 0
	}
	hash ^= b
	hash *= multiplier
	*h = hash
}

func (h *Hash) AddString(v string) {
	h.AddBytes([]byte(v))
}

func (h *Hash) AddBytes(v []byte) {
	hash := *h
	for _, b := range v {
		hash ^= Hash(b)
		hash *= multiplier
	}
	*h = hash
}

func (h *Hash) Add(other Hash) {
	h.addUint32(uint32(other.underlyingHash()))
}

func (h *Hash) underlyingHash() Hash {
	return *h
}
