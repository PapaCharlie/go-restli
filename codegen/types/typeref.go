package types

import (
	"encoding/json"
	"plugin"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

var TyperefBindingPlugin *plugin.Plugin

type TyperefCodeProvider interface {
	ReferencedTypes() utils.IdentifierSet
	GenerateType() Code
	GenerateMarshalRaw(def *Group)
	GenerateUnmarshalRaw(raw Code, def *Group)
	GenerateEquals(other Code, def *Group)
	GenerateComputeHash(h Code, def *Group)
	ZeroValue() Code
	ShouldReference() bool
}

type Typeref struct {
	NamedType
	Type         *PrimitiveType `json:"type"`
	codeProvider TyperefCodeProvider
}

func (r *Typeref) UnmarshalJSON(data []byte) (err error) {
	type _t Typeref
	err = json.Unmarshal(data, (*_t)(r))
	if err != nil {
		return err
	}
	if TyperefBindingPlugin != nil {
		var f plugin.Symbol
		f, err = TyperefBindingPlugin.Lookup("GetTyperefCodeProvider")
		if err != nil {
			return err
		}
		r.codeProvider = f.(func(*Typeref) TyperefCodeProvider)(r)
	}
	return nil
}

func (r *Typeref) InnerTypes() utils.IdentifierSet {
	if r.codeProvider != nil {
		return r.codeProvider.ReferencedTypes()
	} else {
		return nil
	}
}

func (r *Typeref) GenerateCode() Code {
	underlyingType := RestliType{Primitive: r.Type}
	cast := r.Type.Cast(Op("*").Id(r.Receiver()))

	def := Empty()
	utils.AddWordWrappedComment(def, r.Doc).Line()
	if r.codeProvider != nil {
		def.Add(r.codeProvider.GenerateType()).Line().Line()
	} else {
		def.Type().Id(r.Name).Add(r.Type.GoType()).Line().Line()
	}

	AddEquals(def, r.Receiver(), r.Name, func(other Code, def *Group) {
		if r.codeProvider != nil {
			r.codeProvider.GenerateEquals(other, def)
		} else {
			left, right := Op("*").Id(r.Receiver()), Op("*").Add(other)

			if r.Type.IsBytes() {
				def.Return(Qual("bytes", "Equal").Call(left, right))
			} else {
				def.Return(Add(left).Op("==").Add(right))
			}
		}
	})
	AddComputeHash(def, r.Receiver(), r.Name, func(h Code, def *Group) {
		if r.codeProvider != nil {
			r.codeProvider.GenerateComputeHash(h, def)
		} else {
			def.Add(h).Dot(r.Type.HasherName()).Call(r.Type.Cast(Op("*").Id(r.Receiver())))
		}
		def.Return(h)
	})

	AddMarshalRestLi(def, r.Receiver(), r.Name, func(def *Group) {
		var raw Code
		if r.codeProvider != nil {
			raw = utils.Raw
			def.Var().Add(raw).Add(r.Type.GoType())
			def.List(raw, Err()).Op("=").Id(r.Receiver()).Dot(utils.MarshalRaw).Call()
			def.Add(utils.IfErrReturn(Err()))
		} else {
			raw = cast
		}

		def.Add(Writer.Write(underlyingType, Writer, raw))
		def.Return(Nil())
	})
	AddUnmarshalRestli(def, r.Receiver(), r.Name, func(def *Group) {
		def.Var().Add(utils.Raw).Add(r.Type.GoType())
		def.Add(Reader.Read(underlyingType, Reader, utils.Raw))
		def.Add(utils.IfErrReturn(Err())).Line()

		if r.codeProvider != nil {
			def.Return(Id(r.Receiver()).Dot(utils.UnmarshalRaw).Call(utils.Raw))
		} else {
			def.Op("*").Id(r.Receiver()).Op("=").Id(r.Name).Call(utils.Raw)
			def.Return(Nil())
		}
	})

	if r.codeProvider != nil {
		utils.AddFuncOnReceiver(def, r.Receiver(), r.Name, utils.MarshalRaw).
			Params().
			Params(r.Type.GoType(), Error()).
			BlockFunc(r.codeProvider.GenerateMarshalRaw).
			Line().Line()

		utils.AddFuncOnReceiver(def, r.Receiver(), r.Name, utils.UnmarshalRaw).
			Params(Add(utils.Raw).Add(r.Type.GoType())).
			Params(Error()).
			BlockFunc(func(def *Group) { r.codeProvider.GenerateUnmarshalRaw(utils.Raw, def) }).
			Line().Line()
	}

	return def
}
