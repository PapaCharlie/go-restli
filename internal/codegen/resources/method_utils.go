package resources

import (
	"github.com/PapaCharlie/go-restli/internal/codegen/types"
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

type MethodImplementation interface {
	GetMethod() *Method
	GetResource() *Resource
	IsSupported() bool
	FuncName() string
	FuncParamNames() []Code
	FuncParamTypes() []Code
	NonErrorFuncReturnParams() []Code
	GenerateCode() *utils.CodeFile
}

type methodImplementation struct {
	Resource *Resource
	*Method
}

func (m *methodImplementation) GetMethod() *Method {
	return m.Method
}

func (m *methodImplementation) GetResource() *Resource {
	return m.Resource
}

func formatQueryUrl(m MethodImplementation, def *Group, keyWriter func(itemWriter Code, def *Group), returns ...Code) {
	if m.GetMethod().OnEntity {
		def.List(Path, Err()).Op(":=").Id(ResourceEntityPath).Call(m.GetMethod().entityParamNames()...)
	} else {
		def.List(Path, Err()).Op(":=").Id(ResourcePath).Call(m.GetMethod().entityParamNames()...)
	}

	def.Add(utils.IfErrReturn(returns...)).Line()

	encodeQueryParams := Code(Add(QueryParams).Dot(utils.EncodeQueryParams))
	callEncodeQueryParams := func(encoder Code) {
		rawQuery := Id("rawQuery")
		def.Var().Add(rawQuery).String()
		def.List(rawQuery, Err()).Op("=").Add(encoder)
		def.Add(utils.IfErrReturn(returns...))
		def.Add(Path).Op("+=").Lit("?").Op("+").Add(rawQuery)
		def.Line()
	}

	switch m.(type) {
	case *Action:
		def.Add(Path).Op("+=").Lit("?action=" + m.GetMethod().Name)
	case *Finder:
		callEncodeQueryParams(Add(encodeQueryParams).Call())
	case *RestMethod:
		r := m.(*RestMethod)
		hasParams := len(m.GetMethod().Params) > 0
		if r.isBatch() {
			encoder := Func().Params(Add(types.ItemWriter).Func().Params().Add(types.WriterQual)).Params(Err().Error()).
				BlockFunc(func(def *Group) { keyWriter(types.ItemWriter, def) })
			if hasParams {
				callEncodeQueryParams(Add(encodeQueryParams).Call(encoder))
			} else {
				callEncodeQueryParams(Qual(utils.ProtocolPackage, "GenerateBatchKeysParam").Call(encoder))
			}
		} else {
			if hasParams {
				callEncodeQueryParams(Add(encodeQueryParams).Call())
			}
		}
	}

	def.List(Url, Err()).
		Op(":=").
		Id(ClientReceiver).Dot("FormatQueryUrl").
		Call(Lit(m.GetResource().RootResourceName), Path)
	def.Add(utils.IfErrReturn(returns...)).Line()
}

func methodFuncName(m MethodImplementation, withContext bool) string {
	n := m.FuncName()
	if withContext {
		n += WithContext
	}
	return n
}

func addParams(def *Group, names, types []Code) {
	for i, name := range names {
		def.Add(name).Add(types[i])
	}
}

func methodParamNames(m MethodImplementation) []Code {
	return append(m.GetMethod().entityParamNames(), m.FuncParamNames()...)
}

func methodParamTypes(m MethodImplementation) []Code {
	return append(m.GetMethod().entityParamTypes(), m.FuncParamTypes()...)
}

func methodReturnParams(m MethodImplementation) []Code {
	return append(m.NonErrorFuncReturnParams(), Err().Error())
}

func (m *Method) entityParamNames() (params []Code) {
	for _, pk := range m.PathKeys {
		params = append(params, Id(pk.Name))
	}
	return params
}

func (m *Method) entityParamTypes() (params []Code) {
	for _, pk := range m.PathKeys {
		params = append(params, pk.Type.ReferencedType())
	}
	return params
}
