package restlicodec

import (
	"net/url"
)

type finderEncoder struct {
	urlEncoder
}

func NewFinderEncoder() *Encoder {
	return &Encoder{encoder: &finderEncoder{
		urlEncoder{stringEscaper: url.QueryEscape},
	}}
}

func (f *finderEncoder) WriteObjectStart() {
}

func (f *finderEncoder) WriteFieldNameDelimiter() {
	f.Writer.RawByte('=')
}

func (f *finderEncoder) WriteFieldDelimiter() {
	f.Writer.RawByte('&')
}

func (f *finderEncoder) WriteObjectEnd() {
}
