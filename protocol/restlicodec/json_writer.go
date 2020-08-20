package restlicodec

import (
	"github.com/mailru/easyjson/jwriter"
)

type compactJsonWriter struct {
	jwriter.Writer
}

// NewCompactJsonWriter returns a WriteCloser that serializes objects using JSON. This representation has no extraneous
// whitespace and is intended for wire transport.
func NewCompactJsonWriter() WriteCloser {
	return newGenericWriter(new(compactJsonWriter))
}

func (c *compactJsonWriter) writeMapStart() {
	c.Writer.RawByte('{')
}

func (c *compactJsonWriter) writeKey(key string) {
	c.Writer.RawByte('"')
	c.Writer.RawString(key)
	c.Writer.RawByte('"')
}

func (c *compactJsonWriter) writeKeyDelimiter() {
	c.Writer.RawByte(':')
}

func (c *compactJsonWriter) writeEntryDelimiter() {
	c.Writer.RawByte(',')
}

func (c *compactJsonWriter) writeMapEnd() {
	c.Writer.RawByte('}')
}

func (c *compactJsonWriter) writeArrayStart() {
	c.Writer.RawByte('[')
}

func (c *compactJsonWriter) writeArrayItemDelimiter() {
	c.Writer.RawByte(',')
}

func (c *compactJsonWriter) writeArrayEnd() {
	c.Writer.RawByte(']')
}

func (c *compactJsonWriter) WriteInt32(v int32) {
	c.Writer.Int32(v)
}

func (c *compactJsonWriter) WriteInt64(v int64) {
	c.Writer.Int64(v)
}

func (c *compactJsonWriter) WriteFloat32(v float32) {
	c.Writer.Float32(v)
}

func (c *compactJsonWriter) WriteFloat64(v float64) {
	c.Writer.Float64(v)
}

func (c *compactJsonWriter) WriteBool(v bool) {
	c.Writer.Bool(v)
}

func (c *compactJsonWriter) WriteString(v string) {
	c.Writer.String(v)
}

func (c *compactJsonWriter) WriteBytes(v []byte) {
	c.String(string(v))
}

type prettyJsonWriter struct {
	compactJsonWriter
	indent string
}

// NewPrettyJsonWriter returns a WriteCloser that serializes objects using JSON. This representation delimits fields and
// array items using newlines and provides indentation for nested objects. It generates a lot of unnecessary bytes and
// is intended primarily for debugging or human-readability purposes.
func NewPrettyJsonWriter() WriteCloser {
	return newGenericWriter(new(prettyJsonWriter))
}

func (p *prettyJsonWriter) incrementIndent() {
	p.indent += "  "
}

func (p *prettyJsonWriter) decrementIndent() {
	p.indent = p.indent[:len(p.indent)-2]
}

func (p *prettyJsonWriter) writeMapStart() {
	p.incrementIndent()
	p.RawString("{\n")
}

func (p *prettyJsonWriter) writeKey(key string) {
	p.Writer.RawString(p.indent)
	p.compactJsonWriter.writeKey(key)
}

func (p *prettyJsonWriter) writeKeyDelimiter() {
	p.Writer.RawString(": ")
}

func (p *prettyJsonWriter) writeEntryDelimiter() {
	p.Writer.RawString(",\n")
}

func (p *prettyJsonWriter) writeMapEnd() {
	p.decrementIndent()
	p.RawString("\n" + p.indent + "}")
}

func (p *prettyJsonWriter) writeArrayStart() {
	p.incrementIndent()
	p.RawString("[\n" + p.indent)
}

func (p *prettyJsonWriter) writeArrayItemDelimiter() {
	p.Writer.RawString(",\n" + p.indent)
}

func (p *prettyJsonWriter) writeArrayEnd() {
	p.decrementIndent()
	p.RawString("\n" + p.indent + "]")
}
