package restli

import (
	"github.com/dave/jennifer/jen"
	"log"
	"strings"
)

type RestliType interface {
	GoType() *jen.Statement
	UnionFieldName() string
}

type GeneratedType struct {
	ReferenceType
	Doc        string
	Definition jen.Statement
}

type ReferenceType struct {
	NamespacePrefix string
	Namespace       string
	Name            string
}

func (r *ReferenceType) UnionFieldName() string {
	return r.Name
}

func (r *ReferenceType) PackageName() string {
	return strings.Replace(NsJoin(r.NamespacePrefix, r.Namespace), NamespaceSep, "/", -1)
}

func (r *ReferenceType) GoType() *jen.Statement {
	return jen.Qual(r.PackageName(), r.Name)
}

type PrimitiveType struct {
	Type string
	Size int
}

func (t *PrimitiveType) UnionFieldName() string {
	return capitalizeFirstLetter(t.Type)
}

func (t *PrimitiveType) GoType() *jen.Statement {
	switch t.Type {
	case Int:
		return jen.Int()
	case Long:
		return jen.Int64()
	case Float:
		return jen.Float32()
	case Double:
		return jen.Float64()
	case Boolean:
		return jen.Bool()
	case String:
		return jen.String()
	case Bytes:
		return jen.Index().Byte()
	case Fixed:
		return jen.Index(jen.Lit(t.Size)).Byte()
	default:
		log.Panicln("Unknown primitive type", t.Type)
	}
	return nil
}
