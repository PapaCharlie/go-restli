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

	returnType := Index().Add(thisResource.Schema.Model.GoType())

	def := thisResource.addClientFunc(c.Code, f.Doc, func(def *Statement) *Statement {
		return def.Id(FindBy+ExportedIdentifier(f.FinderName)).
			ParamsFunc(func(def *Group) {
				addEntityTypes(def, parentResources)
				def.Id("params").Op("*").Id(f.StructName)
			}).
			Params(Add(returnType), Error())
	})

	def.BlockFunc(func(def *Group) {
		def.List(Id(Path), Err()).Op(":=").Id(ResourcePath).Call(entityParams(parentResources)...)
		IfErrReturn(def, Nil(), Err()).Line()

		def.Id("query").Op(":=").Qual("net/url", "Values").Block()
		def.Id("query").Dot("Set").Call(Lit("q"), Lit(f.FinderName))
		def.Line()

		for _, field := range f.Fields {
			accessor := Id("params").Dot(ExportedIdentifier(field.Name))

			setBlock := def.Empty()
			if field.IsPointer() {
				setBlock.If(Add(accessor).Op("!=").Nil())
			}

			accessor = field.RawAccessor(accessor)
			varName := field.Name + "Str"

			setBlock.BlockFunc(func(def *Group) {
				assignment, hasError := field.Type.RestLiURLEncodeModel(accessor)
				if hasError {
					def.List(Id(varName), Err()).Op(":=").Add(assignment)
					IfErrReturn(def, Nil(), Err())
				} else {
					def.Id(varName).Op(":=").Add(assignment)
				}
				def.Id("query").Dot("Set").Call(Lit(field.Name), Id(varName))
			})
			def.Line()
		}

		def.Id(Path).Op("+=").Lit("?").Op("+").Id("query").Dot("Encode").Call()

		callFormatQueryUrl(def, parentResources, thisResource)
		IfErrReturn(def, Nil(), Err()).Line()

		def.List(Id(Req), Err()).Op(":=").Id(ClientReceiver).Dot("GetRequest").Call(Id(Url), RestLiMethod(protocol.Method_finder))
		IfErrReturn(def, Nil(), Err()).Line()

		def.Id(DoAndDecodeResult).Op(":=").Struct(Id("Elements").Add(returnType)).Block()
		callDoAndDecode(def)
		def.Return(Id(DoAndDecodeResult).Dot("Elements"), Nil())
	})

	return c
}
