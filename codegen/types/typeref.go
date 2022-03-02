package types

import (
	"log"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

const TyperefShouldUsePointer = utils.No

type NativeTypeRef struct {
	Type                    PrimitiveType `json:"type"`
	Ref                     string        `json:"ref"`
	NativePackage           string        `json:"nativePackage"`
	NativeIdentifier        string        `json:"nativeIdentifier"`
	NonReceiverFuncsPackage *string       `json:"nonReceiverFuncsPackage"`

	ShouldReference utils.ShouldUsePointer `json:"-"`
}

func (n *NativeTypeRef) IsCustomStruct() bool {
	return n.NonReceiverFuncsPackage == nil
}

func (n *NativeTypeRef) GoType() *Statement {
	return Qual(n.NativePackage, n.NativeIdentifier)
}

type Typeref struct {
	NamedType
	Type *PrimitiveType `json:"type"`

	native *NativeTypeRef
}

func (r *Typeref) CheckNativeTyperef() {
	var ok bool
	r.native, ok = nativeTyperefs[r.Identifier.String()]
	if !ok {
		return
	}
	if r.native.Type.Type != r.Type.Type {
		log.Panicf("Typeref is defined as a %q, but native typeref \"%s.%s\" is defined as %q",
			r.Type.Type, r.native.NativePackage, r.native.NativeIdentifier, r.native.Type.Type)
	}
}

func (r *Typeref) InnerTypes() utils.IdentifierSet {
	return nil
}

func (r *Typeref) ShouldReference() utils.ShouldUsePointer {
	if r.native != nil {
		return r.native.ShouldReference
	} else {
		return TyperefShouldUsePointer
	}
}

func (r *Typeref) GenerateCode() (def *Statement) {
	def = Empty()

	utils.AddWordWrappedComment(def, r.Doc).Line()

	receiver := Code(Id(r.Receiver()))

	var cast func(Code) *Statement
	typedef := def.Type().Id(r.Name)
	if r.native != nil {
		cast = func(receiver Code) *Statement {
			def := Qual(r.native.NativePackage, r.native.NativeIdentifier)
			if r.native.ShouldReference.ShouldUsePointer() {
				def = Parens(Op("*").Add(def))
			}
			return def.Call(receiver)
		}
		typedef.Add(r.native.GoType())
	} else {
		cast = r.Type.Cast
		typedef.Add(r.Type.GoType())
	}
	typedef.Line().Line()

	AddEquals(def, r.Receiver(), r.Name, r.ShouldReference(), func(other Code, def *Group) {
		switch {
		case r.native != nil && r.native.IsCustomStruct():
			def.Return(cast(receiver).Dot(utils.Equals).Call(cast(other)))
		case r.native != nil:
			def.Return(Qual(*r.native.NonReceiverFuncsPackage, utils.Equals+r.Name).Call(cast(receiver), cast(other)))
		case r.Type.IsBytes():
			def.Return(Qual(utils.EqualsPackage, "Bytes").Call(receiver, other))
		default:
			def.Return(Add(receiver).Op("==").Add(other))
		}
	})

	AddCustomComputeHash(def, r.Receiver(), r.Name, r.ShouldReference(), func(def *Group) {
		if r.ShouldReference().ShouldUsePointer() {
			def.Add(If(Add(receiver).Op("==").Nil()).Block(Return(utils.ZeroHash)))
		}

		switch {
		case r.native != nil && r.native.IsCustomStruct():
			def.Return(cast(receiver).Dot(utils.ComputeHash).Call())
		case r.native != nil:
			def.Return(Qual(*r.native.NonReceiverFuncsPackage, utils.ComputeHash+r.Name).Call(cast(receiver)))
		default:
			h := Id("hash")
			def.Add(h).Op(":=").Add(utils.NewHash).Line()
			def.Add(h).Dot(r.Type.HasherName()).Call(cast(receiver))
			def.Return(h)
		}
	})

	if r.native == nil {
		utils.AddPointer(def, r.Receiver(), r.Name)
	}

	tmp := Code(Id("tmp"))
	underlyingType := RestliType{Primitive: r.Type}
	AddMarshalRestLi(def, r.Receiver(), r.Name, r.ShouldReference(), func(def *Group) {
		var accessor Code
		if r.native != nil {
			accessor = tmp
			set := def.List(tmp, Err()).Op(":=")
			if r.native.IsCustomStruct() {
				set.Add(cast(receiver)).Dot(utils.MarshalRestLi).Call()
			} else {
				set.Qual(*r.native.NonReceiverFuncsPackage, utils.MarshalRestLi+r.Name).Call(cast(receiver))
			}
			def.Add(utils.IfErrReturn(Err())).Line()
		} else {
			accessor = cast(receiver)
		}
		def.Add(Writer.Write(underlyingType, Writer, accessor))
		def.Return(Nil())
	})

	AddUnmarshalerFunc(def, r.Receiver(), r.Identifier, r.ShouldReference())

	AddUnmarshalRestli(def, r.Receiver(), r.Name, func(def *Group) {
		def.Var().Add(tmp).Add(r.Type.GoType())
		def.Add(Reader.Read(underlyingType, Reader, tmp))
		def.Add(utils.IfErrReturn(Err())).Line()

		switch {
		case r.native != nil:
			var unmarshalRestLiPackage string
			if r.native.IsCustomStruct() {
				unmarshalRestLiPackage = r.native.NativePackage
			} else {
				unmarshalRestLiPackage = *r.native.NonReceiverFuncsPackage
			}

			unmarshaled := Id("unmarshaled")
			def.List(unmarshaled, Err()).Op(":=").Qual(unmarshalRestLiPackage, utils.UnmarshalRestLi+r.Name).Call(tmp)
			def.Add(utils.IfErrReturn(Err()))
			if r.native.ShouldReference.ShouldUsePointer() {
				unmarshaled = Op("*").Add(unmarshaled)
			}
			def.Op("*").Add(receiver).Op("=").Id(r.Name).Call(unmarshaled)
		default:
			def.Op("*").Id(r.Receiver()).Op("=").Id(r.Name).Call(tmp)
		}
		def.Return(Nil())
	})

	return def
}

var nativeTyperefs = map[string]*NativeTypeRef{}

func RegisterNativeTyperef(n *NativeTypeRef) {
	if _, ok := nativeTyperefs[n.Ref]; ok {
		log.Panicf("Native typeref for %q defined more than once", n.Ref)
	}
	nativeTyperefs[n.Ref] = n
}
