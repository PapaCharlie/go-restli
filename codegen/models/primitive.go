package models

import (
	"encoding/json"
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

const (
	Int     = Primitive("int32")
	Long    = Primitive("int64")
	Float   = Primitive("float32")
	Double  = Primitive("float64")
	Boolean = Primitive("bool")
	String  = Primitive("string")
)

type Primitive string

func ParsePrimitive(p string) *Primitive {
	var primitive Primitive
	switch p {
	case "int":
		primitive = Int
	case "long":
		primitive = Long
	case "float":
		primitive = Float
	case "double":
		primitive = Double
	case "boolean":
		primitive = Boolean
	case "string":
		primitive = String
	default:
		return nil
	}
	return &primitive
}

func (p *Primitive) UnmarshalJSON(data []byte) error {
	var primitiveType string
	if err := json.Unmarshal(data, &primitiveType); err != nil {
		return errors.WithStack(err)
	}

	parsedPrimitive := ParsePrimitive(primitiveType)
	if parsedPrimitive != nil {
		*p = *parsedPrimitive
		return nil
	} else {
		return errors.Errorf("not a valid primitive type: |%s|", primitiveType)
	}
}

func (p *Primitive) GoType() *jen.Statement {
	return jen.Empty().Add(jen.Id(string(*p)))
}
