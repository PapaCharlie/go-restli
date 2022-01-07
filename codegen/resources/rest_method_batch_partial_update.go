package resources

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func (r *RestMethod) batchPartialUpdateReturnType() Code {
	return Map(r.EntityPathKey.Type.ReferencedType()).Op("*").Add(BatchEntityUpdateResponse)
}

func (r *RestMethod) generateBatchPartialUpdate(def *Group) {
	returns := []Code{Nil(), Err()}

	formatQueryUrl(r, def, func(def *Group) {
		def.For(Add(Key).Op(":=").Range().Add(Entities)).BlockFunc(func(def *Group) {
			r.addEntityId(def, Key, returns)
		})
	}, returns...)

	entityMap := Code(Id("entityMap"))
	def.Add(entityMap).Op(":=").Make(BatchEntities)
	def.For(List(Key, Entity).Op(":=").Range().Add(Entities)).BlockFunc(func(def *Group) {
		var keyMarshaler Code
		if p := r.EntityPathKey.Type.Primitive; p != nil {
			keyMarshaler = p.NewPrimitiveMarshaler(Key)
		} else {
			keyMarshaler = Key
		}

		def.Add(Err().Op("=").Add(entityMap).Dot("Add").Call(keyMarshaler, Entity))
		def.Add(utils.IfErrReturn(returns...))
	}).Line()

	def.Add(Statuses).Op("=").Make(r.batchPartialUpdateReturnType())
	resultsReader := r.batchMethodBoilerplate(def, func(def *Group, keyAccessor, valueReader Code) {
		accessor := Add(Statuses).Index(keyAccessor)
		def.Add(accessor).Op("=").New(BatchEntityUpdateResponse)
		def.Add(Return(Add(accessor).Dot(utils.UnmarshalRestLi).Call(valueReader)))
	})

	def.Err().Op("=").Id(ClientReceiver).Dot("DoBatchPartialUpdateRequest").Call(Ctx, Url, entityMap, resultsReader, r.Resource.createAndReadOnlyFields()).Line()
	def.Return(Statuses, Err())
}
