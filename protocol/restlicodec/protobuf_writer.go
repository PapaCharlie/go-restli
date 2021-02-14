package restlicodec

import (
	"bytes"
	"io"
	"io/ioutil"
)

type protobufWriter struct {
	bytes.Buffer
}

const (
	pbMapOrdinal                = '\x00'
	pbListOrdinal               = '\x01'
	pbStringLiteralOrdinal      = '\x02'
	pbStringReferenceOrdinal    = '\x03'
	pbIntegerOrdinal            = '\x04'
	pbLongOrdinal               = '\x05'
	pbFloatOrdinal              = '\x06'
	pbDoubleOrdinal             = '\x07'
	pbBooleanTrueOrdinal        = '\x08'
	pbBooleanFalseOrdinal       = '\x09'
	pbRawBytesOrdinal           = '\x0A'
	pbNullOrdinal               = '\x0B'
	pbASCIIStringLiteralOrdinal = '\x14'
	pbFixedFloatOrdinal         = '\x15'
	pbFixedDoubleOrdinal        = '\x16'
	pbSerializeFixedFloats      = "serialize_fixed_floats"
	pbSymbolTableParamName      = "symbol-table"
)

var _ PrimitiveWriter = new(protobufWriter) // assert PrimitiveWriter
var _ Writer = NewProtobufWriter()          // assert Writer
var _ WriteCloser = NewProtobufWriter()     // assert WriteCloser
var _ rawWriter = new(protobufWriter)       // assert rawWriter

// NewProtobufWriter returns a WriteCloser that serializes objects using a Protocol buffer encoder modeled after this codec:
// https://github.com/linkedin/rest.li/blob/master/data/src/main/java/com/linkedin/data/codec/ProtobufDataCodec.java
func NewProtobufWriter() WriteCloser {
	return newGenericWriter(new(protobufWriter), nil)
}

// NewProtobufWriterWithExcludedFields returns a WriteCloser as from NewProtobufWriter(), excluding any fields matched by the given PathSpec.
func NewProtobufWriterWithExcludedFields(excludedFields PathSpec) WriteCloser {
	return newGenericWriter(new(protobufWriter), excludedFields)
}

func (c *protobufWriter) writeMapStart() {
	c.WriteOrdinal(pbMapOrdinal)
	// somehow need to output length before it is known
}

func (c *protobufWriter) writeKey(key string) {
	c.WriteString(key)
}

func (c *protobufWriter) writeKeyDelimiter()   {} // no key delimiter
func (c *protobufWriter) writeEntryDelimiter() {} // no entry delimiter
func (c *protobufWriter) writeMapEnd()         {} // no end marker

func (c *protobufWriter) writeArrayStart() {
	c.WriteOrdinal(pbListOrdinal)
	// somehow need to output length before it is known
}

func (c *protobufWriter) writeArrayItemDelimiter() {} // no array item delimiter

func (c *protobufWriter) writeArrayEnd() {} // no end marker

func (c *protobufWriter) WriteBool(v bool) {
	if v {
		c.WriteOrdinal(pbBooleanTrueOrdinal)
	} else {
		c.WriteOrdinal(pbBooleanFalseOrdinal)
	}
}

// WriteInt32 implements PrimitiveWriter
func (c *protobufWriter) WriteInt32(v int32) {
	c.WriteOrdinal(pbIntegerOrdinal)
	c.WriteVarInt(int64(v))
}

// WriteInt64 implements PrimitiveWriter
func (c *protobufWriter) WriteInt64(v int64) {
	c.WriteOrdinal(pbLongOrdinal)
	c.WriteVarInt(v)
}

// WriteFloat32 implements PrimitiveWriter
func (c *protobufWriter) WriteFloat32(v float32) {
	c.WriteOrdinal(pbFloatOrdinal)
	c.WriteVarDouble(float64(v))
}

// WriteFloat64 implements PrimitiveWriter
func (c *protobufWriter) WriteFloat64(v float64) {
	c.WriteOrdinal(pbDoubleOrdinal)
	c.WriteVarDouble(v)
}
func (c *protobufWriter) WriteString(v string) {
	// TODO add string symbol table (pbStringReferenceOrdinal)
	c.WriteOrdinal(pbStringLiteralOrdinal)
	b := []byte(v)
	c.WriteVarInt(int64(len(b)))
	c.Buffer.Write(b)
}

// WriteBytes implements PrimitiveWriter
func (c *protobufWriter) WriteBytes(v []byte) {
	c.WriteOrdinal(pbRawBytesOrdinal)
	c.WriteVarInt(int64(len(v)))
	c.Buffer.Write(v)
}

// WriteOrdinal writes the given ordinal byte to the buffer
func (c *protobufWriter) WriteOrdinal(v byte) {
	c.Buffer.WriteByte(v)
}

func (c *protobufWriter) RawByte(v byte) {
	c.Buffer.WriteByte(v)
}
func (c *protobufWriter) Raw(v []byte, _ error) {
	c.Buffer.Write(v)
}
func (c *protobufWriter) RawString(v string) {
	c.Buffer.WriteString(v)
}
func (c *protobufWriter) BuildBytes(...[]byte) ([]byte, error) {
	return c.Buffer.Bytes(), nil
}
func (c *protobufWriter) ReadCloser() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewReader(c.Buffer.Bytes())), nil
}
func (c *protobufWriter) Size() int {
	return c.Buffer.Len()
}

//  https://developers.google.com/protocol-buffers/docs/encoding#varints
func (c *protobufWriter) WriteVarInt(v int64)      {}
func (c *protobufWriter) WriteVarDouble(v float64) {}
