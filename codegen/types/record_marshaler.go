package types

import (
	"sort"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func AddMarshalRestLi(def *Statement, receiver, typeName string, f func(def *Group)) *Statement {
	utils.AddFuncOnReceiver(def, receiver, typeName, utils.MarshalRestLi).
		Params(Add(Writer).Add(WriterQual)).
		Params(Err().Error()).
		BlockFunc(f).
		Line().Line()

	utils.AddFuncOnReceiver(def, receiver, typeName, "MarshalJSON").
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
	return AddMarshalRestLi(Empty(), r.Receiver(), r.Name, func(def *Group) {
		r.generateMarshaler(def)
	})
}

func (r *Record) generateMarshaler(def *Group) {
	fields := r.SortedFields()

	def.Return(Writer.WriteMap(Writer, func(keyWriter Code, def *Group) {
		writeAllFields(def, fields, func(_ int, f Field) Code { return r.fieldAccessor(f) }, keyWriter)
	}))
}

func (r *Record) GenerateQueryParamMarshaler(finderName *string, isBatchRequest bool) *Statement {
	receiver := r.Receiver()
	var params []Code
	if isBatchRequest {
		params = []Code{Add(utils.EntityIDsEncoder).Op("*").Add(utils.BatchEntityIDsEncoder)}
	}

	return utils.AddFuncOnReceiver(Empty(), receiver, r.Name, utils.EncodeQueryParams).
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
			if isBatchRequest {
				idsIndex = sort.Search(len(fields), func(i int) bool { return fields[i].Name >= utils.EntityIDsParam })
				insertFieldAt(idsIndex, Field{})
			}

			paramNameWriter := Id("paramNameWriter")
			paramNameWriterFunc := Add(paramNameWriter).Func().Params(String()).Add(WriterQual)
			def.Err().Op("=").Add(Writer).Dot("WriteParams").Call(Func().Params(paramNameWriterFunc).Params(Err().Error()).BlockFunc(func(def *Group) {
				for i, f := range fields {
					if i == idsIndex {
						def.BlockFunc(func(def *Group) {
							def.Err().Op("=").Add(utils.EntityIDsEncoder).Dot("Encode").Call(paramNameWriter)
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
			def.Return(Writer.Finalize(), Nil())
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

	serialize := def.Empty()
	if f.IsOptionalOrDefault() {
		serialize.If(Add(accessor).Op("!=").Nil())
	}

	serialize.BlockFunc(func(def *Group) {
		if f.IsOptionalOrDefault() && f.Type.Reference == nil {
			accessor = Op("*").Add(accessor)
		}
		def.Add(Writer.Write(f.Type, Add(keyWriter).Call(Lit(f.Name)), accessor, Err()))
	}).Line()
}
