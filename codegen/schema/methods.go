package schema

import (
	"log"

	. "github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/protocol"
	. "github.com/dave/jennifer/jen"
)

var RestliMethodToHttpMethod = map[string]string{
	"get":            "GET",
	"create":         "POST",
	"delete":         "DELETE",
	"update":         "PUT",
	"partial_update": "POST",

	"batch_get":            "GET",
	"batch_create":         "POST",
	"batch_delete":         "DELETE",
	"batch_update":         "PUT",
	"batch_partial_update": "POST",

	"get_all": "GET",
}

type MethodGenerator func(m Method, parentResources []*Resource, thisResource *Resource) *Statement

// https://github.com/linkedin/rest.li/wiki/Rest.li-User-Guide#resource-methods
func (m *Method) generate(parentResources []*Resource, thisResource *Resource) *Statement {
	switch m.Method {
	case protocol.MethodGet:
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

	def.Params(Op("*").Add(thisResource.Schema.GoType()), Error())

	def.BlockFunc(func(def *Group) {
		def.List(Id("path"), Err()).Op(":=").Id(ClientReceiver).Dot(ResourceEntityPath).Call(entityParams(resources)...)
		IfErrReturn(def, Nil(), Err()).Line()

		def.List(Id(Url), Err()).Op(":=").Id(ClientReceiver).Dot(FormatQueryUrl).Call(Id("path"))
		IfErrReturn(def, Nil(), Err()).Line()
		def.List(Id(Req), Err()).Op(":=").Id(ClientReceiver).Dot("GetRequest").Call(Id("url"), RestLiMethod(protocol.MethodGet))
		IfErrReturn(def, Nil(), Err()).Line()

		def.Id("result").Op(":=").New(thisResource.Schema.GoType())
		def.List(Err()).Op("=").Id(ClientReceiver).Dot("DoAndDecode").Call(Id(Req), Id("result"))
		IfErrReturn(def, Nil(), Err()).Line()
		def.Return(Id("result"), Err())
	})

	return def
}
