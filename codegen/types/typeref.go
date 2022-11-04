package types

import (
	"github.com/PapaCharlie/go-restli/v2/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

const TyperefShouldUsePointer = utils.No

type Typeref struct {
	NamedType
	Type     *PrimitiveType `json:"type"`
	IsCustom bool           `json:"isCustom"`
}

func (r *Typeref) InnerTypes() utils.IdentifierSet {
	return nil
}

func (r *Typeref) ShouldReference() utils.ShouldUsePointer {
	return TyperefShouldUsePointer
}

func (r *Typeref) GenerateCode() (def *Statement) {
	if r.IsCustomTyperef() {
		return nil
	}
	def = Empty()

	underlyingType := RestliType{Primitive: r.Type}
	cast := r.Type.Cast(Id(r.Receiver()))

	utils.AddWordWrappedComment(def, r.Doc).Line()
	def.Type().Id(r.TypeName()).Add(r.Type.GoType()).Line().Line()

	AddEquals(def, r.Receiver(), r.TypeName(), TyperefShouldUsePointer, func(other Code, def *Group) {
		if r.Type.IsBytes() {
			def.Return(Qual("bytes", "Equal").Call(Id(r.Receiver()), other))
		} else {
			def.Return(Id(r.Receiver()).Op("==").Add(other))
		}
	})

	AddComputeHash(def, r.Receiver(), r.TypeName(), TyperefShouldUsePointer, func(h Code, def *Group) {
		def.Add(h).Dot(r.Type.HasherName()).Call(cast)
	})

	utils.AddPointer(def, r.Receiver(), r.TypeName())

	AddMarshalRestLi(def, r.Receiver(), r.TypeName(), TyperefShouldUsePointer, func(def *Group) {
		def.Add(Writer.Write(underlyingType, Writer, cast))
		def.Return(Nil())
	})

	AddUnmarshalRestli(def, r.Receiver(), r.TypeName(), TyperefShouldUsePointer, func(def *Group) {
		tmp := Id("tmp")
		def.Var().Add(tmp).Add(r.Type.GoType())
		def.Add(Reader.Read(underlyingType, Reader, tmp))
		def.Add(utils.IfErrReturn(Err())).Line()

		def.Op("*").Id(r.Receiver()).Op("=").Id(r.TypeName()).Call(tmp)
		def.Return(Nil())
	})

	return def
}
func customTyperefFunc(id utils.Identifier, prefix string) *Statement {
	return Qual(id.PackagePath(), prefix+id.TypeName())
}

func customTyperefEqualsFunc(id utils.Identifier) *Statement {
	return customTyperefFunc(id, utils.Equals)
}

func customTyperefHashFunc(id utils.Identifier) *Statement {
	return customTyperefFunc(id, utils.ComputeHash)
}

func customTyperefMarshalerFunc(id utils.Identifier) *Statement {
	return customTyperefFunc(id, utils.Marshal)
}

func customTyperefUnmarshalerFunc(id utils.Identifier) *Statement {
	return customTyperefFunc(id, utils.Unmarshal)
}

func writeCustomTyperef(writerAccessor, sourceAccessor Code, id utils.Identifier) *Statement {
	return Add(utils.WriteCustomTyperef).Call(writerAccessor, sourceAccessor, customTyperefMarshalerFunc(id))
}

func readCustomTyperef(readerAccessor Code, id utils.Identifier) *Statement {
	return Add(utils.ReadCustomTyperef).Call(readerAccessor, customTyperefUnmarshalerFunc(id))
}

func CallRegisterCustomTyperef(id utils.Identifier) *Statement {
	return Qual(utils.RestLiCodecPackage, "RegisterCustomTyperef").Custom(utils.MultiLineCall,
		customTyperefMarshalerFunc(id),
		customTyperefUnmarshalerFunc(id),
		customTyperefHashFunc(id),
		customTyperefEqualsFunc(id),
	)
}

func AddCustomTyperefEquals(name string, f func(def *Group, left, right Code)) *Statement {
	left, right := Code(Id("left")), Code(Id("right"))

	return Func().Id(utils.Equals + name).Params(List(left, right).Id(name)).Bool().BlockFunc(func(def *Group) {
		f(def, left, right)
	})
}

func AddCustomTyperefComputeHash(name string, f func(def *Group, in, h Code)) *Statement {
	in := Code(Id("in"))
	return Func().Id(utils.ComputeHash + name).Params(Add(in).Id(name)).Add(utils.Hash).BlockFunc(func(def *Group) {
		h := Code(Id("hash"))
		def.Add(h).Op(":=").Add(utils.NewHash)

		f(def, in, h)

		def.Return(h)
	})
}

func AddCustomTyperefMarshal(name string, pt PrimitiveType, f func(def *Group, in, out Code)) *Statement {
	in, out := Code(Id("in")), Code(Id("out"))
	return Func().Id(utils.Marshal+name).Params(Add(in).Id(name)).Params(Add(out).Add(pt.GoType()), Err().Error()).BlockFunc(func(def *Group) {
		f(def, in, out)
	})
}

func AddCustomTyperefUnmarshal(name string, pt PrimitiveType, f func(def *Group, in, out Code)) *Statement {
	in, out := Code(Id("in")), Code(Id("out"))
	return Func().Id(utils.Unmarshal+name).Params(Add(in).Add(pt.GoType())).Params(Add(out).Id(name), Err().Error()).BlockFunc(func(def *Group) {
		f(def, in, out)
	})
}
