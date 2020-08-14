package codegen

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/protocol"
	. "github.com/dave/jennifer/jen"
)

const EncodeQueryParams = "EncodeQueryParams"

func (m *Method) finderFuncName() string {
	return FindBy + ExportedIdentifier(m.Name)
}

func (m *Method) finderStructType() string {
	return FindBy + ExportedIdentifier(m.Name) + "Params"
}

func (r *Resource) finderFuncParams(m *Method, def *Group) {
	m.addEntityTypes(def)
	def.Id(QueryParams).Op("*").Qual(r.PackagePath(), m.finderStructType())
}

func (m *Method) finderMethodCallParams() (params []Code) {
	if len(m.Params) > 0 {
		params = append(params, Id(QueryParams))
	}
	return params
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

	params := &Record{
		NamedType: NamedType{
			Identifier: Identifier{
				Name:      f.finderStructType(),
				Namespace: r.Namespace,
			},
			Doc: fmt.Sprintf("This struct provides the parameters to the %s finder", f.Name),
		},
		Fields: f.Params,
	}
	c.Code.Add(params.generateStruct()).Line().Line()
	c.Code.Add(params.generateQueryParamEncoder(&f.Name))

	r.addClientFuncDeclarations(c.Code, ClientType, f, func(def *Group) {
		formatQueryUrl(r, f, def, Nil(), Err())

		def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver).Dot("GetRequest").Call(Id(ContextVar), Id(UrlVar), RestLiMethod(protocol.Method_finder))
		IfErrReturn(def, Nil(), Err()).Line()

		accessor := Id("elements")
		def.Var().Add(accessor).Struct(Id("Elements").Add(f.finderReturnType()))
		callDoAndDecode(def, Op("&").Add(accessor), Nil())
		def.Return(Add(accessor).Dot("Elements"), Nil())
	})

	return c
}
