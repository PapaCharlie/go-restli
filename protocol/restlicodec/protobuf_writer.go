package restlicodec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type protobufWriter struct {
	excludedFields PathSpec
	scope          []string
	options        ProtobufOptions
	buf            *bytes.Buffer
}

// options for protobufWriter and protobufReader
type ProtobufOptions struct {
	symbolTable       SymbolTable
	fixedWidthFloat32 bool
	fixedWidthFloat64 bool
	expVarintImpl     bool
}

type pbOrdinal byte

func (o pbOrdinal) String() string {
	omap := map[pbOrdinal]string{
		pbMapOrdinal:                "Map",
		pbListOrdinal:               "List",
		pbStringLiteralOrdinal:      "StringLiteral",
		pbStringReferenceOrdinal:    "StringReference",
		pbIntegerOrdinal:            "Integer",
		pbLongOrdinal:               "Long",
		pbFloatOrdinal:              "Float",
		pbDoubleOrdinal:             "Double",
		pbBooleanTrueOrdinal:        "BooleanTrue",
		pbBooleanFalseOrdinal:       "BooleanFalse",
		pbRawBytesOrdinal:           "RawBytes",
		pbNullOrdinal:               "Null",
		pbASCIIStringLiteralOrdinal: "ASCIIStringLiteral",
		pbFixedFloatOrdinal:         "FixedFloat",
		pbFixedDoubleOrdinal:        "FixedDouble",
	}
	n, f := omap[o]
	if !f {
		n = "Unknown"
	}
	return fmt.Sprintf("%x [%v]", byte(o), n)
}

const (
	pbMapOrdinal                pbOrdinal = 0x00
	pbListOrdinal               pbOrdinal = 0x01
	pbStringLiteralOrdinal      pbOrdinal = 0x02
	pbStringReferenceOrdinal    pbOrdinal = 0x03
	pbIntegerOrdinal            pbOrdinal = 0x04
	pbLongOrdinal               pbOrdinal = 0x05
	pbFloatOrdinal              pbOrdinal = 0x06
	pbDoubleOrdinal             pbOrdinal = 0x07
	pbBooleanTrueOrdinal        pbOrdinal = 0x08
	pbBooleanFalseOrdinal       pbOrdinal = 0x09
	pbRawBytesOrdinal           pbOrdinal = 0x0A
	pbNullOrdinal               pbOrdinal = 0x0B
	pbASCIIStringLiteralOrdinal pbOrdinal = 0x14
	pbFixedFloatOrdinal         pbOrdinal = 0x15
	pbFixedDoubleOrdinal        pbOrdinal = 0x16
)

// TODO remove these type assertions
var _ PrimitiveWriter = new(protobufWriter) // assert PrimitiveWriter
var _ Writer = new(protobufWriter)          // assert Writer
var _ WriteCloser = new(protobufWriter)     // assert WriteCloser

// NewProtobufWriter returns a WriteCloser that serializes objects using a Protocol buffer encoder modeled after this codec:
// https://github.com/linkedin/rest.li/blob/master/data/src/main/java/com/linkedin/data/codec/ProtobufDataCodec.java
func NewProtobufWriter() WriteCloser {
	out := new(protobufWriter)
	out.buf = new(bytes.Buffer)
	return out
}

// NewProtobufWriterWithExcludedFields returns a WriteCloser as from NewProtobufWriter(), excluding any fields matched by the given PathSpec.
func NewProtobufWriterWithExcludedFields(excludedFields PathSpec) WriteCloser {
	out := new(protobufWriter)
	out.buf = new(bytes.Buffer)
	out.excludedFields = excludedFields
	return out
}

func (p *protobufWriter) subWriter(key string) *protobufWriter {
	var out protobufWriter
	out = *p
	out.scope = copyAndAppend(p.scope, key)
	return &out
}
func (p *protobufWriter) pushBuffer() *bytes.Buffer {
	old := p.buf
	p.buf = new(bytes.Buffer)
	return old
}
func (p *protobufWriter) swapBuffer(next *bytes.Buffer) *bytes.Buffer {
	prev := p.buf
	p.buf = next
	return prev
}

func (p *protobufWriter) WriteMap(mapWriter MapWriter) (err error) {
	tmpBuffer := p.pushBuffer()
	var count int64
	sub := p.subWriter("")

	defer func() {
		tmpBuffer = p.swapBuffer(tmpBuffer)
		p.WriteOrdinal(pbMapOrdinal)
		p.WriteVarint(count)
		if count > 0 {
			p.WriteRawBytes(tmpBuffer.Bytes())
		}
	}()

	err = mapWriter(func(key string) Writer {
		if p.IsKeyExcluded(key) {
			return noopWriter
		}
		count++
		p.WriteString(key)
		sub.scope[len(sub.scope)-1] = key
		var writer Writer = sub
		return writer
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *protobufWriter) WriteArray(arrayWriter ArrayWriter) (err error) {
	tmpBuffer := p.pushBuffer()
	var count int64
	sub := p.subWriter(WildCard)

	defer func() {
		tmpBuffer = p.swapBuffer(tmpBuffer)
		p.WriteOrdinal(pbListOrdinal)
		p.WriteVarint(count)
		if count > 0 {
			p.WriteRawBytes(tmpBuffer.Bytes())
		}
	}()

	err = arrayWriter(func() Writer {
		return sub
	})
	if err != nil {
		return err
	}
	return nil
}

// IsKeyExcluded implements Writer
func (p *protobufWriter) IsKeyExcluded(key string) bool {
	p.scope = append(p.scope, key)
	excluded := p.excludedFields.Matches(p.scope)
	p.scope = p.scope[:len(p.scope)-1]
	return excluded
}

// SetScope implements Writer
func (p *protobufWriter) SetScope(scope ...string) Writer {
	var out protobufWriter
	out = *p
	out.scope = scope
	return &out
}

func (p *protobufWriter) WriteBool(v bool) {
	if v {
		p.WriteOrdinal(pbBooleanTrueOrdinal)
	} else {
		p.WriteOrdinal(pbBooleanFalseOrdinal)
	}
}

// WriteInt32 implements PrimitiveWriter
func (p *protobufWriter) WriteInt32(v int32) {
	p.WriteOrdinal(pbIntegerOrdinal)
	p.WriteVarint(int64(v))
}

// WriteInt64 implements PrimitiveWriter
func (p *protobufWriter) WriteInt64(v int64) {
	p.WriteOrdinal(pbLongOrdinal)
	p.WriteVarint(v)
}

// WriteFloat32 implements PrimitiveWriter
func (p *protobufWriter) WriteFloat32(v float32) {
	if p.options.fixedWidthFloat32 {
		p.WriteOrdinal(pbFixedFloatOrdinal)
		p.writeFixed32(uint32(math.Float64bits(float64(v))))
	} else {
		p.WriteOrdinal(pbFloatOrdinal)
		p.WriteVarDouble(float64(v))
	}
}

// WriteFloat64 implements PrimitiveWriter
func (p *protobufWriter) WriteFloat64(v float64) {
	if p.options.fixedWidthFloat32 {
		p.WriteOrdinal(pbFixedDoubleOrdinal)
		p.writeFixed64(math.Float64bits(float64(v)))
	} else {
		p.WriteOrdinal(pbDoubleOrdinal)
		p.WriteVarDouble(v)
	}
}

func (p *protobufWriter) WriteString(v string) {
	if p.options.symbolTable != nil {
		id, found := p.options.symbolTable.GetSymbolId(v)
		if found {
			p.WriteOrdinal(pbStringReferenceOrdinal)
			p.WriteVarint(int64(id))
			return
		}
	}
	p.WriteOrdinal(pbStringLiteralOrdinal)
	b := []byte(v)
	p.WriteVarint(int64(len(b)))
	p.buf.Write(b)
}

// WriteBytes implements PrimitiveWriter
func (p *protobufWriter) WriteBytes(v []byte) {
	p.WriteOrdinal(pbRawBytesOrdinal)
	p.WriteVarint(int64(len(v)))
	p.buf.Write(v)
}

// WriteOrdinal writes the given ordinal byte to the buffer
func (p *protobufWriter) WriteOrdinal(v pbOrdinal) {
	p.buf.WriteByte(byte(v))
}

func (p *protobufWriter) RawByte(v byte) {
	p.buf.WriteByte(v)
}
func (p *protobufWriter) Raw(v []byte, _ error) {
	p.buf.Write(v)
}
func (p *protobufWriter) RawString(v string) {
	p.buf.WriteString(v)
}
func (p *protobufWriter) WriteRawBytes(data []byte) {
	p.Raw(data, nil)
}
func (p *protobufWriter) BuildBytes(...[]byte) ([]byte, error) {
	return p.buf.Bytes(), nil
}
func (p *protobufWriter) ReadCloser() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(p.buf.Bytes()))
}
func (p *protobufWriter) Size() int {
	return p.buf.Len()
}
func (p *protobufWriter) Finalize() string {
	data, _ := p.BuildBytes()
	return string(data)
}

//  https://developers.google.com/protocol-buffers/docs/encoding#varints
func (p *protobufWriter) WriteVarint(v int64) {
	if p.options.expVarintImpl {
		// use the embedded varint implementation
		p.writeVarintImpl(uint64(v))
	} else {
		// use the Go built-in varint implementation
		var buf []byte = make([]byte, 4)
		s := binary.PutVarint(buf, 4)
		p.WriteRawBytes(buf[:s])
	}
}
func (p *protobufWriter) WriteUvarint(v uint64) {
	if p.options.expVarintImpl {
		// use the embedded varint implementation
		p.writeVarintImpl(v)
	} else {
		// use the Go built-in varint implementation
		var buf []byte = make([]byte, 16)
		s := binary.PutUvarint(buf, 16)
		p.WriteRawBytes(buf[:s])
	}
}
func (p *protobufWriter) WriteVarDouble(v float64) {
	uintBits := math.Float64bits(v)
	p.WriteUvarint(uintBits)
}
func (p *protobufWriter) writeFixed32(value uint32) {
	p.buf.WriteByte(byte(value & 0xFF))
	p.buf.WriteByte(byte((value >> 8) & 0xFF))
	p.buf.WriteByte(byte((value >> 16) & 0xFF))
	p.buf.WriteByte(byte((value >> 24) & 0xFF))
}
func (p *protobufWriter) writeFixed64(value uint64) {
	p.buf.WriteByte(byte(value & 0xFF))
	p.buf.WriteByte(byte((value >> 8) & 0xFF))
	p.buf.WriteByte(byte((value >> 16) & 0xFF))
	p.buf.WriteByte(byte((value >> 24) & 0xFF))
	p.buf.WriteByte(byte((value >> 32) & 0xFF))
	p.buf.WriteByte(byte((value >> 40) & 0xFF))
	p.buf.WriteByte(byte((value >> 48) & 0xFF))
	p.buf.WriteByte(byte((value >> 56) & 0xFF))
}
func (p *protobufWriter) writeVarintImpl(value uint64) {

}
