package codegen

import (
	"encoding/json"
	"reflect"

	. "github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

type PrimitiveType struct {
	Type        string
	newInstance func() interface{}
	empty       interface{}
}

var PrimitiveTypes = []PrimitiveType{
	{Type: "int32", newInstance: func() interface{} { return new(int32) }, empty: int32(0)},
	{Type: "int64", newInstance: func() interface{} { return new(int64) }, empty: int64(0)},
	{Type: "float32", newInstance: func() interface{} { return new(float32) }, empty: float32(0.0)},
	{Type: "float64", newInstance: func() interface{} { return new(float64) }, empty: float64(0.0)},
	{Type: "bool", newInstance: func() interface{} { return new(bool) }, empty: false},
	{Type: "string", newInstance: func() interface{} { return new(string) }, empty: ""},
	{Type: "bytes", newInstance: func() interface{} { return new([]byte) }, empty: nil},
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

func (p *PrimitiveType) Nil() *Statement {
	return Lit(p.empty)
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
		return Bytes()
	} else {
		return Id(p.Type)
	}
}

func (p *PrimitiveType) getLit(rawJson string) interface{} {
	v := p.newInstance()

	err := json.Unmarshal([]byte(rawJson), v)
	if err != nil {
		Logger.Panicf("(%v) Illegal primitive literal: \"%s\" (%s)", p, rawJson, err)
	}
	return reflect.ValueOf(v).Elem().Interface()
}

func (p *PrimitiveType) encode(accessor *Statement) *Statement {
	return Id(Codec).Dot("Encode" + ExportedIdentifier(p.Type)).Call(accessor)
}

func (p *PrimitiveType) decode(accessor *Statement) *Statement {
	return Id(Codec).Dot("Decode"+ExportedIdentifier(p.Type)).Call(Id("data"), Call(Op("*").Add(p.GoType())).Call(accessor))
}
