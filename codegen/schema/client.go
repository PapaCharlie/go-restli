package schema

import (
	"fmt"
	"log"
	"strings"

	. "github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
)

const (
	ResourcePath       = "ResourcePath"
	ResourceEntityPath = "ResourceEntityPath"
)

func (r *Resource) generateClient(parentResources []*Resource) (c *CodeFile) {
	c = NewCodeFile("client", r.PackagePath(), r.Name)

	c.Code.Const().DefsFunc(func(def *Group) {
		def.Id(ExportedIdentifier(r.Name + "Path")).Op("=").Lit(r.Path)
		if e := r.getEntity(); e != nil {
			def.Id(ExportedIdentifier(r.Name + "EntityPath")).Op("=").Lit(e.Path)
		}
	}).Line().Line()

	AddWordWrappedComment(c.Code, r.Doc).Line()
	c.Code.Type().Id(Client).Struct(Qual(ProtocolPackage, "RestLiClient")).Line().Line()

	resources := make([]*Resource, 0, len(parentResources)+1)
	resources = append(resources, parentResources...)

	addResourcePathFunc(c.Code, ResourcePath, resources, r.Path)

	if e := r.getEntity(); e != nil {
		resources = append(resources, r)
		addResourcePathFunc(c.Code, ResourceEntityPath, resources, e.Path)
	}

	return c
}

func addResourcePathFunc(def *Statement, funcName string, resources []*Resource, path string) {
	addClientFunc(def, funcName).ParamsFunc(func(def *Group) {
		addEntityTypes(def, resources)
	}).Params(String(), Error()).BlockFunc(func(def *Group) {
		def.Var().Id(Path).String()

		for _, resource := range resources {
			if id := resource.getIdentifier(); id != nil {
				assignment, hasError := id.Type.RestLiURLEncodeModel(Id(id.Name))
				if hasError {
					def.List(Id(id.EncodedVariableName()), Err()).Op(":=").Add(assignment)
					IfErrReturn(def, Lit(""), Err())
				} else {
					def.Id(id.EncodedVariableName()).Op(":=").Add(assignment)
				}

				pattern := fmt.Sprintf("{%s}", id.Name)
				idx := strings.Index(path, pattern)
				if idx < 0 {
					log.Panicf("%s does not appear in %s", pattern, path)
				}
				def.Id(Path).Op("+=").Lit(path[:idx]).Op("+").Id(id.EncodedVariableName())
				path = path[idx+len(pattern):]
			}
		}
		def.Line()

		if path != "" {
			def.Id(Path).Op("+=").Lit(path)
		}

		def.Return(Id(Path), Nil())
	}).Line().Line()
}

func addEntityTypes(def *Group, resources []*Resource) {
	for _, r := range resources {
		if id := r.getIdentifier(); id != nil {
			def.Id(id.Name).Add(id.Type.GoType())
		}
	}
}

func entityParams(resources []*Resource) []Code {
	var params []Code
	for _, r := range resources {
		if id := r.getIdentifier(); id != nil {
			params = append(params, Id(id.Name))
		}
	}
	return params
}

func addClientFunc(def *Statement, funcName string) *Statement {
	return AddFuncOnReceiver(def, ClientReceiver, Client, funcName)
}

func rootResourceName(parentResources []*Resource, thisResource *Resource) string {
	var resource *Resource
	if len(parentResources) > 0 {
		resource = parentResources[0]
	} else {
		resource = thisResource
	}
	return resource.Name
}

func callFormatQueryUrl(def *Group, parentResources []*Resource, thisResource *Resource) {
	def.List(Id(Url), Err()).
		Op(":=").
		Id(ClientReceiver).Dot(FormatQueryUrl).
		Call(Lit(rootResourceName(parentResources, thisResource)), Id(Path))
}

func callDoAndDecode(def *Group) {
	def.List(Id("_"), Err()).Op("=").Id(ClientReceiver).Dot(DoAndDecode).Call(Id(Req), Op("&").Id(DoAndDecodeResult))
	IfErrReturn(def, Nil(), Err()).Line()
}
