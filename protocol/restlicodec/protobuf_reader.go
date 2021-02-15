package restlicodec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type protobufReader struct {
	missingFieldsTracker
	buf     *bytes.Buffer
	started bool
	options ProtobufOptions
}

var _ PrimitiveReader = new(protobufReader) // assert PrimitiveReader

// NewProtobufReader returns a Reader that deserializes objects from a Protocol buffer encoded bytestring.
// Ref: https://github.com/linkedin/rest.li/blob/master/data/src/main/java/com/linkedin/data/codec/ProtobufDataCodec.java
func NewProtobufReader(data []byte) Reader {
	p := new(protobufReader)
	p.buf = new(bytes.Buffer)
	p.buf.Write(data)
	return p
}

// AtInputStart implements Reader
func (p *protobufReader) AtInputStart() bool {
	return !p.started
}

// ReadMap implements Reader
func (p *protobufReader) ReadMap(mapReader MapReader) (err error) {
	ord, err := p.ReadOrdinal()
	if err != nil {
		return err
	}
	if pbMapOrdinal != ord {
		return &DeserializationError{
			Err: fmt.Errorf("unexpected token in protobuf stream %v - expected %v", ord, pbMapOrdinal),
		}
	}
	count, err := p.readVarint()
	if err != nil {
		return err
	}
	for i := int64(0); i < count; i++ {
		fieldName, err := p.ReadString()
		if err != nil {
			return err
		}
		err = mapReader(p, fieldName)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadArray implements Reader
func (p *protobufReader) ReadArray(arrayReader ArrayReader) error {
	ord, err := p.ReadOrdinal()
	if err != nil {
		return err
	}
	if pbListOrdinal != ord {
		return &DeserializationError{
			Err: fmt.Errorf("unexpected token in protobuf stream %v - expected %v", ord, pbListOrdinal),
		}
	}
	count, err := p.readVarint()
	if err != nil {
		return err
	}
	for i := int64(0); i < count; i++ {
		err = arrayReader(p)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadInterface implements Reader
func (p *protobufReader) ReadInterface() (interface{}, error) {
	return nil, &DeserializationError{
		Err: fmt.Errorf("this isnt finished"),
	}
}

// ReadRawBytes implements Reader
func (p *protobufReader) ReadRawBytes() ([]byte, error) {
	return nil, &DeserializationError{
		Err: fmt.Errorf("this isnt finished"),
	}
}

// Skip implements Reader
func (p *protobufReader) Skip() error {
	return &DeserializationError{
		Err: fmt.Errorf("this isnt finished"),
	}
}

// ReadBool implements Reader
func (p *protobufReader) ReadBool() (bool, error) {
	v, err := p.ReadOrdinal()
	if err != nil {
		return false, &DeserializationError{Err: err}
	}
	switch v {
	case pbBooleanTrueOrdinal:
		return true, nil
	case pbBooleanFalseOrdinal:
		return false, nil
	default:
		return false, &DeserializationError{
			Err: fmt.Errorf("unexpected token in protobuf stream %v - expected %v or %v", v, pbBooleanTrueOrdinal, pbBooleanFalseOrdinal),
		}
	}
}

// ReadInt32 implements Reader
func (p *protobufReader) ReadInt32() (int32, error) {
	ord, err := p.ReadOrdinal()
	if err != nil {
		return 0, err
	}
	if pbIntegerOrdinal != ord {
		return 0, &DeserializationError{
			Err: fmt.Errorf("unexpected token in protobuf stream %v - expected %v", ord, pbIntegerOrdinal),
		}
	}
	val, err := p.readVarint()
	return int32(val), err
}

// ReadInt64 implements Reader
func (p *protobufReader) ReadInt64() (int64, error) {
	ord, err := p.ReadOrdinal()
	if err != nil {
		return 0, err
	}
	if pbLongOrdinal != ord {
		return 0, &DeserializationError{
			Err: fmt.Errorf("unexpected token in protobuf stream %v - expected %v", ord, pbLongOrdinal),
		}
	}
	val, err := p.readVarint()
	return val, err
}

func (p *protobufReader) readVarint() (int64, error) {
	var br io.ByteReader = p.buf
	val, err := binary.ReadVarint(br)
	if err != nil {
		return 0, err
	}
	return val, nil
}
func (p *protobufReader) readUvarint() (uint64, error) {
	var br io.ByteReader = p.buf
	val, err := binary.ReadUvarint(br)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ReadFloat32 implements Reader
func (p *protobufReader) ReadFloat32() (float32, error) {
	ord, err := p.ReadOrdinal()
	if err != nil {
		return 0, err
	}
	if pbFloatOrdinal != ord {
		return 0, &DeserializationError{
			Err: fmt.Errorf("unexpected token in protobuf stream %v - expected %v", ord, pbFloatOrdinal),
		}
	}
	if p.options.fixedWidthFloat32 {
		val, err := p.readFixed32()
		if err != nil {
			return 0, err
		}
		return math.Float32frombits(val), nil
	}
	val, err := p.readUvarint()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(uint32(val)), nil
}

// ReadFloat64 implements Reader
func (p *protobufReader) ReadFloat64() (float64, error) {
	ord, err := p.ReadOrdinal()
	if err != nil {
		return 0, err
	}
	if pbDoubleOrdinal != ord {
		return 0, &DeserializationError{
			Err: fmt.Errorf("unexpected token in protobuf stream %v - expected %v", ord, pbDoubleOrdinal),
		}
	}
	if p.options.fixedWidthFloat64 {
		val, err := p.readFixed64()
		if err != nil {
			return 0, err
		}
		return math.Float64frombits(val), nil
	}
	val, err := p.readUvarint()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(uint64(val)), nil
}

func (p *protobufReader) readFixed32() (uint32, error) {
	var (
		err   error
		b     byte
		value uint32 = 0
	)
	b, err = p.buf.ReadByte()
	if err != nil {
		return 0, err
	}
	value = value & uint32(b)
	b, err = p.buf.ReadByte()
	if err != nil {
		return 0, err
	}
	value = value & (uint32(b) << 8)
	b, err = p.buf.ReadByte()
	if err != nil {
		return 0, err
	}
	value = value & (uint32(b) << 16)
	b, err = p.buf.ReadByte()
	if err != nil {
		return 0, err
	}
	value = value & (uint32(b) << 24)
	return value, nil
}
func (p *protobufReader) readFixed64() (uint64, error) {
	var (
		err   error
		b     byte
		value uint64 = 0
	)
	b, err = p.buf.ReadByte()
	if err != nil {
		return 0, err
	}
	value = value & uint64(b)
	b, err = p.buf.ReadByte()
	if err != nil {
		return 0, err
	}
	value = value | (uint64(b) << 8)
	b, err = p.buf.ReadByte()
	if err != nil {
		return 0, err
	}
	value = value | (uint64(b) << 16)
	b, err = p.buf.ReadByte()
	if err != nil {
		return 0, err
	}
	value = value | (uint64(b) << 24)
	b, err = p.buf.ReadByte()
	if err != nil {
		return 0, err
	}
	value = value | (uint64(b) << 32)
	b, err = p.buf.ReadByte()
	if err != nil {
		return 0, err
	}
	value = value | (uint64(b) << 40)
	b, err = p.buf.ReadByte()
	if err != nil {
		return 0, err
	}
	value = value | (uint64(b) << 48)
	b, err = p.buf.ReadByte()
	if err != nil {
		return 0, err
	}
	value = value | (uint64(b) << 56)
	return value, nil
}

func (p *protobufReader) ReadString() (string, error) {
	ord, err := p.ReadOrdinal()
	if err != nil {
		return "", err
	}
	if pbStringReferenceOrdinal == ord {
		return p.readStringReference()
	}
	if pbStringLiteralOrdinal == ord {
		return p.readStringLiteral()
	}
	return "", &DeserializationError{
		Err: fmt.Errorf("unexpected token in protobuf stream %v - expected %v or %v", ord, pbStringLiteralOrdinal, pbStringReferenceOrdinal),
	}
}

func (p *protobufReader) readStringLiteral() (string, error) {
	size, err := p.readVarint()
	if err != nil {
		return "", err
	}
	var strBuf bytes.Buffer
	for i := int64(0); i < size; i++ {
		b, e := p.buf.ReadByte()
		if e != nil {
			return "", e
		}
		e = strBuf.WriteByte(b)
		if e != nil {
			return "", e
		}
	}
	return strBuf.String(), nil
}
func (p *protobufReader) readStringReference() (string, error) {
	val, err := p.readVarint()
	if err != nil {
		return "", err
	}
	if p.options.symbolTable == nil {
		return "", &DeserializationError{
			Err: fmt.Errorf("protobuf string dereference without symbolTable"),
		}
	}
	symbolID := int(val)
	symbol, found := p.options.symbolTable.GetSymbolName(symbolID)
	if found {
		return symbol, nil
	}
	return "", &DeserializationError{
		Err: fmt.Errorf("protobuf string symbol not in symbolTable: %v", symbolID),
	}
}

func (p *protobufReader) ReadOrdinal() (pbOrdinal, error) {
	p.started = true
	b, err := p.buf.ReadByte()
	if err != nil {
		return pbOrdinal(b), err
	}
	return pbOrdinal(b), err
}

func (p *protobufReader) PeakOrdinal() (pbOrdinal, error) {
	o, e := p.ReadOrdinal()
	if e == nil {
		p.buf.UnreadByte()
	}
	return o, e
}

func (p *protobufReader) PeakByte() (byte, error) {
	b, err := p.buf.ReadByte()
	if err == nil {
		p.buf.UnreadByte()
	}
	return b, err
}

func (p *protobufReader) ReadBytes() ([]byte, error) {
	b := p.buf.Bytes()
	return b, nil
}
