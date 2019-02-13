package schema

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	. "go-restli/codegen"
	"log"
	"strings"
)

func generateAllActionStructs(packagePrefix string, parentResources []*Resource, thisResource *Resource) (defs []*Statement) {
	fullName := prefixNameWithParentResources(thisResource.Name, parentResources, thisResource)
	defs = append(defs, Const().Id(fullName + "Path").Op("=").Lit(thisResource.Path))

	if thisResource.Simple != nil {
		defs = append(defs, thisResource.Simple.generateActionParamStructs(packagePrefix, parentResources, thisResource)...)
		defs = append(defs, thisResource.Simple.Entity.generateActionParamStructs(packagePrefix, parentResources, thisResource)...)
		newParentResources := make([]*Resource, len(parentResources)+1)
		copy(parentResources, newParentResources)
		newParentResources = append(newParentResources, thisResource)
		for _, r := range thisResource.Simple.Entity.Subresources {
			defs = append(defs, generateAllActionStructs(packagePrefix, parentResources, &r)...)
		}
		return
	}

	if thisResource.Collection != nil {
		defs = append(defs, thisResource.Collection.generateActionParamStructs(packagePrefix, parentResources, thisResource)...)
		defs = append(defs, thisResource.Collection.Entity.generateActionParamStructs(packagePrefix, parentResources, thisResource)...)
		newParentResources := make([]*Resource, len(parentResources)+1)
		copy(parentResources, newParentResources)
		newParentResources = append(newParentResources, thisResource)
		for _, r := range thisResource.Collection.Entity.Subresources {
			defs = append(defs, generateAllActionStructs(packagePrefix, parentResources, &r)...)
		}
		return
	}

	if thisResource.Association != nil {
		defs = append(defs, thisResource.Association.generateActionParamStructs(packagePrefix, parentResources, thisResource)...)
		defs = append(defs, thisResource.Association.Entity.generateActionParamStructs(packagePrefix, parentResources, thisResource)...)
		newParentResources := make([]*Resource, len(parentResources)+1)
		copy(parentResources, newParentResources)
		newParentResources = append(newParentResources, thisResource)
		for _, r := range thisResource.Association.Entity.Subresources {
			defs = append(defs, generateAllActionStructs(packagePrefix, parentResources, &r)...)
		}
		return
	}

	if thisResource.ActionsSet != nil {
		defs = append(defs, thisResource.ActionsSet.generateActionParamStructs(packagePrefix, parentResources, thisResource)...)
		return
	}

	log.Panicln(thisResource, "does not define any resources")
	return
}

func (h *HasActions) generateActionParamStructs(packagePrefix string, parentResources []*Resource, thisResource *Resource) (defs []*Statement) {
	for _, a := range h.Actions {
		defs = append(defs, a.generateActionParamStructs(packagePrefix, parentResources, thisResource, false))
	}
	return defs
}

func (e *Entity) generateActionParamStructs(packagePrefix string, parentResources []*Resource, thisResource *Resource) (defs []*Statement) {
	for _, a := range e.Actions {
		defs = append(defs, a.generateActionParamStructs(packagePrefix, parentResources, thisResource, true))
	}
	return defs
}

func (a *Action) generateActionParamStructs(packagePrefix string, parentResources []*Resource, thisResource *Resource, isOnEntity bool) (def *Statement) {
	fullName := prefixNameWithParentResources(a.Name, parentResources, thisResource)
	structName := fullName + "ActionParams"

	def = Empty()
	def.Const().Id(fullName + "Action").Op("=").Lit(a.Name).Line()

	def.Type().Id(structName).StructFunc(func(def *Group) {
		for _, p := range a.Parameters {
			paramDef := def.Empty()
			AddWordWrappedComment(paramDef, p.Doc).Line()
			paramDef.Id(ExportedIdentifier(p.Name))
			paramDef.Add(p.Type.GoType(packagePrefix)).Tag(JsonTag(p.Name))
		}
	}).Line()

	var queryPath string
	if isOnEntity {
		queryPath = thisResource.getEntity().Path
	} else {
		queryPath = thisResource.Path
	}

	def.Func().Params(Id(ClientReceiver).Op("*").Id(thisResource.clientType())).Id(ExportedIdentifier(a.Name) + "Action")
	def.ParamsFunc(func(def *Group) {
		for _, r := range parentResources {
			if id := r.getIdentifier(); id != nil {
				def.Id(id.Name).Add(id.Type.GoType(packagePrefix))
				queryPath = strings.Replace(queryPath, fmt.Sprintf("{%s}", id.Name), "%s", 1)
			}
		}
		if id := thisResource.getIdentifier(); id != nil {
			def.Id(id.Name).Add(id.Type.GoType(packagePrefix))
			queryPath = strings.Replace(queryPath, fmt.Sprintf("{%s}", id.Name), "%s", 1)
		}
		def.Id("params").Id(structName)
	})

	Req := func() *Statement { return Id("req") }
	Res := func() *Statement { return Id("res") }
	ActionResult := func() *Statement { return Id("actionResult") }
	returns := a.Returns != nil

	def.ParamsFunc(func(def *Group) {
		if returns {
			def.Add(ActionResult()).Add(a.Returns.GoType(packagePrefix))
		}
		def.Err().Error()
	})

	def.BlockFunc(func(def *Group) {
		def.Id("url").Op(":=").Add(Id(ClientReceiver).Dot(HostnameClientField)).Op("+").Qual("fmt", "Sprintf").
			CallFunc(func(def *Group) {
				def.Lit(queryPath + "?action=" + a.Name)
				for _, r := range parentResources {
					if id := r.getIdentifier(); id != nil {
						def.Id(id.Name)
					}
				}
				if isOnEntity {
					def.Id(thisResource.getIdentifier().Name)
				}
			})
		def.Var().Add(Req()).Op("*").Qual(NetHttp, "Request")
		def.List(Req(), Err()).Op("=").Qual(packagePrefix+"/protocol", "RestliPost").Call(Id("url"), Lit(""), Id("params"))
		ifErrReturn(def)

		def.Var().Add(Res()).Op("*").Qual(NetHttp, "Response")
		def.List(Res(), Err()).Op("=").Id(ClientReceiver).Dot("Do").Call(Req())
		ifErrReturn(def)

		def.Err().Op("=").Qual(packagePrefix+"/protocol", "IsErrorResponse").Call(Res())
		ifErrReturn(def)

		if returns {
			def.Id("result").Op(":=").Struct(Id("Value").Add(a.Returns.GoType(packagePrefix))).Block()
			def.Err().Op("=").Qual("encoding/json", "NewDecoder").Call(Res().Dot("Body")).Dot("Decode").Call(Op("&").Id("result"))
			ifErrReturn(def)
			def.Add(ActionResult()).Op("=").Add(Id("result").Dot("Value"))
		}
		def.Return()
	})

	return
}

func prefixNameWithParentResources(name string, parentResources []*Resource, thisResource *Resource) string {
	var names []string
	for _, r := range parentResources {
		names = append(names, ExportedIdentifier(r.Name))
	}
	names = append(names, ExportedIdentifier(thisResource.Name))
	names = append(names, ExportedIdentifier(name))
	return strings.Join(names, "_")
}

func ifErrReturn(c *Group) {
	c.If(Err().Op("!=").Nil()).Block(Return()).Line()
}
