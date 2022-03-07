package types

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

const FixedShouldUsePointer = utils.Yes

type Fixed struct {
	NamedType
	Size int
}

var FixedUnderlyingType = RestliType{Primitive: &BytesPrimitive}

func (f *Fixed) InnerTypes() utils.IdentifierSet {
	return nil
}

func (f *Fixed) ShouldReference() utils.ShouldUsePointer {
	return FixedShouldUsePointer
}

func (f *Fixed) GenerateCode() (def *Statement) {
	def = Empty()
	o := NewObjectCodeGenerator(f.Identifier, FixedShouldUsePointer)

	o.DeclareType(def, f.Doc, Index(Lit(f.Size)).Byte())

	errorMsg := fmt.Sprintf("size of %s must be exactly %d bytes (was %%d)", f.Name, f.Size)
	slice := Index(Op(":"))

	o.Equals(def, func(receiver, other Code, def *Group) {
		def.Return(Qual("bytes", "Equal").Call(
			Add(receiver).Add(slice),
			Add(other).Add(slice)))
	})

	o.ComputeHash(def, func(receiver, h Code, def *Group) {
		def.Add(hash(h, FixedUnderlyingType, false, Add(receiver).Add(slice)))
	})

	utils.AddPointer(def, f.Receiver(), f.Name)

	o.MarshalRestLi(def, func(receiver, writer Code, def *Group) {
		def.Add(WriterUtils.Write(FixedUnderlyingType, writer, Add(receiver).Add(slice)))
		def.Return(Nil())
	})

	o.UnmarshalRestLi(def, func(receiver, reader Code, def *Group) {
		data := Id("data")
		def.Var().Add(data).Index().Byte()
		def.Add(ReaderUtils.Read(FixedUnderlyingType, reader, data))
		def.Add(utils.IfErrReturn(Err())).Line()

		def.If(Len(data).Op("!=").Lit(f.Size)).BlockFunc(func(def *Group) {
			def.Return(Qual("fmt", "Errorf").Call(Lit(errorMsg), Len(data)))
		}).Line()

		def.Copy(Add(receiver).Add(slice), Add(data).Index(Op(":").Lit(f.Size)))
		def.Return(Nil())
	})

	return def
}

func (f *Fixed) getLit(rawJson string) *Statement {
	return f.Qual().Add(getLitBytesValues(rawJson))
}
