package types

import (
	"sort"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	"github.com/PapaCharlie/go-restli/protocol/batchkeyset"
	. "github.com/dave/jennifer/jen"
)

func (r *Record) GenerateMarshalRestLi() *Statement {
	return AddMarshalRestLi(Empty(), r.Receiver(), r.Name, RecordShouldUsePointer, func(_, writer Code, def *Group) {
		r.generateMarshaler(writer, def)
	})
}

func (r *Record) generateMarshaler(writer Code, def *Group) {
	fields := r.SortedFields()

	def.Return(WriterUtils.WriteMap(writer, func(keyWriter Code, def *Group) {
		writeAllFields(def, fields, func(_ int, f Field) Code { return r.fieldAccessor(f) }, keyWriter)
	}))
}

func (r *Record) GenerateQueryParamMarshaler(finderName *string, batchKeyType Code) *Statement {
	receiver := r.Receiver()
	var params []Code
	if batchKeyType != nil {
		params = []Code{Add(utils.BatchKeySet).Qual(utils.BatchKeySetPackage, "BatchKeySet").Index(batchKeyType)}
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
						writeField(def, i, f, func(i int, f Field) Code {
							if i == qIndex {
								return Lit(*finderName)
							} else {
								return r.fieldAccessor(f)
							}
						}, paramNameWriter)
					}
				}
				def.Return(Nil())
			}))

			def.Add(utils.IfErrReturn(Lit(""), Err()))
			def.Return(WriterUtils.Finalize(Writer), Nil())
		})
}

func writeAllFields(def *Group, fields []Field, fieldAccessor func(i int, f Field) Code, keyWriter Code) {
	for i, f := range fields {
		writeField(def, i, f, fieldAccessor, keyWriter)
	}
	def.Return(Nil())
}

func writeField(def *Group, i int, f Field, fieldAccessor func(i int, f Field) Code, keyWriter Code) {
	accessor := fieldAccessor(i, f)

	serialize := func() Code {
		if f.IsOptionalOrDefault() && f.Type.Reference == nil {
			accessor = Op("*").Add(accessor)
		}
		return WriterUtils.Write(f.Type, Add(keyWriter).Call(Lit(f.Name)), accessor, Err())
	}

	if f.IsOptionalOrDefault() {
		def.If(Add(accessor).Op("!=").Nil()).Block(serialize())
	} else {
		def.Add(serialize())
	}
}
