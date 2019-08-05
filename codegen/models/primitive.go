package models

import (
	"encoding/json"
	"log"

	"github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

type PrimitiveModel [2]string

var (
	IntPrimitive     = PrimitiveModel{"int", "int32"}
	LongPrimitive    = PrimitiveModel{"long", "int64"}
	FloatPrimitive   = PrimitiveModel{"float", "float32"}
	DoublePrimitive  = PrimitiveModel{"double", "float64"}
	BooleanPrimitive = PrimitiveModel{"boolean", "bool"}
	StringPrimitive  = PrimitiveModel{"string", "string"}
)

func ParsePrimitiveModel(p string) *PrimitiveModel {
	var primitive PrimitiveModel
	switch p {
	case IntPrimitive[0]:
		primitive = IntPrimitive
	case LongPrimitive[0]:
		primitive = LongPrimitive
	case FloatPrimitive[0]:
		primitive = FloatPrimitive
	case DoublePrimitive[0]:
		primitive = DoublePrimitive
	case BooleanPrimitive[0]:
		primitive = BooleanPrimitive
	case StringPrimitive[0]:
		primitive = StringPrimitive
	default:
		return nil
	}
	return &primitive
}

func (p *PrimitiveModel) UnmarshalJSON(data []byte) error {
	var primitiveType string
	if err := json.Unmarshal(data, &primitiveType); err != nil {
		return errors.WithStack(err)
	}

	parsedPrimitive := ParsePrimitiveModel(primitiveType)
	if parsedPrimitive != nil {
		*p = *parsedPrimitive
		return nil
	} else {
		return errors.Errorf("not a valid primitive type: |%s|", primitiveType)
	}
}

func (p *PrimitiveModel) GoType() *Statement {
	return Id(p[1])
}

func (p *PrimitiveModel) GetLit(rawJson string) interface{} {
	unmarshal := func(v interface{}) interface{} {
		err := json.Unmarshal([]byte(rawJson), &v)
		if err != nil {
			log.Panicf("(%v) Illegal primitive: \"%s\" (%s)", p, rawJson, err)
		}
		return v
	}

	switch *p {
	case IntPrimitive:
		v := new(int32)
		unmarshal(v)
		return *v
	case LongPrimitive:
		v := new(int64)
		unmarshal(v)
		return *v
	case FloatPrimitive:
		v := new(float32)
		unmarshal(v)
		return *v
	case DoublePrimitive:
		v := new(float64)
		unmarshal(v)
		return *v
	case BooleanPrimitive:
		v := new(bool)
		unmarshal(v)
		return *v
	case StringPrimitive:
		v := new(string)
		unmarshal(v)
		return *v
	}

	log.Panicln("Illegal primitive", p)
	return nil
}

func (p *PrimitiveModel) encode(accessor *Statement) *Statement {
	return Id(codegen.Codec).Dot("Encode" + codegen.ExportedIdentifier(p[0])).Call(accessor)
}

func (p *PrimitiveModel) decode(accessor *Statement) *Statement {
	return Id(codegen.Codec).Dot("Decode"+codegen.ExportedIdentifier(p[0])).Call(Id("data"), Call(Op("*").Id(p[1])).Call(accessor))
}
