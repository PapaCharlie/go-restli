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
	Reader     = &reader{Id("reader")}
	ReaderQual = Qual(utils.RestLiCodecPackage, "Reader")
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
		return Err().Op("=").Add(targetAccessor).Dot(utils.UnmarshalRestLi).Call(reader)
	case t.Array != nil:
		return Err().Op("=").Add(d.ReadArray(reader, func(reader Code, def *Group) {
			d.ReadArrayFunc(t, reader, targetAccessor, def)
		}))
	case t.Map != nil:
		return Add(targetAccessor).Op("=").Make(t.GoType()).Line().
			Err().Op("=").
			Add(d.ReadMap(reader, func(reader, key Code, def *Group) {
				_, value := tempIteratorVariableNames(t)
				def.Var().Add(value).Add(t.Map.GoType())
				def.Add(d.Read(*t.Map, reader, value))
				def.Add(utils.IfErrReturn(Err()))
				if t.Map.ShouldReference() {
					value = Op("&").Add(value)
				}
				def.Parens(targetAccessor).Index(key).Op("=").Add(value)
				def.Return(Nil())
			}))
	case t.NativeTyperef != nil:
		return BlockFunc(func(def *Group) {
			raw := Id("raw")
			def.Var().Add(raw).Add(t.NativeTyperef.Primitive.GoType())
			def.List(raw, Err()).Op("=").Add(reader).Dot(t.NativeTyperef.Primitive.ReaderName()).Call()
			def.Add(utils.IfErrReturn(Err()))
			def.List(targetAccessor, Err()).Op("=").Add(t.NativeTyperef.Unmarshaler()).Call(raw)
		})
	default:
		log.Panicf("Illegal restli type: %+v", t)
		return nil
	}
}

func (d *reader) ReadArrayFunc(t RestliType, reader, targetAccessor Code, def *Group) {
	_, item := tempIteratorVariableNames(t)
	def.Var().Add(item).Add(t.Array.GoType())
	def.Add(d.Read(*t.Array, reader, item))
	def.Add(utils.IfErrReturn(Err()))
	if t.Array.ShouldReference() {
		item = Op("&").Add(item)
	}
	def.Add(targetAccessor).Op("=").Append(targetAccessor, item)
	def.Return(Nil())
}

func tempIteratorVariableNames(t RestliType) (Code, Code) {
	var tempName func(t RestliType) string
	tempName = func(t RestliType) string {
		if t.Array != nil {
			if t.Array.IsMapOrArray() {
				return "array" + utils.ExportedIdentifier(tempName(*t.Array))
			} else {
				return ""
			}
		} else {
			if t.Map.IsMapOrArray() {
				return "map" + utils.ExportedIdentifier(tempName(*t.Map))
			} else {
				return ""
			}
		}
	}
	prefix := tempName(t)
	var left, right string
	if t.Array != nil {
		left, right = "index", "item"
	} else {
		left, right = "key", "value"
	}
	if prefix != "" {
		left, right = utils.ExportedIdentifier(left), utils.ExportedIdentifier(right)
	}
	return Id(prefix + left), Id(prefix + right)
}
