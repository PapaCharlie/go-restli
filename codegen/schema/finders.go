package schema

import (
	. "github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/protocol"
	. "github.com/dave/jennifer/jen"
)

const FindBy = "FindBy"

func (f *Finder) generate(parentResources []*Resource, thisResource *Resource) *CodeFile {
	c := NewCodeFile("findBy"+ExportedIdentifier(f.FinderName), thisResource.PackagePath(), thisResource.Name)

	c.Code.Const().Id(ExportedIdentifier(FindBy + ExportedIdentifier(f.FinderName))).Op("=").Lit(f.FinderName).Line()
	c.Code.Add(f.GenerateCode())

	returnType := Index().Add(thisResource.Schema.GoType())

	def := addClientFunc(c.Code, FindBy+ExportedIdentifier(f.FinderName)).
		ParamsFunc(func(def *Group) {
			addEntityTypes(def, parentResources)
			def.Id("params").Op("*").Id(f.StructName)
		}).
		Params(Add(returnType), Error())

	def.BlockFunc(func(def *Group) {
		def.List(Id("path"), Err()).Op(":=").Id(ClientReceiver).Dot(ResourcePath).Call(entityParams(parentResources)...)
		IfErrReturn(def, Nil(), Err()).Line()

		def.Id("query").Op(":=").Qual("net/url", "Values").Block()
		def.Id("query").Dot("Set").Call(Lit("q"), Lit(f.FinderName))
		for _, field := range f.Fields {
			accessor := Id("params").Dot(ExportedIdentifier(field.Name))

			setBlock := def.Empty()
			if field.IsPointer() {
				setBlock.If(Add(accessor).Op("!=").Nil())
			}

			if !field.Type.IsMapOrArray() && field.IsPointer() {
				accessor = Op("*").Add(accessor)
			}
			varName := field.Name + "Str"

			setBlock.BlockFunc(func(def *Group) {
				hasError, assignment := field.Type.RestLiURLEncode(accessor)
				if hasError {
					def.List(Id(varName), Err()).Op(":=").Add(assignment)
					IfErrReturn(def, Nil(), Err())
				} else {
					def.Id(varName).Op(":=").Add(assignment)
				}
				def.Id("query").Dot("Set").Call(Lit(field.Name), Id(varName))
			})
		}
		def.Line()

		def.Id("path").Op("+=").Lit("?").Op("+").Id("query").Dot("Encode").Call()

		def.List(Id("url"), Err()).Op(":=").Id(ClientReceiver).Dot(FormatQueryUrl).Call(Id("path"))
		IfErrReturn(def, Nil(), Err()).Line()

		def.List(Id(Req), Err()).Op(":=").Id(ClientReceiver).Dot("GetRequest").Call(Id("url"), RestLiMethod(protocol.NoMethod))
		IfErrReturn(def, Nil(), Err()).Line()

		def.Id("result").Op(":=").Struct(Id("Elements").Add(returnType)).Block()
		def.Err().Op("=").Id(ClientReceiver).Dot("DoAndDecode").Call(Id(Req), Op("&").Id("result"))
		IfErrReturn(def, Nil(), Err()).Line()
		def.Return(Id("result").Dot("Elements"), Nil())
	})

	return c
}
