package restlicodec

import (
	"net/url"
	"strings"
)

const (
	emptyString = `''`
)

var headerEncodingEscaper = strings.NewReplacer(
	"%", url.QueryEscape("%"),
	",", url.QueryEscape(","),
	"(", url.QueryEscape("("),
	")", url.QueryEscape(")"),
	"'", url.QueryEscape("'"),
	":", url.QueryEscape(":")).Replace

func hexEscape(buf *strings.Builder, c byte) {
	const hexChars = "0123456789ABCDEF"

	buf.WriteByte('%')
	buf.WriteByte(hexChars[c>>4])
	buf.WriteByte(hexChars[c&15])
}
