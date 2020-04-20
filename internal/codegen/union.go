package codegen

import (
	"fmt"
	"strings"

	. "github.com/dave/jennifer/jen"
)

type UnionType []UnionMember

func (u *UnionType) InnerModels() IdentifierSet {
	innerTypes := make(IdentifierSet)
	for _, m := range *u {
		innerTypes.AddAll(m.Type.InnerTypes())
	}
	return innerTypes
}

func (u *UnionType) GoType() *Statement {
	return StructFunc(func(def *Group) {
		for _, m := range *u {
			field := def.Empty()
			field.Id(m.name())
			field.Add(m.Type.PointerType())
			field.Tag(JsonFieldTag(m.Alias, true))
		}
	})
}

func (u *UnionType) FieldGoType() *Statement {
	return StructFunc(func(def *Group) {
		for _, m := range *u {
			field := def.Empty()
			field.Id(m.name())
			field.Add(m.Type.FieldGoType())
			field.Tag(JsonFieldTag(m.Alias, true))
		}
	})
}

func (u *UnionType) validateUnionFields(def *Group, accessor *Statement) {
	isSet := "is" + canonicalizeAccessor(accessor) + "Set"
	def.Id(isSet).Op(":=").False().Line()
	errorMessage := fmt.Sprintf("must specify exactly one member of %s", accessor.GoString())

	for i, t := range *u {
		def.If(Add(accessor).Dot(t.name()).Op("!=").Nil()).
			BlockFunc(func(def *Group) {
				if i == 0 {
					def.Id(isSet).Op("=").True()
				} else {
					def.If(Op("!").Id(isSet)).BlockFunc(func(def *Group) {
						def.Id(isSet).Op("=").True()
					}).Else().BlockFunc(func(def *Group) {
						def.Return(Qual("fmt", "Errorf").Call(Lit(errorMessage)))
					})
				}
			}).Line()
	}
	def.If(Op("!").Id(isSet)).BlockFunc(func(def *Group) {
		def.Return(Qual("fmt", "Errorf").Call(Lit(errorMessage)))
	})
}

type UnionMember struct {
	Type  RestliType
	Alias string
}

func (m *UnionMember) name() string {
	return ExportedIdentifier(m.Alias[strings.LastIndex(m.Alias, ".")+1:])
}
