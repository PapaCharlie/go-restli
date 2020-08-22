package codegen

import (
	"log"

	. "github.com/dave/jennifer/jen"
)

type encoder struct {
	*Statement
}

var Encoder = &encoder{Id("encoder")}

func (e *encoder) WriteObjectStart() *Statement {
	return Add(e).Dot("WriteObjectStart").Call()
}

func (e *encoder) WriteObjectEnd() *Statement {
	return Add(e).Dot("WriteObjectEnd").Call()
}

func (e *encoder) WriteFieldDelimiter() *Statement {
	return Add(e).Dot("WriteFieldDelimiter").Call()
}

func (e *encoder) WriteFieldNameAndDelimiter(fieldName string) *Statement {
	return Add(e).Dot("WriteFieldNameAndDelimiter").Call(Lit(fieldName))
}

func (e *encoder) WriteField(def *Group, fieldName string, t RestliType, fieldAccessor *Statement, returnOnError ...Code) {
	def.Add(e.WriteFieldNameAndDelimiter(fieldName))
	e.Write(def, t, fieldAccessor, returnOnError...)
}

func (e *encoder) Write(def *Group, t RestliType, accessor *Statement, returnOnError ...Code) {
	switch {
	case t.Primitive != nil:
		def.Add(e).Dot(t.Primitive.EncoderName()).Call(accessor)
		return
	case t.Reference != nil:
		def.Err().Op("=").Add(e).Dot("Encodable").Call(accessor)
		IfErrReturn(def, append(append([]Code(nil), returnOnError...), Err())...)
	case t.Array != nil:
		e.ArrayEncoder(def, func(def *Group, indexWriter *Statement) {
			index, item := Id("index"), Id("item")
			def.For(List(index, item).Op(":=").Range().Add(accessor)).BlockFunc(func(def *Group) {
				def.Add(indexWriter).Call(index)
				if t.Array.IsReferenceEncodable() {
					item = Op("&").Add(item)
				}
				e.Write(def, *t.Array, item)
			})
			def.Return(Nil())
		})
		IfErrReturn(def, append(append([]Code(nil), returnOnError...), Err())...)
	case t.Map != nil:
		e.MapEncoder(def, func(def *Group, keyWriter *Statement) {
			key, value := Id("key"), Id("value")
			def.For(List(key, value).Op(":=").Range().Parens(accessor)).BlockFunc(func(def *Group) {
				def.Add(keyWriter).Call(key)
				if t.Map.IsReferenceEncodable() {
					value = Op("&").Add(value)
				}
				e.Write(def, *t.Map, value)
			})
			def.Return(Nil())
		})
		IfErrReturn(def, append(append([]Code(nil), returnOnError...), Err())...)
	default:
		log.Panicf("Illegal restli type: %+v", t)
	}
}

func (e *encoder) ArrayEncoder(def *Group, block func(def *Group, indexWriter *Statement)) {
	indexWriter := Id("indexWriter")
	def.Err().Op("=").Add(e).Dot("Array").Call(Func().Params(Add(indexWriter).Func().Params(Id("index").Int())).Params(Err().Error()).BlockFunc(func(def *Group) {
		block(def, indexWriter)
	}))
}

func (e *encoder) MapEncoder(def *Group, block func(def *Group, keyWriter *Statement)) {
	keyWriter := Id("keyWriter")
	def.Err().Op("=").Add(e).Dot("Map").Call(Func().Params(Add(keyWriter).Func().Params(Id("key").String())).Params(Err().Error()).BlockFunc(func(def *Group) {
		block(def, keyWriter)
	}))

}

func (e *encoder) Finalize() *Statement {
	return Add(e).Dot("Finalize").Call()
}
