package types

import (
	"github.com/PapaCharlie/go-restli/v2/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func AddEquals(def *Statement, receiver, typeName string, pointer utils.ShouldUsePointer, f func(other Code, def *Group)) *Statement {
	other := Id("other")
	rightHandType := Id(typeName)
	if pointer.ShouldUsePointer() {
		rightHandType = Op("*").Add(rightHandType)
	}

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
	return AddEquals(Empty(), r.Receiver(), r.TypeName(), RecordShouldUsePointer, func(other Code, def *Group) {
		if len(r.Fields) == 0 && len(r.Includes) == 0 {
			def.Return(True())
			return
		}

		exp := Empty()
		count := 0
		addAnd := func() {
			if count != 0 {
				exp.Op("&&").Line()
			}
			count++
		}
		for _, i := range r.Includes {
			addAnd()
			if i.IsEmptyRecord() {
				exp.Id(r.Receiver()).Dot(i.TypeName()).Dot(utils.Equals).Call(Add(other).Dot(i.TypeName()))
				// Ignore EmptyRecord as it will never change the outcome of .Equals since it will never have fields
				continue
			} else {
				exp.Id(r.Receiver()).Dot(i.TypeName()).Dot(utils.Equals).Call(Op("&").Add(other).Dot(i.TypeName()))
			}
		}

		for _, f := range r.Fields {
			left, right := r.fieldAccessor(f), fieldAccessor(other, f)
			addAnd()
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
		if !t.Primitive.IsBytes() && !isPointer {
			condition = Add(left).Op("==").Add(right)
		} else {
			if t.Primitive.IsBytes() {
				prefix = "Bytes"
			} else {
				prefix = "Comparable"
			}
			if isPointer {
				pointer = "Pointer"
			}
			condition = equalsFunc(prefix+pointer).Call(left, right)
		}
	case t.Reference != nil && !t.ShouldReference(): // enums and typerefs
		if t.Reference.IsCustomTyperef() {
			f := customTyperefEqualsFunc(*t.Reference)
			if isPointer {
				condition = equalsFunc("GenericPointer").Call(left, right, f)
			} else {
				condition = Add(f).Call(left, right)
			}
		} else {
			if isPointer {
				condition = equalsFunc("ObjectPointer").Call(left, right)
			} else {
				condition = Add(left).Dot(utils.Equals).Call(right)
			}
		}
	case t.Reference != nil:
		if !isPointer {
			right = Op("&").Add(right)
		}
		condition = Add(left).Dot(utils.Equals).Call(right)
	case t.IsMapOrArray():
		innerT, word := t.InnerMapOrArray()
		if isPointer {
			word += "Pointer"
		}

		switch {
		case innerT.Primitive != nil && innerT.Primitive.IsBytes():
			condition = equalsFunc("Bytes"+word).Call(left, right)
		case innerT.Primitive != nil:
			condition = equalsFunc("Comparable"+word).Call(left, right)
		case innerT.Reference != nil && innerT.Reference.IsCustomTyperef():
			condition = equalsFunc("Generic"+word).Call(left, right, customTyperefEqualsFunc(*innerT.Reference))
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
