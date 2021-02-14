package restlicodec

import (
	"bytes"
	"io"
)

type protobufWriter struct {
	bytes.Buffer
}

const (
	pbMapOrdinal                  = '\x00'
	pbListOrdinal                 = '\x01'
	pbString_literalOrdinal       = '\x02'
	pbStringReferenceOrdinal      = '\x03'
	pbIntegerOrdinal              = '\x04'
	pbLongOrdinal                 = '\x05'
	pbFloatOrdinal                = '\x06'
	pbDoubleOrdinal               = '\x07'
	pbBooleanTrueOrdinal          = '\x08'
	pbBooleanFalseOrdinal         = '\x09'
	pbRawBytesOrdinal             = '\x0A'
	pbNullOrdinal                 = '\x0B'
	pbAscii_string_literalOrdinal = '\x14'
	pbFixed_floatOrdinal          = '\x15'
	pbFixed_doubleOrdinal         = '\x16'
	pbSerializeFixedFloats        = "serialize_fixed_floats"
	pbSymbolTableParamName        = "symbol-table"
)

var _ PrimitiveWriter = new(protobufWriter)
var _ Writer = new(protobufWriter)
var _ WriteCloser = new(protobufWriter)
var _ rawWriter = new(protobufWriter)

// NewProtobufWriter returns a WriteCloser that serializes objects using a Protocol buffer encoder.
// https://github.com/linkedin/rest.li/blob/master/data/src/main/java/com/linkedin/data/codec/ProtobufDataCodec.java
func NewProtobufWriter() WriteCloser {
	return newGenericWriter(new(protobufWriter), nil)
}

// PrimitiveWriter

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

func (c *protobufWriter) WriteInt32(v int32) {
	// c.WriteVarInt(int64(v))
}
func (c *protobufWriter) WriteInt64(v int64) {
	// c.WriteVarInt(v)
}
func (c *protobufWriter) WriteFloat32(v float32) {
	// c.WriteVarDouble(float64(v))
}
func (c *protobufWriter) WriteFloat64(v float64) {
	// c.WriteVarDouble(v)
}
func (c *protobufWriter) WriteString(v string) {
	// if v in symbol_table
	// c.WriteOrdinal(pbStringReferenceOrdinal)
	// c.WriteVarInt(symbol_id)

}

// WriteBytes
func (c *protobufWriter) WriteBytes(v []byte) {
	c.WriteOrdinal(pbRawBytesOrdinal)
	c.WriteVarInt(int64(len(v)))
	c.Buffer.Write(v)
}

func (c *protobufWriter) WriteOrdinal(v byte) {
	c.Buffer.WriteByte(v)
}

// The following are exposed directly by jwriter.Writer
func (c *protobufWriter) RawByte(v byte) {
	c.Buffer.WriteByte(v)
}
func (c *protobufWriter) Raw(v []byte, _ error) {
	c.Buffer.Write(v)
}
func (c *protobufWriter) RawString(v string) {
	c.Buffer.WriteString(v)
}
func (c *protobufWriter) BuildBytes(...[]byte) ([]byte, error) {}
func (c *protobufWriter) ReadCloser() (io.ReadCloser, error)   {}
func (c *protobufWriter) Size() int {
	return c.Buffer.Len()
}

//  https://developers.google.com/protocol-buffers/docs/encoding#varints
func (c *protobufWriter) WriteVarInt(v int64)      {}
func (c *protobufWriter) WriteVarDouble(v float64) {}
