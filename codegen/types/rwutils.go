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
	KeyWriter     = Code(Id("keyWriter"))
	KeyWriterFunc = Code(Add(KeyWriter).Func().Params(String()).Add(WriterQual))
	ItemWriter    = Code(Id("itemWriter"))
	Writer        = &writer{Id("writer")}
	WriterQual    = Code(Qual(utils.RestLiCodecPackage, "Writer"))
	WriterParam   = Code(Add(Writer).Add(WriterQual))
)

func (e *writer) WriteMap(writerAccessor Code, writer func(keyWriter Code, def *Group)) Code {
	return Add(writerAccessor).Dot("WriteMap").Call(Func().Params(KeyWriterFunc).Params(Err().Error()).BlockFunc(func(def *Group) {
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
	case t.IsCustomTyperef():
		def := Err().Op("=").Add(writeCustomTyperef(writerAccessor, sourceAccessor, *t.Reference)).Line()
		def.Add(utils.IfErrReturn(returnOnError...))
		return def
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
	case t.IsCustomTyperef():
		v := Id("t")
		return Func().Params(Add(v).Add(t.GoType()), Add(WriterParam)).Error().Block(
			Return(writeCustomTyperef(Writer, v, *t.Reference)),
		)
	case t.Reference != nil:
		q := t.Reference.Qual()
		if t.ShouldReference() {
			q = Parens(Op("*").Add(q))
		}
		return q.Dot(utils.MarshalRestLi)
	case t.IsMapOrArray():
		v := Id("t")
		return Func().
			Params(Add(v).Add(t.GoType()), Add(Writer).Add(WriterQual)).Error().
			Block(Return(e.NestedMarshaler(t, Writer, v)))
	default:
		log.Panicf("Illegal restli type: %+v", t)
		return nil
	}
}

func (e *writer) NestedMarshaler(t RestliType, writerAccessor, sourceAccessor Code) Code {
	innerT, word := t.InnerMapOrArray()
	return Qual(utils.RestLiCodecPackage, "Write"+word).Call(writerAccessor, sourceAccessor, e.MarshalerFunc(innerT))
}

func (e *writer) Finalize() Code {
	return Add(e).Dot("Finalize").Call()
}

type reader struct {
	Code
}

var (
	Reader         = &reader{Id("reader")}
	ReaderQual     = Code(Qual(utils.RestLiCodecPackage, "Reader"))
	ReaderParam    = Code(Add(Reader).Add(ReaderQual))
	FieldParamName = Code(Id("field"))
)

func (d *reader) ReadMap(reader Code, mapReader func(reader, key Code, def *Group)) Code {
	key := Id("key")
	return Add(reader).Dot("ReadMap").Call(Func().Params(Add(d).Add(ReaderQual), Add(key).String()).Params(Err().Error()).BlockFunc(func(def *Group) {
		mapReader(d, key, def)
	}))
}

func (d *reader) ReadRecord(reader Code, requiredFields Code, mapReader func(reader, key Code, def *Group)) Code {
	return Add(reader).Dot("ReadRecord").Call(requiredFields, Func().Params(Add(d).Add(ReaderQual), Add(FieldParamName).String()).Params(Err().Error()).BlockFunc(func(def *Group) {
		mapReader(d, FieldParamName, def)
	}))
}

func (d *reader) ReadArray(reader Code, arrayReader func(reader Code, def *Group)) Code {
	return Add(reader).Dot("ReadArray").Call(Func().Params(Add(d).Add(ReaderQual)).Params(Err().Error()).BlockFunc(func(def *Group) {
		arrayReader(d, def)
	}))
}

func (d *reader) Skip(reader Code) *Statement {
	return Add(reader).Dot("Skip").Call()
}

func (d *reader) Read(t RestliType, reader, targetAccessor Code) Code {
	switch {
	case t.Primitive != nil:
		return List(targetAccessor, Err()).Op("=").Add(reader).Dot(t.Primitive.ReaderName()).Call()
	case t.IsCustomTyperef():
		return List(targetAccessor, Err()).Op("=").Add(readCustomTyperef(reader, *t.Reference))
	case t.Reference != nil:
		return Err().Op("=").Add(targetAccessor).Dot(utils.UnmarshalRestLi).Call(reader)
	case t.IsMapOrArray():
		return List(targetAccessor, Err()).Op("=").Add(d.NestedUnmarshaler(t, reader))
	default:
		log.Panicf("Illegal restli type: %+v", t)
		return nil
	}
}

func (d *reader) UnmarshalerFunc(t RestliType) Code {
	switch {
	case t.Primitive != nil:
		return t.Primitive.UnmarshalerFunc()
	case t.IsCustomTyperef():
		return Func().
			Params(Add(Reader).Add(ReaderQual)).
			Params(t.GoType(), Error()).
			BlockFunc(func(def *Group) {
				def.Return(readCustomTyperef(Reader, *t.Reference))
			})
	case t.Reference != nil:
		return Qual(utils.RestLiCodecPackage, utils.UnmarshalRestLi).Index(t.ReferencedType())
	case t.IsMapOrArray():
		return Func().
			Params(Add(Reader).Add(ReaderQual)).
			Params(t.GoType(), Error()).
			BlockFunc(func(def *Group) {
				def.Return(d.NestedUnmarshaler(t, Reader))
			})
	default:
		log.Panicf("Illegal restli type: %+v", t)
		return nil
	}
}

func (d *reader) NestedUnmarshaler(t RestliType, reader Code) Code {
	innerT, word := t.InnerMapOrArray()
	return Qual(utils.RestLiCodecPackage, "Read"+word).Call(reader, d.UnmarshalerFunc(innerT))
}
