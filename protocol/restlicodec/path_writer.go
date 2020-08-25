package restlicodec

import "net/url"

type PathWriter interface {
	Closer
	Writer
	RawPathSegment(segment string)
}

type pathWriter struct {
	*genericWriter
}

func NewPathWriter() PathWriter {
	return &pathWriter{genericWriter: newGenericWriter(&urlWriter{stringEscaper: url.PathEscape})}
}

func (p *pathWriter) RawPathSegment(segment string) {
	p.RawString(segment)
}
