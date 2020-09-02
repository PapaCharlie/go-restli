package types

import (
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func AddUnmarshalRestli(def *Statement, receiver, typeName string, f func(def *Group)) *Statement {
	utils.AddFuncOnReceiver(def, receiver, typeName, UnmarshalRestLi).
		Params(Add(Reader).Add(ReaderQual)).
		Params(Err().Error()).
		BlockFunc(f).
		Line().Line()

	data := Id("data")
	utils.AddFuncOnReceiver(def, receiver, typeName, "UnmarshalJSON").
		Params(Add(data).Index().Byte()).
		Params(Error()).
		BlockFunc(func(def *Group) {
			def.Add(Reader).Op(":=").Add(NewJsonReader).Call(data)
			def.Return(Id(receiver).Dot(UnmarshalRestLi).Call(Reader))
		}).Line().Line()

	return def
}

func (r *Record) GenerateUnmarshalRestLi() *Statement {
	return AddUnmarshalRestli(Empty(), r.Receiver(), r.Name, func(def *Group) {
		r.generateUnmarshaler(def, nil, nil)
	})
}

func (r *Record) generateUnmarshaler(def *Group, complexKeyKeyAccessor *Statement, complexKeyParamsType *utils.Identifier) {
	fields := r.SortedFields()

	complexKeyParamsIndex := -1
	if complexKeyKeyAccessor != nil {
		fields = append([]Field{{
			Name:       "$params",
			IsOptional: true,
			Type:       RestliType{Reference: complexKeyParamsType},
		}}, fields...)
		complexKeyParamsIndex = 0
	}

	if len(fields) == 0 {
		def.Return(Reader.ReadMap(Reader, func(reader Code, field Code, def *Group) {
			def.Return(Reader.Skip(reader))
		}))
		return
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

	requiredFieldsRemaining := Id("requiredFieldsRemaining")
	def.Add(requiredFieldsRemaining).Op(":=").Map(String()).Bool().Values(DictFunc(func(dict Dict) {
		for _, f := range fields {
			if !f.IsOptionalOrDefault() {
				dict[Lit(f.Name)] = True()
			}
		}
	})).Line()

	def.Err().Op("=").Add(Reader.ReadMap(Reader, func(reader, field Code, def *Group) {
		def.Switch(field).BlockFunc(func(def *Group) {
			for i, f := range fields {
				def.Case(Lit(f.Name)).BlockFunc(func(def *Group) {
					accessor := fieldAccessor(i, f)

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
		def.Add(utils.IfErrReturn(Err()))
		def.Delete(requiredFieldsRemaining, field)
		def.Return(Nil())
	})).Line()

	def.Add(utils.IfErrReturn(Err())).Line()

	def.If(Len(requiredFieldsRemaining).Op("!=").Lit(0)).BlockFunc(func(def *Group) {
		def.Return(Qual("fmt", "Errorf").Call(Lit("required fields not all present: %+v"), requiredFieldsRemaining))
	}).Line()

	if r.hasDefaultValue() {
		def.Id(r.Receiver()).Dot(PopulateLocalDefaultValues).Call()
	}

	def.Return(Nil())
}

// TODO
// func (r *Record) generateQueryParamsUnmarhsaler(def *Group, finderName *string) {
// 	fields := r.SortedFields()
// 	qIndex := -1
// 	if finderName != nil {
// 		qIndex = sort.Search(len(fields), func(i int) bool { return fields[i].Name >= FinderNameParam })
// 		fields = append(fields[:qIndex], append([]Field{{
// 			Name:       FinderNameParam,
// 			IsOptional: false,
// 			Type:       RestliType{Primitive: &StringPrimitive},
// 		}}, fields[qIndex:]...)...)
// 	}
//
// 	finderNameVar := Id("finderName")
// 	if finderName != nil {
// 		def.Var().Add(finderNameVar).String()
// 		def.Line()
// 	}
//
// }
