package models

import "github.com/dave/jennifer/jen"

const TyperefType = "typeref"

type Typeref struct {
	NameAndDoc
	Ref *Model `json:"ref"`
}

func (t *Typeref) InnerModels() (models []*Model) {
	return []*Model{t.Ref}
}

func (t *Typeref) GoType(destinationPackage string) *jen.Statement {
	panic("typerefs cannot be directly referenced!")
}

func (t *Typeref) generateCode(destinationPackage string) (def *jen.Statement) {
	def = jen.Empty()
	addWordWrappedComment(def, t.Doc)
	def.Type().Id(t.Name).Add(t.Ref.GoType(destinationPackage))
	return
}
