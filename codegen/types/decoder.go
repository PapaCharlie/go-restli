package types

import (
	"log"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type readerUtils struct{}

var (
	ReaderUtils      = &readerUtils{}
	ReaderQual  Code = Qual(utils.RestLiCodecPackage, "Reader")
	Reader      Code = Id("reader")
	ReaderParam Code = Add(Reader).Add(ReaderQual)
)

func (d *readerUtils) ReadMap(reader Code, mapReader func(reader, key Code, def *Group)) Code {
	key := Id("key")
	return Add(reader).Dot("ReadMap").Call(Func().Params(ReaderParam, Add(key).String()).Params(Err().Error()).BlockFunc(func(def *Group) {
		mapReader(Reader, key, def)
	}))
}

func (d *readerUtils) ReadRecord(reader Code, requiredFields Code, mapReader func(reader, key Code, def *Group)) Code {
	field := Id("field")
	return Add(reader).Dot("ReadRecord").Call(requiredFields, Func().Params(ReaderParam, Add(field).String()).Params(Err().Error()).BlockFunc(func(def *Group) {
		mapReader(Reader, field, def)
	}))
}

func (d *readerUtils) ReadArray(reader Code, arrayReader func(reader Code, def *Group)) Code {
	return Add(reader).Dot("ReadArray").Call(Func().Params(ReaderParam).Params(Err().Error()).BlockFunc(func(def *Group) {
		arrayReader(Reader, def)
	}))
}

func (d *readerUtils) Skip(reader Code) *Statement {
	return Add(reader).Dot("Skip").Call()
}

func (d *readerUtils) Read(t RestliType, reader, targetAccessor Code) Code {
	switch {
	case t.Primitive != nil:
		return List(targetAccessor, Err()).Op("=").Add(reader).Dot(t.Primitive.ReaderName()).Call()
	case t.Reference != nil:
		return Err().Op("=").Add(targetAccessor).Dot(utils.UnmarshalRestLi).Call(reader)
	case t.IsMapOrArray():
		return List(targetAccessor, Err()).Op("=").Add(d.NestedUnmarshaler(t, reader))
	default:
		log.Panicf("Illegal restli type: %+v", t)
		return nil
	}
}

func (d *readerUtils) UnmarshalerFunc(t RestliType) Code {
	switch {
	case t.Primitive != nil:
		return t.Primitive.UnmarshalerFunc()
	case t.Reference != nil:
		return t.Reference.UnmarshalerFunc()
	case t.IsMapOrArray():
		return Func().
			Params(ReaderParam).
			Params(t.GoType(), Error()).
			BlockFunc(func(def *Group) {
				def.Return(d.NestedUnmarshaler(t, Reader))
			})
	default:
		log.Panicf("Illegal restli type: %+v", t)
		return nil
	}
}

func (d *readerUtils) NestedUnmarshaler(t RestliType, reader Code) Code {
	innerT, word := t.InnerMapOrArray()
	return Qual(utils.RestLiCodecPackage, "Read"+word).Call(reader, d.UnmarshalerFunc(innerT))
}
