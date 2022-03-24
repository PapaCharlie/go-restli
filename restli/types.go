package restli

import (
	"strings"

	"github.com/PapaCharlie/go-restli/restlicodec"
)

type ResourcePath interface {
	RootResource() string
	ResourcePath() (path string, err error)
}

type ResourcePathUnmarshaler[T any] interface {
	NewInstance() T
	UnmarshalResourcePath(segments []restlicodec.Reader) error
}

func UnmarshalResourcePath[T ResourcePathUnmarshaler[T]](segments []restlicodec.Reader) (t T, err error) {
	t = t.NewInstance()
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
