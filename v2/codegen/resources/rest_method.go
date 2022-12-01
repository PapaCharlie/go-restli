package resources

import (
	"fmt"
	"log"

	"github.com/PapaCharlie/go-restli/v2/codegen/types"
	"github.com/PapaCharlie/go-restli/v2/codegen/utils"
	"github.com/PapaCharlie/go-restli/v2/restli"
	. "github.com/dave/jennifer/jen"
)

var restMethodFuncNames = map[restli.Method]string{
	restli.Method_get:            "Get",
	restli.Method_get_all:        "GetAll",
	restli.Method_create:         "Create",
	restli.Method_delete:         "Delete",
	restli.Method_update:         "Update",
	restli.Method_partial_update: "PartialUpdate",

	restli.Method_batch_get:            "BatchGet",
	restli.Method_batch_create:         "BatchCreate",
	restli.Method_batch_delete:         "BatchDelete",
	restli.Method_batch_update:         "BatchUpdate",
	restli.Method_batch_partial_update: "BatchPartialUpdate",
}

var batchQueryParamMethods = map[restli.Method]bool{
	restli.Method_batch_get:            true,
	restli.Method_batch_delete:         true,
	restli.Method_batch_update:         true,
	restli.Method_batch_partial_update: true,
}

var takesReadOnlyFields = map[restli.Method]bool{
	restli.Method_create:       true,
	restli.Method_batch_create: true,
}

var takesCreateAndReadOnlyFields = map[restli.Method]bool{
	restli.Method_update:               true,
	restli.Method_partial_update:       true,
	restli.Method_batch_update:         true,
	restli.Method_batch_partial_update: true,
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
	case restli.Method_batch_delete, restli.Method_batch_get:
		params = append(params, Keys)
	case restli.Method_create, restli.Method_update, restli.Method_partial_update:
		params = append(params, Entity)
	case restli.Method_batch_create, restli.Method_batch_update, restli.Method_batch_partial_update:
		params = append(params, Entities)
	}
	if r.hasParams() {
		params = append(params, QueryParams)
	}
	return params
}

func (r *RestMethod) FuncParamTypes() (params []Code) {
	switch r.restLiMethod() {
	case restli.Method_create, restli.Method_update:
		params = append(params, r.EntityType())
	case restli.Method_partial_update:
		params = append(params, r.PartialEntityUpdateType())
	case restli.Method_batch_create:
		params = append(params, Index().Add(r.EntityType()))
	case restli.Method_batch_delete, restli.Method_batch_get:
		params = append(params, Index().Add(r.EntityKeyType()))
	case restli.Method_batch_update:
		params = append(params, Map(r.EntityKeyType()).Add(r.EntityType()))
	case restli.Method_batch_partial_update:
		params = append(params, Map(r.EntityKeyType()).Add(r.PartialEntityUpdateType()))
	}
	if r.hasParams() {
		params = append(params, Op("*").Qual(r.Resource.PackagePath(), r.queryParamsStructName()))
	}
	return params
}

func (r *RestMethod) NonErrorFuncReturnParam() Code {
	switch r.restLiMethod() {
	case restli.Method_get:
		return Add(Entity).Add(r.EntityType())
	case restli.Method_get_all:
		return Add(Results).Op("*").Add(r.Resource.LocalType(Elements))
	case restli.Method_create, restli.Method_batch_create:
		returns := Op("*")
		if r.ReturnEntity {
			returns.Add(r.Resource.LocalType(CreatedAndReturnedEntity))
		} else {
			returns.Add(r.Resource.LocalType(CreatedEntity))
		}
		if r.restLiMethod() == restli.Method_batch_create {
			return Id("createdEntities").Index().Add(returns)
		} else {
			return Id("createdEntity").Add(returns)
		}
	case restli.Method_batch_get:
		return Add(Results).Op("*").Add(r.Resource.LocalType(BatchEntities))
	case restli.Method_batch_delete, restli.Method_batch_update, restli.Method_batch_partial_update:
		return Add(Results).Op("*").Add(r.Resource.LocalType(BatchResponse))
	case restli.Method_partial_update:
		if r.ReturnEntity {
			return Id("updatedEntity").Add(r.EntityType())
		} else {
			return nil
		}
	default:
		return nil
	}
}

func (r *RestMethod) GenericParams() Code {
	switch r.restLiMethod() {
	case restli.Method_get, restli.Method_update, restli.Method_get_all:
		return r.EntityType()
	case restli.Method_partial_update:
		if r.ReturnEntity {
			return List(r.PartialEntityUpdateType(), r.EntityType())
		} else {
			return r.PartialEntityUpdateType()
		}
	case restli.Method_create, restli.Method_batch_create:
		if r.ReturnEntity {
			return List(r.EntityKeyType(), r.EntityType())
		} else {
			return r.EntityKeyType()
		}
	case restli.Method_batch_get:
		return List(r.EntityKeyType(), r.EntityType())
	case restli.Method_batch_delete, restli.Method_batch_update, restli.Method_batch_partial_update:
		return r.EntityKeyType()
	default:
		return nil
	}
}

func (r *RestMethod) restLiMethod() restli.Method {
	method, ok := restli.MethodNameMapping[r.Name]
	if !ok {
		log.Panicf("Unknown restli method: %s", r.Name)
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

	if r.hasParams() {
		var errReturns []Code
		if r.NonErrorFuncReturnParam() != nil {
			errReturns = append(errReturns, Nil())
		}
		errReturns = append(errReturns, Qual(utils.RestLiPackage, "NilQueryParams"))
		def.If(Add(QueryParams).Op("==").Nil()).Block(Return(errReturns...))
	} else {
		params = append(params, Nil())
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

	call := Qual(utils.RestLiPackage, f)
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

	return Qual(utils.RestLiPackage, name).CallFunc(func(def *Group) {
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

	if r.hasParams() {
		p := &types.Record{
			NamedType: types.NamedType{
				Identifier: utils.Identifier{
					Name:      r.queryParamsStructName(),
					Namespace: r.Resource.Namespace,
				},
				Doc: fmt.Sprintf("%s provides the parameters to the %s method", r.queryParamsStructName(), r.Name),
			},
			Fields:   r.Params,
			Includes: r.includes(),
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
