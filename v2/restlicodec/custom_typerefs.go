package restlicodec

import (
	"log"
	"reflect"
	"sync"

	"github.com/PapaCharlie/go-restli/v2/fnv1a"
)

var customTyperefAdapters sync.Map

type adapter[T any] struct {
	marshaler   GenericMarshaler[T]
	unmarshaler GenericUnmarshaler[T]
	hasher      func(T) fnv1a.Hash
	equals      func(T, T) bool
}

func RegisterCustomTyperef[P, T any](
	marshaler func(T) (P, error),
	unmarshaler func(P) (T, error),
	hasher func(T) fnv1a.Hash,
	equals func(T, T) bool,
) {
	var t *T
	_, loaded := customTyperefAdapters.LoadOrStore(t, &adapter[T]{
		marshaler: func(t T, writer Writer) error {
			return WriteCustomTyperef(writer, t, marshaler)
		},
		unmarshaler: func(reader Reader) (T, error) {
			return ReadCustomTyperef(reader, unmarshaler)
		},
		hasher: hasher,
		equals: equals,
	})
	if loaded {
		log.Panicf("Cannot register custom typeref %q more than once", reflect.TypeOf(t).Elem())
	}
}

func loadAdapter[T any]() *adapter[T] {
	var t *T
	v, ok := customTyperefAdapters.Load(t)
	if !ok {
		log.Panicf("Unregistered custom typeref %q", reflect.TypeOf(t).Elem())
	}
	return v.(*adapter[T])
}

func CustomTyperefMarshaler[T any]() GenericMarshaler[T] {
	return loadAdapter[T]().marshaler
}

func CustomTyperefUnmarshaler[T any]() GenericUnmarshaler[T] {
	return loadAdapter[T]().unmarshaler
}

func CustomTyperefHasher[T any]() func(T) fnv1a.Hash {
	return loadAdapter[T]().hasher
}

func CustomTyperefEquals[T any]() func(T, T) bool {
	return loadAdapter[T]().equals
}

func WriteCustomTyperef[P, T any](writer Writer, t T, marshaler func(T) (P, error)) (err error) {
	p, err := marshaler(t)
	if err != nil {
		return err
	}
	return MarshalRestLi[P](p, writer)
}

func ReadCustomTyperef[P, T any](reader Reader, unmarshaler func(P) (T, error)) (t T, err error) {
	p, err := UnmarshalRestLi[P](reader)
	if err != nil {
		return t, err
	}
	return unmarshaler(p)
}
