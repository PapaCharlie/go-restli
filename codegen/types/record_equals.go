package types

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func (r *Record) GenerateEquals() Code {
	return AddEquals(Empty(), r.Receiver(), r.Name, RecordShouldUsePointer, func(_, other Code, def *Group) {
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
		innerT, word := t.InnerMapOrArray()
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
