package types

import (
	"log"

	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type reader struct {
	Code
}

var (
	Reader     = &reader{Id("reader")}
	ReaderQual = Qual(RestLiCodecPackage, "Reader")
)

func (d *reader) ReadMap(reader Code, mapReader func(reader, key Code, def *Group)) Code {
	key := Id("key")
	return Add(reader).Dot("ReadMap").Call(Func().Params(Add(d).Add(ReaderQual), Add(key).String()).Params(Err().Error()).BlockFunc(func(def *Group) {
		mapReader(d, key, def)
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
		return Err().Op("=").Add(targetAccessor).Dot(UnmarshalRestLi).Call(reader)
	case t.Array != nil:
		return Err().Op("=").Add(d.ReadArray(reader, func(reader Code, def *Group) {
			item := Id(tempReaderVariableName(t))
			def.Var().Add(item).Add(t.Array.GoType())
			def.Add(d.Read(*t.Array, reader, item))
			def.Add(utils.IfErrReturn(Err()))
			if t.Array.ShouldReference() {
				item = Op("&").Add(item)
			}
			def.Add(targetAccessor).Op("=").Append(targetAccessor, item)
			def.Return(Nil())
		}))
	case t.Map != nil:
		return Add(targetAccessor).Op("=").Make(t.GoType()).Line().
			Err().Op("=").
			Add(d.ReadMap(reader, func(reader, key Code, def *Group) {
				value := Id(tempReaderVariableName(t))
				def.Var().Add(value).Add(t.Map.GoType())
				def.Add(d.Read(*t.Map, reader, value))
				def.Add(utils.IfErrReturn(Err()))
				if t.Map.ShouldReference() {
					value = Op("&").Add(value)
				}
				def.Parens(targetAccessor).Index(key).Op("=").Add(value)
				def.Return(Nil())
			}))
	default:
		log.Panicf("Illegal restli type: %+v", t)
		return nil
	}
}

func tempReaderVariableName(t RestliType) string {
	if t.Array != nil {
		if t.Array.IsMapOrArray() {
			return "array" + utils.ExportedIdentifier(tempReaderVariableName(*t.Array))
		} else {
			return "item"
		}
	} else {
		if t.Map.IsMapOrArray() {
			return "map" + utils.ExportedIdentifier(tempReaderVariableName(*t.Map))
		} else {
			return "value"
		}
	}
}
