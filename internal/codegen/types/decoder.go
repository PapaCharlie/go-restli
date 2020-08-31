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

func (d *reader) ReadMap(reader func(key Code, def *Group)) Code {
	key := Id("key")
	return Add(d).Dot("ReadMap").Call(Func().Params(Add(key).String()).Params(Err().Error()).BlockFunc(func(def *Group) {
		reader(key, def)
	}))
}

func (d *reader) ReadArray(creator func(def *Group)) Code {
	return Add(d).Dot("ReadArray").Call(Func().Params().Params(Err().Error()).BlockFunc(func(def *Group) {
		creator(def)
	}))
}

func (d *reader) Skip() *Statement {
	return Add(d).Dot("Skip").Call()
}

func (d *reader) Read(t RestliType, accessor Code) Code {
	switch {
	case t.Primitive != nil:
		return List(accessor, Err()).Op("=").Add(Reader).Dot(t.Primitive.ReaderName()).Call()
	case t.Reference != nil:
		return Err().Op("=").Add(accessor).Dot(UnmarshalRestLi).Call(d)
	case t.Array != nil:
		return Err().Op("=").Add(d.ReadArray(func(def *Group) {
			item := Id(tempReaderVariableName(t))
			def.Var().Add(item).Add(t.Array.GoType())
			def.Add(d.Read(*t.Array, item))
			def.Add(utils.IfErrReturn(Err()))
			if t.Array.ShouldReference() {
				item = Op("&").Add(item)
			}
			def.Add(accessor).Op("=").Append(accessor, item)
			def.Return(Nil())
		}))
	case t.Map != nil:
		return Add(accessor).Op("=").Make(t.GoType()).Line().
			Err().Op("=").
			Add(d.ReadMap(func(key Code, def *Group) {
				value := Id(tempReaderVariableName(t))
				def.Var().Add(value).Add(t.Map.GoType())
				def.Add(d.Read(*t.Map, value))
				def.Add(utils.IfErrReturn(Err()))
				if t.Map.ShouldReference() {
					value = Op("&").Add(value)
				}
				def.Parens(accessor).Index(key).Op("=").Add(value)
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
