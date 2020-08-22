package codegen

import (
	"encoding/json"
	"fmt"
	"log"
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

	populateDefaultValues *Statement
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

func (r *Record) field(f Field) *Statement {
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
		Add(r.generateMarshalingCode()).Line().Line().
		Add(r.generateRestliEncoder()).Line().Line().
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

func (r *Record) generateMarshalingCode() *Statement {
	def := Empty()

	hasDefaultValue := r.generatePopulateDefaultValues(def)

	if hasDefaultValue {
		def.Func().
			Id(r.defaultValuesConstructor()).Params().
			Params(Id(r.Receiver()).Op("*").Id(r.Name))
		def.BlockFunc(func(def *Group) {
			def.Id(r.Receiver()).Op("=").New(Id(r.Name))
			for _, f := range r.Fields {
				if f.Type.Reference == nil {
					continue
				}
				if record, ok := f.Type.Reference.Resolve().(*Record); ok && !f.IsOptionalOrDefault() && record.hasDefaultValue() {
					def.Add(r.field(f)).Op("=").Op("*").Qual(record.PackagePath(), record.defaultValuesConstructor()).Call()
				}
			}
			def.Add(r.populateDefaultValues)
			def.Return()
		}).Line().Line()

		AddUnmarshalJSON(def, r.Receiver(), r.Name, func(def *Group) {
			def.Type().Id("_t").Id(r.Name)
			def.Err().Op("=").Qual(EncodingJson, Unmarshal).Call(Id("data"), Call(Op("*").Id("_t")).Call(Id(r.Receiver())))
			IfErrReturn(def).Line()
			def.Add(r.populateDefaultValues)
			def.Return()
		}).Line().Line()
	}

	return def
}

func (r *Record) generateRestliEncoder() *Statement {
	return AddRestLiEncode(Empty(), r.Receiver(), r.Name, func(def *Group) {
		r.generateEncoder(def, false, nil, nil)
		def.Return(Nil())
	})
}

func (r *Record) generateQueryParamEncoder(finderName *string) *Statement {
	receiver := r.Receiver()
	return AddFuncOnReceiver(Empty(), receiver, r.Name, EncodeQueryParams).
		Params().
		Params(Id("data").String(), Err().Error()).
		BlockFunc(func(def *Group) {
			def.Add(Encoder).Op(":=").Qual(RestLiEncodingPackage, "NewQueryParamsEncoder").Call()

			r.generateEncoder(def, true, finderName, nil)

			def.Return(Encoder.Finalize(), Nil())
		})
}

const finderNameParam = "q"

func (r *Record) generateEncoder(def *Group, forQueryParams bool, finderName *string, complexKeyKeyAccessor *Statement) {
	if finderName != nil && complexKeyKeyAccessor != nil {
		log.Panicln("Cannot provide both a finderName and a complexKeyKeyAccessor")
	}

	var fieldAccessor func(f Field) *Statement
	if complexKeyKeyAccessor != nil {
		fieldAccessor = func(f Field) *Statement {
			return Add(complexKeyKeyAccessor).Dot(f.FieldName())
		}
	} else {
		fieldAccessor = r.field
	}

	fields := append([]Field(nil), r.Fields...)
	sort.Slice(fields, func(i, j int) bool { return fields[i].Name < fields[j].Name })

	qIndex := -1
	if finderName != nil {
		qIndex = sort.Search(len(fields), func(i int) bool { return fields[i].Name >= finderNameParam })
		fields = append(fields[:qIndex], append([]Field{{}}, fields[qIndex:]...)...)
	}
	complexKeyParamsIndex := -1
	const ComplexKeyParams = "Params"
	if complexKeyKeyAccessor != nil {
		fields = append([]Field{{
			Name:       ComplexKeyParams,
			IsOptional: true,
		}}, fields...)
		complexKeyParamsIndex = 0
	}

	def.Add(Encoder.WriteObjectStart())

	if len(fields) == 0 {
		def.Add(Encoder.WriteObjectEnd())
		return
	}

	const needsDelimiterVar = "needsDelimiter"
	needsDelimiterCheckNeeded := len(fields) > 1 && fields[0].IsOptionalOrDefault()
	if needsDelimiterCheckNeeded {
		def.Id(needsDelimiterVar).Op(":=").False()
	}

	var returnOnError []Code
	if forQueryParams {
		returnOnError = append(returnOnError, Lit(""))
	}

	for i, f := range fields {
		var accessor *Statement
		if i == complexKeyParamsIndex {
			accessor = Id(r.Receiver()).Dot(ComplexKeyParams)
		} else {
			accessor = fieldAccessor(f)
		}

		serialize := def.Empty()
		if f.IsOptionalOrDefault() {
			serialize.If(Add(accessor).Op("!=").Nil())
		}

		serialize.BlockFunc(func(def *Group) {
			if i > 0 {
				if needsDelimiterCheckNeeded {
					def.If(Id(needsDelimiterVar)).Block(Encoder.WriteFieldDelimiter())
				} else {
					def.Add(Encoder.WriteFieldDelimiter())
				}
			}

			if i == qIndex {
				def.Add(Encoder.WriteFieldNameAndDelimiter(finderNameParam))
				def.Add(Encoder).Dot("String").Call(Lit(*finderName))
			} else if i == complexKeyParamsIndex {
				def.Add(Encoder.WriteFieldNameAndDelimiter("$params"))
				def.Err().Op("=").Add(Encoder).Dot("Encodable").Call(accessor)
				IfErrReturn(def, Err())
			} else {
				switch {
				case f.Type.Reference == nil && f.IsOptionalOrDefault():
					accessor = Op("*").Add(accessor)
				case f.Type.Reference != nil && !f.IsOptionalOrDefault():
					accessor = Op("&").Add(accessor)
				}
				Encoder.WriteField(def, f.Name, f.Type, accessor, returnOnError...)
			}

			if !f.IsOptionalOrDefault() {
				needsDelimiterCheckNeeded = false
			}

			if needsDelimiterCheckNeeded && i < len(fields)-1 {
				def.Id(needsDelimiterVar).Op("=").True()
			}
		})
		serialize.Line()
	}

	def.Add(Encoder.WriteObjectEnd())
}

func (r *Record) generateRestliDecoder() *Statement {
	return AddRestLiDecode(Empty(), r.Receiver(), r.Name, func(def *Group) {
		r.generateEncoder(def, false, nil, nil)
		def.Return(Nil())
	})
}

func (r *Record) generateDecoder(def *Group, forQueryParams bool, finderName *string, complexKeyKeyAccessor *Statement) {
	if finderName != nil && complexKeyKeyAccessor != nil {
		log.Panicln("Cannot provide both a finderName and a complexKeyKeyAccessor")
	}

	var fieldAccessor func(f Field) *Statement
	if complexKeyKeyAccessor != nil {
		fieldAccessor = func(f Field) *Statement {
			return Add(complexKeyKeyAccessor).Dot(f.FieldName())
		}
	} else {
		fieldAccessor = r.field
	}

	fields := append([]Field(nil), r.Fields...)
	sort.Slice(fields, func(i, j int) bool { return fields[i].Name < fields[j].Name })

	const finderNameParam = "q"
	qIndex := -1
	if finderName != nil {
		qIndex = sort.Search(len(fields), func(i int) bool { return fields[i].Name >= finderNameParam })
		fields = append(fields[:qIndex], append([]Field{{
			Name: "q",
			IsOptional:
		}}, fields[qIndex:]...)...)
	}
	complexKeyParamsIndex := -1
	const ComplexKeyParams = "Params"
	if complexKeyKeyAccessor != nil {
		fields = append([]Field{{
			Name:       ComplexKeyParams,
			IsOptional: true,
		}}, fields...)
		complexKeyParamsIndex = 0
	}

	if len(fields) == 0 {
		def.Return(Decoder.ReadObject(func(field *Statement, def *Group) {
			def.Return(Nil())
		}))
		return
	}

	requiredFieldsRemaining := Id("requiredFieldsRemaining")
	def.Add(requiredFieldsRemaining).Op(":=").Map(String()).Bool().Values(DictFunc(func(dict Dict) {
		for _, f := range fields {
			if !f.IsOptionalOrDefault() {
				dict[Lit(f.Name)] = True()
			}
		}
	}))

	var returnOnError []Code
	if finderName != nil {
		returnOnError = append(returnOnError, Lit(""))
	}

	def.Add(Decoder.ReadObject(func(field *Statement, def *Group) {
		def.Switch(field).BlockFunc(func(def *Group) {
			for i, f := range fields {
				var accessor *Statement
				if i == complexKeyParamsIndex {
					accessor = Id(r.Receiver()).Dot(ComplexKeyParams)
				} else {
					accessor = fieldAccessor(f)
				}

				if i == qIndex {
					def.Add(Encoder.WriteFieldNameAndDelimiter(finderNameParam))
					def.Add(Encoder).Dot("String").Call(Lit(*finderName))
				} else if i == complexKeyParamsIndex {
					def.Add(Encoder.WriteFieldNameAndDelimiter("$params"))
					def.Err().Op("=").Add(Encoder).Dot("Encodable").Call(accessor)
					IfErrReturn(def, Err())
				} else {
					switch {
					case f.Type.Reference == nil && f.IsOptionalOrDefault():
						accessor = Op("*").Add(accessor)
					case f.Type.Reference != nil && !f.IsOptionalOrDefault():
						accessor = Op("&").Add(accessor)
					}
					Encoder.WriteField(def, f.Name, f.Type, accessor, returnOnError...)
				}

			}
		})
	}))

	def.Add(Encoder.WriteObjectEnd())
}

func (r *Record) setDefaultValue(def *Group, name, rawJson string, t *RestliType) {
	def.If(Id(r.Receiver()).Dot(name).Op("==").Nil()).BlockFunc(func(def *Group) {
		switch {
		// Special case for primitives, instead of parsing them from JSON every time, we can leave them as literals
		case t.UnderlyingPrimitive() != nil:
			pt := t.UnderlyingPrimitive()
			def.Id("val").Op(":=").Add(pt.Cast(Lit(pt.getLit(rawJson))))
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

func (r *Record) generatePopulateDefaultValues(def *Statement) bool {
	r.populateDefaultValues = Empty()

	if !r.hasDefaultValue() {
		return false
	}

	AddFuncOnReceiver(def, r.Receiver(), r.Name, PopulateDefaultValues).Params().BlockFunc(func(def *Group) {
		for _, f := range r.Fields {
			if f.DefaultValue != nil {
				r.setDefaultValue(def, f.FieldName(), *f.DefaultValue, &f.Type)
				def.Line()
			}
		}
	}).Line().Line()

	r.populateDefaultValues.Id(r.Receiver()).Dot(PopulateDefaultValues).Call().Line()
	return true
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

	AddRestLiEncode(def, r.Receiver(), r.PartialUpdateStructName(), func(def *Group) {
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

		def.Add(Encoder.WriteObjectStart())

		if len(fields) == 0 {
			def.Add(Encoder.WriteObjectEnd())
			def.Return(Nil())
			return
		}

		needsDelimiter := Id("needsDelimiter")
		def.Add(needsDelimiter).Op(":=").False()
		def.Line()

		hasDeletes := deleteAccessor(fields[0])
		for i := 1; i < len(fields); i++ {
			hasDeletes.Op("||").Add(deleteAccessor(fields[i]))
		}

		def.If(hasDeletes).BlockFunc(func(def *Group) {
			def.Add(Encoder.WriteFieldNameAndDelimiter("$delete"))
			Encoder.ArrayEncoder(def, func(def *Group, indexWriter *Statement) {
				index := Id("index")
				def.Add(index).Op(":=").Lit(0)
				for _, f := range fields {
					accessor := deleteAccessor(f)
					def.If(accessor).BlockFunc(func(def *Group) {
						def.Add(indexWriter).Call(Add(index))
						def.Add(index).Op("++")
						if f.Type.IsReferenceEncodable() {
							accessor = Op("&").Add(accessor)
						}
						def.Add(Encoder).Dot("String").Call(Lit(f.Name))
					}).Line()
				}
				def.Return(Nil())
			})
			IfErrReturn(def, Err()).Line()
			def.Add(needsDelimiter).Op("=").True()
		})
		def.Line()

		hasSets := updateAccessor(fields[0]).Op("!=").Nil()
		for i := 1; i < len(fields); i++ {
			hasSets.Op("||").Add(updateAccessor(fields[i])).Op("!=").Nil()
		}

		def.If(hasSets).BlockFunc(func(def *Group) {
			def.If(needsDelimiter).Block(Encoder.WriteFieldDelimiter()).Line()

			def.Add(Encoder.WriteFieldNameAndDelimiter("$set"))
			def.Add(Encoder.WriteObjectStart())
			needsFirst := len(fields) > 1
			first := Id("first")
			if needsFirst {
				def.Add(first).Op(":=").True()
			}
			def.Line()
			for i, f := range fields {
				accessor := updateAccessor(f)
				def.If(Add(accessor).Op("!=").Nil()).BlockFunc(func(def *Group) {
					def.If(deleteAccessor(f)).BlockFunc(func(def *Group) {
						def.Return(Qual("fmt", "Errorf").Call(Lit(errorMessage(f))))
					})
					if needsFirst {
						if i == 0 {
							def.Add(first).Op("=").False()
						} else {
							def.If(first).Block(Add(first).Op("=").False()).Else().Block(Encoder.WriteFieldDelimiter())
						}
					}
					if f.Type.Reference == nil {
						accessor = Op("*").Add(accessor)
					}
					Encoder.WriteField(def, f.Name, f.Type, accessor)
				})
				def.Line()
			}
			def.Add(Encoder.WriteObjectEnd())
			if innerRecords > 0 {
				def.Add(needsDelimiter).Op("=").True()
			}
		})

		for i, f := range fields {
			isRecord := f.Type.Record() != nil
			if isRecord {
				def.Line()
				accessor := Id(r.Receiver()).Dot(f.FieldName())
				def.If(Add(accessor).Op("!=").Nil()).BlockFunc(func(def *Group) {
					def.If(needsDelimiter).Block(Encoder.WriteFieldDelimiter())
					def.If(Add(deleteAccessor(f)).Op("||").Add(updateAccessor(f).Op("!=").Nil())).BlockFunc(func(def *Group) {
						def.Return(Qual("fmt", "Errorf").Call(Lit(errorMessage(f))))
					})
					Encoder.WriteField(def, f.Name, f.Type, accessor)
					if i < innerRecords-1 {
						def.Add(needsDelimiter).Op("=").True()
					}
				})
			}
		}
		def.Line()

		def.Add(Encoder.WriteObjectEnd())

		def.Return(Nil())
	})

	return def
}
