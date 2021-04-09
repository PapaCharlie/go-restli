package main

import (
	. "github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func GetTyperefCodeProvider(t *Typeref) TyperefCodeProvider {
	if t.Namespace == "extras" && t.Name == "BigDecimal" {
		return &bigDecimalProvider{t}
	}
	if t.Namespace == "extras" && t.Name == "BigDecimal2" {
		return &float64Provider{t}
	}
	return nil
}

var bigFloat Code = Qual("math/big", "Float")

func castBigFloat(id Code) *Statement {
	return Parens(Op("*").Add(bigFloat)).Call(id)
}

type bigDecimalProvider struct{ *Typeref }

func (r *bigDecimalProvider) castReceiver() *Statement {
	return castBigFloat(Id(r.Receiver()))
}

func (r *bigDecimalProvider) ReferencedTypes() utils.IdentifierSet {
	return nil
}

func (r *bigDecimalProvider) GenerateType() Code {
	def := Type().Id(r.Name).Add(bigFloat).Line().Line()
	utils.AddStringer(def, r.Receiver(), r.Name, func(def *Group) {
		def.Return(r.castReceiver().Dot("String").Call())
	})
	return def
}

func (r *bigDecimalProvider) GenerateMarshalRaw(def *Group) {
	def.List(utils.Raw, Err()).Op(":=").Add(r.castReceiver()).Dot("MarshalText").Call()
	def.Add(utils.IfErrReturn(r.Type.ZeroValueLit(), Err()))
	def.Return(String().Call(utils.Raw), Nil())
}

func (r *bigDecimalProvider) GenerateUnmarshalRaw(raw Code, def *Group) {
	def.Return().Add(r.castReceiver()).Dot("UnmarshalText").Call(Index().Byte().Call(raw))
}

func (r *bigDecimalProvider) GenerateEquals(other Code, def *Group) {
	left, right := r.castReceiver(), castBigFloat(other)
	def.Return(Add(left).Dot("Cmp").Call(right).Op("==").Lit(0))
}

func (r *bigDecimalProvider) GenerateComputeHash(h Code, def *Group) {
	def.Add(h).Dot("AddString").Call(Id(r.Receiver()).Dot("String").Call())
}

func (r *bigDecimalProvider) ZeroValue() Code {
	return r.Qual().Values()
}

func (r *bigDecimalProvider) ShouldReference() bool {
	return true
}
