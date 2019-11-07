package schema

import (
	. "github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/protocol"
	. "github.com/dave/jennifer/jen"
)

func (a *Action) generate(parentResources []*Resource, thisResource *Resource, isOnEntity bool) (c *CodeFile) {
	c = NewCodeFile(a.ActionName+"Action", thisResource.PackagePath(), thisResource.Name)

	c.Code.Const().Id(ExportedIdentifier(a.ActionName + "Action")).Op("=").Lit(a.ActionName).Line()

	hasParams := len(a.Fields) > 0
	if hasParams {
		c.Code.Add(a.GenerateCode())
	}

	var resources []*Resource
	resources = append(resources, parentResources...)
	if isOnEntity {
		resources = append(resources, thisResource)
	}

	addClientFunc(c.Code, ExportedIdentifier(a.ActionName)+"Action")
	c.Code.ParamsFunc(func(def *Group) {
		addEntityTypes(def, resources)
		if hasParams {
			def.Id("params").Op("*").Id(a.StructName)
		}
	})

	returns := a.Returns != nil

	c.Code.ParamsFunc(func(def *Group) {
		if returns {
			def.Op("*").Add(a.Returns.GoType())
		}
		def.Error()
	})

	c.Code.BlockFunc(func(def *Group) {
		var pathFunc string
		if isOnEntity {
			pathFunc = ResourceEntityPath
		} else {
			pathFunc = ResourcePath
		}

		var errReturnParams []Code
		if returns {
			errReturnParams = []Code{Nil(), Err()}
		} else {
			errReturnParams = []Code{Err()}
		}

		def.List(Id("path"), Err()).Op(":=").Id(ClientReceiver).Dot(pathFunc).Call(entityParams(resources)...)
		IfErrReturn(def, errReturnParams...).Line()
		def.Id("path").Op("+=").Lit("?action=").Op("+").Id(ExportedIdentifier(a.ActionName + "Action"))

		def.List(Id("url"), Err()).Op(":=").Id(ClientReceiver).Dot(FormatQueryUrl).Call(Id("path"))
		IfErrReturn(def, errReturnParams...).Line()

		req := def.List(Id(Req), Err()).Op(":=").Id(ClientReceiver)
		var params *Statement
		if hasParams {
			params = Id("params")
		} else {
			params = Struct().Block()
		}
		req.Dot("JsonPostRequest").Call(Id("url"), RestLiMethod(protocol.Method_action), params)
		IfErrReturn(def, errReturnParams...).Line()

		if returns {
			def.Id("result").Op(":=").Struct(Id("Value").Add(a.Returns.GoType())).Block()
			def.Err().Op("=").Id(ClientReceiver).Dot("DoAndDecode").Call(Id(Req), Op("&").Id("result"))
			IfErrReturn(def, errReturnParams...).Line()

			def.Return(Op("&").Id("result").Dot("Value"), Nil())
		} else {
			def.Err().Op("=").Id(ClientReceiver).Dot("DoAndIgnore").Call(Id(Req))
			IfErrReturn(def, errReturnParams...).Line()
			def.Return(Nil())
		}
	})

	return c
}
