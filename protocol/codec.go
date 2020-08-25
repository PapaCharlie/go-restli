package protocol

import (
	"net/url"
	"strings"
)

// RawComplexKey is a temporary workaround for the fact that this library does not yet support decoding rest.li URL
// encoded data. As such, CREATE responses for complex keys will grab the header value with the raw complex key and
// return it as-is so that clients can build their own deserialization for the time being
type RawComplexKey string

type RestLiCodec struct {
	encoder func(string) string
	decoder func(string) (string, error)
}

var RestLiUrlPathEncoder = &RestLiCodec{
	encoder: url.PathEscape,
	decoder: url.PathUnescape,
}

var RestLiQueryEncoder = &RestLiCodec{
	encoder: url.QueryEscape,
	decoder: url.QueryUnescape,
}

var RestLiReducedEncoder = &RestLiCodec{
	encoder: strings.NewReplacer(
		"%", url.QueryEscape("%"),
		",", url.QueryEscape(","),
		"(", url.QueryEscape("("),
		")", url.QueryEscape(")"),
		"'", url.QueryEscape("'"),
		":", url.QueryEscape(":")).Replace,
	decoder: url.QueryUnescape,
}
