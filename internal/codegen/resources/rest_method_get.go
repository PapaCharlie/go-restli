package resources

import (
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

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
