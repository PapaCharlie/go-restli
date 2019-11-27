package schema

import (
	"log"

	. "github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/protocol"
	. "github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

type MethodGenerator func(m Method, parentResources []*Resource, thisResource *Resource) *Statement

// https://linkedin.github.io/rest.li/user_guide/restli_server#resource-methods
func (m *Method) generate(parentResources []*Resource, thisResource *Resource) *Statement {
	switch m.Method {
	case protocol.Method_get:
		return m.generateGet(parentResources, thisResource)
	case protocol.Method_update:
		return m.generateUpdate(parentResources, thisResource)
	case protocol.Method_delete:
		return m.generateDelete(parentResources, thisResource)
	default:
		log.Printf("Warning: %s method is not currently implemented", m.Name)
		return nil
	}
}

func (m *Method) addMethodFunc(parentResources []*Resource, thisResource *Resource) *Statement {
	def := Empty()

	AddWordWrappedComment(def, m.Doc).Line()
	addClientFunc(def, strcase.ToCamel(m.Method.String()))
	return def
}

func (m *Method) generateGet(parentResources []*Resource, thisResource *Resource) *Statement {
	def := m.addMethodFunc(parentResources, thisResource)

	resources := append([]*Resource(nil), parentResources...)
	resources = append(resources, thisResource)
	def.ParamsFunc(func(def *Group) {
		addEntityTypes(def, resources)
	})

	def.Params(Op("*").Add(thisResource.Schema.Model.GoType()), Error())

	def.BlockFunc(func(def *Group) {
		def.List(Id("path"), Err()).Op(":=").Id(ClientReceiver).Dot(ResourceEntityPath).Call(entityParams(resources)...)
		IfErrReturn(def, Nil(), Err()).Line()

		def.List(Id(Url), Err()).Op(":=").Id(ClientReceiver).Dot(FormatQueryUrl).Call(Id("path"))
		IfErrReturn(def, Nil(), Err()).Line()
		def.List(Id(Req), Err()).Op(":=").Id(ClientReceiver).Dot("GetRequest").Call(Id("url"), RestLiMethod(protocol.Method_get))
		IfErrReturn(def, Nil(), Err()).Line()

		def.Id("result").Op(":=").New(thisResource.Schema.Model.GoType())
		def.List(Id("_"), Err()).Op("=").Id(ClientReceiver).Dot("DoAndDecode").Call(Id(Req), Id("result"))
		IfErrReturn(def, Nil(), Err()).Line()
		def.Return(Id("result"), Err())
	})

	return def
}

func (m *Method) generateUpdate(parentResources []*Resource, thisResource *Resource) *Statement {
	def := m.addMethodFunc(parentResources, thisResource)

	resources := append([]*Resource(nil), parentResources...)
	resources = append(resources, thisResource)

	paramName := "o"
	def.ParamsFunc(func(def *Group) {
		addEntityTypes(def, resources)
		def.Id(paramName).Op("*").Add(thisResource.Schema.GoType())
	})

	def.Params(Error())

	def.BlockFunc(func(def *Group) {
		def.List(Id("path"), Err()).Op(":=").Id(ClientReceiver).Dot(ResourceEntityPath).Call(entityParams(resources)...)
		IfErrReturn(def, Err()).Line()

		def.List(Id(Url), Err()).Op(":=").Id(ClientReceiver).Dot(FormatQueryUrl).Call(Id("path"))
		IfErrReturn(def, Err()).Line()
		def.List(Id(Req), Err()).Op(":=").Id(ClientReceiver).Dot("JsonPutRequest").Call(Id("url"), RestLiMethod(protocol.Method_update), Id(paramName))
		IfErrReturn(def, Err()).Line()

		def.List(Id("res"), Err()).Op(":=").Id(ClientReceiver).Dot("DoAndIgnore").Call(Id(Req))
		IfErrReturn(def, Err()).Line()

		def.If(Id("res").Dot("StatusCode").Op("/").Lit(100).Op("!=").Lit(2)).BlockFunc(func(def *Group) {
			def.Return(Qual("fmt", "Errorf").Call(Lit("Invalid response code from %s: %d"), Id(Url), Id("res").Dot("StatusCode")))
		})
		def.Return(Nil())
	})

	return def
}

func (m *Method) generateDelete(parentResources []*Resource, thisResource *Resource) *Statement {
	def := m.addMethodFunc(parentResources, thisResource)

	resources := append([]*Resource(nil), parentResources...)
	resources = append(resources, thisResource)

	def.ParamsFunc(func(def *Group) {
		addEntityTypes(def, resources)
	})

	def.Params(Error())

	def.BlockFunc(func(def *Group) {
		def.List(Id("path"), Err()).Op(":=").Id(ClientReceiver).Dot(ResourceEntityPath).Call(entityParams(resources)...)
		IfErrReturn(def, Err()).Line()

		def.List(Id(Url), Err()).Op(":=").Id(ClientReceiver).Dot(FormatQueryUrl).Call(Id("path"))
		IfErrReturn(def, Err()).Line()
		def.List(Id(Req), Err()).Op(":=").Id(ClientReceiver).Dot("DeleteRequest").Call(Id("url"), RestLiMethod(protocol.Method_update))
		IfErrReturn(def, Err()).Line()

		def.List(Id("res"), Err()).Op(":=").Id(ClientReceiver).Dot("DoAndIgnore").Call(Id(Req))
		IfErrReturn(def, Err()).Line()

		def.If(Id("res").Dot("StatusCode").Op("/").Lit(100).Op("!=").Lit(2)).BlockFunc(func(def *Group) {
			def.Return(Qual("fmt", "Errorf").Call(Lit("Invalid response code from %s: %d"), Id(Url), Id("res").Dot("StatusCode")))
		})
		def.Return(Nil())
	})

	return def
}
