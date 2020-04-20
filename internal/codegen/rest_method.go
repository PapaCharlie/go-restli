package codegen

import (
	"github.com/PapaCharlie/go-restli/protocol"
	. "github.com/dave/jennifer/jen"
)

const (
	CreateParam      = "create"
	CreateResponseId = "createResponseId"
	UpdateParam      = "update"
)

// isCreatedEntityIdInHeaders returns true if the Create method is supposed to parse the created record's ID from the
// Location header in the response
func (m *Method) isCreatedEntityIdInHeaders() bool {
	if m.EntityPathKey == nil {
		return false
	}

	return m.EntityPathKey.Type.Primitive != nil || m.EntityPathKey.Type.PrimitiveTyperef() != nil
}

func (m *Method) RestLiMethod() protocol.RestLiMethod {
	return protocol.RestLiMethodNameMapping[m.Name]
}

func (m *Method) restMethodFuncName() string {
	return ExportedIdentifier(m.Name)
}

func (m *Method) restMethodFuncParams(def *Group, resourceSchema *RestliType) {
	switch m.RestLiMethod() {
	case protocol.Method_get:
		m.addEntityTypes(def)
	case protocol.Method_create:
		m.addEntityTypes(def)
		def.Id(CreateParam).Add(resourceSchema.PointerType())
	case protocol.Method_update:
		m.addEntityTypes(def)
		def.Id(UpdateParam).Add(resourceSchema.PointerType())
	case protocol.Method_delete:
		m.addEntityTypes(def)
	}
}

func (m *Method) restMethodFuncReturnParams(def *Group) {
	switch m.RestLiMethod() {
	case protocol.Method_get:
		def.Add(m.Return.PointerType())
		def.Error()
	case protocol.Method_create:
		if m.isCreatedEntityIdInHeaders() {
			def.Id(CreateResponseId).Add(m.EntityPathKey.Type.GoType())
		}
		def.Err().Error()
	case protocol.Method_update:
		def.Error()
	case protocol.Method_delete:
		def.Error()
	}
}

// https://linkedin.github.io/rest.li/user_guide/restli_server#resource-methods
func (r *Resource) GenerateRestMethodCode(m *Method) *Statement {
	switch m.RestLiMethod() {
	case protocol.Method_get:
		return r.generateGet(m)
	case protocol.Method_create:
		return r.generateCreate(m)
	case protocol.Method_update:
		return r.generateUpdate(m)
	case protocol.Method_delete:
		return r.generateDelete(m)
	default:
		Logger.Printf("Warning: %s method is not currently implemented", m.Name)
		return nil
	}
}

func (m *Method) callResourcePath(def *Group) {
	if m.OnEntity {
		def.List(Id(PathVar), Err()).Op(":=").Id(ResourceEntityPath).Call(m.entityParams()...)
	} else {
		def.List(Id(PathVar), Err()).Op(":=").Id(ResourcePath).Call(m.entityParams()...)
	}
}

func (r *Resource) generateGet(m *Method) *Statement {
	def := Empty()
	r.addClientFunc(def, m)

	def.BlockFunc(func(def *Group) {
		m.callResourcePath(def)
		IfErrReturn(def, Nil(), Err()).Line()
		r.callFormatQueryUrl(def)
		IfErrReturn(def, Nil(), Err()).Line()

		def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver).Dot("GetRequest").Call(Id(UrlVar), RestLiMethod(protocol.Method_get))
		IfErrReturn(def, Nil(), Err()).Line()

		def.Id(DoAndDecodeResult).Op(":=").New(m.Return.GoType())
		callDoAndDecode(def)
		def.Return(Id(DoAndDecodeResult), Err())
	})

	return def
}

func (r *Resource) generateCreate(m *Method) *Statement {
	def := Empty()
	r.addClientFunc(def, m)

	// TODO: Support @ReturnEntity annotation

	def.BlockFunc(func(def *Group) {
		var returns []Code
		if m.isCreatedEntityIdInHeaders() {
			returns = append(returns, Id(CreateResponseId))
		}
		returns = append(returns, Err())

		m.callResourcePath(def)
		IfErrReturn(def, returns...).Line()
		r.callFormatQueryUrl(def)
		IfErrReturn(def, returns...).Line()

		def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver).Dot("JsonPostRequest").Call(Id(UrlVar), RestLiMethod(protocol.Method_create), Id(CreateParam))
		IfErrReturn(def, returns...).Line()

		def.List(Id(ResVar), Err()).Op(":=").Id(ClientReceiver).Dot(DoAndIgnore).Call(Id(ReqVar))
		IfErrReturn(def, returns...).Line()

		def.If(Id(ResVar).Dot("StatusCode").Op("/").Lit(100).Op("!=").Lit(2)).BlockFunc(func(def *Group) {
			def.Err().Op("=").Qual("fmt", "Errorf").Call(Lit("Invalid response code from %s: %d"), Id(UrlVar), Id(ResVar).Dot("StatusCode"))
			def.Return(returns...)
		})

		if m.isCreatedEntityIdInHeaders() {
			def.Err().Op("=").Add(m.EntityPathKey.Type.RestLiReducedDecodeModel(
				Id(ResVar).Dot("Header").Dot("Get").Call(Qual(ProtocolPackage, RestLiHeaderID)),
				Op("&").Id(CreateResponseId),
			))
			IfErrReturn(def, returns...)
			def.Return(Id(CreateResponseId), Nil())
		} else {
			def.Return(Nil())
		}
	})

	return def
}

func (r *Resource) generateUpdate(m *Method) *Statement {
	def := Empty()
	r.addClientFunc(def, m)

	def.BlockFunc(func(def *Group) {
		m.callResourcePath(def)
		IfErrReturn(def, Err()).Line()
		r.callFormatQueryUrl(def)
		IfErrReturn(def, Err()).Line()

		def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver).Dot("JsonPutRequest").Call(Id(UrlVar), RestLiMethod(protocol.Method_update), Id(UpdateParam))
		IfErrReturn(def, Err()).Line()

		def.List(Id(ResVar), Err()).Op(":=").Id(ClientReceiver).Dot(DoAndIgnore).Call(Id(ReqVar))
		IfErrReturn(def, Err()).Line()

		def.If(Id(ResVar).Dot("StatusCode").Op("/").Lit(100).Op("!=").Lit(2)).BlockFunc(func(def *Group) {
			def.Return(Qual("fmt", "Errorf").Call(Lit("Invalid response code from %s: %d"), Id(UrlVar), Id(ResVar).Dot("StatusCode")))
		})
		def.Return(Nil())
	})

	return def
}

func (r *Resource) generateDelete(m *Method) *Statement {
	def := Empty()
	r.addClientFunc(def, m)

	def.BlockFunc(func(def *Group) {
		m.callResourcePath(def)
		IfErrReturn(def, Err()).Line()
		r.callFormatQueryUrl(def)
		IfErrReturn(def, Err()).Line()

		def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver).Dot("DeleteRequest").Call(Id(UrlVar), RestLiMethod(protocol.Method_update))
		IfErrReturn(def, Err()).Line()

		def.List(Id(ResVar), Err()).Op(":=").Id(ClientReceiver).Dot(DoAndIgnore).Call(Id(ReqVar))
		IfErrReturn(def, Err()).Line()

		def.If(Id(ResVar).Dot("StatusCode").Op("/").Lit(100).Op("!=").Lit(2)).BlockFunc(func(def *Group) {
			def.Return(Qual("fmt", "Errorf").Call(Lit("Invalid response code from %s: %d"), Id(UrlVar), Id(ResVar).Dot("StatusCode")))
		})
		def.Return(Nil())
	})

	return def
}
