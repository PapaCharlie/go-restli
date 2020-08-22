package restliencoding

import (
	"net/url"
)

type queryParamsEncoder struct {
	urlEncoder
}

func NewQueryParamsEncoder() *Encoder {
	return &Encoder{encoder: &queryParamsEncoder{
		urlEncoder{stringEscaper: url.QueryEscape},
	}}
}

func (f *queryParamsEncoder) WriteObjectStart() {
}

func (f *queryParamsEncoder) WriteFieldNameDelimiter() {
	f.Writer.RawByte('=')
}

func (f *queryParamsEncoder) WriteFieldDelimiter() {
	f.Writer.RawByte('&')
}

func (f *queryParamsEncoder) WriteObjectEnd() {
}
