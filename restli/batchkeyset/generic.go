package batchkeyset

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/fnv1a"
	"github.com/PapaCharlie/go-restli/restli/equals"
	"github.com/PapaCharlie/go-restli/restlicodec"
)

type genericBatchKeySet[T any] struct {
	originalKeys map[fnv1a.HashMapKey][]T
	keyCount     int
	hash         func(T) fnv1a.Hash
	equals       func(left, right T) bool
}

func (s *genericBatchKeySet[T]) AddKey(t T) error {
	h := s.hash(t).MapKey()
	found := false
	for _, key := range s.originalKeys[h] {
		if s.equals(t, key) {
			found = true
			break
		}
	}
	if found {
		return fmt.Errorf("go-restli: Cannot add key %+v twice to BatchKeySet", t)
	}
	s.originalKeys[h] = append(s.originalKeys[h], t)
	s.keyCount++
	return nil
}

func (s *genericBatchKeySet[T]) LocateOriginalKey(key T) (originalKey T, found bool) {
	for _, k := range s.originalKeys[s.hash(key).MapKey()] {
		if s.equals(k, key) {
			found = true
			originalKey = k
			break
		}
	}
	return originalKey, found
}

func (s *genericBatchKeySet[T]) LocateOriginalKeyFromReader(keyReader restlicodec.Reader) (originalKey T, err error) {
	key, err := restlicodec.UnmarshalRestLi[T](keyReader)
	if err != nil {
		return originalKey, err
	}

	originalKey, found := s.LocateOriginalKey(key)
	if !found {
		err = fmt.Errorf("go-restli: Unknown key returned by batch method: %q", keyReader)
	}
	return originalKey, err
}

func (s *genericBatchKeySet[T]) encodeKeys() ([]string, error) {
	encodedKeys := make([]string, 0, s.keyCount)
	for _, keys := range s.originalKeys {
		for _, k := range keys {
			w := restlicodec.NewRestLiQueryParamsWriter()
			err := restlicodec.MarshalRestLi(k, w)
			if err != nil {
				return nil, err
			}
			encodedKeys = append(encodedKeys, w.Finalize())
		}
	}
	return encodedKeys, nil
}

func (s *genericBatchKeySet[T]) Encode(paramNameWriter func(string) restlicodec.Writer) error {
	return encode[T](s, paramNameWriter)
}

func (s *genericBatchKeySet[T]) EncodeQueryParams() (params string, err error) {
	return generateRawQuery[T](s)
}

type ComplexKey[T any] interface {
	restlicodec.Marshaler
	ComputeComplexKeyHash() fnv1a.Hash
	ComplexKeyEquals(other T) bool
}

func NewComplexKeySet[T ComplexKey[T]]() BatchKeySet[T] {
	return &genericBatchKeySet[T]{
		originalKeys: map[fnv1a.HashMapKey][]T{},
		hash:         T.ComputeComplexKeyHash,
		equals:       T.ComplexKeyEquals,
	}
}

type SimpleKey[T any] interface {
	restlicodec.Marshaler
	fnv1a.Hashable
	equals.Comparable[T]
}

func NewSimpleKeySet[T SimpleKey[T]]() BatchKeySet[T] {
	return &genericBatchKeySet[T]{
		originalKeys: map[fnv1a.HashMapKey][]T{},
		hash:         T.ComputeHash,
		equals:       T.Equals,
	}
}

func NewBytesKeySet() BatchKeySet[[]byte] {
	return &genericBatchKeySet[[]byte]{
		originalKeys: map[fnv1a.HashMapKey][][]byte{},
		hash:         fnv1a.HashBytes,
		equals:       equals.Bytes,
	}
}

func NewBatchKeySet[K any]() BatchKeySet[K] {
	var t K
	var set any
	switch any(t).(type) {
	case ComplexKey[K]:
		set = &genericBatchKeySet[K]{
			originalKeys: map[fnv1a.HashMapKey][]K{},
			hash:         func(t K) fnv1a.Hash { return any(t).(ComplexKey[K]).ComputeComplexKeyHash() },
			equals:       func(left, right K) bool { return any(left).(ComplexKey[K]).ComplexKeyEquals(right) },
		}
	case SimpleKey[K]:
		set = &genericBatchKeySet[K]{
			originalKeys: map[fnv1a.HashMapKey][]K{},
			hash:         func(t K) fnv1a.Hash { return any(t).(SimpleKey[K]).ComputeHash() },
			equals:       func(left, right K) bool { return any(left).(SimpleKey[K]).Equals(right) },
		}
	case []byte:
		set = NewBytesKeySet()
	case int32:
		set = NewPrimitiveKeySet[int32]()
	case int64:
		set = NewPrimitiveKeySet[int64]()
	case float32:
		set = NewPrimitiveKeySet[float32]()
	case float64:
		set = NewPrimitiveKeySet[float64]()
	case bool:
		set = NewPrimitiveKeySet[bool]()
	case string:
		set = NewPrimitiveKeySet[string]()
	}
	return set.(BatchKeySet[K])
}
