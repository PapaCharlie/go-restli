package models

import (
	"github.com/dave/jennifer/jen"
)

const ArrayType = "array"

type Array struct {
	Items *Model
}

func (a *Array) GoType(destinationPackage string) *jen.Statement {
	return jen.Index().Add(a.Items.GoType(destinationPackage))
}

func (a *Array) InnerModels() []*Model {
	return []*Model{a.Items}
}
