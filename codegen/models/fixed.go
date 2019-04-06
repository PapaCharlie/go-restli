package models

import (
	"fmt"

	. "github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
)

const FixedModelTypeName = "fixed"

type FixedModel struct {
	NameAndDoc
	Size int
}

func (f *FixedModel) generateCode() (def *Statement) {
	def = Empty()
	AddWordWrappedComment(def, f.Doc).Line()
	def.Type().Id(f.Name).Index(Lit(f.Size)).Byte().Line().Line()

	receiver := ReceiverName(f.Name)
	errorMsg := fmt.Sprintf("size of %s must be exactly %d bytes (was %%d)", f.Name, f.Size)

	AddMarshalJSON(def, receiver, f.Name, func(def *Group) {
		def.Id("bytes").Op(":=").Add(Bytes()).Call(Id(receiver).Index(Op(":")))
		def.Return(Id("bytes").Dot(MarshalJSON).Call())
	}).Line().Line()
	AddUnmarshalJSON(def, receiver, f.Name, func(def *Group) {
		def.Id("bytes").Op(":=").Make(Bytes(), Lit(f.Size))
		def.Err().Op("=").Id("bytes").Dot(UnmarshalJSON).Call(Id("data"))
		IfErrReturn(def)
		def.If(Len(Id("bytes")).Op("!=").Lit(f.Size)).BlockFunc(func(def *Group) {
			def.Err().Op("=").Qual("fmt", "Errorf").Call(Lit(errorMsg), Len(Id("bytes")))
			def.Return()
		})
		def.Copy(Id(receiver).Index(Op(":")), Id("bytes").Index(Op(":").Lit(f.Size)))
		def.Return()
	}).Line().Line()

	AddRestLiEncode(def, receiver, f.Name, func(def *Group) {
		def.Return(Id(Codec).Dot("EncodeBytes").Call(Id(receiver).Index(Op(":"))), Nil())
	}).Line().Line()
	AddRestLiDecode(def, receiver, f.Name, func(def *Group) {
		def.Id("bytes").Op(":=").Make(Bytes(), Lit(f.Size))
		def.Err().Op("=").Id(Codec).Dot("DecodeBytes").Call(Id("data"), Op("&").Id("bytes"))
		IfErrReturn(def)
		def.If(Len(Id("bytes")).Op("!=").Lit(f.Size)).BlockFunc(func(def *Group) {
			def.Err().Op("=").Qual("fmt", "Errorf").Call(Lit(errorMsg), Len(Id("bytes")))
			def.Return()
		})
		def.Copy(Id(receiver).Index(Op(":")), Id("bytes").Index(Op(":").Lit(f.Size)))
		def.Return()
	}).Line().Line()

	return
}
