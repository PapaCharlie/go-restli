package codegen

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/protocol"
	. "github.com/dave/jennifer/jen"
)

const (
	CreateParam = "create"
	UpdateParam = "update"
	QueryParams = "queryParams"
)

// createdEntityIdType returns true if the Create method is supposed to parse the created record's ID from the
// Location header in the response
func (m *Method) createdEntityIdType() *Statement {
	if up := m.EntityPathKey.Type.UnderlyingPrimitive(); up != nil {
		return m.EntityPathKey.Type.GoType()
	} else {
		return RawComplexKey()
	}
}

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
		def.Id(CreateParam).Add(r.ResourceSchema.ReferencedType())
	case protocol.Method_update:
		def.Id(UpdateParam).Add(r.ResourceSchema.ReferencedType())
	case protocol.Method_partial_update:
		def.Id(UpdateParam).Op("*").Add(r.ResourceSchema.Record().PartialUpdateStruct())
	}
	if len(m.Params) > 0 {
		def.Id(QueryParams).Op("*").Qual(r.PackagePath(), m.restMethodQueryParamsStructName())
	}
}

func (m *Method) restMethodFuncReturnParams(def *Group) {
	switch m.RestLiMethod() {
	case protocol.Method_get:
		def.Add(m.Return.ReferencedType())
	case protocol.Method_create:
		def.Add(m.createdEntityIdType())
	}
	def.Error()
}

func (m *Method) restMethodCallParams() (params []Code) {
	switch m.RestLiMethod() {
	case protocol.Method_create:
		params = append(params, Id(CreateParam))
	case protocol.Method_update:
		params = append(params, Id(UpdateParam))
	case protocol.Method_partial_update:
		params = append(params, Id(UpdateParam))
	}
	if len(m.Params) > 0 {
		params = append(params, Id(QueryParams))
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
		c.Code.Add(p.generateQueryParamEncoder(nil)).Line().Line()
	}

	r.addClientFuncDeclarations(c.Code, ClientType, m, func(def *Group) {
		generators[m.RestLiMethod()](r, m, def)
	})

	return c
}

func (m *Method) callResourcePath(def *Group) {
	if m.OnEntity {
		def.List(Id(PathVar), Err()).Op(":=").Id(ResourceEntityPath).Call(m.entityParams()...)
	} else {
		def.List(Id(PathVar), Err()).Op(":=").Id(ResourcePath).Call(m.entityParams()...)
	}
}

func generateGet(r *Resource, m *Method, def *Group) {
	returns := []Code{
		m.Return.ZeroValueReference(),
		Err(),
	}

	formatQueryUrl(r, m, def, returns...)

	def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver).Dot("GetRequest").Call(Id(ContextVar), Id(UrlVar), RestLiMethod(protocol.Method_get))
	IfErrReturn(def, returns...).Line()

	result := Id("getResult")
	def.Var().Add(result).Add(m.Return.GoType())
	callDoAndDecode(def, Op("&").Add(result), m.Return.ZeroValueReference())

	if m.Return.ShouldReference() {
		result = Op("&").Add(result)
	}
	def.Return(result, Nil())
}

func generateCreate(r *Resource, m *Method, def *Group) {
	primitiveReturnType := m.EntityPathKey.Type.UnderlyingPrimitive() != nil
	// TODO: Support @ReturnEntity annotation
	var returns []Code
	if primitiveReturnType {
		returns = append(returns, m.EntityPathKey.Type.UnderlyingPrimitiveZeroValueLit())
	} else {
		returns = append(returns, Lit(""))
	}
	returns = append(returns, Err())

	formatQueryUrl(r, m, def, returns...)

	def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver).Dot("JsonPostRequest").Call(Id(ContextVar), Id(UrlVar), RestLiMethod(protocol.Method_create), Id(CreateParam))
	IfErrReturn(def, returns...).Line()

	def.List(Id(ResVar), Err()).Op(":=").Id(ClientReceiver).Dot(DoAndIgnore).Call(Id(ReqVar))
	IfErrReturn(def, returns...).Line()

	def.If(Id(ResVar).Dot("StatusCode").Op("/").Lit(100).Op("!=").Lit(2)).BlockFunc(func(def *Group) {
		def.Err().Op("=").Qual("fmt", "Errorf").Call(Lit("Invalid response code from %s: %d"), Id(UrlVar), Id(ResVar).Dot("StatusCode"))
		def.Return(returns...)
	}).Line()

	if primitiveReturnType {
		accessor := Id(m.EntityPathKey.Name)
		def.Var().Add(accessor).Add(m.EntityPathKey.Type.GoType())

		def.Err().Op("=").Add(m.EntityPathKey.Type.RestLiReducedDecodeModel(
			Id(ResVar).Dot("Header").Dot("Get").Call(Qual(ProtocolPackage, RestLiHeaderID)),
			accessor,
		))

		IfErrReturn(def, returns...)

		def.Return(accessor, Nil())
	} else {
		def.Return(RawComplexKey().Call(Id(ResVar).Dot("Header").Dot("Get").Call(Qual(ProtocolPackage, RestLiHeaderID))), Nil())
	}
}

func generateUpdate(r *Resource, m *Method, def *Group) {
	formatQueryUrl(r, m, def, Err())

	def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver).Dot("JsonPutRequest").Call(Id(ContextVar), Id(UrlVar), RestLiMethod(protocol.Method_update), Id(UpdateParam))
	IfErrReturn(def, Err()).Line()

	def.List(Id(ResVar), Err()).Op(":=").Id(ClientReceiver).Dot(DoAndIgnore).Call(Id(ReqVar))
	IfErrReturn(def, Err()).Line()

	def.If(Id(ResVar).Dot("StatusCode").Op("/").Lit(100).Op("!=").Lit(2)).BlockFunc(func(def *Group) {
		def.Return(Qual("fmt", "Errorf").Call(Lit("Invalid response code from %s: %d"), Id(UrlVar), Id(ResVar).Dot("StatusCode")))
	})
	def.Return(Nil())
}

func generatePartialUpdate(r *Resource, m *Method, def *Group) {
	formatQueryUrl(r, m, def, Err())

	def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver).Dot("PartialUpdateRequest").Call(
		Id(ContextVar),
		Id(UrlVar),
		Id(UpdateParam),
	)
	IfErrReturn(def, Err()).Line()

	def.List(Id(ResVar), Err()).Op(":=").Id(ClientReceiver).Dot(DoAndIgnore).Call(Id(ReqVar))
	IfErrReturn(def, Err()).Line()

	def.If(Id(ResVar).Dot("StatusCode").Op("/").Lit(100).Op("!=").Lit(2)).BlockFunc(func(def *Group) {
		def.Return(Qual("fmt", "Errorf").Call(Lit("Invalid response code from %s: %d"), Id(UrlVar), Id(ResVar).Dot("StatusCode")))
	})
	def.Return(Nil())
}

func generateDelete(r *Resource, m *Method, def *Group) {
	formatQueryUrl(r, m, def, Err())

	def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver).Dot("DeleteRequest").Call(Id(ContextVar), Id(UrlVar), RestLiMethod(protocol.Method_update))
	IfErrReturn(def, Err()).Line()

	def.List(Id(ResVar), Err()).Op(":=").Id(ClientReceiver).Dot(DoAndIgnore).Call(Id(ReqVar))
	IfErrReturn(def, Err()).Line()

	def.If(Id(ResVar).Dot("StatusCode").Op("/").Lit(100).Op("!=").Lit(2)).BlockFunc(func(def *Group) {
		def.Return(Qual("fmt", "Errorf").Call(Lit("Invalid response code from %s: %d"), Id(UrlVar), Id(ResVar).Dot("StatusCode")))
	})
	def.Return(Nil())
}

func formatQueryUrl(r *Resource, m *Method, def *Group, returns ...Code) {
	m.callResourcePath(def)
	IfErrReturn(def, returns...).Line()

	r.callFormatQueryUrl(def)
	IfErrReturn(def, returns...).Line()

	if len(m.Params) > 0 {
		def.List(Id(UrlVar).Dot("RawQuery"), Err()).Op("=").Id(QueryParams).Dot(EncodeQueryParams).Call()
		IfErrReturn(def, returns...)
	}
}
