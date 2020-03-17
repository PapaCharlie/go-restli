package codegen

import (
	. "github.com/dave/jennifer/jen"
)

type Typeref struct {
	NamedType
	Ref RestliType
}

func (r *Typeref) InnerTypes() IdentifierSet {
	return r.Ref.InnerTypes()
}

func (r *Typeref) GenerateCode() (def *Statement) {
	def = Empty()

	if ref := r.Ref.Reference; ref != nil {
		// TODO
		Logger.Printf("Warning: type references to non-primitive types are not yet supported (%s)", r.Identifier)
		return def
	}

	AddWordWrappedComment(def, r.Doc).Line()
	def.Type().Id(r.Name).Add(r.Ref.GoType()).Line().Line()

	if pt := r.Ref.Primitive; pt != nil {
		AddRestLiEncode(def, r.Receiver(), r.Name, func(def *Group) {
			def.Return(pt.encode(pt.Cast(Op("*").Id(r.Receiver()))), Nil())
		}).Line().Line()
		AddRestLiDecode(def, r.Receiver(), r.Name, func(def *Group) {
			def.Return(pt.decode(Id(r.Receiver())))
		}).Line().Line()

		return def
	}

	if union := r.Ref.Union; union != nil {
		AddRestLiEncode(def, r.Receiver(), r.Name, func(def *Group) {
			def.Err().Op("=").Id(r.Receiver()).Dot(ValidateUnionFields).Call()
			def.If(Err().Op("!=").Nil()).Block(Return()).Line()
			def.Var().Id("buf").Qual("strings", "Builder")
			r.Ref.WriteToBuf(def, Id(r.Receiver()))
			def.Id("data").Op("=").Id("buf").Dot("String").Call()
			def.Return()
		}).Line().Line()

		AddFuncOnReceiver(def, r.Receiver(), r.Name, ValidateUnionFields).
			Params().
			Params(Error()).
			BlockFunc(func(def *Group) {
				union.validateUnionFields(def, Id(r.Receiver()))
				def.Line().Return(Nil())
			})

		return def
	}

	Logger.Panicf("Illegal typeref type %+v defined in %s", r.Ref, r.GetSourceFile())
	return nil
}

func (r *Typeref) isPrimitive() bool {
	switch {
	case r.Ref.Primitive != nil:
		return true
	case r.Ref.Reference != nil:
		if ref, ok := r.Ref.Reference.Resolve().(*Typeref); ok {
			return ref.isPrimitive()
		}
	}
	return false
}
