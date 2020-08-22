package restliencoding

import (
	"io"

	"github.com/mailru/easyjson/jwriter"
)

type compactJsonEncoder struct {
	jwriter.Writer
}

func NewCompactJsonEncoder() *Encoder {
	return &Encoder{encoder: &compactJsonEncoder{}}
}

func (c *compactJsonEncoder) WriteObjectStart() {
	c.Writer.RawByte('{')
}

func (c *compactJsonEncoder) WriteFieldName(name string) {
	c.Writer.RawByte('"')
	c.Writer.RawString(name)
	c.Writer.RawByte('"')
}

func (c *compactJsonEncoder) WriteFieldNameDelimiter() {
	c.Writer.RawByte(':')
}

func (c *compactJsonEncoder) WriteFieldDelimiter() {
	c.Writer.RawByte(',')
}

func (c *compactJsonEncoder) WriteObjectEnd() {
	c.Writer.RawByte('}')
}

func (c *compactJsonEncoder) WriteMapStart() {
	c.WriteObjectStart()
}

func (c *compactJsonEncoder) WriteMapKey(key string) {
	c.WriteFieldName(key)
}

func (c *compactJsonEncoder) WriteMapKeyDelimiter() {
	c.WriteFieldNameDelimiter()
}

func (c *compactJsonEncoder) WriteMapEntryDelimiter() {
	c.WriteFieldDelimiter()
}

func (c *compactJsonEncoder) WriteMapEnd() {
	c.WriteObjectEnd()
}

func (c *compactJsonEncoder) WriteArrayStart() {
	c.Writer.RawByte('[')
}

func (c *compactJsonEncoder) WriteArrayItemDelimiter() {
	c.Writer.RawByte(',')
}

func (c *compactJsonEncoder) WriteArrayEnd() {
	c.Writer.RawByte(']')
}

func (c *compactJsonEncoder) SubEncoder() encoder {
	return c
}

func (c *compactJsonEncoder) Bytes(bytes []byte) {
	c.Writer.String(string(bytes))
}

func (c *compactJsonEncoder) Finalize() string {
	data, _ := c.BuildBytes()
	return string(data)
}

func (c *compactJsonEncoder) ReadCloser() io.ReadCloser {
	rc, _ := c.Writer.ReadCloser()
	return rc
}
