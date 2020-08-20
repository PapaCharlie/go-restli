package codegen

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"

	. "github.com/dave/jennifer/jen"
)

var (
	emptyMapRegex   = regexp.MustCompile("{ *}")
	emptyArrayRegex = regexp.MustCompile("\\[ *]")
)

type Record struct {
	NamedType
	Fields []Field
}

func (r *Record) InnerTypes() IdentifierSet {
	innerTypes := make(IdentifierSet)
	for _, f := range r.Fields {
		innerTypes.AddAll(f.Type.InnerTypes())
	}

	return innerTypes
}

func (r *Record) PartialUpdateStructName() string {
	return r.Name + PartialUpdate
}

func (r *Record) PartialUpdateStruct() *Statement {
	return Qual(r.PackagePath(), r.PartialUpdateStructName())
}

type Field struct {
	Type         RestliType
	Name         string
	Doc          string
	IsOptional   bool
	DefaultValue *string
}

func (r *Record) field(f Field) Code {
	return Id(r.Receiver()).Dot(f.FieldName())
}

func (f *Field) IsOptionalOrDefault() bool {
	return f.IsOptional || f.DefaultValue != nil
}

func (f *Field) FieldName() string {
	return ExportedIdentifier(f.Name)
}

func (r *Record) SortedFields() (fields []Field) {
	fields = append(fields, r.Fields...)
	sort.Slice(fields, func(i, j int) bool { return fields[i].Name < fields[j].Name })
	return fields
}

func (r *Record) GenerateCode() *Statement {
	return Empty().
		Add(r.generateStruct()).Line().Line().
		Add(r.generateDefaultValuesCode()).Line().Line().
		Add(r.generateMarshalRestLi()).Line().Line().
		Add(r.generateUnmarshalRestLi()).Line().Line().
		Add(r.generatePartialUpdateStruct()).Line()
}

func (r *Record) generateStruct() *Statement {
	return AddWordWrappedComment(Empty(), r.Doc).Line().
		Type().Id(r.Name).
		StructFunc(func(def *Group) {
			for _, f := range r.Fields {
				field := def.Empty()
				AddWordWrappedComment(field, f.Doc).Line()
				field.Id(f.FieldName())

				if f.IsOptionalOrDefault() {
					field.Add(f.Type.PointerType())
				} else {
					field.Add(f.Type.GoType())
				}

				field.Tag(JsonFieldTag(f.Name, f.IsOptionalOrDefault()))
			}
		})
}

func (r *Record) generateDefaultValuesCode() Code {
	if !r.hasDefaultValue() {
		return Empty()
	}

	def := Empty()

	def.Func().
		Id(r.defaultValuesConstructor()).Params().
		Params(Id(r.Receiver()).Op("*").Id(r.Name)).
		BlockFunc(func(def *Group) {
			def.Id(r.Receiver()).Op("=").New(Id(r.Name))
			for _, f := range r.Fields {
				if f.Type.Reference == nil {
					continue
				}
				if record, ok := f.Type.Reference.Resolve().(*Record); ok && !f.IsOptionalOrDefault() && record.hasDefaultValue() {
					def.Add(r.field(f)).Op("=").Op("*").Qual(record.PackagePath(), record.defaultValuesConstructor()).Call()
				}
			}
			def.Id(r.Receiver()).Dot(PopulateLocalDefaultValues).Call()
			def.Return()
		}).Line().Line()

	AddFuncOnReceiver(def, r.Receiver(), r.Name, PopulateLocalDefaultValues).Params().BlockFunc(func(def *Group) {
		for _, f := range r.Fields {
			if f.DefaultValue != nil {
				r.setDefaultValue(def, f.FieldName(), *f.DefaultValue, &f.Type)
				def.Line()
			}
		}
	}).Line().Line()

	return def
}

func (r *Record) generateMarshalRestLi() *Statement {
	return AddMarshalRestLi(Empty(), r.Receiver(), r.Name, func(def *Group) {
		r.generateMarshaler(def, nil)
	})
}

func (r *Record) generateUnmarshalRestLi() *Statement {
	return AddRestLiDecode(Empty(), r.Receiver(), r.Name, func(def *Group) {
		r.generateUnmarshaler(def, nil, nil)
	})
}

const finderNameParam = "q"

func (r *Record) generateQueryParamMarshaler(finderName *string) *Statement {
	receiver := r.Receiver()
	return AddFuncOnReceiver(Empty(), receiver, r.Name, EncodeQueryParams).
		Params().
		Params(Id("data").String(), Err().Error()).
		BlockFunc(func(def *Group) {
			def.Add(Writer).Op(":=").Qual(RestLiCodecPackage, "NewRestLiQueryParamsWriter").Call()

			fields := r.SortedFields()

			qIndex := -1
			if finderName != nil {
				qIndex = sort.Search(len(fields), func(i int) bool { return fields[i].Name >= finderNameParam })
				fields = append(fields[:qIndex], append([]Field{{
					Type:       RestliType{Primitive: &StringPrimitive},
					Name:       finderNameParam,
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

			def.Add(IfErrReturn(Lit(""), Err()))
			def.Return(Writer.Finalize(), Nil())
		})
}

const ComplexKeyParamsField = "Params"

func (r *Record) generateMarshaler(def *Group, complexKeyKeyAccessor *Statement) {
	fields := r.SortedFields()

	complexKeyParamsIndex := -1
	if complexKeyKeyAccessor != nil {
		fields = append([]Field{{
			Name:       "$params",
			IsOptional: true,
			Type:       RestliType{Reference: new(Identifier)},
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

func (r *Record) generateUnmarshaler(def *Group, complexKeyKeyAccessor *Statement, complexKeyParamsType *Identifier) {
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
		def.Return(Reader.ReadMap(func(field Code, def *Group) {
			def.Return(Reader.Skip())
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

	def.Err().Op("=").Add(Reader.ReadMap(func(field Code, def *Group) {
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

					def.Add(Reader.Read(f.Type, accessor))
				})
			}
			def.Default().BlockFunc(func(def *Group) {
				def.Err().Op("=").Add(Reader.Skip())
			})
		})
		def.Add(IfErrReturn(Err()))
		def.Delete(requiredFieldsRemaining, field)
		def.Return(Nil())
	})).Line()

	def.Add(IfErrReturn(Err())).Line()

	def.If(Len(requiredFieldsRemaining).Op("!=").Lit(0)).BlockFunc(func(def *Group) {
		def.Return(Qual("fmt", "Errorf").Call(Lit("required fields not all present: %+v"), requiredFieldsRemaining))
	}).Line()

	if r.hasDefaultValue() {
		def.Id(r.Receiver()).Dot(PopulateLocalDefaultValues).Call()
	}

	def.Return(Nil())
}

func (r *Record) generateQueryParamsUnmarhsaler(def *Group, finderName *string) {
	fields := r.SortedFields()
	qIndex := -1
	if finderName != nil {
		qIndex = sort.Search(len(fields), func(i int) bool { return fields[i].Name >= finderNameParam })
		fields = append(fields[:qIndex], append([]Field{{
			Name:       finderNameParam,
			IsOptional: false,
			Type:       RestliType{Primitive: &StringPrimitive},
		}}, fields[qIndex:]...)...)
	}

	finderNameVar := Id("finderName")
	if finderName != nil {
		def.Var().Add(finderNameVar).String()
		def.Line()
	}

}

func (r *Record) setDefaultValue(def *Group, name, rawJson string, t *RestliType) {
	def.If(Id(r.Receiver()).Dot(name).Op("==").Nil()).BlockFunc(func(def *Group) {
		switch {
		// Special case for primitives, instead of parsing them from JSON every time, we can leave them as literals
		case t.UnderlyingPrimitive() != nil:
			pt := t.UnderlyingPrimitive()
			def.Id("val").Op(":=").Add(pt.Cast(pt.getLit(rawJson)))
			def.Id(r.Receiver()).Dot(name).Op("= &").Id("val")
			return
		// If the default value for an array is the empty array, we can leave it as nil since that will behave
		// identically to an empty slice
		case t.Array != nil && emptyArrayRegex.MatchString(rawJson):
			return
		// For convenience, we create empty maps of the right type if the default value is the empty map
		case t.Map != nil && emptyMapRegex.MatchString(rawJson):
			def.Id(r.Receiver()).Dot(name).Op("= &").Add(t.GoType()).Values()
			return
		// Enum values can also be added as literals
		case t.Reference != nil:
			if enum, ok := t.Reference.Resolve().(*Enum); ok {
				var v string
				err := json.Unmarshal([]byte(rawJson), &v)
				if err != nil {
					Logger.Panicln("illegal enum", err)
				}
				def.Id("val").Op(":=").Qual(enum.PackagePath(), enum.SymbolIdentifier(v))
				def.Id(r.Receiver()).Dot(name).Op("= &").Id("val")
				return
			}

			if fixed, ok := t.Reference.Resolve().(*Fixed); ok {
				def.Id("val").Op(":=").Add(fixed.getLit(rawJson))
				def.Id(r.Receiver()).Dot(name).Op("= &").Id("val")
				return
			}
		}

		field := Op("&").Id(r.Receiver()).Dot(name)

		def.Err().Op(":=").Qual(EncodingJson, Unmarshal).Call(Index().Byte().Call(Lit(rawJson)), field)
		def.If(Err().Op("!=").Nil()).Block(Qual("log", "Panicln").Call(Lit("Illegal default value"), Err()))
	})
}

func (r *Record) hasDefaultValue() bool {
	for _, f := range r.Fields {
		if f.DefaultValue != nil {
			return true
		}
	}
	return false
}

func (r *Record) defaultValuesConstructor() string {
	return "New" + r.Name + "WithDefaultValues"
}

func (r *Record) generatePartialUpdateStruct() *Statement {
	def := Empty()

	const (
		DeleteField = "Delete"
		UpdateField = "Update"
	)

	// Generate the struct
	AddWordWrappedComment(def, fmt.Sprintf(
		"%s is used to represent a partial update on %s. Toggling the value of a field\n"+
			"in Delete represents selecting it for deletion in a partial update, while\n"+
			"setting the value of a field in Update represents setting that field in the\n"+
			"current struct. Other fields in this struct represent record fields that can\n"+
			"themselves be partially updated.",
		r.PartialUpdateStructName(), r.Name,
	)).Line()

	def.Type().Id(r.PartialUpdateStructName()).StructFunc(func(def *Group) {
		def.Id(DeleteField).StructFunc(func(def *Group) {
			for _, f := range r.Fields {
				def.Id(f.FieldName()).Bool()
			}
		})
		def.Id(UpdateField).StructFunc(func(def *Group) {
			for _, f := range r.Fields {
				def.Id(f.FieldName()).Add(Op("*").Add(f.Type.GoType()))
			}
		})
		for _, f := range r.Fields {
			if record := f.Type.Record(); record != nil {
				def.Id(f.FieldName()).Op("*").Add(record.PartialUpdateStruct())
			}
		}
	}).Line().Line()

	AddMarshalRestLi(def, r.Receiver(), r.PartialUpdateStructName(), func(def *Group) {
		fields := r.SortedFields()
		innerRecords := 0
		for _, f := range fields {
			if f.Type.Record() != nil {
				innerRecords++
			}
		}

		deleteAccessor := func(f Field) *Statement { return Id(r.Receiver()).Dot(DeleteField).Dot(f.FieldName()) }
		updateAccessor := func(f Field) *Statement { return Id(r.Receiver()).Dot(UpdateField).Dot(f.FieldName()) }
		errorMessage := func(f Field) string {
			return fmt.Sprintf("cannot both delete and update %q of %q", f.Name, r.Name)
		}

		if len(fields) == 0 {
			def.Return(Writer.WriteMap(Writer, func(keyWriter Code, def *Group) {
				def.Return(Nil())
			}))
			return
		}

		def.Return(Writer.WriteMap(Writer, func(keyWriter Code, def *Group) {
			hasDeletes := deleteAccessor(fields[0])
			for i := 1; i < len(fields); i++ {
				hasDeletes.Op("||").Add(deleteAccessor(fields[i]))
			}

			def.If(hasDeletes).BlockFunc(func(def *Group) {
				deleteType := RestliType{Primitive: &StringPrimitive}
				def.Err().Op("=").Add(Writer.WriteArray(Add(keyWriter).Call(Lit("$delete")), func(itemWriter Code, def *Group) {
					for _, f := range fields {
						accessor := deleteAccessor(f)
						def.If(accessor).BlockFunc(func(def *Group) {
							def.Add(Writer.Write(deleteType, Add(itemWriter).Call(), Lit(f.Name)))
						}).Line()
					}
					def.Add(Return(Nil()))
				}))
				def.Add(IfErrReturn(Err()))
			})
			def.Line()

			hasSets := updateAccessor(fields[0]).Op("!=").Nil()
			for i := 1; i < len(fields); i++ {
				hasSets.Op("||").Add(updateAccessor(fields[i])).Op("!=").Nil()
			}

			def.If(hasSets).BlockFunc(func(def *Group) {
				def.Err().Op("=").Add(Writer.WriteMap(Add(keyWriter).Call(Lit("$set")), func(keyWriter Code, def *Group) {
					for _, f := range fields {
						accessor := updateAccessor(f)
						def.If(Add(accessor).Op("!=").Nil()).BlockFunc(func(def *Group) {
							def.If(deleteAccessor(f)).BlockFunc(func(def *Group) {
								def.Return(Qual("fmt", "Errorf").Call(Lit(errorMessage(f))))
							})
							if f.Type.Reference == nil {
								accessor = Op("*").Add(accessor)
							}
							def.Add(Writer.Write(f.Type, Add(keyWriter).Call(Lit(f.Name)), accessor, Err()))
						})
						def.Line()
					}
					def.Return(Nil())
				}))
				def.Add(IfErrReturn(Err()))
			})

			for _, f := range fields {
				if f.Type.Record() != nil {
					def.Line()
					accessor := Id(r.Receiver()).Dot(f.FieldName())
					def.If(Add(accessor).Op("!=").Nil()).BlockFunc(func(def *Group) {
						def.If(Add(deleteAccessor(f)).Op("||").Add(updateAccessor(f).Op("!=").Nil())).BlockFunc(func(def *Group) {
							def.Return(Qual("fmt", "Errorf").Call(Lit(errorMessage(f))))
						})
						def.Add(Writer.Write(f.Type, Add(keyWriter).Call(Lit(f.Name)), accessor, Err()))
					})
				}
			}
			def.Line()

			def.Return(Nil())
		}))
	})

	return def
}
