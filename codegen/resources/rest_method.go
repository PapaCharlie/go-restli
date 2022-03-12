package resources

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/PapaCharlie/go-restli/codegen/utils"
	"github.com/PapaCharlie/go-restli/protocol"
	. "github.com/dave/jennifer/jen"
)

var restMethodFuncNames = map[protocol.RestLiMethod]string{
	protocol.Method_get:            "Get",
	protocol.Method_get_all:        "GetAll",
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

var batchQueryParamMethods = map[protocol.RestLiMethod]bool{
	protocol.Method_batch_get:            true,
	protocol.Method_batch_delete:         true,
	protocol.Method_batch_update:         true,
	protocol.Method_batch_partial_update: true,
}

var takesReadOnlyFields = map[protocol.RestLiMethod]bool{
	protocol.Method_create:       true,
	protocol.Method_batch_create: true,
}

var takesCreateAndReadOnlyFields = map[protocol.RestLiMethod]bool{
	protocol.Method_update:               true,
	protocol.Method_partial_update:       true,
	protocol.Method_batch_update:         true,
	protocol.Method_batch_partial_update: true,
}

type RestMethod struct{ methodImplementation }

func (r *RestMethod) EntityType() *Statement {
	return r.Resource.ResourceSchema.ReferencedType()
}

func (r *RestMethod) PartialEntityUpdateType() *Statement {
	return Op("*").Add(r.Resource.ResourceSchema.Record().PartialUpdateStruct())
}

func (r *RestMethod) EntityKeyType() *Statement {
	return r.Resource.LastSegment().PathKey.GoType()
}

func (r *RestMethod) FuncName() string {
	return restMethodFuncNames[r.restLiMethod()]
}

func (r *RestMethod) FuncParamNames() (params []Code) {
	switch r.restLiMethod() {
	case protocol.Method_batch_delete, protocol.Method_batch_get:
		params = append(params, Keys)
	case protocol.Method_create, protocol.Method_update, protocol.Method_partial_update:
		params = append(params, Entity)
	case protocol.Method_batch_create, protocol.Method_batch_update, protocol.Method_batch_partial_update:
		params = append(params, Entities)
	}
	if len(r.Params) > 0 {
		params = append(params, QueryParams)
	}
	return params
}

func (r *RestMethod) FuncParamTypes() (params []Code) {
	switch r.restLiMethod() {
	case protocol.Method_create, protocol.Method_update:
		params = append(params, r.EntityType())
	case protocol.Method_partial_update:
		params = append(params, r.PartialEntityUpdateType())
	case protocol.Method_batch_create:
		params = append(params, Index().Add(r.EntityType()))
	case protocol.Method_batch_delete, protocol.Method_batch_get:
		params = append(params, Index().Add(r.EntityKeyType()))
	case protocol.Method_batch_update:
		params = append(params, Map(r.EntityKeyType()).Add(r.EntityType()))
	case protocol.Method_batch_partial_update:
		params = append(params, Map(r.EntityKeyType()).Add(r.PartialEntityUpdateType()))
	}
	if len(r.Params) > 0 {
		params = append(params, Op("*").Qual(r.Resource.PackagePath(), r.queryParamsStructName()))
	}
	return params
}

func (r *RestMethod) NonErrorFuncReturnParam() Code {
	switch r.restLiMethod() {
	case protocol.Method_get:
		return Add(Entity).Add(r.EntityType())
	case protocol.Method_get_all:
		return Add(Results).Op("*").Add(r.Resource.LocalType(Elements))
	case protocol.Method_create, protocol.Method_batch_create:
		returns := Op("*")
		if r.ReturnEntity {
			returns.Add(r.Resource.LocalType(CreatedAndReturnedEntity))
		} else {
			returns.Add(r.Resource.LocalType(CreatedEntity))
		}
		if r.restLiMethod() == protocol.Method_batch_create {
			return Id("createdEntities").Index().Add(returns)
		} else {
			return Id("createdEntity").Add(returns)
		}
	case protocol.Method_batch_get:
		return Add(Results).Op("*").Add(r.Resource.LocalType(BatchEntities))
	case protocol.Method_batch_delete, protocol.Method_batch_update, protocol.Method_batch_partial_update:
		return Add(Results).Op("*").Add(r.Resource.LocalType(BatchResponse))
	default:
		return nil
	}
}

func (r *RestMethod) GenericParams() Code {
	switch r.restLiMethod() {
	case protocol.Method_get, protocol.Method_update, protocol.Method_get_all:
		return r.EntityType()
	case protocol.Method_partial_update:
		return r.PartialEntityUpdateType()
	case protocol.Method_create, protocol.Method_batch_create:
		if r.ReturnEntity {
			return List(r.EntityKeyType(), r.EntityType())
		} else {
			return r.EntityKeyType()
		}
	case protocol.Method_batch_get:
		return List(r.EntityKeyType(), r.EntityType())
	case protocol.Method_batch_delete, protocol.Method_batch_update, protocol.Method_batch_partial_update:
		return r.EntityKeyType()
	default:
		return nil
	}
}

func (r *RestMethod) restLiMethod() protocol.RestLiMethod {
	method, ok := protocol.RestLiMethodNameMapping[r.Name]
	if !ok {
		utils.Logger.Panicf("Unknown restli method: %s", r.Name)
	}
	return method
}

func (r *RestMethod) usesBatchQueryParams() bool {
	return batchQueryParamMethods[r.restLiMethod()]
}

func (r *RestMethod) queryParamsStructName() string {
	return restMethodFuncNames[r.restLiMethod()] + "Params"
}

func (r *RestMethod) clientMethodGenerator(def *Group) {
	params := r.FuncParamNames()

	if len(r.Params) == 0 {
		params = append(params, Nil())
	} else {
		def.If(Add(QueryParams).Op("==").Nil()).Block(Return(Nil(), Qual(utils.ProtocolPackage, "NilQueryParams")))
	}

	f := r.FuncName()
	if r.ReturnEntity {
		f += "WithReturnEntity"
	}

	if takesReadOnlyFields[r.restLiMethod()] {
		params = append(params, r.Resource.readOnlyFields())
	} else if takesCreateAndReadOnlyFields[r.restLiMethod()] {
		params = append(params, r.Resource.createAndReadOnlyFields())
	}

	call := Qual(utils.ProtocolPackage, f)
	if p := r.GenericParams(); p != nil {
		call.Index(p)
	}

	declareRpStruct(r, def)
	def.Return(call.Call(
		append([]Code{RestLiClientReceiver, Ctx, Rp}, params...)...,
	))
}

func (r *RestMethod) RegisterMethod(server, resource, segments Code) Code {
	name := "Register" + r.FuncName()
	if r.ReturnEntity {
		name += "WithReturnEntity"
	}

	return Qual(utils.ProtocolPackage, name).CallFunc(func(def *Group) {
		def.Add(server)
		def.Add(segments)

		if takesReadOnlyFields[r.restLiMethod()] {
			def.Add(r.Resource.readOnlyFields())
		} else if takesCreateAndReadOnlyFields[r.restLiMethod()] {
			def.Add(r.Resource.createAndReadOnlyFields())
		}

		def.Line().Func().Params(registerParams(r)...).Params(methodReturnParams(r)...).
			BlockFunc(func(def *Group) {
				def.Return(resource).Dot(r.FuncName()).Call(splatRpAndParams(r)...)
			})
	})
}

// https://linkedin.github.io/rest.li/user_guide/restli_server#resource-methods
func (r *RestMethod) GenerateCode() *utils.CodeFile {
	c := r.Resource.NewCodeFile(r.Name)

	if len(r.Params) > 0 || r.PagingSupported {
		p := &types.Record{
			NamedType: types.NamedType{
				Identifier: utils.Identifier{
					Name:      r.queryParamsStructName(),
					Namespace: r.Resource.Namespace,
				},
				Doc: fmt.Sprintf("%s provides the parameters to the %s method", r.queryParamsStructName(), r.Name),
			},
			Fields: r.Params,
		}
		if r.PagingSupported {
			addPagingContextFields(p)
		}
		c.Code.Add(p.GenerateStruct()).Line().Line()
		var batchKeyType *types.RestliType
		if r.usesBatchQueryParams() {
			batchKeyType = &r.Resource.LastSegment().PathKey.Type
		}
		c.Code.
			Add(p.GenerateQueryParamMarshaler(nil, batchKeyType)).Line().Line().
			Add(p.GenerateQueryParamUnmarshaler(batchKeyType)).Line().Line().
			Add(p.GeneratePopulateDefaultValues()).Line().Line()
	}

	r.Resource.addClientFuncDeclarations(c.Code, ClientType, r, func(def *Group) {
		r.clientMethodGenerator(def)
	})

	return c
}
