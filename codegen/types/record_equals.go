package types

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func AddEquals(def *Statement, receiver, typeName string, pointer utils.ShouldUsePointer, f func(other Code, def *Group)) *Statement {
	other := Id("other")
	otherInterface := Id("otherInterface")
	rightHandType := Id(typeName)
	if pointer.ShouldUsePointer() {
		rightHandType = Op("*").Add(rightHandType)
	}
	utils.AddFuncOnReceiver(def, receiver, typeName, utils.EqualsInterface, pointer).
		Params(Add(otherInterface).Interface()).Bool().
		BlockFunc(func(def *Group) {
			ok := Id("ok")
			def.List(other, ok).Op(":=").Add(otherInterface).Assert(rightHandType)
			def.If(Op("!").Add(ok)).Block(Return(False())).Line()
			def.Return(Id(receiver).Dot(utils.Equals).Call(other))
		}).Line().Line()
	return utils.AddFuncOnReceiver(def, receiver, typeName, utils.Equals, pointer).
		Params(Add(other).Add(rightHandType)).Bool().
		BlockFunc(func(def *Group) {
			if pointer.ShouldUsePointer() {
				def.If(Id(receiver).Op("==").Add(other)).Block(Return(True()))
				def.If(Id(receiver).Op("==").Nil().Op("||").Add(other).Op("==").Nil()).Block(Return(False())).Line()
			}
			f(other, def)
		}).Line().Line()
}

func (r *Record) GenerateEquals() Code {
	return AddEquals(Empty(), r.Receiver(), r.Name, RecordShouldUsePointer, func(other Code, def *Group) {
		if len(r.Fields) == 0 {
			def.Return(True())
			return
		}

		exp := Empty()
		for i, f := range r.Fields {
			left, right := r.fieldAccessor(f), fieldAccessor(other, f)
			if i != 0 {
				exp.Op("&&").Line()
			}
			exp.Add(equalsCondition(f.Type, f.IsOptionalOrDefault(), left, right))
		}
		def.Return(exp)
	})
}

func equalsCondition(t RestliType, isPointer bool, left, right Code) Code {
	var condition Code
	switch {
	case t.Primitive != nil:
		var prefix, pointer string
		if t.Primitive.IsBytes() {
			prefix = "Bytes"
		} else {
			prefix = "Comparable"
		}
		if isPointer {
			pointer = "Pointer"
		}
		condition = equalsFunc(prefix+pointer).Call(left, right)
	case t.Reference != nil && !t.ShouldReference(): // enums and typerefs
		if isPointer {
			condition = equalsFunc("ObjectPointer").Call(left, right)
		} else {
			condition = Add(left).Dot(utils.Equals).Call(right)
		}
	case t.Reference != nil:
		if !isPointer {
			right = Op("&").Add(right)
		}
		condition = Add(left).Dot(utils.Equals).Call(right)
	case t.IsMapOrArray():
		var innerT RestliType
		var word string
		if t.Array != nil {
			innerT = *t.Array
			word = "Array"
		} else {
			innerT = *t.Map
			word = "Map"
		}
		if isPointer {
			word += "Pointer"
		}

		switch {
		case innerT.Primitive != nil && innerT.Primitive.IsBytes():
			condition = equalsFunc("Bytes"+word).Call(left, right)
		case innerT.Primitive != nil:
			condition = equalsFunc("Comparable"+word).Call(left, right)
		case innerT.Reference != nil:
			condition = equalsFunc("Object"+word).Call(left, right)
		case innerT.IsMapOrArray():
			l, r := Id("left"), Id("right")
			condition = equalsFunc("Generic"+word).Call(left, right, Func().Params(List(l, r).Add(innerT.GoType())).Bool().Block(
				Return(equalsCondition(innerT, false, l, r)),
			))
		}
	}

	return condition
}

func equalsFunc(name string) *Statement {
	return Qual(utils.EqualsPackage, name)
}
