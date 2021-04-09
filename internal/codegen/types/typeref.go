package types

import (
	"encoding/json"

	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type Typeref struct {
	NamedType
	Type       *PrimitiveType             `json:"type"`
	Properties map[string]json.RawMessage `json:"properties"`
}

func (r *Typeref) InnerTypes() utils.IdentifierSet {
	return nil
}

func (r *Typeref) GenerateCode() (def *Statement) {
	def = Empty()

	underlyingType := RestliType{Primitive: r.Type}
	cast := r.Type.Cast(Op("*").Id(r.Receiver()))

	utils.AddWordWrappedComment(def, r.Doc).Line()
	def.Type().Id(r.Name).Add(r.Type.GoType()).Line().Line()

	AddEquals(def, r.Receiver(), r.Name, func(other Code, def *Group) {
		left, right := Op("*").Id(r.Receiver()), Op("*").Add(other)

		if r.Type.IsBytes() {
			def.Return(Qual("bytes", "Equal").Call(left, right))
		} else {
			def.Return(Add(left).Op("==").Add(right))
		}
	})
	AddComputeHash(def, r.Receiver(), r.Name, func(h Code, def *Group) {
		def.Add(h).Dot(r.Type.HasherName()).Call(cast)
		def.Return(h)
	})
	AddMarshalRestLi(def, r.Receiver(), r.Name, func(def *Group) {
		def.Add(Writer.Write(underlyingType, Writer, cast))
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
