package restlicodec

type prettyJsonEncoder struct {
	compactJsonEncoder
	indent string
}

func NewPrettyJsonEncoder() *Encoder {
	return &Encoder{encoder: &prettyJsonEncoder{}}
}

func (p *prettyJsonEncoder) incrementIndent() {
	p.indent += "  "
}

func (p *prettyJsonEncoder) decrementIndent() {
	p.indent = p.indent[:len(p.indent)-2]
}

func (p *prettyJsonEncoder) WriteObjectStart() {
	p.incrementIndent()
	p.RawString("{\n")
}

func (p *prettyJsonEncoder) WriteFieldName(name string) {
	p.Writer.RawString(p.indent)
	p.compactJsonEncoder.WriteFieldName(name)
}

func (p *prettyJsonEncoder) WriteFieldNameDelimiter() {
	p.Writer.RawString(": ")
}

func (p *prettyJsonEncoder) WriteFieldDelimiter() {
	p.Writer.RawString(",\n")
}

func (p *prettyJsonEncoder) WriteObjectEnd() {
	p.decrementIndent()
	p.RawString("\n" + p.indent + "}")
}

func (p *prettyJsonEncoder) WriteMapStart() {
	p.WriteObjectStart()
}

func (p *prettyJsonEncoder) WriteMapKey(key string) {
	p.WriteFieldName(key)
}

func (p *prettyJsonEncoder) WriteMapKeyDelimiter() {
	p.WriteFieldNameDelimiter()
}

func (p *prettyJsonEncoder) WriteMapEntryDelimiter() {
	p.WriteFieldDelimiter()
}

func (p *prettyJsonEncoder) WriteMapEnd() {
	p.WriteObjectEnd()
}

func (p *prettyJsonEncoder) WriteArrayStart() {
	p.incrementIndent()
	p.RawString("[\n" + p.indent)
}

func (p *prettyJsonEncoder) WriteArrayItemDelimiter() {
	p.Writer.RawString(",\n" + p.indent)
}

func (p *prettyJsonEncoder) WriteArrayEnd() {
	p.decrementIndent()
	p.RawString("\n" + p.indent + "]")
}

func (p *prettyJsonEncoder) SubEncoder() encoder {
	return p
}
