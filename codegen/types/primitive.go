package types

import (
	"encoding/json"
	"reflect"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

type PrimitiveType struct {
	Type        string
	newInstance func() interface{}
}

var (
	Int32Primitive   = PrimitiveType{Type: "int32", newInstance: func() interface{} { return new(int32) }}
	Int64Primitive   = PrimitiveType{Type: "int64", newInstance: func() interface{} { return new(int64) }}
	Float32Primitive = PrimitiveType{Type: "float32", newInstance: func() interface{} { return new(float32) }}
	Float64Primitive = PrimitiveType{Type: "float64", newInstance: func() interface{} { return new(float64) }}
	BoolPrimitive    = PrimitiveType{Type: "bool", newInstance: func() interface{} { return new(bool) }}
	StringPrimitive  = PrimitiveType{Type: "string", newInstance: func() interface{} { return new(string) }}
	BytesPrimitive   = PrimitiveType{Type: "bytes", newInstance: func() interface{} { return new([]byte) }}
)

var PrimitiveTypes = []PrimitiveType{
	Int32Primitive,
	Int64Primitive,
	Float32Primitive,
	Float64Primitive,
	BoolPrimitive,
	StringPrimitive,
	BytesPrimitive,
}

func (p *PrimitiveType) UnmarshalJSON(data []byte) error {
	var primitiveType string
	if err := json.Unmarshal(data, &primitiveType); err != nil {
		return errors.WithStack(err)
	}

	for _, pt := range PrimitiveTypes {
		if primitiveType == pt.Type {
			*p = pt
			return nil
		}
	}

	return errors.Errorf("Unknown type: %s", primitiveType)
}

func (p *PrimitiveType) IsBytes() bool {
	return p.Type == "bytes"
}

func (p *PrimitiveType) Cast(accessor *Statement) *Statement {
	var cast *Statement
	if p.IsBytes() {
		cast = Index().Byte()
	} else {
		cast = Id(p.Type)
	}
	return cast.Call(accessor)
}

func (p *PrimitiveType) GoType() *Statement {
	if p.IsBytes() {
		return Index().Byte()
	} else {
		return Id(p.Type)
	}
}

func getLitBytesValues(rawJson string) *Statement {
	var v string
	if err := json.Unmarshal([]byte(rawJson), &v); err != nil {
		utils.Logger.Panicf("(%v) Illegal primitive literal: \"%s\" (%s)", BytesPrimitive, rawJson, err)
	}
	return ValuesFunc(func(def *Group) {
		for _, c := range v {
			def.LitRune(c)
		}
	})
}

func (p *PrimitiveType) getLit(rawJson string) *Statement {
	if p.IsBytes() {
		return Index().Byte().Add(getLitBytesValues(rawJson))
	} else {
		v := p.newInstance()

		err := json.Unmarshal([]byte(rawJson), v)
		if err != nil {
			utils.Logger.Panicf("(%v) Illegal primitive literal: \"%s\" (%s)", p, rawJson, err)
		}
		return Lit(reflect.ValueOf(v).Elem().Interface())
	}
}

func (p *PrimitiveType) zeroValueLit() *Statement {
	if p.IsBytes() {
		return Nil()
	} else {
		return Lit(reflect.ValueOf(p.newInstance()).Elem().Interface())
	}
}

func (p *PrimitiveType) exportedName() string {
	return utils.ExportedIdentifier(p.Type)
}

func (p *PrimitiveType) WriterName() string {
	return "Write" + p.exportedName()
}

func (p *PrimitiveType) ReaderName() string {
	return "Read" + p.exportedName()
}

func (p *PrimitiveType) HasherName() string {
	return "Add" + p.exportedName()
}

func (p *PrimitiveType) NewPrimitiveMarshaler(accessor Code) Code {
	return Qual(utils.RestLiCodecPackage, "MarshalerFunc").Call(Func().Params(Add(Writer).Add(WriterQual)).Error().BlockFunc(func(def *Group) {
		def.Add(Writer.Write(RestliType{Primitive: p}, Writer, accessor))
		def.Return(Nil())
	}))
}

func (p *PrimitiveType) NewPrimitiveUnmarshaler(accessor Code) Code {
	return Qual(utils.RestLiCodecPackage, "UnmarshalerFunc").Call(Func().Params(Add(Reader).Add(ReaderQual)).Params(Err().Error()).BlockFunc(func(def *Group) {
		def.Add(Reader.Read(RestliType{Primitive: p}, Reader, accessor))
		def.Return(Err())
	}))
}
