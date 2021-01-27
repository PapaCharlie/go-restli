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
	KeyWriter       = Code(Id("keyWriter"))
	ItemWriter      = Code(Id("itemWriter"))
	Writer          = &writer{Id("writer")}
	WriterQual Code = Qual(utils.RestLiCodecPackage, "Writer")
)

func (e *writer) WriteMap(writerAccessor Code, writer func(keyWriter Code, def *Group)) Code {
	keyWriterFunc := Add(KeyWriter).Func().Params(String()).Add(WriterQual)
	return Add(writerAccessor).Dot("WriteMap").Call(Func().Params(keyWriterFunc).Params(Err().Error()).BlockFunc(func(def *Group) {
		writer(KeyWriter, def)
	}))
}

func (e *writer) WriteArray(writerAccessor Code, writer func(itemWriter Code, def *Group)) Code {
	itemWriterFunc := Add(ItemWriter).Func().Params().Add(WriterQual)
	return Add(writerAccessor).Dot("WriteArray").Call(Func().Params(itemWriterFunc).Params(Err().Error()).BlockFunc(func(def *Group) {
		writer(ItemWriter, def)
	}))
}

func (e *writer) Write(t RestliType, writerAccessor, sourceAccessor Code, returnOnError ...Code) Code {
	switch {
	case t.Primitive != nil:
		return Add(writerAccessor).Dot(t.Primitive.WriterName()).Call(sourceAccessor)
	case t.Reference != nil:
		def := Err().Op("=").Add(sourceAccessor).Dot(utils.MarshalRestLi).Call(writerAccessor).Line()
		def.Add(utils.IfErrReturn(returnOnError...))
		return def
	case t.Array != nil:
		def := Err().Op("=").Add(e.WriteArray(writerAccessor, func(itemWriter Code, def *Group) {
			_, item := tempIteratorVariableNames(t)
			def.For(List(Id("_"), item).Op(":=").Range().Add(sourceAccessor)).BlockFunc(func(def *Group) {
				def.Add(e.Write(*t.Array, Add(itemWriter).Call(), item, Err()))
			})
			def.Return(Nil())
		})).Line()
		def.Add(utils.IfErrReturn(returnOnError...))
		return def
	case t.Map != nil:
		def := Err().Op("=").Add(e.WriteMap(writerAccessor, func(keyWriter Code, def *Group) {
			key, value := tempIteratorVariableNames(t)
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

func (e *writer) IsKeyExcluded(writerAccessor, key Code) Code {
	return Add(writerAccessor).Dot("IsKeyExcluded").Call(key)
}
