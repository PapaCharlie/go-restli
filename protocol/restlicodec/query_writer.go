package restlicodec

import (
	"strings"
)

type RestLiQueryParamsWriter interface {
	Closer
	WriteParams(paramsWriter MapWriter) error
}

type queryParamsWriter struct {
	*genericWriter
}

var unescapedQueryCharacters = func() map[byte]struct{} {
	const chars = `!$*-./0123456789;?@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz~`
	m := make(map[byte]struct{}, len(chars))
	for i := range chars {
		m[chars[i]] = struct{}{}
	}
	return m
}()

// Ror2QUeryEscape query-escapes the given string using Rest.li's query escaper (same as Ror2PathEscape). Using the same
// technique of generating all the UTF8 characters, a handful of characters were found to be escaped differently than
// how a normal query encoder would do it.
func Ror2QueryEscape(s string) string {
	buf := new(strings.Builder)
	for _, c := range []byte(s) {
		if c == ' ' {
			hexEscape(buf, ' ')
		} else if _, ok := unescapedQueryCharacters[c]; ok {
			buf.WriteByte(c)
		} else {
			hexEscape(buf, c)
		}
	}
	return buf.String()
}

func NewRestLiQueryParamsWriter() RestLiQueryParamsWriter {
	return &queryParamsWriter{genericWriter: newGenericWriter(&ror2Writer{stringEscaper: Ror2QueryEscape}, nil)}
}

func (w *queryParamsWriter) WriteParams(paramsWriter MapWriter) (err error) {
	first := true
	return paramsWriter(func(paramName string) Writer {
		if first {
			first = false
		} else {
			w.RawByte('&')
		}
		w.RawString(paramName)
		w.RawByte('=')
		return w
	})
}
