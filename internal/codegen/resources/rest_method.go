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
		params = append(params, KeysParams)
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
		return r.genericMethodImplementation("DoPartialUpdateRequest", UpdateParam)
	case protocol.Method_delete:
		return r.genericMethodImplementation("DoDeleteRequest")
	case protocol.Method_batch_get:
		// TODO: ComplexKeys are sometimes returned without the $params object. Ideally, the batch get would look like
		//  this:
		//    BatchGet(keys []*CK) (map[*CK]Entity, error)
		//  But the map would be such that each *CK key is actually the original pointer given in the keys slice. This
		//  would make it far easier and more consistent to use. A batch get that doesn't behave like this would be
		//  impossible to use because the caller would have to manually cross-reference the map's keys with their own
		//  keys using the generated Equals method. This problem can be solved a few different ways:
		//    1. Generate a "Less" method and make all objects of the same type co-sortable. It's generally convenient
		//       to have this capability and it allows the batch get implementation to binary search for the
		//       original *CK pointer after sorting the given keys slice. The problem is that maps can't be sorted so
		//       this wouldn't work for any complex keys with maps in them (unless an arbitrary map sorting mechanism is
		//       implemented, e.g. based on the number of elements in the map and then which keys appear in the map,
		//       etc...)
		//    2. Generate a "Hash" method and implement a generic hash table that would provide quick lookups for the
		//       original pointer. Definitely feasible but generating this hash function would be non-trivial, though
		//       fascinating! Interestingly a "Hasher" interface would actually look very similar to the existing
		//       Reader/Writer interfaces from the restlicodec package. The other difficulty would be to implement an
		//       actual hashtable from scratch. Nothing new or fancy, just tricky to get perfectly right.
		//    3. A naive approach that simply uses the generated Equals method and does linear lookups on the given keys
		//       slice. This can serve as an initial implementation as it would respect the API's declared behavior, it
		//       just wouldn't be very fast for large batches.
		if r.EntityPathKey.Type.UnderlyingPrimitive() == nil {
			return nil
		} else {
			// TODO (cont.): raw primitives or typerefs to primitives don't get referenced so they work as-is with Go's
			//  built-in map type, and no special effort needs to be made there
			return r.generateBatchGet
		}
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

func (r *RestMethod) generateGet(def *Group) {
	returns := []Code{
		r.Return.ZeroValueReference(),
		Err(),
	}

	formatQueryUrl(r, def, nil, returns...)

	var result Code
	if r.Return.ShouldReference() {
		def.Add(Entity).Op("=").New(r.Return.GoType())
		result = Entity
	} else {
		result = Op("&").Add(result)
	}

	def.Err().Op("=").Id(ClientReceiver).Dot("DoGetRequest").Call(Ctx, Url, Entity)
	def.Add(utils.IfErrReturn(returns...)).Line()

	def.Return(Entity, Nil())
}

func (r *RestMethod) generateCreate(def *Group) {
	returns := []Code{r.EntityPathKey.Type.ZeroValueReference()}
	if r.ReturnEntity {
		returns = append(returns, r.Return.ZeroValueReference())
	}
	returns = append(returns, Err())

	formatQueryUrl(r, def, nil, returns...)

	id := Id(r.EntityPathKey.Name)
	def.Var().Add(id).Add(r.EntityPathKey.Type.GoType())

	var idUnmarshaler Code
	if p := r.EntityPathKey.Type.Primitive; p != nil {
		idUnmarshaler = p.NewPrimitiveUnmarshaler(id)
	} else {
		idUnmarshaler = Op("&").Add(id)
	}

	var returnEntityUnmarshaler Code
	if r.ReturnEntity {
		if r.Return.ShouldReference() {
			def.Add(ReturnedEntity).Op("=").New(r.Return.GoType())
			returnEntityUnmarshaler = ReturnedEntity
		} else {
			returnEntityUnmarshaler = Op("&").Add(ReturnedEntity)
		}
	} else {
		returnEntityUnmarshaler = Nil()
	}

	def.Err().Op("=").Id(ClientReceiver).Dot("DoCreateRequest").Call(Ctx, Url, CreateParam, idUnmarshaler, returnEntityUnmarshaler)
	def.Add(utils.IfErrReturn(returns...)).Line()

	if r.EntityPathKey.Type.ShouldReference() {
		id = Op("&").Add(id)
	}

	if r.ReturnEntity {
		def.Return(id, ReturnedEntity, Nil())
	} else {
		def.Return(id, Nil())
	}
}

func (r *RestMethod) generateBatchGet(def *Group) {
	returns := []Code{Nil(), Err()}

	formatQueryUrl(r, def, func(itemWriter Code, def *Group) {
		item := Id("item")
		def.For(List(Id("_"), item).Op(":=").Range().Add(KeysParams)).BlockFunc(func(def *Group) {
			def.Add(types.Writer.Write(r.EntityPathKey.Type, Add(itemWriter).Call(), item, Err()))
		})
		def.Return(Nil())
	}, returns...)

	def.Add(Entities).Op("=").Make(r.batchGetReturnType())
	rawKey := Id("rawKey")
	resultsReader := Func().Params(Add(types.Reader).Add(types.ReaderQual), Add(rawKey).String()).Params(Err().Error()).BlockFunc(func(def *Group) {
		v := Id("v")
		def.Var().Add(v).Add(r.EntityPathKey.Type.GoType())
		keyReader := Id("keyReader")
		def.List(keyReader, Err()).Op(":=").Add(types.NewRor2Reader).Call(rawKey)
		def.Add(utils.IfErrReturn(Err()))
		def.Add(types.Reader.Read(r.EntityPathKey.Type, keyReader, v))
		def.Add(utils.IfErrReturn(Err())).Line()

		accessor := Add(Entities).Index(v)
		if r.Return.ShouldReference() {
			def.Add(accessor).Op("=").New(r.Return.GoType())
		}
		def.Add(types.Reader.Read(*r.Return, types.Reader, accessor))
		def.Add(Return(Err()))
	})
	def.Err().Op("=").Id(ClientReceiver).Dot("DoBatchGetRequest").Call(Ctx, Url, resultsReader)

	def.Add(utils.IfErrReturn(returns...)).Line()

	def.Return()
}

func (r *RestMethod) genericMethodImplementation(doFuncName string, args ...Code) func(*Group) {
	return func(def *Group) {
		formatQueryUrl(r, def, nil, Err())

		def.Return(Id(ClientReceiver).Dot(doFuncName).Call(append([]Code{Ctx, Url}, args...)...))
	}
}
