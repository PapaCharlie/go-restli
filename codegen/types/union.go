package types

import (
	"fmt"
	"strings"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

const unionReceiver = "u"
const UnionShouldUsePointer = utils.Yes

type StandaloneUnion struct {
	NamedType
	Union UnionType `json:"Union"`
}

func (u *StandaloneUnion) InnerTypes() utils.IdentifierSet {
	return u.Union.InnerModels()
}

func (u *StandaloneUnion) ShouldReference() utils.ShouldUsePointer {
	return UnionShouldUsePointer
}

func (u *StandaloneUnion) GenerateCode() *Statement {
	def := Empty()

	utils.AddWordWrappedComment(def, u.Doc).Line().
		Type().Id(u.Name).
		Add(u.Union.GoType()).
		Line().Line()

	AddEquals(def, unionReceiver, u.Name, UnionShouldUsePointer, func(other Code, def *Group) {
		for _, m := range u.Union.Members {
			def.If(Op("!").Add(equalsCondition(m.Type, true, Id(unionReceiver).Dot(m.name()), Add(other).Dot(m.name())))).Block(Return(False()))
		}
		def.Return(True())
	})

	AddComputeHash(def, unionReceiver, u.Name, UnionShouldUsePointer, func(h Code, def *Group) {
		for _, m := range u.Union.Members {
			def.Add(hash(h, m.Type, true, Id(unionReceiver).Dot(m.name()))).Line()
		}
	})

	utils.AddFuncOnReceiver(def, unionReceiver, u.Name, utils.ValidateUnionFields, UnionShouldUsePointer).
		Params().
		Params(Error()).
		BlockFunc(func(def *Group) {
			u.Union.validateUnionFields(def, unionReceiver, u.Name)
		}).Line().Line()

	AddMarshalRestLi(def, unionReceiver, u.Name, UnionShouldUsePointer, func(def *Group) {
		def.Return(Writer.WriteMap(Writer, func(keyWriter Code, def *Group) {
			u.Union.validateAllMembers(def, unionReceiver, u.Name, func(def *Group, m UnionMember) {
				fieldAccessor := Id(unionReceiver).Dot(m.name())
				if m.Type.Reference == nil {
					fieldAccessor = Op("*").Add(fieldAccessor)
				}
				def.Add(Writer.Write(m.Type, Add(keyWriter).Call(Lit(m.Alias)), fieldAccessor, Err()))
			})
			def.Return(Nil())
		}))
	})

	AddUnmarshalRestli(def, unionReceiver, u.Identifier, UnionShouldUsePointer, func(def *Group) {
		u.Union.decode(def, unionReceiver, u.Name)
	})

	return def
}

type UnionType struct {
	HasNull bool
	Members []UnionMember
}

func (u *UnionType) InnerModels() utils.IdentifierSet {
	innerTypes := make(utils.IdentifierSet)
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
			field.Tag(utils.JsonFieldTag(m.Alias, true))
		}
	})
}

func (u *UnionType) validateUnionFields(def *Group, receiver string, typeName string) {
	u.validateAllMembers(def, receiver, typeName, func(*Group, UnionMember) {
		// nothing to do when simply validating
	})
	def.Return(Nil())
}

func (u *UnionType) decode(def *Group, receiver string, typeName string) {
	wasSet := Id("wasSet")
	def.Add(wasSet).Op(":=").False()

	errorMessage := u.errorMessage(typeName)

	decode := Reader.ReadMap(Reader, func(reader, key Code, def *Group) {
		def.If(wasSet).Block(
			Return(errorMessage),
		).Else().Block(
			Add(wasSet).Op("=").True(),
		)
		def.Switch(key).BlockFunc(func(def *Group) {
			for _, m := range u.Members {
				fieldAccessor := Id(receiver).Dot(m.name())
				def.Case(Lit(m.Alias)).BlockFunc(func(def *Group) {
					def.Add(fieldAccessor).Op("=").New(m.Type.GoType())

					if m.Type.Reference == nil {
						fieldAccessor = Op("*").Add(fieldAccessor)
					}

					def.Add(Reader.Read(m.Type, reader, fieldAccessor))
				})
			}
		})
		def.Return(Err())
	})

	if u.HasNull {
		def.Return(decode)
	} else {
		def.Err().Op("=").Add(decode)
		def.Add(utils.IfErrReturn(Err()))
		def.If(Op("!").Add(wasSet)).Block(
			Return(errorMessage),
		)
		def.Return(Nil())
	}
}

func (u *UnionType) errorMessage(typeName string) *Statement {
	if u.HasNull {
		return Qual("errors", "New").Call(Lit(fmt.Sprintf("must specify at most one union member of %s", typeName)))
	} else {
		return Qual("errors", "New").Call(Lit(fmt.Sprintf("must specify exactly one union member of %s", typeName)))
	}
}

func (u *UnionType) validateAllMembers(def *Group, receiver string, typeName string, f func(def *Group, m UnionMember)) {
	isSet := "isSet"
	def.Id(isSet).Op(":=").False().Line()

	errorMessage := u.errorMessage(typeName)

	for i, m := range u.Members {
		def.If(Id(receiver).Dot(m.name()).Op("!=").Nil()).BlockFunc(func(def *Group) {
			if i == 0 {
				def.Id(isSet).Op("=").True()
			} else {
				def.If(Op("!").Id(isSet)).BlockFunc(func(def *Group) {
					def.Id(isSet).Op("=").True()
				}).Else().BlockFunc(func(def *Group) {
					def.Return(errorMessage)
				})
			}
			f(def, m)
		}).Line()
	}

	if !u.HasNull {
		def.If(Op("!").Id(isSet)).BlockFunc(func(def *Group) {
			def.Return(errorMessage)
		})
	}
}

type UnionMember struct {
	Type  RestliType
	Alias string
}

func (m *UnionMember) name() string {
	return utils.ExportedIdentifier(m.Alias[strings.LastIndex(m.Alias, ".")+1:])
}
