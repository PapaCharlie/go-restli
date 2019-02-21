package models

import (
	"encoding/json"
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
	"log"
)

type Primitive string

const (
	Int     = Primitive("int32")
	Long    = Primitive("int64")
	Float   = Primitive("float32")
	Double  = Primitive("float64")
	Boolean = Primitive("bool")
	String  = Primitive("string")
	Bytes   = Primitive("[]byte")
)

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
	case "bytes":
		primitive = Bytes
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

func (p *Primitive) GetLit(rawJson string) interface{} {
	unmarshall := func(v interface{}) interface{} {
		err := json.Unmarshal([]byte(rawJson), &v)
		if err != nil {
			log.Panicln("Illegal primitive", err)
		}
		return v
	}

	switch *p {
	case Int:
		v := new(int32)
		unmarshall(v)
		return *v
	case Long:
		v := new(int64)
		unmarshall(v)
		return *v
	case Float:
		v := new(float32)
		unmarshall(v)
		return *v
	case Double:
		v := new(float64)
		unmarshall(v)
		return *v
	case Boolean:
		v := new(bool)
		unmarshall(v)
		return *v
	case String:
		v := new(string)
		unmarshall(v)
		return *v
	case Bytes:
		v := new([]byte)
		unmarshall(v)
		return *v
	}
	log.Panicln("Illegal primitive", p)
	return nil
}
