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

func (p *Primitive) UnmarshalJSON(data []byte) error {
	var primitiveType string
	if err := json.Unmarshal(data, &primitiveType); err != nil {
		return errors.WithStack(err)
	}

	switch primitiveType {
	case "int":
		*p = Int
	case "long":
		*p = Long
	case "float":
		*p = Float
	case "double":
		*p = Double
	case "boolean":
		*p = Boolean
	case "string":
		*p = String
	default:
		return errors.Errorf("not a valid primitive type: |%s|", primitiveType)
	}
	return nil
}

func (p *Primitive) GoType() *jen.Statement {
	return jen.Empty().Add(jen.Id(string(*p)))
}
