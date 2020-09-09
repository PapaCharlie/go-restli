package types

import (
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func AddEquals(def *Statement, receiver, typeName string, f func(other Code, def *Group)) *Statement {
	other := Id("other")
	otherInterface := Id("otherInterface")
	return utils.AddFuncOnReceiver(def, receiver, typeName, Equals).
		Params(Add(otherInterface).Interface()).Bool().
		BlockFunc(func(def *Group) {
			ok := Id("ok")
			def.List(other, ok).Op(":=").Add(otherInterface).Assert(Op("*").Id(typeName))
			def.If(Op("!").Add(ok)).Block(Return(False())).Line()
			def.If(Id(receiver).Op("==").Add(other)).Block(Return(True()))
			def.If(Id(receiver).Op("==").Nil().Op("||").Add(other).Op("==").Nil()).Block(Return(False())).Line()
			f(other, def)
		}).Line().Line()
}

func (r *Record) GenerateEquals() Code {
	return AddEquals(Empty(), r.Receiver(), r.Name, func(other Code, def *Group) {
		for _, f := range r.SortedFields() {
			left, right := Id(r.Receiver()).Dot(f.FieldName()), Add(other).Dot(f.FieldName())
			def.Add(equals(f.Type, f.IsOptionalOrDefault(), left, right)).Line()
		}
		def.Return(True())
	})
}

func equals(t RestliType, isPointer bool, left, right Code) Code {
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
			def.If(Op("!").Add(left).Dot(Equals).Call(right)).Block(Return(False()))
		case t.Array != nil:
			def.If(Len(left).Op("!=").Len(right)).Block(Return(False())).Line()
			index, item := Id("index"), Id("item")
			def.For().List(index, item).Op(":=").Range().Add(left).BlockFunc(func(def *Group) {
				def.Add(equals(*t.Array, t.Array.ShouldReference(), item, Parens(right).Index(index)))
			})
		case t.Map != nil:
			def.If(Len(left).Op("!=").Len(right)).Block(Return(False())).Line()
			key, value := Id("key"), Id("value")
			def.For().List(key, value).Op(":=").Range().Add(left).BlockFunc(func(def *Group) {
				def.Add(equals(*t.Map, t.Map.ShouldReference(), value, Parens(right).Index(key)))
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
