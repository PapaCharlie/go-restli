package batchkeyset

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/fnv1a"
	"github.com/PapaCharlie/go-restli/protocol/equals"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

type genericBatchKeySet[T any] struct {
	originalKeys map[fnv1a.HashMapKey][]T
	keyCount     int
	marshaler    restlicodec.GenericMarshaler[T]
	unmarshaler  restlicodec.GenericUnmarshaler[T]
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

func (s *genericBatchKeySet[T]) LocateOriginalKey(keyReader restlicodec.Reader) (originalKey T, err error) {
	returnedKey, err := s.unmarshaler(keyReader)
	if err != nil {
		return originalKey, err
	}

	found := false
	for _, key := range s.originalKeys[s.hash(returnedKey).MapKey()] {
		if s.equals(returnedKey, key) {
			found = true
			originalKey = key
			break
		}
	}
	if !found {
		err = fmt.Errorf("go-restli: Unknown key returned by batch method: %q", keyReader)
	}
	return originalKey, err
}

func (s *genericBatchKeySet[T]) MarshalKey(writer restlicodec.Writer, t T) error {
	return s.marshaler(t, writer)
}

func (s *genericBatchKeySet[T]) encodeKeys() ([]string, error) {
	encodedKeys := make([]string, 0, s.keyCount)
	for _, keys := range s.originalKeys {
		for _, k := range keys {
			w := restlicodec.NewRestLiQueryParamsWriter()
			err := s.marshaler(k, w)
			if err != nil {
				return nil, err
			}
			encodedKeys = append(encodedKeys, w.Finalize())
		}
	}
	return encodedKeys, nil
}

type ComplexKey[T any] interface {
	restlicodec.Marshaler
	ComputeComplexKeyHash() fnv1a.Hash
	ComplexKeyEquals(other T) bool
}

func NewComplexKeySet[T ComplexKey[T]](unmarshaler restlicodec.GenericUnmarshaler[T]) BatchKeySet[T] {
	return &genericBatchKeySet[T]{
		originalKeys: map[fnv1a.HashMapKey][]T{},
		marshaler:    T.MarshalRestLi,
		unmarshaler:  unmarshaler,
		hash:         T.ComputeComplexKeyHash,
		equals:       T.ComplexKeyEquals,
	}
}

type SimpleKey[T any] interface {
	restlicodec.Marshaler
	fnv1a.Hashable
	equals.Equatable[T]
}

func NewSimpleKeySet[T SimpleKey[T]](unmarshaler restlicodec.GenericUnmarshaler[T]) BatchKeySet[T] {
	return &genericBatchKeySet[T]{
		originalKeys: map[fnv1a.HashMapKey][]T{},
		marshaler:    T.MarshalRestLi,
		unmarshaler:  unmarshaler,
		hash:         T.ComputeHash,
		equals:       T.Equals,
	}
}

func NewBytesKeySet() BatchKeySet[[]byte] {
	return &genericBatchKeySet[[]byte]{
		originalKeys: map[fnv1a.HashMapKey][][]byte{},
		marshaler: func(v []byte, writer restlicodec.Writer) error {
			writer.WriteBytes(v)
			return nil
		},
		unmarshaler: restlicodec.Reader.ReadBytes,
		hash:        fnv1a.HashBytes,
		equals:      equals.Bytes,
	}
}
