package types

import (
	"log"

	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type writer struct {
	Code
}

var (
	Writer          = &writer{Id("writer")}
	WriterQual Code = Qual(RestLiCodecPackage, "Writer")
)

func (e *writer) WriteMap(writerAccessor Code, writer func(keyWriter Code, def *Group)) Code {
	keyWriter := Id("keyWriter")
	keyWriterFunc := Add(keyWriter).Func().Params(String()).Add(WriterQual)
	return Add(writerAccessor).Dot("WriteMap").Call(Func().Params(keyWriterFunc).Params(Err().Error()).BlockFunc(func(def *Group) {
		writer(keyWriter, def)
	}))
}

func (e *writer) WriteArray(writerAccessor Code, writer func(itemWriter Code, def *Group)) Code {
	itemWriter := Id("itemWriter")
	itemWriterFunc := Add(itemWriter).Func().Params().Add(WriterQual)
	return Add(writerAccessor).Dot("WriteArray").Call(Func().Params(itemWriterFunc).Params(Err().Error()).BlockFunc(func(def *Group) {
		writer(itemWriter, def)
	}))
}

func (e *writer) Write(t RestliType, writerAccessor, sourceAccessor Code, returnOnError ...Code) Code {
	switch {
	case t.Primitive != nil:
		return Add(writerAccessor).Dot(t.Primitive.WriterName()).Call(sourceAccessor)
	case t.Reference != nil:
		def := Err().Op("=").Add(sourceAccessor).Dot(MarshalRestLi).Call(writerAccessor).Line()
		def.Add(utils.IfErrReturn(returnOnError...))
		return def
	case t.Array != nil:
		def := Err().Op("=").Add(e.WriteArray(writerAccessor, func(itemWriter Code, def *Group) {
			item := Id("item")
			def.For(List(Id("_"), item).Op(":=").Range().Add(sourceAccessor)).BlockFunc(func(def *Group) {
				def.Add(e.Write(*t.Array, Add(itemWriter).Call(), item, Err()))
			})
			def.Return(Nil())
		})).Line()
		def.Add(utils.IfErrReturn(returnOnError...))
		return def
	case t.Map != nil:
		def := Err().Op("=").Add(e.WriteMap(writerAccessor, func(keyWriter Code, def *Group) {
			key, value := Id("key"), Id("value")
			def.For(List(key, value).Op(":=").Range().Parens(sourceAccessor)).BlockFunc(func(def *Group) {
				def.Add(e.Write(*t.Map, Add(keyWriter).Call(key), value, Err()))
			})
			def.Return(Nil())
		})).Line()
		def.Add(utils.IfErrReturn(returnOnError...))
		return def
	default:
		log.Panicf("Illegal restli type: %+v", t)
		return nil
	}
}

func (e *writer) Finalize() Code {
	return Add(e).Dot("Finalize").Call()
}
