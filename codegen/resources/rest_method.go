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

var batchMethods = map[protocol.RestLiMethod]bool{
	protocol.Method_batch_get:            true,
	protocol.Method_batch_create:         true,
	protocol.Method_batch_delete:         true,
	protocol.Method_batch_update:         true,
	protocol.Method_batch_partial_update: true,
}

type RestMethod struct{ methodImplementation }

func (r *RestMethod) IsSupported() bool {
	return r.generator() != nil
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
		params = append(params, r.Resource.ResourceSchema.ReferencedType())
	case protocol.Method_partial_update:
		params = append(params, Op("*").Add(r.Resource.ResourceSchema.Record().PartialUpdateStruct()))
	case protocol.Method_batch_create:
		params = append(params, Index().Add(r.Return.ReferencedType()))
	case protocol.Method_batch_delete, protocol.Method_batch_get:
		params = append(params, Index().Add(r.EntityPathKey.Type.ReferencedType()))
	case protocol.Method_batch_update:
		params = append(params, Map(r.EntityPathKey.Type.ReferencedType()).Add(r.Return.ReferencedType()))
	case protocol.Method_batch_partial_update:
		params = append(params, Map(r.EntityPathKey.Type.ReferencedType()).Add(Op("*").Add(r.Resource.ResourceSchema.Record().PartialUpdateStruct())))
	}
	if len(r.Params) > 0 {
		params = append(params, Op("*").Qual(r.Resource.PackagePath(), r.queryParamsStructName()))
	}
	return params
}

func (r *RestMethod) NonErrorFuncReturnParam() Code {
	switch r.restLiMethod() {
	case protocol.Method_get:
		return Add(Entity).Add(r.Return.ReferencedType())
	case protocol.Method_create, protocol.Method_batch_create:
		var returns Code
		if r.ReturnEntity {
			returns = Qual(utils.ProtocolPackage, "CreatedAndReturnedEntity").Index(List(r.EntityPathKey.Type.ReferencedType(), r.Return.ReferencedType()))
		} else {
			returns = Qual(utils.ProtocolPackage, "CreatedEntity").Index(r.EntityPathKey.Type.ReferencedType())
		}
		if r.restLiMethod() == protocol.Method_batch_create {
			returns = Add(CreatedEntities).Index().Op("*").Add(returns)
		} else {
			returns = Add(CreatedEntity).Op("*").Add(returns)
		}
		return returns
	case protocol.Method_batch_get:
		return Add(Entities).Add(Map(r.EntityPathKey.Type.ReferencedType()).Add(r.Return.ReferencedType()))
	case protocol.Method_batch_delete, protocol.Method_batch_update, protocol.Method_batch_partial_update:
		return Add(Statuses).Add(Map(r.EntityPathKey.Type.ReferencedType()).Op("*").Add(BatchEntityUpdateResponse))
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

func (r *RestMethod) isBatch() bool {
	return batchMethods[r.restLiMethod()]
}

func (r *RestMethod) queryParamsStructName() string {
	return restMethodFuncNames[r.restLiMethod()] + "Params"
}

func (r *RestMethod) generator() func(*Group) {
	return func(def *Group) {
		params := r.FuncParamNames()

		if len(r.Params) == 0 {
			params = append(params, Nil())
		} else {
			def.If(Add(QueryParams).Op("==").Nil()).Block(Return(Nil(), Qual(utils.ProtocolPackage, "NilQueryParams")))
		}

		receiver := Id(ClientReceiver)
		if r.GetResource().IsCollection {
			receiver.Dot(CollectionClient)
		} else {
			receiver.Dot(SimpleClient)
		}

		f := r.FuncName()
		if r.ReturnEntity {
			f += "WithReturnEntity"
		}

		declareRpStruct(r, def)
		def.Return(Add(receiver).Dot(f).Call(
			append([]Code{Ctx, Rp}, params...)...,
		))
	}
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
		var batchKeyType Code
		if r.isBatch() {
			batchKeyType = r.EntityPathKey.Type.GoType()
		}
		c.Code.Add(p.GenerateQueryParamMarshaler(nil, batchKeyType)).Line().Line()
	}

	r.Resource.addClientFuncDeclarations(c.Code, ClientType, r, func(def *Group) {
		r.generator()(def)
	})

	return c
}
