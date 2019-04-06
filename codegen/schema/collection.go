package schema

import (
	. "github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
)

func (c *Collection) generateGet(code *CodeFile, parentResources []*Resource, thisResource *Resource, m Method) {
	var resources []*Resource
	resources = append(resources, parentResources...)
	resources = append(resources, thisResource)

	queryPath := thisResource.Path
	queryPath = buildQueryPath(resources, thisResource.Path) + "/%s"

	AddWordWrappedComment(code.Code, m.Doc).Line()
	AddClientFunc(code.Code, RestliMethodNameMapping[m.Method])
	code.Code.ParamsFunc(func(def *Group) {
		addEntityParams(def, resources)
	})

	result := m.Method + "Result"
	code.Code.ParamsFunc(func(def *Group) {
		def.Id(result).Add(thisResource.Schema.GoType())
		def.Err().Error()
	})

	code.Code.BlockFunc(func(def *Group) {
		encodeEntitySegments(def, resources)

		def.Id(Url).Op(":=").Id(ClientReceiver).Dot(HostnameClientField).Op("+").Qual("fmt", "Sprintf").
			CallFunc(func(def *Group) {
				def.Lit(queryPath)
				for _, r := range resources {
					if id := r.getIdentifier(); id != nil {
						def.Id(id.Name + "Str")
					}
				}
			})
		def.List(Id(Req), Err()).Op(":=").Id(ClientReceiver).Dot("GetRequest").Call(Id("url"), Lit(""))
		IfErrReturn(def).Line()

		def.List(Id(Res), Err()).Op(":=").Id(ClientReceiver).Dot("Do").Call(Id(Req))
		IfErrReturn(def).Line()

		def.Err().Op("=").Qual(EncodingJson, "NewDecoder").Call(Id(Res).Dot("Body")).Dot("Decode").Call(Op("&").Id(result))
		def.Return()
	})

	return
}
