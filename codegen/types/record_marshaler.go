package types

import (
	"sort"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	"github.com/PapaCharlie/go-restli/restli/batchkeyset"
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

func (r *Record) GenerateMarshalRestLi() *Statement {
	return AddMarshalRestLi(Empty(), r.Receiver(), r.Name, RecordShouldUsePointer, func(def *Group) {
		r.generateMarshaler(def)
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

	return utils.AddFuncOnReceiver(Empty(), receiver, r.Name, utils.EncodeQueryParams, RecordShouldUsePointer).
		Params(params...).
		Params(Id("rawQuery").String(), Err().Error()).
		BlockFunc(func(def *Group) {
			def.Add(Writer).Op(":=").Qual(utils.RestLiCodecPackage, "NewRestLiQueryParamsWriter").Call()

			fields := r.SortedFields()

			insertFieldAt := func(index int, field Field) {
				fields = append(fields[:index], append([]Field{field}, fields[index:]...)...)
			}
			qIndex := -1
			if finderName != nil {
				qIndex = sort.Search(len(fields), func(i int) bool { return fields[i].Name >= utils.FinderNameParam })
				insertFieldAt(qIndex, Field{
					Type:       RestliType{Primitive: &StringPrimitive},
					Name:       utils.FinderNameParam,
					IsOptional: false,
				})
			}
			idsIndex := -1
			if batchKeyType != nil {
				idsIndex = sort.Search(len(fields), func(i int) bool { return fields[i].Name >= batchkeyset.EntityIDsField })
				insertFieldAt(idsIndex, Field{})
			}

			paramNameWriter := Id("paramNameWriter")
			paramNameWriterFunc := Add(paramNameWriter).Func().Params(String()).Add(WriterQual)
			def.Err().Op("=").Add(Writer).Dot("WriteParams").Call(Func().Params(paramNameWriterFunc).Params(Err().Error()).BlockFunc(func(def *Group) {
				for i, f := range fields {
					if i == idsIndex {
						def.BlockFunc(func(def *Group) {
							def.Err().Op("=").Add(utils.BatchKeySet).Dot("Encode").Call(paramNameWriter)
							def.Add(utils.IfErrReturn(Err()))
						}).Line()
					} else {
						var accessor Code
						if i == qIndex {
							accessor = Lit(*finderName)
						} else {
							accessor = r.fieldAccessor(f)
						}
						writeField(def, f, accessor, paramNameWriter)
					}
				}
				def.Return(Nil())
			}))

			def.Add(utils.IfErrReturn(Lit(""), Err()))
			def.Return(Writer.Finalize(), Nil())
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
		if f.IsOptionalOrDefault() && f.Type.Reference == nil {
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
