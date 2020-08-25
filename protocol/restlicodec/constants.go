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
