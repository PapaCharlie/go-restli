package codegen

import (
	"encoding/json"
	"fmt"
	"regexp"

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
	validateUnionFields   *Statement
}

func (r *Record) InnerTypes() IdentifierSet {
	innerTypes := make(IdentifierSet)
	for _, f := range r.Fields {
		innerTypes.AddAll(f.Type.InnerTypes())
	}

	return innerTypes
}

func (r *Record) PartialName() string {
	return r.Name + "_Partial"
}

func (r *Record) ProjectorName() string {
	return r.Name + "_Projector"
}

func (r *Record) PatchName() string {
	return r.Name + "_Patch"
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
		Add(r.generatePartialStruct()).Line().Line().
		Add(r.generateProjectorStruct()).Line().Line().
		Add(r.generatePatchStruct()).Line()
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
	hasUnionField := r.generateValidateUnionFields(def)

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
	}

	if hasUnionField {
		AddMarshalJSON(def, r.Receiver(), r.Name, func(def *Group) {
			def.Add(r.validateUnionFields)
			def.Type().Id("_t").Id(r.Name)
			def.Return(Qual(EncodingJson, Marshal).Call(Call(Op("*").Id("_t")).Call(Id(r.Receiver()))))
		}).Line().Line()
	}

	if hasDefaultValue || hasUnionField {
		AddUnmarshalJSON(def, r.Receiver(), r.Name, func(def *Group) {
			def.Type().Id("_t").Id(r.Name)
			def.Err().Op("=").Qual(EncodingJson, Unmarshal).Call(Id("data"), Call(Op("*").Id("_t")).Call(Id(r.Receiver())))
			IfErrReturn(def).Line()
			def.Add(r.populateDefaultValues, r.validateUnionFields)
			def.Return()
		}).Line().Line()
	}
	r.generateInitializeUnionFields(def)

	return def
}

func (r *Record) generatePartialStruct() *Statement {
	return AddWordWrappedComment(
		Empty(),
		fmt.Sprintf(
			"%s is %s where all fields can be optionally set. "+
				"Intended for use with patch requests and projections",
			r.PartialName(),
			r.Name,
		),
	).Line().Type().Id(r.PartialName()).
		StructFunc(func(def *Group) {
			for _, f := range r.Fields {
				field := def.Empty()
				AddWordWrappedComment(field, f.Doc).Line()
				field.Id(f.FieldName())

				// Unions, maps and arrays have in-built nil values, no need for a pointer
				if !f.Type.IsUnion() && !f.Type.IsMapOrArray() {
					field.Add(f.Type.PointerType())
				} else {
					field.Add(f.Type.GoType())
				}

				field.Tag(JsonFieldTag(f.Name, true))
			}
		})
}

// GenerateSelected gathers the selected fields and returns a []string
func GenerateSelected(t *RestliType, def *Group, accessor *Statement, name string) {
	switch {
	case t.IsUnion():
		for _, f := range *t.Union {
			GenerateSelected(&f.Type, def, accessor.Clone().Dot(f.name()), name+"."+f.Alias)
		}
	case t.IsRecord():
		def.For(List(Id("_"), Id(FieldVar)).Op(":=").Range().Add(accessor).Dot(Selected).Call()).
			BlockFunc(func(def *Group) {
				def.Id(SelectedVar).Op("=").Append(Id(SelectedVar), Lit(name+".").Op("+").Id(FieldVar))
			})
	default:
		def.If(accessor).BlockFunc(func(def *Group) {
			def.Id(SelectedVar).Op("=").Append(Id(SelectedVar), Lit(name))
		}).Line()
	}
}

func ProjectorGoType(t *RestliType) *Statement {
	switch {
	case t.IsUnion():
		return StructFunc(func(def *Group) {
			for _, m := range *t.Union {
				field := def.Empty()
				field.Id(m.name())
				field.Add(ProjectorGoType(&m.Type))
				field.Tag(JsonFieldTag(m.Alias, true))
			}
		})
	case t.IsRecord():
		return Qual(t.Reference.PackagePath(), t.Reference.Name+"_Projector")
	default:
		return Bool()
	}
}

func (r *Record) generateProjectorStruct() *Statement {
	def := Empty()
	// Generate the struct
	AddWordWrappedComment(def,
		fmt.Sprintf(
			"%s is used to represent a selection of the fields from %s. "+
				"Toggling the value of a field represents selecting it for a projection",
			r.ProjectorName(),
			r.Name,
		),
	).Line().Type().Id(r.ProjectorName()).
		StructFunc(func(def *Group) {
			for _, f := range r.Fields {
				field := def.Empty()
				AddWordWrappedComment(field, f.Doc).Line()
				field.Id(f.FieldName())
				field.Add(ProjectorGoType(&f.Type))
				field.Tag(JsonFieldTag(f.Name, true))
			}
		}).Line().Line()
	// Add the selected method
	receiver := ReceiverName(r.Name)
	AddFuncOnReceiver(def, receiver, r.ProjectorName(), Selected).Params().Index().String().
		BlockFunc(func(def *Group) {
			def.Var().Id(SelectedVar).Index().String()
			for _, f := range r.Fields {
				GenerateSelected(&f.Type, def, Id(receiver).Dot(f.FieldName()), f.Name)
			}
			def.Return(Id(SelectedVar))
		}).Line().Line()
	// MarshalJSON calls gatherFields with no prefix
	AddMarshalJSON(def, receiver, r.ProjectorName(), func(def *Group) {
		def.Return(Qual(EncodingJson, Marshal).Call(Id(receiver).Dot(Selected).Call()))
	}).Line().Line()
	// String calls gatherFields with no prefix then strings.Join with a comma
	AddStringer(def, receiver, r.ProjectorName(), func(def *Group) {
		def.Return(Qual("fmt", "Sprint").Call(Id(receiver).Dot(Selected).Call()))
	}).Line().Line()
	return def
}

func (r *Record) generatePatchStruct() *Statement {
	return AddWordWrappedComment(
		Empty(),
		fmt.Sprintf("%s applies a patch to *%s", r.PatchName(), r.Name),
	).Line().Type().Id(r.PatchName()).
		StructFunc(func(def *Group) {
			def.Id("Update").
				Qual(r.PackagePath(), r.PartialName()).
				Tag(JsonFieldTag("$set", true))
			def.Id("Delete").
				Qual(r.PackagePath(), r.ProjectorName()).
				Tag(JsonFieldTag("$delete", true))
		})
}

func (r *Record) generateRestliEncoder() *Statement {
	return AddRestLiEncode(Empty(), r.Receiver(), r.Name, func(def *Group) {
		def.Add(r.populateDefaultValues, r.validateUnionFields)

		const needsCommaVar = "needsComma"
		usesNeedsComma := false
		for _, f := range r.Fields[:len(r.Fields)-1] {
			usesNeedsComma = usesNeedsComma || f.IsOptionalOrDefault()
		}

		if usesNeedsComma {
			def.Id(needsCommaVar).Op(":=").False()
		}
		def.Var().Id("buf").Qual("strings", "Builder")
		def.Id("buf").Dot("WriteByte").Call(LitRune('('))

		for i, f := range r.Fields {
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
					writeComma := Id("buf").Dot("WriteByte").Call(LitRune(','))
					if r.Fields[i-1].IsOptionalOrDefault() && usesNeedsComma {
						def.If(Id(needsCommaVar)).Block(writeComma)
					} else {
						def.Add(writeComma)
					}
				}

				accessor := r.field(f)
				if f.IsPointer() && f.Type.Reference == nil && f.Type.Union == nil {
					accessor = Op("*").Add(accessor)
				}

				def.Id("buf").Dot("WriteString").Call(Lit(f.Name + ":"))
				f.Type.WriteToBuf(def, accessor)

				if usesNeedsComma && f.IsOptionalOrDefault() && i < len(r.Fields)-1 {
					def.Id(needsCommaVar).Op("=").True()
				}
			})
			serialize.Line()
		}
		def.Id("buf").Dot("WriteByte").Call(LitRune(')'))

		def.Id("data").Op("=").Id("buf").Dot("String").Call()
		def.Return()
	})
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

func (r *Record) generateValidateUnionFields(def *Statement) bool {
	r.validateUnionFields = Empty()

	hasUnion := false
	for _, f := range r.Fields {
		if f.Type.Union != nil {
			hasUnion = true
			break
		}
	}
	if !hasUnion {
		return false
	}

	AddFuncOnReceiver(def, r.Receiver(), r.Name, ValidateUnionFields).
		Params().
		Params(Error()).
		BlockFunc(func(def *Group) {
			for _, f := range r.Fields {
				if union := f.Type.Union; union != nil {
					def.BlockFunc(func(def *Group) {
						if f.IsPointer() {
							def.If(Id(r.Receiver()).Dot(f.FieldName()).Op("==").Nil()).
								Block(Return(Nil())).Line()
						}

						union.validateUnionFields(def, Id(r.Receiver()).Dot(f.FieldName()))
					})
				}
			}
			def.Return(Nil())
		}).Line().Line()

	r.validateUnionFields.Err().Op("=").Id(r.Receiver()).Dot(ValidateUnionFields).Call().Line()
	r.validateUnionFields.If(Err().Op("!=").Nil()).Block(Return()).Line()

	return true
}

func (r *Record) generateInitializeUnionFields(def *Statement) {
	for _, f := range r.Fields {
		if union := f.Type.Union; union != nil && f.IsPointer() {
			AddFuncOnReceiver(def, r.Receiver(), r.Name, "Initialize"+f.FieldName()).
				Params().
				Block(Id(r.Receiver()).Dot(f.FieldName()).Op("=").New(union.GoType()))
		}
	}
}

func (r *Record) defaultValuesConstructor() string {
	return "New" + r.Name + "WithDefaultValues"
}
