package types

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/v2/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

const FixedShouldUsePointer = utils.Yes

type Fixed struct {
	NamedType
	Size int
}

var FixedUnderlyingType = RestliType{Primitive: &BytesPrimitive}

func (f *Fixed) ReferencedTypes() utils.IdentifierSet {
	return nil
}

func (f *Fixed) ShouldReference() utils.ShouldUsePointer {
	return FixedShouldUsePointer
}

func (f *Fixed) GenerateCode() (def *Statement) {
	def = Empty()
	utils.AddWordWrappedComment(def, f.Doc).Line()
	def.Type().Id(f.TypeName()).Index(Lit(f.Size)).Byte().Line().Line()

	receiver := f.Receiver()
	errorMsg := fmt.Sprintf("size of %s must be exactly %d bytes (was %%d)", f.TypeName(), f.Size)
	slice := Index(Op(":"))

	AddEquals(def, receiver, f.TypeName(), FixedShouldUsePointer, func(other Code, def *Group) {
		def.Return(Qual("bytes", "Equal").Call(
			Id(receiver).Add(slice),
			Add(other).Add(slice)))
	})

	AddComputeHash(def, receiver, f.TypeName(), FixedShouldUsePointer, func(h Code, def *Group) {
		def.Add(hash(h, FixedUnderlyingType, false, Id(receiver).Add(slice)))
	})

	utils.AddPointer(def, f.Receiver(), f.TypeName())

	AddMarshalRestLi(def, receiver, f.TypeName(), FixedShouldUsePointer, func(def *Group) {
		def.Add(Writer.Write(FixedUnderlyingType, Writer, Id(receiver).Add(slice)))
		def.Return(Nil())
	})

	AddUnmarshalRestli(def, receiver, f.TypeName(), FixedShouldUsePointer, func(def *Group) {
		data := Id("data")
		def.Var().Add(data).Index().Byte()
		def.Add(Reader.Read(FixedUnderlyingType, Reader, data))
		def.Add(utils.IfErrReturn(Err())).Line()

		def.If(Len(data).Op("!=").Lit(f.Size)).BlockFunc(func(def *Group) {
			def.Return(Qual("fmt", "Errorf").Call(Lit(errorMsg), Len(data)))
		}).Line()

		def.Copy(Id(receiver).Add(slice), Add(data).Index(Op(":").Lit(f.Size)))
		def.Return(Nil())
	})

	return def
}

func (f *Fixed) getLit(rawJson string) *Statement {
	return f.Qual().Add(getLitBytesValues(rawJson))
}
