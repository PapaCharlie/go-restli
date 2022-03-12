package batchkeyset

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

type primitiveKeySet[T restlicodec.ComparablePrimitive] struct {
	originalKeys map[T]struct{}
}

func (s *primitiveKeySet[T]) AddKey(t T) error {
	if _, ok := s.originalKeys[t]; ok {
		return fmt.Errorf("go-restli: Cannot add key %+v twice to BatchKeySet", t)
	}
	s.originalKeys[t] = struct{}{}
	return nil
}

func (s *primitiveKeySet[T]) LocateOriginalKey(key T) (originalKey T, found bool) {
	_, found = s.originalKeys[key]
	return key, found
}

func (s *primitiveKeySet[T]) LocateOriginalKeyFromReader(keyReader restlicodec.Reader) (originalKey T, err error) {
	originalKey, err = restlicodec.UnmarshalRestLi[T](keyReader)
	if err != nil {
		return originalKey, err
	}

	_, ok := s.LocateOriginalKey(originalKey)
	if !ok {
		err = fmt.Errorf("go-restli: Unknown key returned by batch method: %q", keyReader)
	}
	return originalKey, err
}

func (s *primitiveKeySet[T]) Encode(paramNameWriter func(string) restlicodec.Writer) error {
	return encode[T](s, paramNameWriter)
}

func (s *primitiveKeySet[T]) EncodeQueryParams() (params string, err error) {
	return generateRawQuery[T](s)
}

func (s *primitiveKeySet[T]) encodeKeys() ([]string, error) {
	encodedKeys := make([]string, 0, len(s.originalKeys))
	for k := range s.originalKeys {
		w := restlicodec.NewRestLiQueryParamsWriter()
		_ = restlicodec.MarshalRestLi(k, w)
		encodedKeys = append(encodedKeys, w.Finalize())
	}
	return encodedKeys, nil
}

func NewPrimitiveKeySet[T restlicodec.ComparablePrimitive]() BatchKeySet[T] {
	return &primitiveKeySet[T]{originalKeys: map[T]struct{}{}}
}
