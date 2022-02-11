package types

import (
	"log"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type reader struct {
	Code
}

var (
	Reader      = &reader{Id("reader")}
	ReaderQual  = Qual(utils.RestLiCodecPackage, "Reader")
	ReaderParam = Add(Reader).Add(ReaderQual)
)

func (d *reader) ReadMap(reader Code, mapReader func(reader, key Code, def *Group)) Code {
	key := Id("key")
	return Add(reader).Dot("ReadMap").Call(Func().Params(Add(d).Add(ReaderQual), Add(key).String()).Params(Err().Error()).BlockFunc(func(def *Group) {
		mapReader(d, key, def)
	}))
}

func (d *reader) ReadRecord(reader Code, requiredFields Code, mapReader func(reader, key Code, def *Group)) Code {
	field := Id("field")
	return Add(reader).Dot("ReadRecord").Call(requiredFields, Func().Params(Add(d).Add(ReaderQual), Add(field).String()).Params(Err().Error()).BlockFunc(func(def *Group) {
		mapReader(d, field, def)
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
	case t.Reference != nil:
		return t.Reference.UnmarshalerFunc()
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
	var readFunc Code
	var innerT RestliType
	if t.Map != nil {
		readFunc = utils.ReadMap
		innerT = *t.Map
	} else {
		readFunc = utils.ReadArray
		innerT = *t.Array
	}

	return Add(readFunc).Call(reader, d.UnmarshalerFunc(innerT))
}
