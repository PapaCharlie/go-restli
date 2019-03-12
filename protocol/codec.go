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
	encoder: url.PathEscape,
	decoder: url.PathUnescape,
}

var RestLiReducedEncoder = RestLiCodec{
	encoder: strings.NewReplacer(
		",", url.PathEscape(","),
		"(", url.PathEscape("("),
		")", url.PathEscape(")"),
		"'", url.PathEscape("'"),
		":", url.PathEscape(":")).Replace,
	decoder: url.PathUnescape,
}

type RestLiEncodable interface {
	RestLiEncode(codec RestLiCodec) (data string, err error)
	RestLiDecode(codec RestLiCodec, data string) (err error)
}
