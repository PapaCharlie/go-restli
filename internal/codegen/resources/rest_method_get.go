package resources

import (
	. "github.com/dave/jennifer/jen"
)

func (r *RestMethod) generateGet(def *Group) {
	formatQueryUrl(r, def, nil, r.Return.ZeroValueReference(), Err())

	var result Code
	if r.Return.ShouldReference() {
		def.Add(Entity).Op("=").New(r.Return.GoType())
		result = Entity
	} else {
		result = Op("&").Add(result)
	}

	def.Err().Op("=").Id(ClientReceiver).Dot("DoGetRequest").Call(Ctx, Url, Entity)
	def.Return(Entity, Err())
}
