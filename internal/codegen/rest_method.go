package codegen

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/protocol"
	. "github.com/dave/jennifer/jen"
)

var (
	CreateParam = Code(Id("create"))
	UpdateParam = Code(Id("update"))
	QueryParams = Code(Id("queryParams"))
)

func (m *Method) RestLiMethod() protocol.RestLiMethod {
	return protocol.RestLiMethodNameMapping[m.Name]
}

var restMethodFuncNames = map[protocol.RestLiMethod]string{
	protocol.Method_get:            "Get",
	protocol.Method_create:         "Create",
	protocol.Method_delete:         "Delete",
	protocol.Method_update:         "Update",
	protocol.Method_partial_update: "PartialUpdate",

	protocol.Method_batch_get:            "BatchGet",
	protocol.Method_batch_create:         "BatchCreate",
	protocol.Method_batch_delete:         "BatchDelete",
	protocol.Method_batch_update:         "BatchUpdate",
	protocol.Method_batch_partial_update: "BatchPartialUpdate",
}

func (m *Method) restMethodFuncName() string {
	return restMethodFuncNames[m.RestLiMethod()]
}

func (m *Method) restMethodQueryParamsStructName() string {
	return restMethodFuncNames[m.RestLiMethod()] + "Params"
}

func (r *Resource) restMethodFuncParams(m *Method, def *Group) {
	m.addEntityTypes(def)
	switch m.RestLiMethod() {
	case protocol.Method_create:
		def.Add(CreateParam).Add(r.ResourceSchema.ReferencedType())
	case protocol.Method_update:
		def.Add(UpdateParam).Add(r.ResourceSchema.ReferencedType())
	case protocol.Method_partial_update:
		def.Add(UpdateParam).Op("*").Add(r.ResourceSchema.Record().PartialUpdateStruct())
	}
	if len(m.Params) > 0 {
		def.Add(QueryParams).Op("*").Qual(r.PackagePath(), m.restMethodQueryParamsStructName())
	}
}

func (m *Method) restMethodFuncReturnParams(def *Group) {
	switch m.RestLiMethod() {
	case protocol.Method_get:
		def.Add(m.Return.ReferencedType())
	case protocol.Method_create:
		def.Add(m.EntityPathKey.Type.ReferencedType())
		if m.ReturnEntity {
			def.Add(m.Return.ReferencedType())
		}
	}
	def.Error()
}

func (m *Method) restMethodCallParams() (params []Code) {
	switch m.RestLiMethod() {
	case protocol.Method_create:
		params = append(params, CreateParam)
	case protocol.Method_update:
		params = append(params, UpdateParam)
	case protocol.Method_partial_update:
		params = append(params, UpdateParam)
	}
	if len(m.Params) > 0 {
		params = append(params, QueryParams)
	}

	return params
}

var generators = map[protocol.RestLiMethod]func(*Resource, *Method, *Group){
	protocol.Method_get:            generateGet,
	protocol.Method_create:         generateCreate,
	protocol.Method_update:         generateUpdate,
	protocol.Method_partial_update: generatePartialUpdate,
	protocol.Method_delete:         generateDelete,
}

func isMethodSupported(m protocol.RestLiMethod) bool {
	_, ok := generators[m]
	return ok
}

// https://linkedin.github.io/rest.li/user_guide/restli_server#resource-methods
func (r *Resource) GenerateRestMethodCode(m *Method) *CodeFile {
	c := r.NewCodeFile(m.Name)

	if len(m.Params) > 0 {
		p := &Record{
			NamedType: NamedType{
				Identifier: Identifier{
					Name:      m.restMethodQueryParamsStructName(),
					Namespace: r.Namespace,
				},
				Doc: fmt.Sprintf("This struct provides the parameters to the %s method", m.Name),
			},
			Fields: m.Params,
		}
		c.Code.Add(p.generateStruct()).Line().Line()
		c.Code.Add(p.generateQueryParamMarshaler(nil)).Line().Line()
	}

	r.addClientFuncDeclarations(c.Code, ClientType, m, func(def *Group) {
		generators[m.RestLiMethod()](r, m, def)
	})

	return c
}

func (m *Method) callResourcePath(def *Group) {
	if m.OnEntity {
		def.List(PathVar, Err()).Op(":=").Id(ResourceEntityPath).Call(m.entityParams()...)
	} else {
		def.List(PathVar, Err()).Op(":=").Id(ResourcePath).Call(m.entityParams()...)
	}
}

func generateGet(r *Resource, m *Method, def *Group) {
	returns := []Code{
		m.Return.ZeroValueReference(),
		Err(),
	}

	formatQueryUrl(r, m, def, returns...)

	result := Id("getResult")
	def.Var().Add(result).Add(m.Return.GoType())

	def.Err().Op("=").Id(ClientReceiver).Dot("DoGetRequest").Call(ContextVar, UrlVar, Op("&").Add(result))
	def.Add(IfErrReturn(returns...)).Line()

	if m.Return.ShouldReference() {
		result = Op("&").Add(result)
	}
	def.Return(result, Nil())
}

func generateCreate(r *Resource, m *Method, def *Group) {
	returns := []Code{m.EntityPathKey.Type.ZeroValueReference()}
	if m.ReturnEntity {
		returns = append(returns, m.Return.ZeroValueReference())
	}
	returns = append(returns, Err())

	formatQueryUrl(r, m, def, returns...)

	id := Id(m.EntityPathKey.Name)
	def.Var().Add(id).Add(m.EntityPathKey.Type.GoType())

	var idUnmarshaler Code
	if p := m.EntityPathKey.Type.Primitive; p != nil {
		idUnmarshaler = p.NewPrimitiveUnmarshaler(id)
	} else {
		idUnmarshaler = Op("&").Add(id)
	}

	returnEntity := Id("returnEntity")
	var returnEntityUnmarshaler Code
	if m.ReturnEntity {
		def.Var().Add(returnEntity).Add(m.Return.GoType())
		returnEntityUnmarshaler = Op("&").Add(returnEntity)
	} else {
		returnEntityUnmarshaler = Nil()
	}

	def.Err().Op("=").Id(ClientReceiver).Dot("DoCreateRequest").Call(ContextVar, UrlVar, CreateParam, idUnmarshaler, returnEntityUnmarshaler)
	def.Add(IfErrReturn(returns...)).Line()

	if m.EntityPathKey.Type.ShouldReference() {
		id = Op("&").Add(id)
	}
	if m.Return.ShouldReference() {
		returnEntity = Op("&").Add(returnEntity)
	}

	if m.ReturnEntity {
		def.Return(id, returnEntity, Nil())
	} else {
		def.Return(id, Nil())
	}
}

func generateUpdate(r *Resource, m *Method, def *Group) {
	formatQueryUrl(r, m, def, Err())

	def.Err().Op("=").Id(ClientReceiver).Dot("DoUpdateRequest").Call(ContextVar, UrlVar, UpdateParam)
	def.Return(Err())
}

func generatePartialUpdate(r *Resource, m *Method, def *Group) {
	formatQueryUrl(r, m, def, Err())

	def.Err().Op("=").Id(ClientReceiver).Dot("DoPartialUpdateRequest").Call(ContextVar, UrlVar, UpdateParam)
	def.Return(Err())
}

func generateDelete(r *Resource, m *Method, def *Group) {
	formatQueryUrl(r, m, def, Err())

	def.Err().Op("=").Id(ClientReceiver).Dot("DoDeleteRequest").Call(ContextVar, UrlVar)
	def.Add(IfErrReturn(Err())).Line()

	def.Return(Nil())
}

func formatQueryUrl(r *Resource, m *Method, def *Group, returns ...Code) {
	m.callResourcePath(def)
	def.Add(IfErrReturn(returns...)).Line()

	if m.MethodType != ACTION && len(m.Params) > 0 {
		rawQuery := Id("rawQuery")
		def.Var().Add(rawQuery).String()
		def.List(rawQuery, Err()).Op("=").Add(QueryParams).Dot(EncodeQueryParams).Call()
		def.Add(IfErrReturn(returns...))
		def.Add(PathVar).Op("+=").Lit("?").Op("+").Add(rawQuery)
		def.Line()
	}

	r.callFormatQueryUrl(def)
	def.Add(IfErrReturn(returns...)).Line()
}
