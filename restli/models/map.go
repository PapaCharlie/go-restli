package models

import "github.com/dave/jennifer/jen"

const MapType = "map"

type Map struct {
	Values *Model
}

func (m *Map) GoType(destinationPackage string) *jen.Statement {
	return jen.Map(jen.String()).Add(m.Values.GoType(destinationPackage))
}

func (m *Map) InnerModels() []*Model {
	return []*Model{m.Values}
}
