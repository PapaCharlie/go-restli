package codegen

import (
	"github.com/PapaCharlie/go-restli/protocol"
	. "github.com/dave/jennifer/jen"
)

const (
	CreateParam = "create"
	UpdateParam = "update"
)

// isCreatedEntityIdInHeaders returns true if the Create method is supposed to parse the created record's ID from the
// Location header in the response
func (m *Method) isCreatedEntityIdInHeaders() bool {
	if m.EntityPathKey == nil {
		return false
	}

	return m.EntityPathKey.Type.UnderlyingPrimitive() != nil
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

func (r *Resource) restMethodFuncParams(m *Method, def *Group) {
	switch m.RestLiMethod() {
	case protocol.Method_get:
		m.addEntityTypes(def)
	case protocol.Method_create:
		m.addEntityTypes(def)
		def.Id(CreateParam).Add(r.ResourceSchema.ReferencedType())
	case protocol.Method_update:
		m.addEntityTypes(def)
		def.Id(UpdateParam).Add(r.ResourceSchema.ReferencedType())
	case protocol.Method_partial_update:
		m.addEntityTypes(def)
		def.Id(UpdateParam).Add(Op("*").Add(r.ResourceSchema.Record().PartialUpdateStruct()))
	case protocol.Method_delete:
		m.addEntityTypes(def)
	}
}

func (m *Method) restMethodFuncReturnParams(def *Group) {
	switch m.RestLiMethod() {
	case protocol.Method_get:
		def.Add(m.Return.ReferencedType())
		def.Error()
	case protocol.Method_create:
		if m.isCreatedEntityIdInHeaders() {
			def.Add(m.EntityPathKey.Type.GoType())
		}
		def.Error()
	case protocol.Method_update:
		def.Error()
	case protocol.Method_partial_update:
		def.Error()
	case protocol.Method_delete:
		def.Error()
	}
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
func (r *Resource) GenerateRestMethodCode(m *Method) *Statement {
	return r.addClientFuncDeclarations(Empty(), ClientType, m, func(def *Group) {
		generators[m.RestLiMethod()](r, m, def)
	})
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

	m.callResourcePath(def)
	IfErrReturn(def, returns...).Line()
	r.callFormatQueryUrl(def)
	IfErrReturn(def, returns...).Line()

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
	// TODO: Support @ReturnEntity annotation
	var returns []Code
	if m.isCreatedEntityIdInHeaders() {
		returns = append(returns, m.EntityPathKey.Type.ZeroValueReference())
	}
	returns = append(returns, Err())

	m.callResourcePath(def)
	IfErrReturn(def, returns...).Line()
	r.callFormatQueryUrl(def)
	IfErrReturn(def, returns...).Line()

	def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver).Dot("JsonPostRequest").Call(Id(ContextVar), Id(UrlVar), RestLiMethod(protocol.Method_create), Id(CreateParam))
	IfErrReturn(def, returns...).Line()

	def.List(Id(ResVar), Err()).Op(":=").Id(ClientReceiver).Dot(DoAndIgnore).Call(Id(ReqVar))
	IfErrReturn(def, returns...).Line()

	def.If(Id(ResVar).Dot("StatusCode").Op("/").Lit(100).Op("!=").Lit(2)).BlockFunc(func(def *Group) {
		def.Err().Op("=").Qual("fmt", "Errorf").Call(Lit("Invalid response code from %s: %d"), Id(UrlVar), Id(ResVar).Dot("StatusCode"))
		def.Return(returns...)
	}).Line()

	if m.isCreatedEntityIdInHeaders() {
		accessor := Id(m.EntityPathKey.Name)
		def.Var().Add(accessor).Add(m.EntityPathKey.Type.GoType())

		def.Err().Op("=").Add(m.EntityPathKey.Type.RestLiReducedDecodeModel(
			Id(ResVar).Dot("Header").Dot("Get").Call(Qual(ProtocolPackage, RestLiHeaderID)),
			accessor,
		))

		IfErrReturn(def, returns...)

		if m.EntityPathKey.Type.ShouldReference() {
			accessor = Op("&").Add(accessor)
		}

		def.Return(accessor, Nil())
	} else {
		def.Return(Nil())
	}
}

func generateUpdate(r *Resource, m *Method, def *Group) {
	m.callResourcePath(def)
	IfErrReturn(def, Err()).Line()
	r.callFormatQueryUrl(def)
	IfErrReturn(def, Err()).Line()

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
	m.callResourcePath(def)
	IfErrReturn(def, Err()).Line()
	r.callFormatQueryUrl(def)
	IfErrReturn(def, Err()).Line()

	def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver).Dot("JsonPostRequest").Call(
		Id(ContextVar),
		Id(UrlVar),
		RestLiMethod(protocol.Method_partial_update),
		Op("&").Struct(
			Id("Patch").Add(Op("*").Add(r.ResourceSchema.Record().PartialUpdateStruct())).Tag(JsonFieldTag("patch", false)),
		).Values(Dict{Id("Patch"): Id(UpdateParam)}),
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
	m.callResourcePath(def)
	IfErrReturn(def, Err()).Line()
	r.callFormatQueryUrl(def)
	IfErrReturn(def, Err()).Line()

	def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver).Dot("DeleteRequest").Call(Id(ContextVar), Id(UrlVar), RestLiMethod(protocol.Method_update))
	IfErrReturn(def, Err()).Line()

	def.List(Id(ResVar), Err()).Op(":=").Id(ClientReceiver).Dot(DoAndIgnore).Call(Id(ReqVar))
	IfErrReturn(def, Err()).Line()

	def.If(Id(ResVar).Dot("StatusCode").Op("/").Lit(100).Op("!=").Lit(2)).BlockFunc(func(def *Group) {
		def.Return(Qual("fmt", "Errorf").Call(Lit("Invalid response code from %s: %d"), Id(UrlVar), Id(ResVar).Dot("StatusCode")))
	})
	def.Return(Nil())
}
