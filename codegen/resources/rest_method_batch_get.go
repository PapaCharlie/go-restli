package resources

import (
	"github.com/PapaCharlie/go-restli/codegen/types"
	. "github.com/dave/jennifer/jen"
)

func (r *RestMethod) batchGetReturnType() Code {
	return Map(r.EntityPathKey.Type.ReferencedType()).Add(r.Return.ReferencedType())
}

func (r *RestMethod) generateBatchGet(def *Group) {
	returns := []Code{Nil(), Err()}
	key := Code(Id("key"))

	formatQueryUrl(r, def, func(def *Group) {
		def.For(List(Id("_"), key).Op(":=").Range().Add(Keys)).BlockFunc(func(def *Group) {
			r.addEntityId(def, key, returns)
		})
	}, returns...)

	def.Add(Entities).Op("=").Make(r.batchGetReturnType())
	resultsReader := r.batchMethodBoilerplate(def, func(def *Group, keyAccessor, valueReader Code) {
		accessor := Add(Entities).Index(keyAccessor)
		if r.Return.ShouldReference() {
			def.Add(accessor).Op("=").New(r.Return.GoType())
		}
		def.Add(types.Reader.Read(*r.Return, valueReader, accessor))
		def.Return(Err())
	})
	def.Err().Op("=").Id(ClientReceiver).Dot("DoBatchGetRequest").Call(Ctx, Url, resultsReader).Line()
	def.Return(Entities, Err())
}
