package types

import (
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type ComplexKey struct {
	NamedType
	Key    utils.Identifier
	Params utils.Identifier
}

func (ck *ComplexKey) InnerTypes() utils.IdentifierSet {
	return utils.NewIdentifierSet(ck.Key, ck.Params)
}

func (ck *ComplexKey) GenerateCode() *Statement {
	def := utils.AddWordWrappedComment(Empty(), ck.Doc).Line().
		Type().Id(ck.Name).
		StructFunc(func(def *Group) {
			def.Id(ComplexKeyParamsField).Op("*").Add(ck.Params.Qual()).Tag(utils.JsonFieldTag("$params", false))
			def.Add(ck.Key.Qual())
		}).Line().Line()

	record := &Record{
		NamedType: ck.NamedType,
		Fields:    utils.TypeRegistry.Resolve(ck.Key).(*Record).Fields,
	}
	receiver := Id(record.Receiver())

	AddEquals(def, record.Receiver(), ck.Name, func(other Code, def *Group) {
		def.Add(equals(RestliType{Reference: &ck.Key}, false,
			Add(receiver).Add(ck.KeyAccessor()),
			Add(other).Add(ck.KeyAccessor()))).Line()
		def.Add(equals(RestliType{Reference: &ck.Params}, true,
			Add(receiver).Dot(ComplexKeyParamsField),
			Add(other).Dot(ComplexKeyParamsField))).Line()
		def.Return(True())
	})

	AddMarshalRestLi(def, record.Receiver(), ck.Name, func(def *Group) {
		record.generateMarshaler(def, Add(receiver).Add(ck.KeyAccessor()))
	})

	AddUnmarshalRestli(def, record.Receiver(), ck.Name, func(def *Group) {
		record.generateUnmarshaler(def, Add(receiver).Add(ck.KeyAccessor()), &ck.Params)
	})

	return def
}

func (ck *ComplexKey) KeyAccessor() Code {
	return Dot(ck.Key.Name)
}
