package types

import (
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type Typeref struct {
	NamedType
	Type *PrimitiveType `json:"type"`
}

func (r *Typeref) InnerTypes() utils.IdentifierSet {
	return nil
}

func (r *Typeref) GenerateCode() (def *Statement) {
	def = Empty()

	underlyingType := RestliType{Primitive: r.Type}

	utils.AddWordWrappedComment(def, r.Doc).Line()
	def.Type().Id(r.Name).Add(r.Type.GoType()).Line().Line()

	AddEquals(def, r.Receiver(), r.Name, func(other Code, def *Group) {
		def.Return(Op("*").Id(r.Receiver()).Op("== *").Add(other))
	})
	AddMarshalRestLi(def, r.Receiver(), r.Name, func(def *Group) {
		def.Add(Writer.Write(underlyingType, Writer, r.Type.Cast(Op("*").Id(r.Receiver()))))
		def.Return(Nil())
	})
	AddUnmarshalRestli(def, r.Receiver(), r.Name, func(def *Group) {
		tmp := Id("tmp")
		def.Var().Add(tmp).Add(r.Type.GoType())
		def.Add(Reader.Read(underlyingType, Reader, tmp))
		def.Add(utils.IfErrReturn(Err())).Line()

		def.Op("*").Id(r.Receiver()).Op("=").Id(r.Name).Call(tmp)
		def.Return(Nil())
	})

	return def
}
