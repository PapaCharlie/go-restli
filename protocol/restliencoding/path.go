package restliencoding

import "net/url"

type PathEncoder struct {
	*Encoder
}

func NewPathEncoder() *PathEncoder {
	return &PathEncoder{Encoder: &Encoder{encoder: &urlEncoder{
		stringEscaper: url.PathEscape,
	}}}
}

func (e *PathEncoder) RawPathSegment(segment string) {
	e.encoder.(*urlEncoder).Writer.RawString(segment)
}
