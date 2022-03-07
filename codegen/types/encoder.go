package types

import (
	"log"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type writerUtils struct{}

var (
	KeyWriter        = Code(Id("keyWriter"))
	ItemWriter       = Code(Id("itemWriter"))
	WriterUtils      = &writerUtils{}
	WriterQual  Code = Qual(utils.RestLiCodecPackage, "Writer")
	Writer      Code = Id("writer")
	WriterParam Code = Add(Writer).Add(WriterQual)
)

func (e *writerUtils) WriteMap(writerAccessor Code, writer func(keyWriter Code, def *Group)) Code {
	keyWriterFunc := Add(KeyWriter).Func().Params(String()).Add(WriterQual)
	return Add(writerAccessor).Dot("WriteMap").Call(Func().Params(keyWriterFunc).Params(Err().Error()).BlockFunc(func(def *Group) {
		writer(KeyWriter, def)
	}))
}

func (e *writerUtils) WriteArray(writerAccessor Code, writer func(itemWriter Code, def *Group)) Code {
	itemWriterFunc := Add(ItemWriter).Func().Params().Add(WriterQual)
	return Add(writerAccessor).Dot("WriteArray").Call(Func().Params(itemWriterFunc).Params(Err().Error()).BlockFunc(func(def *Group) {
		writer(ItemWriter, def)
	}))
}

func (e *writerUtils) Write(t RestliType, writerAccessor, sourceAccessor Code, returnOnError ...Code) Code {
	switch {
	case t.Primitive != nil:
		return Add(writerAccessor).Dot(t.Primitive.WriterName()).Call(sourceAccessor)
	case t.Reference != nil:
		def := Err().Op("=").Add(sourceAccessor).Dot(utils.MarshalRestLi).Call(writerAccessor).Line()
		def.Add(utils.IfErrReturn(returnOnError...))
		return def
	case t.IsMapOrArray():
		def := Err().Op("=").Add(e.NestedMarshaler(t, writerAccessor, sourceAccessor)).Line()
		def.Add(utils.IfErrReturn(returnOnError...))
		return def
	default:
		log.Panicf("Illegal restli type: %+v", t)
		return nil
	}
}

func (e *writerUtils) MarshalerFunc(t RestliType) Code {
	switch {
	case t.Primitive != nil:
		return t.Primitive.MarshalerFunc()
	case t.Reference != nil:
		return t.Reference.Qual().Dot(utils.MarshalRestLi)
	case t.IsMapOrArray():
		v := Id("t")
		return Func().
			Params(Add(v).Add(t.GoType()), WriterParam).Error().
			Block(Return(e.NestedMarshaler(t, Writer, v)))
	default:
		log.Panicf("Illegal restli type: %+v", t)
		return nil
	}
}

func (e *writerUtils) NestedMarshaler(t RestliType, writerAccessor, sourceAccessor Code) Code {
	innerT, word := t.InnerMapOrArray()
	writeFunc := func(variant string) *Statement {
		return Qual(utils.RestLiCodecPackage, "Write"+variant+word)
	}

	switch {
	case innerT.Primitive != nil:
		return writeFunc("Primitive").Call(writerAccessor, sourceAccessor, innerT.Primitive.MarshalerFunc())
	case innerT.Reference != nil:
		return writeFunc("Object").Call(writerAccessor, sourceAccessor)
	default:
		return writeFunc("").Call(writerAccessor, sourceAccessor, e.MarshalerFunc(innerT))
	}
}

func (e *writerUtils) Finalize(writerAccessor Code) Code {
	return Add(writerAccessor).Dot("Finalize").Call()
}

func (e *writerUtils) IsKeyExcluded(writerAccessor, key Code) Code {
	return Add(writerAccessor).Dot("IsKeyExcluded").Call(key)
}
