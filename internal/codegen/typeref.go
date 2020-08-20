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

	underlyingType := RestliType{Primitive: r.Type}

	AddWordWrappedComment(def, r.Doc).Line()
	def.Type().Id(r.Name).Add(r.Type.GoType()).Line().Line()

	AddMarshalRestLi(def, r.Receiver(), r.Name, func(def *Group) {
		def.Add(Writer.Write(underlyingType, Writer, r.Type.Cast(Op("*").Id(r.Receiver()))))
		def.Return(Nil())
	}).Line().Line()
	AddRestLiDecode(def, r.Receiver(), r.Name, func(def *Group) {
		tmp := Id("tmp")
		def.Var().Add(tmp).Add(r.Type.GoType())
		def.Add(Reader.Read(underlyingType, tmp))
		def.Add(IfErrReturn(Err())).Line()

		def.Op("*").Id(r.Receiver()).Op("=").Id(r.Name).Call(tmp)
		def.Return(Nil())
	}).Line().Line()

	return def
}
