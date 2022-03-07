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
	Package                 string        `json:"package"`
	Name                    string        `json:"name"`
	NonReceiverFuncsPackage *string       `json:"nonReceiverFuncsPackage,omitempty"`

	ShouldReference utils.ShouldUsePointer `json:"-"`
}

func (n *NativeTypeRef) IsCustomStruct() bool {
	return n.NonReceiverFuncsPackage == nil
}

func (n *NativeTypeRef) GoType() *Statement {
	return Qual(n.Package, n.Name)
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
			r.Type.Type, r.native.Package, r.native.Name, r.native.Type.Type)
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
	o := NewObjectCodeGenerator(r.Identifier, r.ShouldReference())

	var cast func(Code) *Statement
	typedef := Empty()
	if r.native != nil {
		cast = func(receiver Code) *Statement {
			def := Qual(r.native.Package, r.native.Name)
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
	o.DeclareType(def, r.Doc, typedef)

	o.Equals(def, func(receiver, other Code, def *Group) {
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

	o.CustomComputeHash(def, func(receiver Code, def *Group) {
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
	o.MarshalRestLi(def, func(receiver, writer Code, def *Group) {
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
		def.Add(WriterUtils.Write(underlyingType, writer, accessor))
		def.Return(Nil())
	})

	o.UnmarshalRestLi(def, func(receiver, reader Code, def *Group) {
		def.Var().Add(tmp).Add(r.Type.GoType())
		def.Add(ReaderUtils.Read(underlyingType, reader, tmp))
		def.Add(utils.IfErrReturn(Err())).Line()

		switch {
		case r.native != nil:
			var unmarshalRestLiPackage string
			if r.native.IsCustomStruct() {
				unmarshalRestLiPackage = r.native.Package
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
