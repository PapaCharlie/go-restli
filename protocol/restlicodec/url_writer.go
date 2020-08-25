package restlicodec

import "github.com/mailru/easyjson/jwriter"

type urlWriter struct {
	jwriter.Writer
	stringEscaper func(string) string
}

func NewHeaderWriter() WriteCloser {
	return newGenericWriter(&urlWriter{stringEscaper: headerEncodingEscaper})
}

func (u *urlWriter) writeMapStart() {
	u.Writer.RawByte('(')
}

func (u *urlWriter) writeKey(key string) {
	u.Writer.RawString(key)
}

func (u *urlWriter) writeKeyDelimiter() {
	u.Writer.RawByte(':')
}

func (u *urlWriter) writeEntryDelimiter() {
	u.Writer.RawByte(',')
}

func (u *urlWriter) writeMapEnd() {
	u.Writer.RawByte(')')
}

func (u *urlWriter) writeArrayStart() {
	u.Writer.RawString("List(")
}

func (u *urlWriter) writeArrayItemDelimiter() {
	u.Writer.RawString(",")
}

func (u *urlWriter) writeArrayEnd() {
	u.Writer.RawByte(')')
}

func (u *urlWriter) WriteInt32(v int32) {
	u.Writer.Int32(v)
}

func (u *urlWriter) WriteInt64(v int64) {
	u.Writer.Int64(v)
}

func (u *urlWriter) WriteFloat32(v float32) {
	u.Writer.Float32(v)
}

func (u *urlWriter) WriteFloat64(v float64) {
	u.Writer.Float64(v)
}

func (u *urlWriter) WriteBool(v bool) {
	u.Writer.Bool(v)
}

func (u *urlWriter) WriteString(v string) {
	if len(v) == 0 {
		u.Writer.RawString(emptyString)
	} else {
		u.Writer.RawString(u.stringEscaper(v))
	}
}

func (u *urlWriter) WriteBytes(v []byte) {
	u.WriteString(string(v))
}
