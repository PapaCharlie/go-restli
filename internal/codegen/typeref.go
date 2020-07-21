package codegen

import (
	"log"

	. "github.com/dave/jennifer/jen"
)

type Typeref struct {
	NamedType
	Primitive *PrimitiveType `json:"primitive"`
	Reference *Identifier    `json:"reference"`
	Array     *RestliType    `json:"array"`
	Map       *RestliType    `json:"map"`
	Union     *UnionType     `json:"union"`
}

func (r *Typeref) InnerTypes() IdentifierSet {
	switch {
	case r.Primitive != nil:
		return nil
	case r.Reference != nil:
		return NewIdentifierSet(*r.Reference)
	case r.Array != nil:
		return r.Array.InnerTypes()
	case r.Map != nil:
		return r.Map.InnerTypes()
	case r.Union != nil:
		return r.Union.InnerModels()
	default:
		log.Panicf("Unknown reference type: %+v", r)
		return nil
	}
}

func (r *Typeref) refGoType() *Statement {
	switch {
	case r.Primitive != nil:
		return r.Primitive.GoType()
	case r.Reference != nil:
		return r.Reference.Qual()
	case r.Array != nil:
		return r.Array.GoType()
	case r.Map != nil:
		return r.Map.GoType()
	case r.Union != nil:
		return r.Union.GoType()
	default:
		log.Panicf("Unknown reference type: %+v", r)
		return nil
	}
}

func (r *Typeref) GenerateCode() (def *Statement) {
	def = Empty()

	AddWordWrappedComment(def, r.Doc).Line()
	def.Type().Id(r.Name).Add(r.refGoType()).Line().Line()

	if pt := r.underlyingPrimitiveType(); pt != nil {
		AddRestLiEncode(def, r.Receiver(), r.Name, func(def *Group) {
			writeStringToBuf(def, pt.encode(pt.Cast(Op("*").Id(r.Receiver()))))
			def.Return(Nil())
		}).Line().Line()
		AddRestLiDecode(def, r.Receiver(), r.Name, func(def *Group) {
			def.Return(pt.decode(Id(r.Receiver())))
		}).Line().Line()

		return def
	}

	if ref := r.Reference; ref != nil {
		AddRestLiEncode(def, r.Receiver(), r.Name, func(def *Group) {
			def.Return(Parens(Op("*").Add(r.Reference.Qual()).Dot(RestLiEncode).Params(Id(Codec), Id("buf"))))
		}).Line().Line()
	}

	if m := r.Array; m != nil {
		AddRestLiEncode(def, r.Receiver(), r.Name, func(def *Group) {
			writeArrayToBuf(def, Id(r.Receiver()), r.Array)
		}).Line().Line()
	}

	if m := r.Map; m != nil {
		AddRestLiEncode(def, r.Receiver(), r.Name, func(def *Group) {
			writeMapToBuf(def, Id(r.Receiver()), r.Map)
		}).Line().Line()
	}

	if union := r.Union; union != nil {
		AddFuncOnReceiver(def, r.Receiver(), r.Name, ValidateUnionFields).
			Params().
			Params(Error()).
			BlockFunc(func(def *Group) {
				union.validateUnionFields(def, r.Receiver(), r.Name)
			})

		AddRestLiEncode(def, r.Receiver(), r.Name, func(def *Group) {
			union.encode(def, r.Receiver(), r.Name)
		}).Line().Line()

		return def
	}

	return nil
}

func (r *Typeref) underlyingPrimitiveType() *PrimitiveType {
	switch {
	case r.Primitive != nil:
		return r.Primitive
	case r.Reference != nil:
		if ref, ok := r.Reference.Resolve().(*Typeref); ok {
			return ref.underlyingPrimitiveType()
		}
	}
	return nil
}
