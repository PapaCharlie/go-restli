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

var batchMethodsWithMapInput = map[protocol.RestLiMethod]bool{
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

func (r *RestMethod) NonErrorFuncReturnParams() []Code {
	switch r.restLiMethod() {
	case protocol.Method_get:
		return []Code{Add(Entity).Add(r.Return.ReferencedType())}
	case protocol.Method_create:
		returns := []Code{Add(EntityKey).Add(r.EntityPathKey.Type.ReferencedType())}
		if r.ReturnEntity {
			returns = append(returns, Add(ReturnedEntity).Add(r.Return.ReferencedType()))
		}
		return returns
	case protocol.Method_batch_create:
		returns := Add(CreatedEntities).Index().Op("*")
		if r.ReturnEntity {
			returns.Qual(utils.ProtocolPackage, "CreatedAndReturnedEntity").Index(List(r.EntityPathKey.Type.ReferencedType(), r.Return.ReferencedType()))
		} else {
			returns.Qual(utils.ProtocolPackage, "CreatedEntity").Index(r.EntityPathKey.Type.ReferencedType())
		}
		return []Code{returns}
	case protocol.Method_batch_get:
		return []Code{Add(Entities).Add(Map(r.EntityPathKey.Type.ReferencedType()).Add(r.Return.ReferencedType()))}
	case protocol.Method_batch_delete, protocol.Method_batch_update, protocol.Method_batch_partial_update:
		return []Code{Add(Statuses).Add(Map(r.EntityPathKey.Type.ReferencedType()).Op("*").Add(BatchEntityUpdateResponse))}
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

func (r *RestMethod) usesBatchMapInput() bool {
	return batchMethodsWithMapInput[r.restLiMethod()]
}

func (r *RestMethod) queryParamsStructName() string {
	return restMethodFuncNames[r.restLiMethod()] + "Params"
}

func (r *RestMethod) generator() func(*Group) {
	switch r.restLiMethod() {
	case protocol.Method_get:
		return r.genericMethodImplementation(
			"DoGetRequest",
			[]Code{r.Return.ZeroValueReference()},
			types.Reader.UnmarshalerFunc(*r.Return),
		)
	case protocol.Method_create:
		f := "DoCreateRequest"
		returns := []Code{r.EntityPathKey.Type.ZeroValueReference()}
		params := []Code{
			Entity,
			r.Resource.readOnlyFields(),
			types.Reader.UnmarshalerFunc(r.EntityPathKey.Type),
		}
		if r.ReturnEntity {
			f += "WithReturnEntity"
			returns = append(returns, r.Return.ZeroValueReference())
			params = append(params, types.Reader.UnmarshalerFunc(*r.Return))
		}
		return r.genericMethodImplementation(f, returns, params...)
	case protocol.Method_update:
		return r.genericMethodImplementation("DoUpdateRequest", nil, Entity)
	case protocol.Method_partial_update:
		return r.genericMethodImplementation(
			"DoPartialUpdateRequest",
			nil,
			Entity,
			r.Resource.createAndReadOnlyFields(),
		)
	case protocol.Method_delete:
		return r.genericMethodImplementation("DoDeleteRequest", nil)
	case protocol.Method_batch_create:
		f := "DoBatchCreateRequest"
		params := []Code{
			Entities,
			r.Resource.readOnlyFields(),
			types.Reader.UnmarshalerFunc(r.EntityPathKey.Type),
		}
		if r.ReturnEntity {
			f += "WithReturnEntity"
			params = append(params, types.Reader.UnmarshalerFunc(*r.Return))
		}
		return r.genericMethodImplementation(
			f,
			[]Code{Nil()},
			params...,
		)
	case protocol.Method_batch_delete:
		return r.genericMethodImplementation(
			"DoBatchDeleteRequest",
			[]Code{Nil()},
			utils.BatchKeySet,
		)
	case protocol.Method_batch_get:
		return r.genericMethodImplementation(
			"DoBatchGetRequest",
			[]Code{Nil()},
			utils.BatchKeySet,
			types.Reader.UnmarshalerFunc(*r.Return),
		)
	case protocol.Method_batch_update:
		return r.genericMethodImplementation(
			"DoBatchUpdateRequest",
			[]Code{Nil()},
			utils.BatchKeySet,
			Entities,
			r.Resource.createAndReadOnlyFields(),
		)
	case protocol.Method_batch_partial_update:
		return r.genericMethodImplementation(
			"DoBatchPartialUpdateRequest",
			[]Code{Nil()},
			utils.BatchKeySet,
			Entities,
			r.Resource.createAndReadOnlyFields(),
		)
	default:
		return nil
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

func (r *RestMethod) genericMethodImplementation(doFuncName string, extraReturnValues []Code, args ...Code) func(*Group) {
	return func(def *Group) {
		returns := append([]Code(nil), extraReturnValues...)
		returns = append(returns, Err())
		formatQueryUrl(r, def, returns...)

		def.Return(Qual(utils.ProtocolPackage, doFuncName).Call(
			append([]Code{RestLiClientReceiver, Ctx, Url}, args...)...,
		))
	}
}
