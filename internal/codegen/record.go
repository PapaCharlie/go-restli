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

func (f *Field) IsPointer() bool {
	return f.IsOptionalOrDefault() && !f.Type.IsMapOrArray()
}

func (f *Field) IsOptionalOrDefault() bool {
	return f.IsOptional || f.DefaultValue != nil
}

func (f *Field) FieldName() string {
	return ExportedIdentifier(f.Name)
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

				if f.IsPointer() {
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
				if record, ok := f.Type.Reference.Resolve().(*Record); ok && !f.IsPointer() && record.hasDefaultValue() {
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
		r.generateEncoder(def, nil, nil)
		def.Return(Nil())
	})
}

func (r *Record) generateEncoder(def *Group, finderName *string, complexKeyParamsType *Identifier) {
	if finderName != nil && complexKeyParamsType != nil {
		log.Panicln("Cannot provide both a finderName and a complexKeyParamType")
	}

	var nameDelimiter string
	var fieldDelimiter rune
	if finderName != nil {
		nameDelimiter = "="
		fieldDelimiter = '&'
	} else {
		nameDelimiter = ":"
		fieldDelimiter = ','
	}

	if finderName == nil {
		def.Id("buf").Dot("WriteByte").Call(LitRune('('))
	}
	def.Line()

	fields := append([]Field(nil), r.Fields...)

	const finderNameParam = "q"
	qIndex := -1
	if finderName != nil {
		sort.Slice(fields, func(i, j int) bool { return fields[i].Name < fields[j].Name })
		qIndex = sort.Search(len(fields), func(i int) bool { return fields[i].Name >= finderNameParam })
		fields = append(fields[:qIndex], append([]Field{{}}, fields[qIndex:]...)...)
	}
	complexKeyParamsIndex := -1
	if complexKeyParamsType != nil {
		fields = append([]Field{{
			Name:       "Params",
			IsOptional: true,
		}}, fields...)
		complexKeyParamsIndex = 0
	}

	if len(r.Fields) == 0 {
		return
	}

	const needsDelimiterVar = "needsDelimiter"
	needsDelimiterCheckNeeded := len(fields) > 1 && fields[0].IsOptionalOrDefault()
	if needsDelimiterCheckNeeded {
		def.Id(needsDelimiterVar).Op(":=").False()
	}

	for i, f := range fields {
		serialize := def.Empty()
		if f.IsOptionalOrDefault() {
			if f.Type.IsMapOrArray() {
				serialize.If(Len(r.field(f)).Op(">").Lit(0))
			} else {
				serialize.If(r.field(f).Op("!=").Nil())
			}
		}

		serialize.BlockFunc(func(def *Group) {
			if i > 0 {
				writeDelimiter := Id("buf").Dot("WriteByte").Call(LitRune(fieldDelimiter))
				if needsDelimiterCheckNeeded {
					def.If(Id(needsDelimiterVar)).Block(writeDelimiter)
				} else {
					def.Add(writeDelimiter)
				}
			}

			if i == qIndex {
				def.Id("buf").Dot("WriteString").Call(Lit(finderNameParam + nameDelimiter + *finderName))
			} else if i == complexKeyParamsIndex {
				def.Id("buf").Dot("WriteString").Call(Lit("$params:"))
				def.Err().Op(":=").Add(r.field(f)).Dot(RestLiEncode).Call(Id(Codec), Id("buf"))
				IfErrReturn(def, Err())
			} else {
				accessor := r.field(f)
				if f.IsPointer() && f.Type.Reference == nil {
					accessor = Op("*").Add(accessor)
				}

				def.Id("buf").Dot("WriteString").Call(Lit(f.Name + nameDelimiter))
				f.Type.WriteToBuf(def, accessor)
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
	if finderName == nil {
		def.Id("buf").Dot("WriteByte").Call(LitRune(')'))
	}
}

func (r *Record) setDefaultValue(def *Group, name, rawJson string, t *RestliType) {
	def.If(Id(r.Receiver()).Dot(name).Op("==").Nil()).BlockFunc(func(def *Group) {
		switch {
		// Special case for primitives, instead of parsing them from JSON every time, we can leave them as literals
		case t.Primitive != nil:
			def.Id("val").Op(":=").Lit(t.Primitive.getLit(rawJson))
			def.Id(r.Receiver()).Dot(name).Op("= &").Id("val")
			return
		// If the default value for an array is the empty array, we can leave it as nil since that will behave
		// identically to an empty slice
		case t.Array != nil && emptyArrayRegex.MatchString(rawJson):
			return
		// For convenience, we create empty maps of the right type if the default value is the empty map
		case t.Map != nil && emptyMapRegex.MatchString(rawJson):
			def.Id(r.Receiver()).Dot(name).Op("=").Make(t.GoType(), Lit(0))
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

	AddMarshalJSON(def, r.Receiver(), r.PartialUpdateStructName(), func(def *Group) {
		partialUpdate := Id("partialUpdate")
		rawMessage := Qual(EncodingJson, "RawMessage")
		def.Var().Add(partialUpdate).StructFunc(func(def *Group) {
			def.Id(DeleteField).Index().String().Tag(JsonFieldTag("$delete", true))
			def.Id(UpdateField).Map(String()).Add(rawMessage).Tag(JsonFieldTag("$set", true))
			for _, f := range r.Fields {
				if record := f.Type.Record(); record != nil {
					def.Id(f.FieldName()).Op("*").Add(record.PartialUpdateStruct()).Tag(JsonFieldTag(f.Name, true))
				}
			}
		})
		def.Add(partialUpdate).Dot(UpdateField).Op("=").Make(Map(String()).Add(rawMessage))
		def.Line()

		for _, f := range r.Fields {
			errorMessage := fmt.Sprintf("Cannot both delete and update %q of %q", f.Name, r.Name)

			isRecord := f.Type.Record() != nil
			def.BlockFunc(func(def *Group) {
				isDeleted := Id("isDeleted")
				isUpdated := Id("isUpdated")
				toInit := []Code{isDeleted}
				if isRecord {
					toInit = append(toInit, isUpdated)
				}
				def.Var().List(toInit...).Bool()

				def.If(Id(r.Receiver()).Dot(DeleteField).Dot(f.FieldName())).BlockFunc(func(def *Group) {
					def.Add(partialUpdate).Dot(DeleteField).Op("=").Append(Add(partialUpdate).Dot(DeleteField), Lit(f.Name))
					def.Add(isDeleted).Op("=").True()
				})

				updateField := Id(r.Receiver()).Dot(UpdateField).Dot(f.FieldName())
				def.If(Add(updateField).Op("!=").Nil()).BlockFunc(func(def *Group) {
					def.If(isDeleted).BlockFunc(func(def *Group) {
						def.Return(Nil(), Qual("fmt", "Errorf").Call(Lit(errorMessage)))
					})
					if isRecord {
						def.Add(isUpdated).Op("=").True()
					}

					def.List(Id("data"), Err()).Op(":=").Qual(EncodingJson, Marshal).Call(updateField)
					IfErrReturn(def, Nil(), Err())

					def.Add(partialUpdate).Dot(UpdateField).Index(Lit(f.Name)).Op("=").Id("data")
				})

				if isRecord {
					recordField := Id(r.Receiver()).Dot(f.FieldName())
					def.If(Add(recordField).Op("!=").Nil()).BlockFunc(func(def *Group) {
						def.If(Add(isDeleted).Op("||").Add(isUpdated)).BlockFunc(func(def *Group) {
							def.Return(Nil(), Qual("fmt", "Errorf").Call(Lit(errorMessage)))
						})
						def.Add(partialUpdate).Dot(f.FieldName()).Op("=").Add(recordField)
					})
				}
			}).Line()
		}

		def.Return(Qual(EncodingJson, Marshal).Call(partialUpdate))
	})

	return def
}
