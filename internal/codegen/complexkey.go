package codegen

import . "github.com/dave/jennifer/jen"

type ComplexKey struct {
	NamedType
	Key    RestliType
	Params RestliType
}

func (ck *ComplexKey) InnerTypes() IdentifierSet {
	innerTypes := make(IdentifierSet)
	innerTypes.AddAll(ck.Key.InnerTypes())
	innerTypes.AddAll(ck.Params.InnerTypes())
	return innerTypes
}

func (ck *ComplexKey) GenerateCode() *Statement {
	def := AddWordWrappedComment(Empty(), ck.Doc).Line().
		Type().Id(ck.Name).
		StructFunc(func(def *Group) {
			def.Add(ck.Key.GoType())
			def.Id("Params").Add(ck.Params.PointerType()).Tag(JsonFieldTag("$params", false))
		}).Line().Line()

	receiver := "ck"
	return AddRestLiEncode(def, receiver, ck.Name, func(def *Group) {
		def.List(Id("encoded"), Err()).Op(":=").Id(receiver).Dot(ck.Key.Reference.Name).Dot(RestLiEncode).Call(Id(Codec))
		IfErrReturn(def, Lit(""), Err()).Line()

		def.If(Id(receiver).Dot("Params").Op("!=").Nil()).BlockFunc(func(def *Group) {
			const encodedParamsVar = "encodedParams"
			def.Var().Id(encodedParamsVar).String()
			def.List(Id(encodedParamsVar), Err()).Op("=").Id(receiver).Dot("Params").Dot(RestLiEncode).Call(Id(Codec))
			IfErrReturn(def, Lit(""), Err())
			def.Id("encoded").Op("=").Lit("($params:").Op("+").Id(encodedParamsVar).Op("+").Lit(",").Op("+").Id("encoded").Index(Lit(1), Empty())
		}).Line()

		def.Return(Id("encoded"), Nil())
	})
}
