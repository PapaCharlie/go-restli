package types

import (
	"sort"

	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func AddMarshalRestLi(def *Statement, receiver, typeName string, f func(def *Group)) *Statement {
	utils.AddFuncOnReceiver(def, receiver, typeName, MarshalRestLi).
		Params(Add(Writer).Add(WriterQual)).
		Params(Err().Error()).
		BlockFunc(f).
		Line().Line()

	utils.AddFuncOnReceiver(def, receiver, typeName, "MarshalJSON").
		Params().
		Params(Id("data").Index().Byte(), Err().Error()).
		BlockFunc(func(def *Group) {
			def.Add(Writer).Op(":=").Qual(RestLiCodecPackage, "NewCompactJsonWriter").Call()
			def.Err().Op("=").Id(receiver).Dot(MarshalRestLi).Call(Writer)
			def.Add(utils.IfErrReturn(Nil(), Err()))
			def.Return(Index().Byte().Call(Add(Writer.Finalize())), Nil())
		}).Line().Line()

	return def
}

func (r *Record) GenerateMarshalRestLi() *Statement {
	return AddMarshalRestLi(Empty(), r.Receiver(), r.Name, func(def *Group) {
		r.generateMarshaler(def, nil)
	})
}

func (r *Record) generateMarshaler(def *Group, complexKeyKeyAccessor *Statement) {
	fields := r.SortedFields()

	complexKeyParamsIndex := -1
	if complexKeyKeyAccessor != nil {
		fields = append([]Field{{
			Name:       "$params",
			IsOptional: true,
			Type:       RestliType{Reference: new(utils.Identifier)},
		}}, fields...)
		complexKeyParamsIndex = 0
	}

	var fieldAccessor func(i int, f Field) Code
	if complexKeyKeyAccessor != nil {
		fieldAccessor = func(i int, f Field) Code {
			if i == complexKeyParamsIndex {
				return Id(r.Receiver()).Dot(ComplexKeyParamsField)
			} else {
				return Add(complexKeyKeyAccessor).Dot(f.FieldName())
			}
		}
	} else {
		fieldAccessor = func(_ int, f Field) Code { return r.field(f) }
	}

	def.Return(Writer.WriteMap(Writer, func(keyWriter Code, def *Group) {
		writeAllFields(def, fields, fieldAccessor, keyWriter)
	}))
}

func (r *Record) GenerateQueryParamMarshaler(finderName *string, isBatchRequest bool) *Statement {
	receiver := r.Receiver()
	var params []Code
	entityIDsWriter := Code(Id("entityIDsWriter"))
	if isBatchRequest {
		params = []Code{Add(entityIDsWriter).Qual(RestLiCodecPackage, "ArrayWriter")}
	}

	return utils.AddFuncOnReceiver(Empty(), receiver, r.Name, EncodeQueryParams).
		Params(params...).
		Params(Id("rawQuery").String(), Err().Error()).
		BlockFunc(func(def *Group) {
			def.Add(Writer).Op(":=").Qual(RestLiCodecPackage, "NewRestLiQueryParamsWriter").Call()

			fields := r.SortedFields()

			insertFieldAt := func(index int, field Field) {
				fields = append(fields[:index], append([]Field{field}, fields[index:]...)...)
			}
			qIndex := -1
			if finderName != nil {
				qIndex = sort.Search(len(fields), func(i int) bool { return fields[i].Name >= FinderNameParam })
				insertFieldAt(qIndex, Field{
					Type:       RestliType{Primitive: &StringPrimitive},
					Name:       FinderNameParam,
					IsOptional: false,
				})
			}
			idsIndex := -1
			if isBatchRequest {
				idsIndex = sort.Search(len(fields), func(i int) bool { return fields[i].Name >= EntityIDsParam })
				insertFieldAt(idsIndex, Field{})
			}

			paramNameWriter := Id("paramNameWriter")
			paramNameWriterFunc := Add(paramNameWriter).Func().Params(String()).Add(WriterQual)
			def.Err().Op("=").Add(Writer).Dot("WriteParams").Call(Func().Params(paramNameWriterFunc).Params(Err().Error()).BlockFunc(func(def *Group) {
				for i, f := range fields {
					if i == idsIndex {
						def.Err().Op("=").Add(paramNameWriter).Call(Lit(EntityIDsParam)).Dot("WriteArray").Call(entityIDsWriter)
					} else {
						writeField(def, i, f, func(i int, f Field) Code {
							if i == qIndex {
								return Lit(*finderName)
							} else {
								return r.field(f)
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
