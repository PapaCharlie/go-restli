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
		addEntityParams(def, resources)
	}).Params(String(), Error()).BlockFunc(func(def *Group) {
		for _, resource := range resources {
			if id := resource.getIdentifier(); id != nil {
				hasError, assignment := id.Type.restLiURLEncode(Id(id.Name))
				if hasError {
					def.List(Id(id.Name+"Str"), Err()).Op(":=").Add(assignment)
					IfErrReturn(def, Lit(""), Err())
				} else {
					def.Id(id.Name + "Str").Op(":=").Add(assignment)
				}
			}
		}
		def.Line()

		def.Var().Id("path").String()

		for _, resource := range resources {
			if id := resource.getIdentifier(); id != nil {
				pattern := fmt.Sprintf("{%s}", id.Name)
				idx := strings.Index(path, pattern)
				if idx < 0 {
					log.Panicf("%s does not appear in %s", pattern, path)
				}
				def.Id("path").Op("+=").Lit(path[:idx]).Op("+").Id(id.Name + "Str")
				path = path[idx+len(pattern):]
			}
		}

		if path != "" {
			def.Id("path").Op("+=").Lit(path)
		}

		def.Return(Id("path"), Nil())
	}).Line().Line()
}

func addEntityParams(def *Group, resources []*Resource) {
	for _, r := range resources {
		if id := r.getIdentifier(); id != nil {
			def.Id(id.Name).Add(id.Type.GoType())
		}
	}
}
