package resources

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/internal/codegen/types"
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
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
	case protocol.Method_create:
		params = append(params, CreateParam)
	case protocol.Method_update, protocol.Method_partial_update:
		params = append(params, UpdateParam)
	case protocol.Method_batch_get:
		params = append(params, Keys)
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
	case protocol.Method_batch_get:
		params = append(params, Index().Add(r.EntityPathKey.Type.ReferencedType()))
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
	case protocol.Method_batch_get:
		return []Code{Add(Entities).Add(r.batchGetReturnType())}
	default:
		return nil
	}
}

func (r *RestMethod) batchGetReturnType() Code {
	return Map(r.EntityPathKey.Type.ReferencedType()).Add(r.Return.ReferencedType())
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
	switch r.restLiMethod() {
	case protocol.Method_get:
		return r.generateGet
	case protocol.Method_create:
		return r.generateCreate
	case protocol.Method_update:
		return r.genericMethodImplementation("DoUpdateRequest", UpdateParam)
	case protocol.Method_partial_update:
		return r.generatePartialUpdate
	case protocol.Method_delete:
		return r.genericMethodImplementation("DoDeleteRequest")
	case protocol.Method_batch_get:
		return r.generateBatchGet
	default:
		return nil
	}
}

// https://linkedin.github.io/rest.li/user_guide/restli_server#resource-methods
func (r *RestMethod) GenerateCode() *utils.CodeFile {
	c := r.Resource.NewCodeFile(r.Name)

	if len(r.Params) > 0 {
		p := &types.Record{
			NamedType: types.NamedType{
				Identifier: utils.Identifier{
					Name:      r.queryParamsStructName(),
					Namespace: r.Resource.Namespace,
				},
				Doc: fmt.Sprintf("This struct provides the parameters to the %s method", r.Name),
			},
			Fields: r.Params,
		}
		c.Code.Add(p.GenerateStruct()).Line().Line()
		c.Code.Add(p.GenerateQueryParamMarshaler(nil, r.isBatch())).Line().Line()
	}

	r.Resource.addClientFuncDeclarations(c.Code, ClientType, r, func(def *Group) {
		r.generator()(def)
	})

	return c
}

func (r *RestMethod) genericMethodImplementation(doFuncName string, args ...Code) func(*Group) {
	return func(def *Group) {
		formatQueryUrl(r, def, nil, Err())

		def.Return(Id(ClientReceiver).Dot(doFuncName).Call(append([]Code{Ctx, Url}, args...)...))
	}
}
