package schema

import (
	"fmt"
	. "github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
	"log"
	"strings"
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
		def.Var().Id("path").String()

		for _, resource := range resources {
			if id := resource.getIdentifier(); id != nil {
				hasError, assignment := id.Type.RestLiURLEncode(Id(id.Name))
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
				def.Id("path").Op("+=").Lit(path[:idx]).Op("+").Id(id.EncodedVariableName())
				path = path[idx+len(pattern):]
			}
		}
		def.Line()

		if path != "" {
			def.Id("path").Op("+=").Lit(path)
		}

		def.Return(Id("path"), Nil())
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
