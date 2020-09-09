package resources

import (
	. "github.com/dave/jennifer/jen"
)

func (r *RestMethod) generatePartialUpdate(def *Group) {
	formatQueryUrl(r, def, nil, Err())

	def.Return(Id(ClientReceiver).Dot("DoPartialUpdateRequest").Call(
		Ctx,
		Url,
		UpdateParam,
		r.Resource.createAndReadOnlyFields(),
	))
}
