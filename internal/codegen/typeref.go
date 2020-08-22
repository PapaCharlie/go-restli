package codegen

import (
	. "github.com/dave/jennifer/jen"
)

type Typeref struct {
	NamedType
	Type *PrimitiveType `json:"type"`
}

func (r *Typeref) InnerTypes() IdentifierSet {
	return nil
}

func (r *Typeref) GenerateCode() (def *Statement) {
	def = Empty()

	AddWordWrappedComment(def, r.Doc).Line()
	def.Type().Id(r.Name).Add(r.Type.GoType()).Line().Line()

	AddRestLiEncode(def, r.Receiver(), r.Name, func(def *Group) {
		Encoder.Write(def, RestliType{Primitive: r.Type}, r.Type.Cast(Op("*").Id(r.Receiver())))
		def.Return(Nil())
	}).Line().Line()
	AddRestLiDecode(def, r.Receiver(), r.Name, func(def *Group) {
		def.Return(r.Type.decode(Id(Codec), Id(r.Receiver())))
	}).Line().Line()

	return def
}
