package types

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func AddComputeHash(def *Statement, receiver, typeName string, pointer utils.ShouldUsePointer, f func(h Code, def *Group)) *Statement {
	return AddCustomComputeHash(def, receiver, typeName, pointer, func(def *Group) {
		if pointer.ShouldUsePointer() {
			def.Add(If(Id(receiver).Op("==").Nil()).Block(Return(utils.ZeroHash)))
		}
		h := Id("hash")
		def.Add(h).Op(":=").Add(utils.NewHash).Line()
		f(h, def)
		def.Return(h)
	})
}

func AddCustomComputeHash(def *Statement, receiver, typeName string, pointer utils.ShouldUsePointer, f func(def *Group)) *Statement {
	return utils.AddFuncOnReceiver(def, receiver, typeName, utils.ComputeHash, pointer).
		Params().Params(utils.Hash).
		BlockFunc(f).Line().Line()
}

func (r *Record) GenerateComputeHash() Code {
	return AddComputeHash(Empty(), r.Receiver(), r.Name, RecordShouldUsePointer, func(h Code, def *Group) {
		for _, f := range r.Fields {
			def.Add(hash(h, f.Type, f.IsOptionalOrDefault(), r.fieldAccessor(f))).Line()
		}
	})
}

func hash(h Code, t RestliType, isPointer bool, accessor Code) Code {
	hasher := func(accessor Code) Code {
		def := Empty()
		switch {
		case t.Primitive != nil:
			def.Add(h).Dot(t.Primitive.HasherName()).Call(accessor)
		case t.Reference != nil:
			def.Add(h).Dot("Add").Call(Add(accessor).Dot(utils.ComputeHash).Call())
		case t.IsMapOrArray():
			var innerT RestliType
			var add, addHashable Code

			if t.Array != nil {
				innerT = *t.Array
				add, addHashable = utils.AddArray, utils.AddHashableArray
			} else {
				innerT = *t.Map
				add, addHashable = utils.AddMap, utils.AddHashableMap
			}

			switch {
			case innerT.Primitive != nil:
				def.Add(add).Call(h, accessor, innerT.Primitive.HasherQual())
			case innerT.Reference != nil:
				def.Add(addHashable).Call(h, accessor)
			default:
				inner := Id("inner")
				def.Add(add).Call(h, accessor, Func().
					Params(Add(h).Add(utils.Hash), Add(inner).Add(innerT.GoType())).
					BlockFunc(func(def *Group) {
						def.Add(hash(h, innerT, innerT.ShouldReference(), inner))
					}))
			}
		}
		return def
	}

	if isPointer {
		def := If(Add(accessor).Op("!=").Nil())
		if t.Reference == nil {
			accessor = Op("*").Add(accessor)
		}
		return def.Block(hasher(accessor))
	} else {
		return hasher(accessor)
	}
}
