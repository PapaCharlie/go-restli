package models

import (
	"log"

	. "github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
)

const TyperefModelTypeName = "typeref"

type TyperefModel struct {
	NameAndDoc
	Ref *Model `json:"ref"`
}

func (t *TyperefModel) InnerModels() (models []*Model) {
	return []*Model{t.Ref}
}

func (t *TyperefModel) generateCode() (def *Statement) {
	def = Empty()
	AddWordWrappedComment(def, t.Doc).Line()
	def.Type().Id(t.Name).Add(t.Ref.GoType()).Line().Line()

	if t.Ref.Primitive == nil && t.Ref.Bytes == nil {
		log.Panicln("illegal non-primitive typeref type", t)
	}

	receiver := ReceiverName(t.Name)

	var accessor *Statement
	var encoder func(*Statement) *Statement
	var decoder func(*Statement) *Statement

	if t.Ref.Bytes != nil {
		accessor = Bytes().Call(Op("*").Id(receiver))
		encoder = t.Ref.Bytes.encode
		decoder = t.Ref.Bytes.decode
	} else {
		accessor = Id(t.Ref.Primitive[1]).Call(Op("*").Id(receiver))
		encoder = t.Ref.Primitive.encode
		decoder = t.Ref.Primitive.decode
	}

	AddRestLiEncode(def, receiver, t.Name, func(def *Group) {
		def.Return(encoder(accessor), Nil())
	}).Line().Line()
	AddRestLiDecode(def, receiver, t.Name, func(def *Group) {
		def.Return(decoder(Id(receiver)))
	}).Line().Line()

	return def
}
