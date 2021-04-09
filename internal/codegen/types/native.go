package types

import (
	"encoding/json"

	. "github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

type NativeTyperef struct {
	Primitive              *PrimitiveType `json:"primitive"`
	TypePackage            string         `json:"typePackage"`
	TypeName               string         `json:"typeName"`
	ObjectFunctionsPackage string         `json:"objectFunctionsPackage"`
}

func (n *NativeTyperef) UnmarshalJSON(data []byte) error {
	type _t NativeTyperef
	err := json.Unmarshal(data, (*_t)(n))
	if err != nil {
		return err
	}
	if n.TypePackage == "" {
		return errors.New("go-restli: Native typeref must declare a type package")
	}
	if n.TypeName == "" {
		return errors.New("go-restli: Native typeref must declare a type name")
	}
	if n.ObjectFunctionsPackage == "" {
		return errors.New("go-restli: Native typeref must declare an object functions package")
	}
	return nil
}

func (n *NativeTyperef) GoType() *Statement {
	return Qual(n.TypePackage, n.TypeName)
}

func (n *NativeTyperef) objectFunction(kind string) *Statement {
	return Qual(n.ObjectFunctionsPackage, kind+n.TypeName)
}

func (n *NativeTyperef) Marshaler() *Statement {
	return n.objectFunction("Marshal")
}

func (n *NativeTyperef) Unmarshaler() *Statement {
	return n.objectFunction("Unmarshal")
}

func (n *NativeTyperef) Equals() *Statement {
	return n.objectFunction("Equals")
}

func (n *NativeTyperef) ComputeHash() *Statement {
	return n.objectFunction("ComputeHash")
}

func (n *NativeTyperef) ZeroValue() *Statement {
	return n.objectFunction("ZeroValue")
}
