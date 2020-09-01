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

func (r *Record) GenerateQueryParamMarshaler(finderName *string) *Statement {
	receiver := r.Receiver()
	return utils.AddFuncOnReceiver(Empty(), receiver, r.Name, EncodeQueryParams).
		Params().
		Params(Id("data").String(), Err().Error()).
		BlockFunc(func(def *Group) {
			def.Add(Writer).Op(":=").Qual(RestLiCodecPackage, "NewRestLiQueryParamsWriter").Call()

			fields := r.SortedFields()

			qIndex := -1
			if finderName != nil {
				qIndex = sort.Search(len(fields), func(i int) bool { return fields[i].Name >= FinderNameParam })
				fields = append(fields[:qIndex], append([]Field{{
					Type:       RestliType{Primitive: &StringPrimitive},
					Name:       FinderNameParam,
					IsOptional: false,
				}}, fields[qIndex:]...)...)
			}

			paramNameWriter := Id("paramNameWriter")
			paramNameWriterFunc := Add(paramNameWriter).Func().Params(String()).Add(WriterQual)
			def.Err().Op("=").Add(Writer).Dot("WriteParams").Call(Func().Params(paramNameWriterFunc).Params(Err().Error()).BlockFunc(func(def *Group) {
				writeAllFields(def, fields, func(i int, f Field) Code {
					if i == qIndex {
						return Lit(*finderName)
					} else {
						return r.field(f)
					}
				}, paramNameWriter)
			}))

			def.Add(utils.IfErrReturn(Lit(""), Err()))
			def.Return(Writer.Finalize(), Nil())
		})
}

func writeAllFields(def *Group, fields []Field, fieldAccessor func(i int, f Field) Code, keyWriter Code) {
	for i, f := range fields {
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
	def.Return(Nil())
}
