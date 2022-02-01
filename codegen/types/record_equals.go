package types

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func AddEquals(def *Statement, receiver, typeName string, f func(other Code, def *Group)) *Statement {
	return AddCustomEquals(def, receiver, typeName, func(other Code, def *Group) {
		def.If(Id(receiver).Op("==").Nil().Op("||").Add(other).Op("==").Nil()).Block(Return(False())).Line()
		f(other, def)
	})
}

func AddCustomEquals(def *Statement, receiver, typeName string, f func(other Code, def *Group)) *Statement {
	other := Id("other")
	otherInterface := Id("otherInterface")
	rightHandType := Op("*").Id(typeName)
	utils.AddFuncOnReceiver(def, receiver, typeName, utils.EqualsInterface).
		Params(Add(otherInterface).Interface()).Bool().
		BlockFunc(func(def *Group) {
			ok := Id("ok")
			def.List(other, ok).Op(":=").Add(otherInterface).Assert(rightHandType)
			def.If(Op("!").Add(ok)).Block(Return(False())).Line()
			def.Return(Id(receiver).Dot(utils.Equals).Call(other))
		}).Line().Line()
	return utils.AddFuncOnReceiver(def, receiver, typeName, utils.Equals).
		Params(Add(other).Add(rightHandType)).Bool().
		BlockFunc(func(def *Group) {
			f(other, def)
		}).Line().Line()
}

func (r *Record) GenerateEquals() Code {
	return AddEquals(Empty(), r.Receiver(), r.Name, func(other Code, def *Group) {
		for _, f := range r.SortedFields() {
			left, right := r.fieldAccessor(f), fieldAccessor(other, f)
			def.Add(equals(f.Type, f.IsOptionalOrDefault(), left, right)).Line()
		}
		def.Return(True())
	})
}

func equals(t RestliType, isPointer bool, left, right Code) Code {
	allocateNewRight := func(def *Group, t RestliType, right Code) Code {
		if t.Typeref() != nil {
			ref := Id("ref")
			def.Add(ref).Op(":=").Add(right)
			return ref
		} else {
			return right
		}
	}
	check := func(left, right Code) Code {
		def := Empty()
		switch {
		case t.Primitive != nil:
			if t.Primitive.IsBytes() {
				def.If(Op("!").Qual("bytes", "Equal").Call(left, right)).Block(Return(False()))
			} else {
				def.If(Add(left).Op("!=").Add(right)).Block(Return(False()))
			}
		case t.Reference != nil:
			if !isPointer {
				right = Op("&").Add(right)
			}
			def.If(Op("!").Add(left).Dot(utils.Equals).Call(right)).Block(Return(False()))
		case t.Array != nil:
			def.If(Len(left).Op("!=").Len(right)).Block(Return(False())).Line()
			index, item := tempIteratorVariableNames(t)
			def.For().List(index, item).Op(":=").Range().Add(left).BlockFunc(func(def *Group) {
				def.Add(equals(*t.Array, t.Array.ShouldReference(), item,
					allocateNewRight(def, *t.Array, Parens(right).Index(index))))
			})
		case t.Map != nil:
			def.If(Len(left).Op("!=").Len(right)).Block(Return(False())).Line()
			key, value := tempIteratorVariableNames(t)
			def.For().List(key, value).Op(":=").Range().Add(left).BlockFunc(func(def *Group) {
				def.Add(equals(*t.Map, t.Map.ShouldReference(), value,
					allocateNewRight(def, *t.Map, Parens(right).Index(key))))
			})
		}
		return def
	}
	if isPointer {
		return If(Add(left).Op("!=").Add(right)).BlockFunc(func(def *Group) {
			def.If(Add(left).Op("==").Nil().Op("||").Add(right).Op("==").Nil()).Block(Return(False()))

			if t.Reference == nil {
				left, right = Op("*").Add(left), Op("*").Add(right)
			}
			def.Add(check(left, right))
		})
	} else {
		return check(left, right)
	}
}
