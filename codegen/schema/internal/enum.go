package internal

import (
	"encoding/json"
	"fmt"

	. "github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
)

const EnumModelTypeName = "enum"

type EnumModel struct {
	Identifier
	Doc string

	Symbols    []string
	SymbolDocs map[string]string
}

func (e *EnumModel) CopyWithAlias(alias string) ComplexType {
	eCopy := *e
	eCopy.Name = alias
	return &eCopy
}

func (e *EnumModel) UnmarshalJSON(data []byte) error {
	t := &struct {
		Identifier
		typeField
		docField
		Symbols    []string          `json:"symbols"`
		SymbolDocs map[string]string `json:"symbolDocs"`
	}{}
	t.Namespace = currentNamespace // default to the current namespace if none is specified
	if err := json.Unmarshal(data, t); err != nil {
		return err
	}
	if t.Type != EnumModelTypeName {
		return &WrongTypeError{Expected: EnumModelTypeName, Actual: string(data)}
	}
	e.Identifier = t.Identifier
	e.Doc = t.Doc
	e.Symbols = t.Symbols
	e.SymbolDocs = t.SymbolDocs
	return nil
}

func (e *EnumModel) GenerateCode() (def *Statement) {
	def = Empty()
	AddWordWrappedComment(def, e.Doc).Line()
	def.Type().Id(e.Name).Int().Line()

	def.Const().DefsFunc(func(def *Group) {
		def.Id("_" + e.SymbolIdentifier("unknown")).Op("=").Id(e.Name).Call(Iota())
		for _, symbol := range e.Symbols {
			def.Add(AddWordWrappedComment(Empty(), e.SymbolDocs[symbol]))
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

	def.Func().Id("All" + e.Name + "Values").Params().Index().Id(e.Name).BlockFunc(func(def *Group) {
		def.Return(Index().Id(e.Name).ValuesFunc(func(def *Group) {
			for _, s := range e.Symbols {
				def.Id(e.SymbolIdentifier(s))
			}
		}))
	}).Line().Line()

	def.Func().Id(getter).Params(Id("val").String()).Params(Id(receiver).Id(e.Name), Err().Error())
	def.BlockFunc(func(def *Group) {
		def.List(Id(receiver), Id("ok")).Op(":=").Id(values).Index(Id("val"))
		def.If(Op("!").Id("ok")).BlockFunc(func(def *Group) {
			def.Err().Op("=").Qual("fmt", "Errorf").Call(Lit(fmt.Sprintf("unknown %s: %%s", e.Name)), Id("val"))
		})
		def.Return()
	}).Line().Line()

	AddStringer(def, receiver, e.Name, func(def *Group) {
		def.Return(Id(strings).Index(Op("*").Id(receiver)))
	}).Line().Line()

	AddMarshalJSON(def, receiver, e.Name, func(def *Group) {
		def.Id("val").Op(":=").Id(receiver).Dot("String").Call()
		def.If(Id("val").Op("==").Lit("")).BlockFunc(func(def *Group) {
			def.Return(Nil(), Qual("fmt", "Errorf").Call(Lit(fmt.Sprintf("illegal %s: %%s", e.Name)), Id(receiver)))
		})
		def.Return(Index().Byte().Call(Lit(`"`).Op("+").Id("val").Op("+").Lit(`"`)), Nil())
	}).Line().Line()

	AddUnmarshalJSON(def, receiver, e.Name, func(def *Group) {
		def.Var().Id("str").String()
		def.Err().Op("=").Qual(EncodingJson, Unmarshal).Call(Id("data"), Op("&").Id("str"))
		IfErrReturn(def)
		def.Line()

		def.List(Id("val"), Err()).Op(":=").Id(getter).Call(Id("str"))
		IfErrReturn(def)
		def.Op("*").Id(receiver).Op("=").Id("val")
		def.Return()
	}).Line().Line()

	AddRestLiEncode(def, receiver, e.Name, func(def *Group) {
		def.Id("data").Op("=").Id(receiver).Dot("String").Call()
		def.Return()
	}).Line().Line()
	AddRestLiDecode(def, receiver, e.Name, func(def *Group) {
		def.List(Op("*").Id(receiver), Err()).Op("=").Id(getter).Call(Id("data"))
		def.Return()
	}).Line().Line()

	return def
}

func (e *EnumModel) SymbolIdentifier(symbol string) string {
	return ExportedIdentifier(e.Name + "_" + symbol)
}
