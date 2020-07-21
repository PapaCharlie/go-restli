package codegen

import (
	"fmt"
	"strings"

	. "github.com/dave/jennifer/jen"
)

const unionReceiver = "u"

type StandaloneUnion struct {
	NamedType
	Union UnionType `json:"Union"`
}

func (u *StandaloneUnion) InnerTypes() IdentifierSet {
	return u.Union.InnerModels()
}

func (u *StandaloneUnion) GenerateCode() *Statement {
	def := Empty()

	AddWordWrappedComment(def, u.Doc).Line().
		Type().Id(u.Name).
		Add(u.Union.GoType()).
		Line().Line()

	AddFuncOnReceiver(def, unionReceiver, u.Name, ValidateUnionFields).
		Params().
		Params(Error()).
		BlockFunc(func(def *Group) {
			u.Union.validateUnionFields(def, unionReceiver, u.Name)
		}).Line().Line()

	AddRestLiEncode(def, unionReceiver, u.Name, func(def *Group) {
		u.Union.encode(def, unionReceiver, u.Name)
	}).Line().Line()

	return def
}

type UnionType struct {
	HasNull bool
	Members []UnionMember
}

func (u *UnionType) InnerModels() IdentifierSet {
	innerTypes := make(IdentifierSet)
	for _, m := range u.Members {
		innerTypes.AddAll(m.Type.InnerTypes())
	}
	return innerTypes
}

func (u *UnionType) GoType() *Statement {
	return StructFunc(func(def *Group) {
		for _, m := range u.Members {
			field := def.Empty()
			field.Id(m.name())
			field.Add(m.Type.PointerType())
			field.Tag(JsonFieldTag(m.Alias, true))
		}
	})
}

func (u *UnionType) validateUnionFields(def *Group, receiver string, typeName string) {
	u.validateAllMembers(def, receiver, typeName, func(*Group, UnionMember) {
		// nothing to do when simply validating
	})
}

func (u *UnionType) encode(def *Group, receiver string, typeName string) {
	u.validateAllMembers(def, receiver, typeName, func(def *Group, m UnionMember) {
		writeStringToBuf(def, Lit("("+m.Alias+":"))
		fieldAccessor := Id(receiver).Dot(m.name())
		if !(m.Type.Reference != nil || m.Type.IsMapOrArray()) {
			fieldAccessor = Op("*").Add(fieldAccessor)
		}
		m.Type.WriteToBuf(def, fieldAccessor)
		def.Id("buf").Dot("WriteByte").Call(LitRune(')'))
	})
}

func (u *UnionType) validateAllMembers(def *Group, receiver string, typeName string, f func(def *Group, m UnionMember)) {
	isSet := "isSet"
	def.Id(isSet).Op(":=").False().Line()

	var errorMessage string
	if u.HasNull {
		errorMessage = fmt.Sprintf("must specify at most one union member of %s", typeName)
	} else {
		errorMessage = fmt.Sprintf("must specify exactly one union member of %s", typeName)
	}

	for i, m := range u.Members {
		def.If(Id(receiver).Dot(m.name()).Op("!=").Nil()).BlockFunc(func(def *Group) {
			if i == 0 {
				def.Id(isSet).Op("=").True()
			} else {
				def.If(Op("!").Id(isSet)).BlockFunc(func(def *Group) {
					def.Id(isSet).Op("=").True()
				}).Else().BlockFunc(func(def *Group) {
					def.Return(Qual("fmt", "Errorf").Call(Lit(errorMessage)))
				})
			}
			f(def, m)
		}).Line()
	}

	if !u.HasNull {
		def.If(Op("!").Id(isSet)).BlockFunc(func(def *Group) {
			def.Return(Qual("fmt", "Errorf").Call(Lit(errorMessage)))
		})
	}

	def.Return(Nil())
}

type UnionMember struct {
	Type  RestliType
	Alias string
}

func (m *UnionMember) name() string {
	return ExportedIdentifier(m.Alias[strings.LastIndex(m.Alias, ".")+1:])
}
