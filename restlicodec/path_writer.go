package restlicodec

import (
	"strings"
)

// Ror2PathWriter is an ROR2 Writer that is intended to construct URLs for entities that are ROR2 encoded.
type Ror2PathWriter interface {
	Writer
	RawPathSegment(segment string)
}

type ror2PathWriter struct {
	*genericWriter
}

var unescapedPathCharacters = func() map[byte]struct{} {
	const chars = `!$&*+-.0123456789=@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz~`
	m := make(map[byte]struct{}, len(chars))
	for i := range chars {
		m[chars[i]] = struct{}{}
	}
	return m
}()

// Ror2PathEscape path-escapes the given string using Rest.li's rather opinionated path escaper. The list of unescaped
// characters was generated directly from the Java code by enumerating all UTF-8 characters and escaping them. Turns out
// the only differences are that '!' and '*' aren't escaped (while url.PathEscape does) and ':' isn't escaped by
// url.PathEscape but the Rest.li escaper escapes it.
func Ror2PathEscape(s string) string {
	buf := new(strings.Builder)
	for _, c := range []byte(s) {
		if _, ok := unescapedPathCharacters[c]; ok {
			buf.WriteByte(c)
		} else {
			hexEscape(buf, c)
		}
	}
	return buf.String()
}

func NewRor2PathWriter() Ror2PathWriter {
	return &ror2PathWriter{genericWriter: newGenericWriter(&ror2Writer{stringEscaper: Ror2PathEscape}, nil)}
}

func (p *ror2PathWriter) RawPathSegment(segment string) {
	p.RawString(segment)
}
