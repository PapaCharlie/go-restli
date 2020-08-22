package codegen

import (
	. "github.com/dave/jennifer/jen"
)

type ComplexKey struct {
	NamedType
	Key    Identifier
	Params Identifier
}

func (ck *ComplexKey) InnerTypes() IdentifierSet {
	return NewIdentifierSet(ck.Key, ck.Params)
}

func (ck *ComplexKey) GenerateCode() *Statement {
	def := AddWordWrappedComment(Empty(), ck.Doc).Line().
		Type().Id(ck.Name).
		StructFunc(func(def *Group) {
			def.Add(ck.Key.Qual())
			def.Id(ComplexKeyParams).Op("*").Add(ck.Params.Qual()).Tag(JsonFieldTag("$params", false))
		}).Line().Line()

	record := &Record{
		NamedType: ck.NamedType,
		Fields:    TypeRegistry.Resolve(ck.Key).(*Record).Fields,
	}

	return AddRestLiEncode(def, record.Receiver(), ck.Name, func(def *Group) {
		record.generateEncoder(def, false, nil, Id(record.Receiver()).Dot(ck.Key.Name))
		def.Return(Nil())
	})
}
