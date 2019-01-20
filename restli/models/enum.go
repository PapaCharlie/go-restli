package models

import (
	"encoding/json"
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
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

func (e *Enum) generateCode(destinationPackage string) (def *jen.Statement) {
	def = jen.Empty()
	addWordWrappedComment(def, e.Doc)
	def.Type().Id(e.Name).String().Line()

	var values []jen.Code
	for symbol, doc := range e.Symbols {
		values = append(values, jen.Id(symbol).Op("=").Id(e.Name).Call(jen.Lit(symbol)).Comment(doc))
	}

	def.Const().Defs(values...)
	return
}
