package types

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func AddComputeHash(def *Statement, receiver, typeName string, f func(h Code, def *Group)) *Statement {
	h := Id("hash")
	return utils.AddFuncOnReceiver(def, receiver, typeName, utils.ComputeHash).
		Params().Params(Add(h).Add(utils.Hash)).
		BlockFunc(func(def *Group) {
			def.Add(If(Id(receiver).Op("==").Nil()).Block(Return(h)))
			f(h, def)
		}).Line().Line()
}

func (r *Record) GenerateComputeHash() Code {
	return AddComputeHash(Empty(), r.Receiver(), r.Name, func(h Code, def *Group) {
		def.Add(h).Op("=").Add(utils.NewHash).Line()
		for _, f := range r.SortedFields() {
			def.Add(hash(h, f.Type, f.IsOptionalOrDefault(), r.fieldAccessor(f))).Line()
		}
		def.Return(h)
	})
}

func hash(h Code, t RestliType, isPointer bool, accessor Code) Code {
	def := Empty()
	if isPointer {
		def.If(Add(accessor).Op("!=").Nil())
		if t.Reference == nil {
			accessor = Op("*").Add(accessor)
		}
	}
	return def.BlockFunc(func(def *Group) {
		switch {
		case t.Primitive != nil:
			def.Add(h).Dot(t.Primitive.HasherName()).Call(accessor)
		case t.Reference != nil:
			def.Add(h).Dot("Add").Call(Add(accessor).Dot(utils.ComputeHash).Call())
		case t.Array != nil:
			_, item := tempIteratorVariableNames(t)
			def.For().List(Id("_"), item).Op(":=").Range().Add(accessor).BlockFunc(func(def *Group) {
				def.Add(hash(h, *t.Array, t.Array.ShouldReference(), item))
			})
		case t.Map != nil:
			hashSum := Id("hashSum")
			def.Var().Add(hashSum).Add(utils.Hash)
			key, value := tempIteratorVariableNames(t)
			def.For().List(key, value).Op(":=").Range().Add(accessor).BlockFunc(func(def *Group) {
				kvHash := Id("hvHash")
				def.Add(kvHash).Op(":=").Add(utils.NewHash)
				def.Add(hash(kvHash, RestliType{Primitive: &StringPrimitive}, false, key)).Line()
				def.Add(hash(kvHash, *t.Map, t.Map.ShouldReference(), value)).Line()
				def.Add(hashSum).Op("+=").Add(kvHash)
			})
			def.Add(h).Dot("Add").Call(hashSum)
		case t.NativeTyperef != nil:
			def.Add(h).Dot("Add").Call(t.NativeTyperef.ComputeHash().Call(accessor))
		}
	})
}
