package protocol

import (
	"reflect"
	"strings"

	"github.com/PapaCharlie/go-restli/fnv1a"
	"github.com/PapaCharlie/go-restli/protocol/equals"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

type RestLiObject[T any] interface {
	equals.Equatable[T]
	fnv1a.Hashable
	restlicodec.Marshaler
}

type ResourcePath interface {
	RootResource() string
	ResourcePath() (path string, err error)
}

type ResourcePathUnmarshaler interface {
	UnmarshalResourcePath(segments []restlicodec.Reader) error
}

func UnmarshalResourcePath[T ResourcePathUnmarshaler](segments []restlicodec.Reader) (t T, err error) {
	v := reflect.New(reflect.TypeOf(t).Elem())
	t = v.Interface().(T)
	return t, t.UnmarshalResourcePath(segments)
}

type ResourcePathString string

func (s ResourcePathString) RootResource() string {
	root, _, _ := strings.Cut(string(s[1:]), "/")
	return root
}

func (s ResourcePathString) ResourcePath() (string, error) {
	return string(s), nil
}

type QueryParamsEncoder interface {
	EncodeQueryParams() (string, error)
}

type QueryParamsString string

func (q QueryParamsString) EncodeQueryParams() (string, error) {
	return string(q), nil
}
