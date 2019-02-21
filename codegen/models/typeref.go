package models

import (
	"github.com/dave/jennifer/jen"
	. "go-restli/codegen"
)

const TyperefType = "typeref"

type Typeref struct {
	NameAndDoc
	Ref *Model `json:"ref"`
}

func (t *Typeref) InnerModels() (models []*Model) {
	return []*Model{t.Ref}
}

func (t *Typeref) GoType(packagePrefix string) *jen.Statement {
	panic("typerefs cannot be directly referenced!")
}

func (t *Typeref) generateCode() (def *jen.Statement) {
	def = jen.Empty()
	AddWordWrappedComment(def, t.Doc).Line()
	def.Type().Id(t.Name).Add(t.Ref.GoType())
	return
}
