package codegen

import (
	"fmt"
	"github.com/PapaCharlie/go-restli/protocol"
	. "github.com/dave/jennifer/jen"
)

const EncodeFinderParams = "EncodeFinderParams"

type FinderParams Record

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

	params := &FinderParams{
		NamedType: NamedType{
			Identifier: Identifier{
				Name:      f.finderStructType(),
				Namespace: r.Namespace,
			},
			Doc: fmt.Sprintf("This struct provides the parameters to the %s finder", f.Name),
		},
		Fields: f.Params,
	}
	c.Code.Add(params.GenerateCode(f)).Line().Line()

	AddWordWrappedComment(c.Code, f.Doc).Line()
	r.addClientFunc(c.Code, f)

	c.Code.BlockFunc(func(def *Group) {
		def.List(Id(PathVar), Err()).Op(":=").Id(ResourcePath).Call(f.entityParams()...)
		IfErrReturn(def, Nil(), Err()).Line()

		def.List(Id("query"), Err()).Op(":=").Id("params").Dot(EncodeFinderParams).Call()
		IfErrReturn(def, Nil(), Err()).Line()

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

func (p *FinderParams) GenerateCode(f *Method) *Statement {
	def := Empty()
	def.Add((*Record)(p).generateStruct()).Line().Line()

	receiver := (*Record)(p).Receiver()
	return AddFuncOnReceiver(def, receiver, p.Name, EncodeFinderParams).
		Params().
		Params(Id("query").Qual("net/url", "Values"), Err().Error()).
		BlockFunc(func(def *Group) {
			def.Id(Codec).Op(":=").Qual(ProtocolPackage, RestLiUrlEncoder).Line()

			def.Id("query").Op("=").Make(Qual("net/url", "Values"))
			def.Id("query").Dot("Set").Call(Lit("q"), Lit(f.Name))
			def.Line()

			def.Var().Id("buf").Qual("strings", "Builder")

			for _, field := range f.Params {
				accessor := Id(receiver).Dot(ExportedIdentifier(field.Name))

				setBlock := def.Empty()
				if field.IsPointer() {
					setBlock.If(Add(accessor).Op("!=").Nil())
				}

				if field.IsPointer() && field.Type.Reference == nil && field.Type.Union == nil {
					accessor = Op("*").Add(accessor)
				}

				setBlock.BlockFunc(func(def *Group) {
					field.Type.WriteToBuf(def, accessor)
					def.Id("query").Dot("Set").Call(Lit(field.Name), Id("buf").Dot("String").Call())
					Id("buf").Dot("Reset").Call()
				})
				def.Line()
			}

			def.Return(Id("query"), Err())
		})
}
