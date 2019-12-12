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

	resources := append([]*Resource(nil), parentResources...)
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

		def.List(Id(Path), Err()).Op(":=").Id(ClientReceiver).Dot(pathFunc).Call(entityParams(resources)...)
		IfErrReturn(def, errReturnParams...).Line()
		def.Id(Path).Op("+=").Lit("?action=").Op("+").Id(ExportedIdentifier(a.ActionName + "Action"))

		callFormatQueryUrl(def, parentResources, thisResource)
		IfErrReturn(def, errReturnParams...).Line()

		req := def.List(Id(Req), Err()).Op(":=").Id(ClientReceiver)
		var params *Statement
		if hasParams {
			params = Id("params")
		} else {
			params = Struct().Block()
		}
		req.Dot("JsonPostRequest").Call(Id(Url), RestLiMethod(protocol.Method_action), params)
		IfErrReturn(def, errReturnParams...).Line()

		if returns {
			def.Id(DoAndDecodeResult).Op(":=").Struct(Id("Value").Add(a.Returns.GoType())).Block()
			callDoAndDecode(def)
			def.Return(Op("&").Id(DoAndDecodeResult).Dot("Value"), Nil())
		} else {
			def.List(Id("_"), Err()).Op("=").Id(ClientReceiver).Dot(DoAndIgnore).Call(Id(Req))
			IfErrReturn(def, errReturnParams...).Line()
			def.Return(Nil())
		}
	})

	return c
}
