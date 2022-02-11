package fnv1a

import (
	"math"
	"sort"
)

const (
	initialHash = hash(2166136261)
	multiplier  = 16777619
	mask        = hash(0xFF)
)

func ZeroHash() Hash {
	return new(hash)
}

func NewHash() Hash {
	h := initialHash
	return &h
}

type Hashable interface {
	ComputeHash() Hash
}

// Hash provides the set of functions needed to compute a running fnv1a of the supported rest.li primitives
type Hash interface {
	AddInt32(v int32)
	AddInt64(v int64)
	AddFloat32(v float32)
	AddFloat64(v float64)
	AddBool(v bool)
	AddString(v string)
	AddBytes(v []byte)
	Add(other Hash)
	Equals(other Hash) bool
	MapKey() HashMapKey

	add(other hash)
	underlying() hash
}

type (
	HashMapKey uint32
	hash       HashMapKey
)

func (h *hash) addUint32(v uint32) {
	hV := *h
	hV ^= hash(v) & mask
	hV *= multiplier
	hV ^= hash(v>>8) & mask
	hV *= multiplier
	hV ^= hash(v>>16) & mask
	hV *= multiplier
	hV ^= hash(v>>24) & mask
	hV *= multiplier
	*h = hV
}

func (h *hash) addUint64(v uint64) {
	hV := *h
	hV ^= hash(v) & mask
	hV *= multiplier
	hV ^= hash(v>>8) & mask
	hV *= multiplier
	hV ^= hash(v>>16) & mask
	hV *= multiplier
	hV ^= hash(v>>24) & mask
	hV *= multiplier
	hV ^= hash(v>>32) & mask
	hV *= multiplier
	hV ^= hash(v>>40) & mask
	hV *= multiplier
	hV ^= hash(v>>48) & mask
	hV *= multiplier
	hV ^= hash(v>>56) & mask
	hV *= multiplier
	*h = hV
}

func (h *hash) AddInt32(v int32) {
	h.addUint32(uint32(v))
}

func (h *hash) AddInt64(v int64) {
	h.addUint64(uint64(v))
}

func (h *hash) AddFloat32(v float32) {
	h.addUint32(math.Float32bits(v))
}

func (h *hash) AddFloat64(v float64) {
	h.addUint64(math.Float64bits(v))
}

func (h *hash) AddBool(v bool) {
	hV := *h
	var b hash
	if v {
		b = 1
	} else {
		b = 0
	}
	hV ^= b
	hV *= multiplier
	*h = hV
}

func (h *hash) AddString(v string) {
	h.AddBytes([]byte(v))
}

func (h *hash) AddBytes(v []byte) {
	hV := *h
	for _, b := range v {
		hV ^= hash(b)
		hV *= multiplier
	}
	*h = hV
}

func (h *hash) Add(other Hash) {
	h.add(other.underlying())
}

func (h *hash) Equals(other Hash) bool {
	return *h == other.underlying()
}

func (h *hash) MapKey() HashMapKey {
	return HashMapKey(*h)
}

func (h *hash) add(other hash) {
	h.addUint32(uint32(other))
}

func (h *hash) underlying() hash {
	return *h
}

func AddArray[T any](h Hash, elements []T, hasher func(Hash, T)) {
	for _, e := range elements {
		hasher(h, e)
	}
}

func AddHashableArray[T Hashable](h Hash, elements []T) {
	AddArray(h, elements, func(hash Hash, t T) {
		hash.Add(t.ComputeHash())
	})
}

func AddMap[T any](h Hash, elements map[string]T, hasher func(Hash, T)) {
	kvHashes := make([]hash, len(elements))
	i := 0
	for k, v := range elements {
		kvHash := &kvHashes[i]
		i++

		kvHash.AddString(k)
		hasher(kvHash, v)
	}

	// Because order matters when hashing, the kvHashes are computed separately and inserted in ascending order. Two
	// identical maps will produce the same kvHashes and therefore add them to the total hash in the same order, meaning
	// they will hash to the same value.
	sort.Slice(kvHashes, func(i, j int) bool {
		return kvHashes[i] < kvHashes[j]
	})

	for _, kvHash := range kvHashes {
		h.add(kvHash.underlying())
	}
}

func AddHashableMap[T Hashable](h Hash, elements map[string]T) {
	AddMap(h, elements, func(hash Hash, t T) {
		hash.Add(t.ComputeHash())
	})
}

func HashInt32(v int32) Hash {
	h := initialHash
	h.AddInt32(v)
	return &h
}

func HashInt64(v int64) Hash {
	h := initialHash
	h.AddInt64(v)
	return &h
}

func HashFloat32(v float32) Hash {
	h := initialHash
	h.AddFloat32(v)
	return &h
}

func HashFloat64(v float64) Hash {
	h := initialHash
	h.AddFloat64(v)
	return &h
}

func HashBool(v bool) Hash {
	h := initialHash
	h.AddBool(v)
	return &h
}

func HashString(v string) Hash {
	h := initialHash
	h.AddString(v)
	return &h
}

func HashBytes(v []byte) Hash {
	h := initialHash
	h.AddBytes(v)
	return &h
}
