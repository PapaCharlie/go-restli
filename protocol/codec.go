package protocol

import (
	"net/url"
	"strings"
)

type RestLiCodec struct {
	encoder func(string) string
	decoder func(string) (string, error)
}

var RestLiUrlEncoder = RestLiCodec{
	encoder: url.QueryEscape,
	decoder: url.QueryUnescape,
}

var RestLiReducedEncoder = RestLiCodec{
	encoder: strings.NewReplacer(
		",", url.QueryEscape(","),
		"(", url.QueryEscape("("),
		")", url.QueryEscape(")"),
		"'", url.QueryEscape("'"),
		":", url.QueryEscape(":")).Replace,
	decoder: url.QueryUnescape,
}

type RestLiEncodable interface {
	RestLiEncode(*RestLiCodec, *strings.Builder) error
	// RestLiDecode(*RestLiCodec, string) error
}
