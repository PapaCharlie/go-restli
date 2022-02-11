package batchkeyset

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

type primitiveKeySet[T restlicodec.ComparablePrimitive] struct {
	originalKeys map[T]struct{}
	marshaler    restlicodec.PrimitiveMarshaler[T]
	unmarshaler  restlicodec.GenericUnmarshaler[T]
}

func (s *primitiveKeySet[T]) AddKey(t T) error {
	if _, ok := s.originalKeys[t]; ok {
		return fmt.Errorf("go-restli: Cannot add key %+v twice to BatchKeySet", t)
	}
	s.originalKeys[t] = struct{}{}
	return nil
}

func (s *primitiveKeySet[T]) LocateOriginalKey(keyReader restlicodec.Reader) (originalKey T, err error) {
	originalKey, err = s.unmarshaler(keyReader)
	if err != nil {
		return originalKey, err
	}

	_, ok := s.originalKeys[originalKey]
	if !ok {
		err = fmt.Errorf("go-restli: Unknown key returned by batch method: %q", keyReader)
	}
	return originalKey, err
}

func (s *primitiveKeySet[T]) MarshalKey(writer restlicodec.Writer, t T) error {
	s.marshaler(writer, t)
	return nil
}

func (s *primitiveKeySet[T]) encodeKeys() ([]string, error) {
	encodedKeys := make([]string, 0, len(s.originalKeys))
	for k := range s.originalKeys {
		w := restlicodec.NewRestLiQueryParamsWriter()
		s.marshaler(w, k)
		encodedKeys = append(encodedKeys, w.Finalize())
	}
	return encodedKeys, nil
}

func NewPrimitiveKeySet[T restlicodec.ComparablePrimitive](
	marshaler restlicodec.PrimitiveMarshaler[T],
	unmarshaler restlicodec.GenericUnmarshaler[T],
) BatchKeySet[T] {
	return &primitiveKeySet[T]{
		originalKeys: map[T]struct{}{},
		marshaler:    marshaler,
		unmarshaler:  unmarshaler,
	}
}
