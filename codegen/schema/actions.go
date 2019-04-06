package schema

import (
	. "github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
)

func (a *Action) generateActionParamStructs(parentResources []*Resource, thisResource *Resource, isOnEntity bool) (c *CodeFile) {
	c = NewCodeFile(a.ActionName, thisResource.PackagePath(), thisResource.Name)

	c.Code.Const().Id(ExportedIdentifier(a.ActionName + "Action")).Op("=").Lit(a.ActionName).Line()
	c.Code.Add(a.GenerateCode())

	var resources []*Resource
	resources = append(resources, parentResources...)
	if isOnEntity {
		resources = append(resources, thisResource)
	}

	var queryPath string
	if isOnEntity {
		queryPath = thisResource.getEntity().Path
	} else {
		queryPath = thisResource.Path
	}
	queryPath = buildQueryPath(resources, queryPath)

	AddClientFunc(c.Code, ExportedIdentifier(a.ActionName)+"Action")
	c.Code.ParamsFunc(func(def *Group) {
		addEntityParams(def, resources)
		def.Id("params").Op("*").Id(a.StructName)
	})

	returns := a.Returns != nil

	c.Code.ParamsFunc(func(def *Group) {
		if returns {
			def.Id(ActionResult).Add(a.Returns.GoType())
		}
		def.Err().Error()
	})

	c.Code.BlockFunc(func(def *Group) {
		encodeEntitySegments(def, resources)

		def.Id(Url).Op(":=").Id(ClientReceiver).Dot(HostnameClientField).Op("+").Qual("fmt", "Sprintf").
			CallFunc(func(def *Group) {
				def.Lit(queryPath + "?action=" + a.ActionName)
				for _, r := range resources {
					if id := r.getIdentifier(); id != nil {
						def.Id(id.Name + "Str")
					}
				}
			})
		def.List(Id(Req), Err()).Op(":=").Id(ClientReceiver).Dot("PostRequest").Call(Id("url"), Lit(""), Id("params"))
		IfErrReturn(def).Line()

		var resDef *Statement
		if returns {
			resDef = def.List(Id(Res), Err()).Op(":=")
		} else {
			resDef = def.List(Id("_"), Err()).Op("=")
		}
		resDef.Id(ClientReceiver).Dot("Do").Call(Id(Req))
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
