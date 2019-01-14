package restli

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"log"
)

type RestliType interface {
	GoType() *jen.Statement
	UnionFieldName() string
}

type GeneratedType struct {
	NamespacePrefix string
	Namespace       string
	Name            string
	Doc             string
	Definition      jen.Statement
}

func (g GeneratedType) GoType() *jen.Statement {
	return getQual(NsJoin(g.NamespacePrefix, g.Namespace), g.Name)
}

type TyperefType struct {
	NamespacePrefix string
	Namespace       string
	Name            string
	Doc             string
}

func (t *TyperefType) GoType() *jen.Statement {
	return getQual(NsJoin(t.NamespacePrefix, t.Namespace), t.Name)
}

func (t *TyperefType) UnionFieldName() string {
	return t.Name
}

type ReferenceType struct {
	NamespacePrefix string
	Namespace       string
	Name            string
}

func (r *ReferenceType) UnionFieldName() string {
	return r.Name
}

func (r *ReferenceType) GoType() *jen.Statement {
	fmt.Println(NsJoin(r.NamespacePrefix, r.Namespace), r.Name)
	return getQual(NsJoin(r.NamespacePrefix, r.Namespace), r.Name)
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
		return jen.Index(jen.Lit(t.Size).Byte())
	default:
		log.Panicln("Unknown primitive type", t.Type)
	}
	return nil
}
