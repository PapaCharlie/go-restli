package types

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type Fixed struct {
	NamedType
	Size int
}

var FixedUnderlyingType = RestliType{Primitive: &BytesPrimitive}

func (f *Fixed) InnerTypes() utils.IdentifierSet {
	return nil
}

func (f *Fixed) GenerateCode() (def *Statement) {
	def = Empty()
	utils.AddWordWrappedComment(def, f.Doc).Line()
	def.Type().Id(f.Name).Index(Lit(f.Size)).Byte().Line().Line()

	receiver := f.Receiver()
	errorMsg := fmt.Sprintf("size of %s must be exactly %d bytes (was %%d)", f.Name, f.Size)

	AddEquals(def, receiver, f.Name, func(other Code, def *Group) {
		def.Return(Qual("bytes", "Equal").Call(
			Id(receiver).Index(Op(":")),
			Add(other).Index(Op(":"))))
	})

	AddMarshalRestLi(def, receiver, f.Name, func(def *Group) {
		def.Add(Writer.Write(FixedUnderlyingType, Writer, Id(receiver).Index(Op(":"))))
		def.Return(Nil())
	})
	AddUnmarshalRestli(def, receiver, f.Name, func(def *Group) {
		data := Id("data")
		def.Var().Add(data).Index().Byte()
		def.Add(Reader.Read(FixedUnderlyingType, data))
		def.Add(utils.IfErrReturn(Err())).Line()

		def.If(Len(data).Op("!=").Lit(f.Size)).BlockFunc(func(def *Group) {
			def.Return(Qual("fmt", "Errorf").Call(Lit(errorMsg), Len(data)))
		}).Line()

		def.Copy(Id(receiver).Index(Op(":")), Add(data).Index(Op(":").Lit(f.Size)))
		def.Return(Nil())
	})

	return def
}

func (f *Fixed) getLit(rawJson string) *Statement {
	return f.Qual().Add(getLitBytesValues(rawJson))
}
