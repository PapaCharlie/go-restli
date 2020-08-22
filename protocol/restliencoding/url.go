package restliencoding

import (
	"io"
	"net/url"
	"strings"

	"github.com/mailru/easyjson/jwriter"
)

const emptyString = `''`

type urlEncoder struct {
	jwriter.Writer
	stringEscaper func(string) string
}

var headerEncodingEscaper = strings.NewReplacer(
	"%", url.QueryEscape("%"),
	",", url.QueryEscape(","),
	"(", url.QueryEscape("("),
	")", url.QueryEscape(")"),
	"'", url.QueryEscape("'"),
	":", url.QueryEscape(":")).Replace

func NewHeaderEncoder() *Encoder {
	return &Encoder{encoder: &urlEncoder{
		stringEscaper: headerEncodingEscaper,
	}}
}

func (u *urlEncoder) WriteObjectStart() {
	u.Writer.RawByte('(')
}

func (u *urlEncoder) WriteFieldName(name string) {
	u.Writer.RawString(name)
}

func (u *urlEncoder) WriteFieldNameDelimiter() {
	u.Writer.RawByte(':')
}

func (u *urlEncoder) WriteFieldDelimiter() {
	u.Writer.RawByte(',')
}

func (u *urlEncoder) WriteObjectEnd() {
	u.Writer.RawByte(')')
}

func (u *urlEncoder) WriteMapStart() {
	u.WriteObjectStart()
}

func (u *urlEncoder) WriteMapKey(key string) {
	u.WriteFieldName(key)
}

func (u *urlEncoder) WriteMapKeyDelimiter() {
	u.WriteFieldNameDelimiter()
}

func (u *urlEncoder) WriteMapEntryDelimiter() {
	u.WriteFieldDelimiter()
}

func (u *urlEncoder) WriteMapEnd() {
	u.WriteObjectEnd()
}

func (u *urlEncoder) WriteArrayStart() {
	u.Writer.RawString("List(")
}

func (u *urlEncoder) WriteArrayItemDelimiter() {
	u.Writer.RawString(",")
}

func (u *urlEncoder) WriteArrayEnd() {
	u.Writer.RawByte(')')
}

func (u *urlEncoder) String(v string) {
	if len(v) == 0 {
		u.Writer.RawString(emptyString)
	} else {
		u.Writer.RawString(u.stringEscaper(v))
	}
}

func (u *urlEncoder) Bytes(v []byte) {
	u.String(string(v))
}

func (u *urlEncoder) Finalize() string {
	data, _ := u.BuildBytes()
	return string(data)
}

func (u *urlEncoder) ReadCloser() io.ReadCloser {
	rc, _ := u.Writer.ReadCloser()
	return rc
}

func (u *urlEncoder) SubEncoder() encoder {
	return u
}
