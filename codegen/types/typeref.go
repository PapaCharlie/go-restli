package types

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

const TyperefShouldUsePointer = utils.No

type Typeref struct {
	NamedType
	Type *PrimitiveType `json:"type"`
}

func (r *Typeref) InnerTypes() utils.IdentifierSet {
	return nil
}

func (r *Typeref) ShouldReference() utils.ShouldUsePointer {
	return TyperefShouldUsePointer
}

func (r *Typeref) GenerateCode() (def *Statement) {
	def = Empty()

	underlyingType := RestliType{Primitive: r.Type}
	cast := r.Type.Cast(Id(r.Receiver()))

	utils.AddWordWrappedComment(def, r.Doc).Line()
	def.Type().Id(r.Name).Add(r.Type.GoType()).Line().Line()

	AddEquals(def, r.Receiver(), r.Name, TyperefShouldUsePointer, func(other Code, def *Group) {
		if r.Type.IsBytes() {
			def.Return(Qual("bytes", "Equal").Call(Id(r.Receiver()), other))
		} else {
			def.Return(Id(r.Receiver()).Op("==").Add(other))
		}
	})

	AddComputeHash(def, r.Receiver(), r.Name, TyperefShouldUsePointer, func(h Code, def *Group) {
		def.Add(h).Dot(r.Type.HasherName()).Call(cast)
	})

	utils.AddPointer(def, r.Receiver(), r.Name)

	AddMarshalRestLi(def, r.Receiver(), r.Name, TyperefShouldUsePointer, func(def *Group) {
		def.Add(Writer.Write(underlyingType, Writer, cast))
		def.Return(Nil())
	})

	AddUnmarshalRestli(def, r.Receiver(), r.Identifier, TyperefShouldUsePointer, func(def *Group) {
		tmp := Id("tmp")
		def.Var().Add(tmp).Add(r.Type.GoType())
		def.Add(Reader.Read(underlyingType, Reader, tmp))
		def.Add(utils.IfErrReturn(Err())).Line()

		def.Op("*").Id(r.Receiver()).Op("=").Id(r.Name).Call(tmp)
		def.Return(Nil())
	})

	return def
}
