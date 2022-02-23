package protocol

import (
	"strings"
)

type ResourcePath interface {
	RootResource() string
	ResourcePath() (path string, err error)
}

type ResourcePathString string

func (s ResourcePathString) RootResource() string {
	root, _, _ := strings.Cut(string(s[1:]), "/")
	return root
}

func (s ResourcePathString) ResourcePath() (string, error) {
	return string(s), nil
}

type QueryParams interface {
	EncodeQueryParams() (string, error)
}

type QueryParamsString string

func (q QueryParamsString) EncodeQueryParams() (string, error) {
	return string(q), nil
}
