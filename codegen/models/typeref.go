package models

import (
	"encoding/json"
	"log"

	. "github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
)

const TyperefModelTypeName = "typeref"

type TyperefModel struct {
	Identifier
	Doc string
	Ref *Model
}

func (r *TyperefModel) UnmarshalJSON(data []byte) error {
	t := &struct {
		typeField
		docField
		Identifier
		Ref *Model `json:"ref"`
	}{}
	if err := json.Unmarshal(data, t); err != nil {
		return err
	}
	if t.Type != TyperefModelTypeName {
		return &WrongTypeError{Expected: TyperefModelTypeName, Actual: t.Type}
	}
	r.Identifier = t.Identifier
	r.Doc = t.Doc
	r.Ref = t.Ref
	return nil
}

func (r *TyperefModel) innerModels() []*Model {
	return []*Model{r.Ref}
}

func (r *TyperefModel) GenerateCode() (def *Statement) {
	def = Empty()
	AddWordWrappedComment(def, r.Doc).Line()
	def.Type().Id(r.Name).Add(r.Ref.GoType()).Line().Line()

	var accessor *Statement
	var encoder func(*Statement) *Statement
	var decoder func(*Statement) *Statement

	if bytes, ok := r.Ref.BuiltinType.(*BytesModel); ok {
		accessor = Bytes().Call(Op("*").Id(r.receiver()))
		encoder = bytes.encode
		decoder = bytes.decode
	}
	if primitive, ok := r.Ref.BuiltinType.(*PrimitiveModel); ok {
		accessor = Id(primitive[1]).Call(Op("*").Id(r.receiver()))
		encoder = primitive.encode
		decoder = primitive.decode
	}
	if accessor == nil {
		log.Panicln("Illegal typeref type:", r.Ref)
	}

	AddRestLiEncode(def, r.receiver(), r.Name, func(def *Group) {
		def.Return(encoder(accessor), Nil())
	}).Line().Line()
	AddRestLiDecode(def, r.receiver(), r.Name, func(def *Group) {
		def.Return(decoder(Id(r.receiver())))
	}).Line().Line()

	return def
}
