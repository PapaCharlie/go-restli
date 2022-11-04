package types

import (
	"github.com/PapaCharlie/go-restli/v2/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func AddMarshalRestLi(def *Statement, receiver, typeName string, pointer utils.ShouldUsePointer, f func(def *Group)) *Statement {
	utils.AddFuncOnReceiver(def, receiver, typeName, utils.MarshalRestLi, pointer).
		Params(WriterParam).
		Params(Err().Error()).
		BlockFunc(f).
		Line().Line()

	utils.AddFuncOnReceiver(def, receiver, typeName, "MarshalJSON", pointer).
		Params().
		Params(Id("data").Index().Byte(), Err().Error()).
		BlockFunc(func(def *Group) {
			def.Add(Writer).Op(":=").Qual(utils.RestLiCodecPackage, "NewCompactJsonWriter").Call()
			def.Err().Op("=").Id(receiver).Dot(utils.MarshalRestLi).Call(Writer)
			def.Add(utils.IfErrReturn(Nil(), Err()))
			def.Return(Index().Byte().Call(Add(Writer.Finalize())), Nil())
		}).Line().Line()

	return def
}

const MarshalFields = "MarshalFields"

func (r *Record) GenerateMarshalFields() *Statement {
	return utils.AddFuncOnReceiver(Empty(), r.Receiver(), r.TypeName(), MarshalFields, RecordShouldUsePointer).
		Params(KeyWriterFunc).
		Params(Err().Error()).
		BlockFunc(func(def *Group) {
			for _, i := range r.Includes {
				def.Err().Op("=").Id(r.Receiver()).Dot(i.TypeName()).Dot(MarshalFields).Call(KeyWriter)
				def.Add(utils.IfErrReturn(Err()))
			}
			fields := r.SortedFields()
			writeAllFields(def, fields, func(_ int, f Field) Code { return r.fieldAccessor(f) }, KeyWriter)
		}).
		Line().Line()
}

func (r *Record) GenerateMarshalRestLi() *Statement {
	def := r.GenerateMarshalFields()
	return AddMarshalRestLi(def, r.Receiver(), r.TypeName(), RecordShouldUsePointer, func(def *Group) {
		def.Return(Add(Writer).Dot("WriteMap").Call(Id(r.Receiver()).Dot(MarshalFields)))
	})
}

func (r *Record) generateMarshaler(def *Group) {
	fields := r.SortedFields()

	def.Return(Writer.WriteMap(Writer, func(keyWriter Code, def *Group) {
		writeAllFields(def, fields, func(_ int, f Field) Code { return r.fieldAccessor(f) }, keyWriter)
	}))
}

func (r *Record) GenerateQueryParamMarshaler(finderName *string, batchKeyType *RestliType) *Statement {
	receiver := r.Receiver()
	var params []Code
	if batchKeyType != nil {
		params = []Code{Add(utils.BatchKeySet).Qual(utils.BatchKeySetPackage, "BatchKeySet").Index(batchKeyType.ReferencedType())}
	}

	def := r.GenerateMarshalFields()

	return utils.AddFuncOnReceiver(def, receiver, r.TypeName(), utils.EncodeQueryParams, RecordShouldUsePointer).
		Params(params...).
		Params(Id("rawQuery").String(), Err().Error()).
		BlockFunc(func(def *Group) {
			paramNameWriter := Id("paramNameWriter")
			paramNameWriterFunc := Add(paramNameWriter).Func().Params(String()).Add(WriterQual)
			def.Return(Qual(utils.RestLiCodecPackage, "BuildQueryParams").Call(Func().Params(paramNameWriterFunc).Params(Err().Error()).BlockFunc(func(def *Group) {
				if batchKeyType != nil {
					def.Err().Op("=").Add(utils.BatchKeySet).Dot("Encode").Call(paramNameWriter)
					def.Add(utils.IfErrReturn(Err()))
				}

				if finderName != nil {
					def.Add(Writer.Write(
						RestliType{Primitive: &StringPrimitive},
						Add(paramNameWriter).Call(Lit(utils.FinderNameParam)),
						Lit(*finderName),
						Err(),
					))
				}

				def.Return(Id(r.Receiver()).Dot(MarshalFields).Call(paramNameWriter))
			})))
		})
}

func writeAllFields(def *Group, fields []Field, fieldAccessor func(i int, f Field) Code, keyWriter Code) {
	for i, f := range fields {
		writeField(def, f, fieldAccessor(i, f), keyWriter)
	}
	def.Return(Nil())
}

func writeField(def *Group, f Field, accessor Code, keyWriter Code) {
	serialize := func() Code {
		if f.IsOptionalOrDefault() && (f.Type.Reference == nil || f.Type.IsCustomTyperef()) {
			accessor = Op("*").Add(accessor)
		}
		return Writer.Write(f.Type, Add(keyWriter).Call(Lit(f.Name)), accessor, Err())
	}

	if f.IsOptionalOrDefault() {
		def.If(Add(accessor).Op("!=").Nil()).Block(serialize())
	} else {
		def.Add(serialize())
	}
}
