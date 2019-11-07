package schema

import (
	"log"

	. "github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/protocol"
	. "github.com/dave/jennifer/jen"
)

type MethodGenerator func(m Method, parentResources []*Resource, thisResource *Resource) *Statement

// https://linkedin.github.io/rest.li/user_guide/restli_server#resource-methods
func (m *Method) generate(parentResources []*Resource, thisResource *Resource) *Statement {
	switch m.Method {
	case protocol.Method_get:
		return m.generateGet(parentResources, thisResource)
	default:
		log.Printf("Warning: %s method is not currently implemented", m.Name)
		return nil
	}
}

func (m *Method) generateGet(parentResources []*Resource, thisResource *Resource) *Statement {
	var resources []*Resource
	resources = append(resources, parentResources...)
	resources = append(resources, thisResource)

	def := Empty()

	AddWordWrappedComment(def, m.Doc).Line()
	addClientFunc(def, m.Name)
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
		def.List(Err()).Op("=").Id(ClientReceiver).Dot("DoAndDecode").Call(Id(Req), Id("result"))
		IfErrReturn(def, Nil(), Err()).Line()
		def.Return(Id("result"), Err())
	})

	return def
}
