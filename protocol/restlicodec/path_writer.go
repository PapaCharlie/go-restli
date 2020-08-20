package restlicodec

import "net/url"

// Ror2PathWriter is an ROR2 Writer that is intended to construct URLs for entities that are ROR2 encoded.
type Ror2PathWriter interface {
	Closer
	Writer
	RawPathSegment(segment string)
}

type ror2PathWriter struct {
	*genericWriter
}

func NewRor2PathWriter() Ror2PathWriter {
	return &ror2PathWriter{genericWriter: newGenericWriter(&ror2Writer{stringEscaper: url.PathEscape})}
}

func (p *ror2PathWriter) RawPathSegment(segment string) {
	p.RawString(segment)
}
