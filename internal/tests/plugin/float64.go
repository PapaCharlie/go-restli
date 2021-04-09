package main

import (
	. "github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type float64Provider struct{ *Typeref }

func (r *float64Provider) castReceiver() *Statement {
	return Float64().Call(Op("*").Id(r.Receiver()))
}

func (r *float64Provider) ReferencedTypes() utils.IdentifierSet {
	return nil
}

func (r *float64Provider) GenerateType() Code {
	return Type().Id(r.Name).Float64()
}

func (r *float64Provider) GenerateMarshalRaw(def *Group) {
	def.Return(r.castReceiver(), Nil())
}

func (r *float64Provider) GenerateUnmarshalRaw(raw Code, def *Group) {
	def.Op("*").Id(r.Receiver()).Op("=").Add(r.Qual()).Call(raw)
	def.Return(Nil())
}

func (r *float64Provider) GenerateEquals(other Code, def *Group) {
	def.Return(Op("*").Id(r.Receiver()).Op("==").Op("*").Add(other))
}

func (r *float64Provider) GenerateComputeHash(h Code, def *Group) {
	def.Add(h).Dot("AddFloat64").Call(r.castReceiver())
}

func (r *float64Provider) ZeroValue() Code {
	return Lit(float64(0))
}

func (r *float64Provider) ShouldReference() bool {
	return false
}
