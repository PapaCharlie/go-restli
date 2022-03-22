package types

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	"github.com/PapaCharlie/go-restli/restli/batchkeyset"
	. "github.com/dave/jennifer/jen"
)

func AddUnmarshalRestli(def *Statement, receiver, typeName string, f func(def *Group)) *Statement {
	utils.AddFuncOnReceiver(def, receiver, typeName, utils.UnmarshalRestLi, utils.Yes).
		Params(ReaderParam).
		Params(Err().Error()).
		BlockFunc(f).
		Line().Line()

	data := Id("data")
	const unmarshalJSON = "UnmarshalJSON"
	utils.AddFuncOnReceiver(def, receiver, typeName, unmarshalJSON, utils.Yes).
		Params(Add(data).Index().Byte()).
		Params(Error()).
		BlockFunc(func(def *Group) {
			def.Return(Qual(utils.RestLiCodecPackage, unmarshalJSON).Call(data, Id(receiver)))
		}).Line().Line()

	return def
}

func (r *Record) GenerateUnmarshalRestLi() *Statement {
	requiredFields, def := r.generateRequiredFields(nil)
	return AddUnmarshalRestli(def, r.Receiver(), r.Name, func(def *Group) {
		r.generateUnmarshaler(def, requiredFields, nil)
	})
}

func (r *Record) generateUnmarshaler(def *Group, requiredFields Code, batchKeyType *RestliType) {
	if len(r.Fields) == 0 {
		def.Return(Reader.ReadMap(Reader, func(reader Code, field Code, def *Group) {
			def.Return(Reader.Skip(reader))
		}))
		return
	}

	ids := Code(Id(batchkeyset.EntityIDsField))
	def.Err().Op("=").Add(Reader.ReadRecord(Reader, requiredFields, func(reader, field Code, def *Group) {
		def.Switch(field).BlockFunc(func(def *Group) {
			if batchKeyType != nil {
				def.Case(Lit(batchkeyset.EntityIDsField)).BlockFunc(func(def *Group) {
					def.Add(Reader.Read(RestliType{Array: batchKeyType}, reader, ids))
				})
			}
			for _, f := range r.Fields {
				def.Case(Lit(f.Name)).BlockFunc(func(def *Group) {
					accessor := r.fieldAccessor(f)

					if f.IsOptionalOrDefault() {
						def.Add(accessor).Op("=").New(f.Type.GoType())
						if f.Type.Reference == nil {
							accessor = Op("*").Add(accessor)
						}
					}

					def.Add(Reader.Read(f.Type, reader, accessor))
				})
			}
			def.Default().BlockFunc(func(def *Group) {
				def.Err().Op("=").Add(Reader.Skip(reader))
			})
		})
		def.Return(Err())
	}))
	if batchKeyType != nil {
		def.Add(utils.IfErrReturn(Nil(), Err())).Line()
	} else {
		def.Add(utils.IfErrReturn(Err())).Line()
	}

	if r.hasDefaultValue() {
		def.Id(r.Receiver()).Dot(utils.PopulateLocalDefaultValues).Call()
	}

	if batchKeyType != nil {
		def.Return(ids, Err())
	} else {
		def.Return(Err())
	}
}

func (r *Record) GenerateQueryParamUnmarshaler(batchKeyType *RestliType) Code {
	requiredFields, def := r.generateRequiredFields(batchKeyType)
	ids := Code(Id(batchkeyset.EntityIDsField))

	var params []Code
	if batchKeyType != nil {
		params = append(params, Add(ids).Index().Add(batchKeyType.GoType()))
	}
	params = append(params, Err().Error())

	return utils.AddFuncOnReceiver(def, r.Receiver(), r.Name, "DecodeQueryParams", RecordShouldUsePointer).
		Params(Add(Reader).Qual(utils.RestLiCodecPackage, "QueryParamsReader")).
		Params(params...).
		BlockFunc(func(def *Group) {
			r.generateUnmarshaler(def, requiredFields, batchKeyType)
		})
}

func (r *Record) generateRequiredFields(batchKeyType *RestliType) (requiredFields Code, def *Statement) {
	hasRequiredFields := batchKeyType != nil
	for _, f := range r.Fields {
		if !f.IsOptionalOrDefault() {
			hasRequiredFields = true
			break
		}
	}

	if hasRequiredFields {
		requiredFields = Code(Id("_" + r.Name + "RequiredFields"))
		def = Var().Add(requiredFields).Op("=").Add(utils.RequiredFields).ValuesFunc(func(def *Group) {
			for _, f := range r.Fields {
				if !f.IsOptionalOrDefault() {
					def.Line().Lit(f.Name)
				}
			}
			if batchKeyType != nil {
				def.Line().Lit(batchkeyset.EntityIDsField)
			}
			def.Line()
		}).Line().Line()
	} else {
		requiredFields = Nil()
		def = Empty()
	}

	return requiredFields, def
}
