package codegen

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/protocol"
	. "github.com/dave/jennifer/jen"
)

func (m *Method) finderFuncName() string {
	return FindBy + ExportedIdentifier(m.Name)
}

func (m *Method) finderStructType() string {
	return FindBy + ExportedIdentifier(m.Name) + "Params"
}

func (m *Method) finderFuncParams(def *Group) {
	m.addEntityTypes(def)
	def.Id("params").Op("*").Id(m.finderStructType())
}

func (m *Method) finderReturnType() Code {
	return Index().Add(m.Return.PointerType())
}

func (m *Method) finderFuncReturnParams(def *Group) {
	def.Add(m.finderReturnType())
	def.Error()
}

func (r *Resource) GenerateFinderCode(f *Method) *CodeFile {
	c := r.NewCodeFile("findBy" + ExportedIdentifier(f.Name))

	c.Code.Const().Id(ExportedIdentifier(FindBy + ExportedIdentifier(f.Name))).Op("=").Lit(f.Name).Line()

	record := &Record{
		NamedType: NamedType{
			Identifier: Identifier{
				Name:      f.finderStructType(),
				Namespace: r.Namespace,
			},
			Doc: fmt.Sprintf("This struct provides the parameters to the %s finder", f.Name),
		},
		Fields: f.Params,
	}
	c.Code.Add(record.GenerateCode())

	AddWordWrappedComment(c.Code, f.Doc).Line()
	r.addClientFunc(c.Code, f)

	c.Code.BlockFunc(func(def *Group) {
		def.List(Id(PathVar), Err()).Op(":=").Id(ResourcePath).Call(f.entityParams()...)
		IfErrReturn(def, Nil(), Err()).Line()

		def.Id("query").Op(":=").Qual("net/url", "Values").Block()
		def.Id("query").Dot("Set").Call(Lit("q"), Lit(f.Name))
		def.Line()

		for _, field := range f.Params {
			accessor := Id("params").Dot(ExportedIdentifier(field.Name))

			setBlock := def.Empty()
			if field.IsPointer() {
				setBlock.If(Add(accessor).Op("!=").Nil())
			}

			if field.IsPointer() && field.Type.Reference == nil && field.Type.Union == nil {
				accessor = Op("*").Add(accessor)
			}
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

		def.Id(PathVar).Op("+=").Lit("?").Op("+").Id("query").Dot("Encode").Call()

		r.callFormatQueryUrl(def)
		IfErrReturn(def, Nil(), Err()).Line()

		def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver).Dot("GetRequest").Call(Id(UrlVar), RestLiMethod(protocol.Method_finder))
		IfErrReturn(def, Nil(), Err()).Line()

		def.Id(DoAndDecodeResult).Op(":=").Struct(Id("Elements").Add(f.finderReturnType())).Block()
		callDoAndDecode(def)
		def.Return(Id(DoAndDecodeResult).Dot("Elements"), Nil())
	})

	return c
}
