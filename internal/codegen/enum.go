package codegen

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
)

type Enum struct {
	NamedType
	Symbols     []string
	SymbolToDoc map[string]string
}

func (e *Enum) InnerTypes() IdentifierSet {
	return nil
}

func (e *Enum) GenerateCode() (def *Statement) {
	def = Empty()
	AddWordWrappedComment(def, e.Doc).Line()
	def.Type().Id(e.Name).Int().Line()

	unknownEnum := Id("_" + e.SymbolIdentifier("unknown"))

	def.Const().DefsFunc(func(def *Group) {
		def.Add(unknownEnum).Op("=").Id(e.Name).Call(Iota())
		for _, symbol := range e.Symbols {
			def.Add(AddWordWrappedComment(Empty(), e.SymbolToDoc[symbol]))
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

	receiver := ReceiverName(e.Name)
	getter := "Get" + e.Name + "FromString"
	accessor := Op("*").Id(receiver)
	getEnumString := func(def *Group) (*Statement, *Statement) {
		val, ok := Id("val"), Id("ok")
		def.List(val, ok).Op(":=").Id(strings).Index(accessor)
		return val, ok
	}

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

	AddStringer(def, receiver, e.Name, func(def *Group) {
		val, ok := getEnumString(def)
		def.If(Op("!").Add(ok)).Block(
			Return(Lit("$UNKNOWN$")),
		)
		def.Return(val)
	}).Line().Line()

	AddFuncOnReceiver(def, receiver, e.Name, "Pointer").
		Params().
		Op("*").Id(e.Name).
		BlockFunc(func(def *Group) {
			def.Return(Id(receiver))
		}).Line().Line()

	AddUnmarshalJSON(def, receiver, e.Name, func(def *Group) {
		value := Id("value")
		def.Var().Add(value).String()
		def.Err().Op("=").Qual(EncodingJson, Unmarshal).Call(Id("data"), Op("&").Add(value))
		IfErrReturn(def)
		def.Line()

		def.Op("*").Id(receiver).Op("=").Id(values).Index(value)
		def.Return()
	}).Line().Line()

	AddRestLiEncode(def, receiver, e.Name, func(def *Group) {
		val, ok := getEnumString(def)
		def.If(Op("!").Add(ok)).Block(
			Return(Qual("fmt", "Errorf").Call(Lit(fmt.Sprintf("illegal %q: %%d", e.Identifier)), Int().Call(accessor))),
		)
		def.List(Encoder).Dot("String").Call(val)
		def.Return(Nil())
	}).Line().Line()
	AddRestLiDecode(def, receiver, e.Name, func(def *Group) {
		def.Var().Id("value").String()
		def.Err().Op("=").Id(Codec).Dot("DecodeString").Call(Id("data"), Op("&").Id("value"))
		IfErrReturn(def, Err())
		def.Line()

		def.Op("*").Id(receiver).Op("=").Id(values).Index(Id("value"))
		def.Return(Nil())
	}).Line().Line()

	return def
}

func (e *Enum) SymbolIdentifier(symbol string) string {
	return ExportedIdentifier(e.Name + "_" + symbol)
}
