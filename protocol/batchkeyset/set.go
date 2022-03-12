package batchkeyset

import (
	"sort"

	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

const EntityIDsField = "ids"

type batchKeySet[T any] interface {
	encodeKeys() ([]string, error)
	AddKey(T) error
}

type BatchKeySet[T any] interface {
	batchKeySet[T]
	LocateOriginalKey(key T) (originalKey T, found bool)
	LocateOriginalKeyFromReader(keyReader restlicodec.Reader) (originalKey T, err error)
	Encode(paramNameWriter func(string) restlicodec.Writer) error
	EncodeQueryParams() (params string, err error)
}

func generateRawQuery[T any](set batchKeySet[T]) (string, error) {
	writer := restlicodec.NewRestLiQueryParamsWriter()
	err := writer.WriteParams(func(paramNameWriter func(key string) restlicodec.Writer) error {
		return encode(set, paramNameWriter)
	})
	if err != nil {
		return "", err
	}
	return writer.Finalize(), nil
}

func encode[T any](set batchKeySet[T], paramNameWriter func(string) restlicodec.Writer) (err error) {
	encodedKeys, err := set.encodeKeys()
	if err != nil {
		return err
	}
	sort.Strings(encodedKeys)

	return paramNameWriter(EntityIDsField).WriteArray(func(itemWriter func() restlicodec.Writer) error {
		for _, k := range encodedKeys {
			itemWriter().WriteRawBytes([]byte(k))
		}
		return nil
	})
}

func AddAllKeys[T any](set BatchKeySet[T], keys ...T) (err error) {
	for _, k := range keys {
		err = set.AddKey(k)
		if err != nil {
			return err
		}
	}
	return err
}

func AddAllMapKeys[K comparable, V any](set BatchKeySet[K], entities map[K]V) (err error) {
	for k := range entities {
		err = set.AddKey(k)
		if err != nil {
			return err
		}
	}
	return err
}
