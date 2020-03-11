package codegen

import (
	"strings"
	"unicode"

	"github.com/bored-engineer/go-restli/protocol"
	. "github.com/dave/jennifer/jen"
)

const UpdateParam = "update"
const PartialUpdateSetParam = "set"
const PartialUpdateDeleteParam = "delete"

func (m *Method) RestLiMethod() protocol.RestLiMethod {
	return protocol.RestLiMethodNameMapping[m.Name]
}

func (m *Method) restMethodFuncName() string {
	// Converts Method_partial_update to MethodPartialUpdate
	name := m.Name
	for {
		idx := strings.IndexRune(name, '_')
		if idx == -1 {
			break
		}
		name = name[:idx] + string(unicode.ToUpper(rune(name[idx+1]))) + name[idx+2:]
	}
	return ExportedIdentifier(name)
}

func (m *Method) restMethodFuncParams(def *Group, resourceSchema *RestliType) {
	switch m.RestLiMethod() {
	case protocol.Method_get:
		m.addEntityTypes(def)
	case protocol.Method_update:
		m.addEntityTypes(def)
		def.Id(UpdateParam).Add(resourceSchema.PointerType())
	case protocol.Method_partial_update:
		m.addEntityTypes(def)
		def.Id(PartialUpdateSetParam).Add(resourceSchema.PointerType())
		def.Id(PartialUpdateDeleteParam).Add(Index().String())
	case protocol.Method_delete:
		m.addEntityTypes(def)
	}
}

func (m *Method) restMethodFuncReturnParams(def *Group) {
	switch m.RestLiMethod() {
	case protocol.Method_get:
		def.Add(m.Return.PointerType())
		def.Error()
	case protocol.Method_update:
		def.Error()
	case protocol.Method_partial_update:
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
	case protocol.Method_update:
		return r.generateUpdate(m)
	case protocol.Method_partial_update:
		return r.generatePartialUpdate(m)
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

func (r *Resource) generatePartialUpdate(m *Method) *Statement {
	def := Empty()
	r.addClientFunc(def, m)

	def.BlockFunc(func(def *Group) {
		m.callResourcePath(def)
		IfErrReturn(def, Err()).Line()
		r.callFormatQueryUrl(def)
		IfErrReturn(def, Err()).Line()

		def.Id(PatchVar).Op(":=").Struct(
			Id("Set").Add(m.Return.PointerType()).Tag(JsonFieldTag("$set", true)),
			Id("Delete").Add(Index().String()).Tag(JsonFieldTag("$delete", true)),
		).Values(Id(PartialUpdateSetParam), Id(PartialUpdateDeleteParam))

		def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver).Dot("JsonPostRequest").Call(Id(UrlVar), RestLiMethod(protocol.Method_partial_update), Op("&").Struct(
			Id("Patch").Add(Interface()).Tag(JsonFieldTag("patch", true)),
		).Values(Id(PatchVar)))
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
