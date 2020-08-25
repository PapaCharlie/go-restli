package codegen

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
)

type Fixed struct {
	NamedType
	Size int
}

var FixedUnderlyingType = RestliType{Primitive: &BytesPrimitive}

func (f *Fixed) InnerTypes() IdentifierSet {
	return nil
}

func (f *Fixed) GenerateCode() (def *Statement) {
	def = Empty()
	AddWordWrappedComment(def, f.Doc).Line()
	def.Type().Id(f.Name).Index(Lit(f.Size)).Byte().Line().Line()

	receiver := ReceiverName(f.Name)
	errorMsg := fmt.Sprintf("size of %s must be exactly %d bytes (was %%d)", f.Name, f.Size)

	AddMarshalRestLi(def, receiver, f.Name, func(def *Group) {
		def.Add(Writer.Write(FixedUnderlyingType, Writer, Id(receiver).Index(Op(":"))))
		def.Return(Nil())
	}).Line().Line()
	AddRestLiDecode(def, receiver, f.Name, func(def *Group) {
		bytes := Id("bytes")
		def.Var().Add(bytes).Index().Byte()
		def.Add(Reader.Read(FixedUnderlyingType, bytes))
		def.Add(IfErrReturn(Err())).Line()

		def.If(Len(Id("bytes")).Op("!=").Lit(f.Size)).BlockFunc(func(def *Group) {
			def.Return(Qual("fmt", "Errorf").Call(Lit(errorMsg), Len(Id("bytes"))))
		}).Line()

		def.Copy(Id(receiver).Index(Op(":")), Id("bytes").Index(Op(":").Lit(f.Size)))
		def.Return(Nil())
	}).Line().Line()

	return def
}

func (f *Fixed) getLit(rawJson string) *Statement {
	return f.Qual().Add(getLitBytesValues(rawJson))
}
