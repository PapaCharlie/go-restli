package restlicodec

import (
	"math"

	"github.com/mailru/easyjson/jwriter"
)

type ror2Writer struct {
	jwriter.Writer
	stringEscaper func(string) string
}

// NewRor2HeaderWriter returns a new WriteCloser that serializes objects using the rest.li protocol 2.0 object and
// array representation (ROR2), whose spec is defined here:
// https://linkedin.github.io/rest.li/spec/protocol#restli-protocol-20-object-and-listarray-representation
// This specific WriteCloser uses the "reduced" URL encoding instead of the full URL encoding, i.e. it only escapes the
// following characters using url.QueryEscape:
//   % , ( ) ' :
func NewRor2HeaderWriter() WriteCloser {
	return newGenericWriter(&ror2Writer{stringEscaper: headerEncodingEscaper}, nil)
}

// NewRor2HeaderWriterWithExcludedFields returns a new WriteCloser that serializes objects using the rest.li protocol
// 2.0 object and array representation (ROR2), whose spec is defined here:
// https://linkedin.github.io/rest.li/spec/protocol#restli-protocol-20-object-and-listarray-representation
// This specific WriteCloser uses the "reduced" URL encoding instead of the full URL encoding, i.e. it only escapes the
// following characters using url.QueryEscape:
//   % , ( ) ' :
// Any fields matched by the given PathSpec are excluded from serialization
func NewRor2HeaderWriterWithExcludedFields(excludedFields PathSpec) WriteCloser {
	return newGenericWriter(&ror2Writer{stringEscaper: headerEncodingEscaper}, excludedFields)
}

func (u *ror2Writer) writeMapStart() {
	u.Writer.RawByte('(')
}

func (u *ror2Writer) writeKey(key string) {
	u.Writer.RawString(key)
}

func (u *ror2Writer) writeKeyDelimiter() {
	u.Writer.RawByte(':')
}

func (u *ror2Writer) writeEntryDelimiter() {
	u.Writer.RawByte(',')
}

func (u *ror2Writer) writeMapEnd() {
	u.Writer.RawByte(')')
}

func (u *ror2Writer) writeEmptyMap() {
	u.writeMapStart()
	u.writeMapEnd()
}

func (u *ror2Writer) writeArrayStart() {
	u.Writer.RawString("List(")
}

func (u *ror2Writer) writeArrayItemDelimiter() {
	u.Writer.RawString(",")
}

func (u *ror2Writer) writeArrayEnd() {
	u.Writer.RawByte(')')
}

func (u *ror2Writer) writeEmptyArray() {
	u.writeArrayStart()
	u.writeArrayEnd()
}

func (u *ror2Writer) WriteInt(v int) {
	u.Writer.Int(v)
}

func (u *ror2Writer) WriteInt32(v int32) {
	u.Writer.Int32(v)
}

func (u *ror2Writer) WriteInt64(v int64) {
	u.Writer.Int64(v)
}

func (u *ror2Writer) WriteFloat32(v float32) {
	u.WriteFloat64(float64(v))
}

func (u *ror2Writer) WriteFloat64(v float64) {
	switch {
	case v > math.MaxFloat64:
		u.Writer.RawString("Infinity")
	case v < -math.MaxFloat64:
		u.Writer.RawString("-Infinity")
	case v != v:
		u.Writer.RawString("NaN")
	default:
		u.Writer.Float64(v)
	}
}

func (u *ror2Writer) WriteBool(v bool) {
	u.Writer.Bool(v)
}

func (u *ror2Writer) WriteString(v string) {
	if len(v) == 0 {
		u.Writer.RawString(emptyString)
	} else {
		u.Writer.RawString(u.stringEscaper(v))
	}
}

func (u *ror2Writer) WriteBytes(v []byte) {
	u.WriteString(string(v))
}
