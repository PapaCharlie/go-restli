package schema

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	. "go-restli/codegen"
	"log"
	"strings"
)

func generateAllActionStructs(parentResources []*Resource, thisResource *Resource) (code []*CodeFile) {
	if thisResource.Simple != nil {
		code = append(code, thisResource.Simple.generateActionParamStructs(parentResources, thisResource)...)
		code = append(code, thisResource.Simple.Entity.generateActionParamStructs(parentResources, thisResource)...)
		newParentResources := make([]*Resource, len(parentResources)+1)
		copy(parentResources, newParentResources)
		newParentResources = append(newParentResources, thisResource)
		for _, r := range thisResource.Simple.Entity.Subresources {
			code = append(code, generateAllActionStructs(parentResources, &r)...)
		}
		return
	}

	if thisResource.Collection != nil {
		code = append(code, thisResource.Collection.generateActionParamStructs(parentResources, thisResource)...)
		code = append(code, thisResource.Collection.Entity.generateActionParamStructs(parentResources, thisResource)...)
		newParentResources := make([]*Resource, len(parentResources)+1)
		copy(parentResources, newParentResources)
		newParentResources = append(newParentResources, thisResource)
		for _, r := range thisResource.Collection.Entity.Subresources {
			code = append(code, generateAllActionStructs(parentResources, &r)...)
		}
		return
	}

	if thisResource.Association != nil {
		code = append(code, thisResource.Association.generateActionParamStructs(parentResources, thisResource)...)
		code = append(code, thisResource.Association.Entity.generateActionParamStructs(parentResources, thisResource)...)
		newParentResources := make([]*Resource, len(parentResources)+1)
		copy(parentResources, newParentResources)
		newParentResources = append(newParentResources, thisResource)
		for _, r := range thisResource.Association.Entity.Subresources {
			code = append(code, generateAllActionStructs(parentResources, &r)...)
		}
		return
	}

	if thisResource.ActionsSet != nil {
		code = append(code, thisResource.ActionsSet.generateActionParamStructs(parentResources, thisResource)...)
		return
	}

	log.Panicln(thisResource, "does not define any resources")
	return
}

func (h *HasActions) generateActionParamStructs(parentResources []*Resource, thisResource *Resource) (code []*CodeFile) {
	for _, a := range h.Actions {
		code = append(code, a.generateActionParamStructs(parentResources, thisResource, false))
	}
	return code
}

func (e *Entity) generateActionParamStructs(parentResources []*Resource, thisResource *Resource) (code []*CodeFile) {
	for _, a := range e.Actions {
		code = append(code, a.generateActionParamStructs(parentResources, thisResource, true))
	}
	return code
}

func (a *Action) generateActionParamStructs(parentResources []*Resource, thisResource *Resource, isOnEntity bool) (c *CodeFile) {
	c = NewCodeFile(a.ActionName, thisResource.PackagePath(), thisResource.Name)

	c.Code.Const().Id(ExportedIdentifier(a.ActionName + "Action")).Op("=").Lit(a.ActionName).Line()
	c.Code.Add(a.GenerateCode())

	var queryPath string
	if isOnEntity {
		queryPath = thisResource.getEntity().Path
	} else {
		queryPath = thisResource.Path
	}

	c.Code.Func().Params(Id(ClientReceiver).Op("*").Id(Client)).Id(ExportedIdentifier(a.ActionName) + "Action")
	c.Code.ParamsFunc(func(def *Group) {
		for _, r := range parentResources {
			if id := r.getIdentifier(); id != nil {
				def.Id(id.Name).Add(id.Type.GoType())
				queryPath = strings.Replace(queryPath, fmt.Sprintf("{%s}", id.Name), "%s", 1)
			}
		}
		if id := thisResource.getIdentifier(); isOnEntity && id != nil {
			def.Id(id.Name).Add(id.Type.GoType())
			queryPath = strings.Replace(queryPath, fmt.Sprintf("{%s}", id.Name), "%s", 1)
		}
		def.Id("params").Id(a.StructName)
	})

	returns := a.Returns != nil

	c.Code.ParamsFunc(func(def *Group) {
		if returns {
			def.Id(ActionResult).Add(a.Returns.GoType())
		}
		def.Err().Error()
	})

	c.Code.BlockFunc(func(def *Group) {
		def.Id(Url).Op(":=").Id(ClientReceiver).Dot(HostnameClientField).Op("+").Qual("fmt", "Sprintf").
			CallFunc(func(def *Group) {
				def.Lit(queryPath + "?action=" + a.ActionName)
				for _, r := range parentResources {
					if id := r.getIdentifier(); id != nil {
						def.Id(id.Name)
					}
				}
				if isOnEntity {
					def.Id(thisResource.getIdentifier().Name)
				}
			})
		def.List(Id(Req), Err()).Op(":=").Qual(GetRestLiProtocolPackage(), "RestliPost").Call(Id("url"), Lit(""), Id("params"))
		IfErrReturn(def).Line()

		var resDef *Statement
		if returns {
			resDef = def.List(Id(Res), Err()).Op(":=")
		} else {
			resDef = def.List(Id("_"), Err()).Op("=")
		}
		resDef.Qual(GetRestLiProtocolPackage(), "RestliDo").Call(Id(ClientReceiver).Dot(Client), Id(Req))
		IfErrReturn(def).Line()

		if returns {
			def.Id("result").Op(":=").Struct(Id("Value").Add(a.Returns.GoType())).Block()
			def.Err().Op("=").Qual(EncodingJson, "NewDecoder").Call(Id(Res).Dot("Body")).Dot("Decode").Call(Op("&").Id("result"))
			IfErrReturn(def).Line()
			def.Id(ActionResult).Op("=").Id("result").Dot("Value")
		}
		def.Return()
	})

	return
}
