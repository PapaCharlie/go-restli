package models

import (
	"encoding/json"
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

func (e *Enum) generateCode(packagePrefix string) (def *jen.Statement) {
	def = jen.Empty()
	AddWordWrappedComment(def, e.Doc).Line()
	def.Type().Id(e.Name).String().Line()

	var values []jen.Code
	for symbol, doc := range e.Symbols {
		def := jen.Id(symbol).Op("=")
		def.Id(e.Name).Call(jen.Lit(symbol))
		AddWordWrappedComment(def, doc)
		values = append(values, def)
	}

	def.Const().Defs(values...)
	return
}
