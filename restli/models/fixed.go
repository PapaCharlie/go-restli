package models

import (
	"github.com/dave/jennifer/jen"
)

const FixedType = "fixed"

type Fixed struct {
	NameAndDoc
	Size int
}

func (f *Fixed) goType() *jen.Statement {
	return jen.Index(jen.Lit(f.Size)).Byte()
}

func (f *Fixed) GoType() *jen.Statement {
	return f.goType()
}

func (f *Fixed) GenerateCode() (def *jen.Statement) {
	def = jen.Empty()
	addWordWrappedComment(def, f.Doc)
	def.Type().Id(f.Name).Add(f.goType())
	return
}
