package resources

import (
	. "github.com/dave/jennifer/jen"
)

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

	def.Err().Op("=").Id(ClientReceiver).Dot("DoCreateRequest").Call(
		Ctx,
		Url,
		CreateParam,
		r.Resource.readOnlyFields(),
		idUnmarshaler,
		returnEntityUnmarshaler,
	)

	if r.EntityPathKey.Type.ShouldReference() {
		id = Op("&").Add(id)
	}

	if r.ReturnEntity {
		def.Return(id, ReturnedEntity, Err())
	} else {
		def.Return(id, Err())
	}
}
