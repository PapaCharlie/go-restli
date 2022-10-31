package types

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	"github.com/PapaCharlie/go-restli/restli/batchkeyset"
	. "github.com/dave/jennifer/jen"
)

const NewInstance = "NewInstance"

func AddNewInstance(def *Statement, receiver, typeName string) {
	utils.AddFuncOnReceiver(def, receiver, typeName, NewInstance, utils.Yes).Params().Op("*").Id(typeName).Block(
		Return(New(Id(typeName))),
	).Line().Line()
}

const UnmarshalField = "UnmarshalField"

var Found = Code(Id("found"))

func AddUnmarshalRestli(def *Statement, receiver, typeName string, pointer utils.ShouldUsePointer, f func(def *Group)) *Statement {
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

	if pointer.ShouldUsePointer() {
		AddNewInstance(def, receiver, typeName)
	}

	return def
}

func (r *Record) GenerateUnmarshalField(batchKeyType *RestliType) *Statement {
	params := []Code{ReaderParam, Add(FieldParamName).String()}
	if batchKeyType != nil {
		params = append(params, Add(Ids).Add((&RestliType{Array: batchKeyType}).PointerType()))
	}
	return utils.AddFuncOnReceiver(Empty(), r.Receiver(), r.TypeName(), UnmarshalField, RecordShouldUsePointer).
		Params(params...).
		Params(Add(Found).Bool(), Err().Error()).
		BlockFunc(func(def *Group) {
			if len(r.Fields) == 0 && len(r.Includes) == 0 && batchKeyType == nil {
				def.Return(False(), Nil())
				return
			}

			for _, i := range r.Includes {
				def.List(Found, Err()).Op("=").Id(r.Receiver()).Dot(i.TypeName()).Dot(UnmarshalField).Call(Reader, FieldParamName)
				def.Add(utils.IfErrReturn(Found, Err()))
				def.If(Found).Block(Return(Found, Nil()))
			}
			def.Switch(FieldParamName).BlockFunc(func(def *Group) {
				if batchKeyType != nil {
					def.Case(Lit(batchkeyset.EntityIDsField)).BlockFunc(func(def *Group) {
						def.Add(Found).Op("=").True()
						def.Add(Reader.Read(RestliType{Array: batchKeyType}, Reader, Op("*").Add(Ids)))
					})
				}
				for _, f := range r.Fields {
					def.Case(Lit(f.Name)).BlockFunc(func(def *Group) {
						def.Add(Found).Op("=").True()
						r.readField(def, f, Reader)
					})
				}
			})
			def.Return(Found, Err())
		}).
		Line().Line()
}

func (r *Record) GenerateUnmarshalRestLi() *Statement {
	requiredFields, def := r.generateRequiredFields(nil)
	def.Add(r.GenerateUnmarshalField(nil))
	return AddUnmarshalRestli(def, r.Receiver(), r.TypeName(), RecordShouldUsePointer, func(def *Group) {
		r.generateUnmarshaler(def, requiredFields, nil)
	})
}

func (r *Record) readField(def *Group, f Field, reader Code) {
	accessor := r.fieldAccessor(f)

	if f.IsOptionalOrDefault() {
		def.Add(accessor).Op("=").New(f.Type.GoType())
		if f.Type.Reference == nil || f.Type.IsCustomTyperef() {
			accessor = Op("*").Add(accessor)
		}
	}

	def.Add(Reader.Read(f.Type, reader, accessor))
}

func (r *Record) generateUnmarshaler(def *Group, requiredFields Code, batchKeyType *RestliType) {
	def.Err().Op("=").Add(Reader.ReadRecord(Reader, requiredFields, func(reader, field Code, def *Group) {
		callParams := []Code{reader, field}
		if batchKeyType != nil {
			callParams = append(callParams, Op("&").Add(Ids))
		}
		def.List(Found, Err()).Op(":=").Id(r.Receiver()).Dot(UnmarshalField).Call(callParams...)
		def.Add(utils.IfErrReturn(Err()))
		def.If(Op("!").Add(Found)).Block(Err().Op("=").Add(Reader.Skip(reader)))
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
		def.Return(Ids, Nil())
	} else {
		def.Return(Nil())
	}
}

var Ids = Code(Id(batchkeyset.EntityIDsField))

func (r *Record) GenerateQueryParamUnmarshaler(batchKeyType *RestliType) Code {
	requiredFields, def := r.generateRequiredFields(batchKeyType)

	var params []Code
	if batchKeyType != nil {
		params = append(params, Add(Ids).Index().Add(batchKeyType.ReferencedType()))
	}
	params = append(params, Err().Error())

	AddNewInstance(def, r.Receiver(), r.TypeName())
	def.Add(r.GenerateUnmarshalField(batchKeyType))

	return utils.AddFuncOnReceiver(def, r.Receiver(), r.TypeName(), "DecodeQueryParams", RecordShouldUsePointer).
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

	const suffix = "RequiredFields"
	requiredFields = Code(Id(r.TypeName() + suffix))
	def = Var().Add(requiredFields).Op("=").
		Add(utils.NewRequiredFields).CustomFunc(utils.MultiLineCall, func(def *Group) {
		for _, i := range r.Includes {
			def.Qual(i.PackagePath(), i.TypeName()+suffix)
		}
	})

	if hasRequiredFields {
		const suffix = "RequiredFields"
		requiredFields = Code(Id(r.TypeName() + suffix))
		def.Dot("Add").CustomFunc(utils.MultiLineCall, func(def *Group) {
			for _, f := range r.Fields {
				if !f.IsOptionalOrDefault() {
					def.Lit(f.Name)
				}
			}
			if batchKeyType != nil {
				def.Qual(utils.BatchKeySetPackage, "EntityIDsField")
			}
		})
	}

	def.Line().Line()

	return requiredFields, def
}
