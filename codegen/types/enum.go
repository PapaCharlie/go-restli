package types

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type Enum struct {
	NamedType
	Symbols     []string
	SymbolToDoc map[string]string
}

func (e *Enum) InnerTypes() utils.IdentifierSet {
	return nil
}

func (e *Enum) GenerateCode() Code {
	def := Empty()
	utils.AddWordWrappedComment(def, e.Doc).Line()
	def.Type().Id(e.Name).Int().Line()

	unknownEnum := Id("_" + e.SymbolIdentifier("unknown"))

	def.Const().DefsFunc(func(def *Group) {
		def.Add(unknownEnum).Op("=").Id(e.Name).Call(Iota())
		for _, symbol := range e.Symbols {
			def.Add(utils.AddWordWrappedComment(Empty(), e.SymbolToDoc[symbol]))
			def.Id(e.SymbolIdentifier(symbol))
		}
	}).Line()

	values := "_" + e.Name + "_values"
	def.Var().Id(values).Op("=").Map(String()).Id(e.Name).Values(DictFunc(func(dict Dict) {
		for _, s := range e.Symbols {
			dict[Lit(s)] = Id(e.SymbolIdentifier(s))
		}
	})).Line()

	strings := "_" + e.Name + "_strings"
	def.Var().Id(strings).Op("=").Map(Id(e.Name)).String().Values(DictFunc(func(dict Dict) {
		for _, s := range e.Symbols {
			dict[Id(e.SymbolIdentifier(s))] = Lit(s)
		}
	})).Line().Line()

	receiver := utils.ReceiverName(e.Name)
	getter := "Get" + e.Name + "FromString"
	accessor := Op("*").Id(receiver)
	getEnumString := func(def *Group) (*Statement, *Statement) {
		val, ok := Id("val"), Id("ok")
		def.List(val, ok).Op(":=").Id(strings).Index(accessor)
		return val, ok
	}

	AddEquals(def, receiver, e.Name, func(other Code, def *Group) {
		def.Return(Op("*").Id(receiver).Op("== *").Add(other))
	})

	AddComputeHash(def, receiver, e.Name, func(h Code, def *Group) {
		def.Add(h).Op("=").Add(utils.Hash).Call(Op("*").Id(receiver))
		def.Return(h)
	})

	def.Func().Id("All" + e.Name + "Values").Params().Index().Id(e.Name).BlockFunc(func(def *Group) {
		def.Return(Index().Id(e.Name).ValuesFunc(func(def *Group) {
			for _, s := range e.Symbols {
				def.Id(e.SymbolIdentifier(s))
			}
		}))
	}).Line().Line()

	def.Func().Id(getter).Params(Id("val").String()).Params(Id(receiver).Id(e.Name), Err().Error()).
		BlockFunc(func(def *Group) {
			def.List(Id(receiver), Id("ok")).Op(":=").Id(values).Index(Id("val"))
			def.If(Op("!").Id("ok")).BlockFunc(func(def *Group) {
				def.Err().Op("=").Qual("fmt", "Errorf").Call(Lit(fmt.Sprintf("unknown %q: %%s", e.Identifier)), Id("val"))
			})
			def.Return()
		}).Line().Line()

	utils.AddStringer(def, receiver, e.Name, func(def *Group) {
		val, ok := getEnumString(def)
		def.If(Op("!").Add(ok)).Block(
			Return(Lit("$UNKNOWN$")),
		)
		def.Return(val)
	}).Line().Line()

	def.Func().
		Params(Id(receiver).Id(e.Name)).
		Id("Pointer").Params().
		Op("*").Id(e.Name).
		BlockFunc(func(def *Group) {
			def.Return(Op("&").Id(receiver))
		}).Line().Line()

	AddMarshalRestLi(def, receiver, e.Name, func(def *Group) {
		val, ok := getEnumString(def)
		def.If(Op("!").Add(ok)).Block(
			Return(Qual("fmt", "Errorf").Call(Lit(fmt.Sprintf("illegal %q: %%d", e.Identifier)), Int().Call(accessor))),
		)
		def.List(Writer).Dot("WriteString").Call(val)
		def.Return(Nil())
	})
	AddUnmarshalRestli(def, receiver, e.Name, func(def *Group) {
		value := Id("value")
		def.Var().Add(value).String()
		def.Add(Reader.Read(RestliType{Primitive: &StringPrimitive}, Reader, value))
		def.Add(utils.IfErrReturn(Err()))
		def.Line()

		def.Op("*").Id(receiver).Op("=").Id(values).Index(Id("value"))
		def.Return(Nil())
	})

	return def
}

func (e *Enum) SymbolIdentifier(symbol string) string {
	return utils.ExportedIdentifier(e.Name + "_" + symbol)
}

func (e *Enum) zeroValueLit() *Statement {
	return e.Qual().Call(Lit(0))
}

func (e *Enum) isValidSymbol(v string) bool {
	for _, symbol := range e.Symbols {
		if symbol == v {
			return true
		}
	}
	return false
}
