package types

import (
	"log"

	"github.com/PapaCharlie/go-restli/codegen/utils"
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
	case t.IsMapOrArray():
		def := Err().Op("=").Add(e.NestedMarshaler(t, writerAccessor, sourceAccessor)).Line()
		def.Add(utils.IfErrReturn(returnOnError...))
		return def
	default:
		log.Panicf("Illegal restli type: %+v", t)
		return nil
	}
}

func (e *writer) MarshalerFunc(t RestliType) Code {
	switch {
	case t.Primitive != nil:
		return t.Primitive.MarshalerFunc()
	case t.Reference != nil:
		return t.Reference.Qual().Dot(utils.MarshalRestLi)
	case t.IsMapOrArray():
		v := Id("t")
		return Func().
			Params(Add(v).Add(t.GoType()), Add(Writer).Add(WriterQual)).Error().
			BlockFunc(func(def *Group) {
				def.Return(e.NestedMarshaler(t, Writer, v))
			})
	default:
		log.Panicf("Illegal restli type: %+v", t)
		return nil
	}
}

func (e *writer) NestedMarshaler(t RestliType, writerAccessor, sourceAccessor Code) Code {
	var innerT RestliType
	if t.Array != nil {
		innerT = *t.Array
	} else {
		innerT = *t.Map
	}

	switch {
	case innerT.Primitive != nil:
		return writeFunc(t, "Primitive").Call(writerAccessor, sourceAccessor, innerT.Primitive.MarshalerFunc())
	case innerT.Reference != nil:
		return writeFunc(t, "Object").Call(writerAccessor, sourceAccessor)
	default:
		return writeFunc(t, "").Call(writerAccessor, sourceAccessor, e.MarshalerFunc(innerT))
	}

}

func writeFunc(t RestliType, writeFuncType string) *Statement {
	var variant string
	if t.Array != nil {
		variant = "Array"
	} else {
		variant = "Map"
	}
	return Qual(utils.RestLiCodecPackage, "Write"+writeFuncType+variant)
}

func (e *writer) Finalize() Code {
	return Add(e).Dot("Finalize").Call()
}

func (e *writer) IsKeyExcluded(writerAccessor, key Code) Code {
	return Add(writerAccessor).Dot("IsKeyExcluded").Call(key)
}
