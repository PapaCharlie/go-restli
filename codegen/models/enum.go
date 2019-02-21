package models

import (
	"encoding/json"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
	. "go-restli/codegen"
)

const EnumType = "enum"

type Enum struct {
	NameAndDoc
	Symbols map[string]string
}

func (e *Enum) UnmarshalJSON(data []byte) error {
	enum := &struct {
		NameAndDoc
		Symbols    []string          `json:"symbols"`
		SymbolDocs map[string]string `json:"symbolDocs"`
	}{}

	if err := json.Unmarshal(data, enum); err != nil {
		return errors.Wrap(err, "Could not unmarshal enum")
	}

	e.NameAndDoc = enum.NameAndDoc
	if e.Symbols == nil {
		e.Symbols = make(map[string]string)
	}

	for _, s := range enum.Symbols {
		e.Symbols[s] = enum.SymbolDocs[s]
	}
	return nil
}

func (e *Enum) generateCode() (def *jen.Statement) {
	def = jen.Empty()
	AddWordWrappedComment(def, e.Doc).Line()
	def.Type().Id(e.Name).Int().Line()

	var symbols []string

	def.Const().DefsFunc(func(def *jen.Group) {
		def.Id("_" + e.SymbolIdentifier("unknown")).Op("=").Id(e.Name).Call(jen.Iota())
		for symbol, symbolDoc := range e.Symbols {
			symbols = append(symbols, symbol)
			def.Add(AddWordWrappedComment(jen.Empty(), symbolDoc))
			def.Id(e.SymbolIdentifier(symbol))
		}
	}).Line()

	values := "_" + e.Name + "_values"
	def.Var().Id(values).Op("=").Map(jen.String()).Id(e.Name).Values(jen.DictFunc(func(dict jen.Dict) {
		for _, s := range symbols {
			dict[jen.Lit(s)] = jen.Id(e.SymbolIdentifier(s))
		}
	})).Line()

	strings := "_" + e.Name + "_strings"
	def.Var().Id(strings).Op("=").Map(jen.Id(e.Name)).String().Values(jen.DictFunc(func(dict jen.Dict) {
		for _, s := range symbols {
			dict[jen.Id(e.SymbolIdentifier(s))] = jen.Lit(s)
		}
	})).Line().Line()

	receiver := PrivateIdentifier(e.Name[:1])
	getter := "Get" + e.Name + "FromString"

	def.Func().Id(getter).Params(jen.Id("val").String()).Params(jen.Id(receiver).Id(e.Name), jen.Err().Error())
	def.BlockFunc(func(def *jen.Group) {
		def.List(jen.Id(receiver), jen.Id("ok")).Op(":=").Id(values).Index(jen.Id("val"))
		def.If(jen.Op("!").Id("ok")).BlockFunc(func(def *jen.Group) {
			jen.Err().Op("=").Qual("fmt", "Errorf").Call(jen.Lit(fmt.Sprintf("illegal %s: %%s", e.Name)), jen.Id(receiver))
		})
		def.Return()
	}).Line().Line()

	AddStringer(def, receiver, e.Name, func(def *jen.Group) {
		def.Return(jen.Id(strings).Index(jen.Op("*").Id(receiver)))
	}).Line().Line()

	AddMarshalJSON(def, receiver, e.Name, func(def *jen.Group) {
		def.Id("val").Op(":=").Id(receiver).Dot("String").Call()
		def.If(jen.Id("val").Op("==").Lit("")).BlockFunc(func(def *jen.Group) {
			def.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit(fmt.Sprintf("illegal %s: %%s", e.Name)), jen.Id(receiver)))
		})
		def.Return(jen.Index().Byte().Call(jen.Lit(`"`).Op("+").Id("val").Op("+").Lit(`"`)), jen.Nil())
	}).Line().Line()

	AddUnmarshalJSON(def, receiver, e.Name, func(def *jen.Group) {
		def.Var().Id("str").String()
		def.Err().Op("=").Qual(EncodingJson, Unmarshal).Call(jen.Id("data"), jen.Op("&").Id("str"))
		IfErrReturn(def)
		def.Line()

		def.List(jen.Id("val"), jen.Err()).Op(":=").Id(getter).Call(jen.Id("str"))
		IfErrReturn(def)
		def.Op("*").Id(receiver).Op("=").Id("val")
		def.Return()
	}).Line().Line()

	return
}

func (e *Enum) SymbolIdentifier(symbol string) string {
	return ExportedIdentifier(e.Name + "_" + symbol)
}
