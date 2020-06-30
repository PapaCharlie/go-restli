package codegen

import (
	"encoding/json"
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
	validateUnionFields   *Statement
}

func (r *Record) InnerTypes() IdentifierSet {
	innerTypes := make(IdentifierSet)
	for _, f := range r.Fields {
		innerTypes.AddAll(f.Type.InnerTypes())
	}

	return innerTypes
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
		Add(r.generatePatchStructs()).Line()
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

func (r *Record) generatePatchStructs() *Statement {
	// TODO
	return Empty()
}

func (r *Record) generateRestliEncoder() *Statement {
	return AddRestLiEncode(Empty(), r.Receiver(), r.Name, func(def *Group) {
		r.generateEncoder(def, nil)
	})
}

func (r *Record) generateEncoder(def *Group, finderName *string) {
	if len(r.Fields) == 0 {
		def.Return(Lit(""), Nil())
		return
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

	def.Add(r.validateUnionFields)

	def.Var().Id("buf").Qual("strings", "Builder")
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

	bufLenCheckNeeded := len(fields) > 1 && fields[0].IsOptionalOrDefault()

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
				if bufLenCheckNeeded {
					def.If(Id("buf").Dot("Len").Call().Op("!=").Lit(0)).Block(writeDelimiter)
				} else {
					def.Add(writeDelimiter)
				}
			}

			if i == qIndex {
				def.Id("buf").Dot("WriteString").Call(Lit(finderNameParam + nameDelimiter + *finderName))
			} else {
				accessor := r.field(f)
				if f.IsPointer() && f.Type.Reference == nil && f.Type.Union == nil {
					accessor = Op("*").Add(accessor)
				}

				def.Id("buf").Dot("WriteString").Call(Lit(f.Name + nameDelimiter))
				f.Type.WriteToBuf(def, accessor)
			}
		})
		serialize.Line()

		if !f.IsOptionalOrDefault() {
			bufLenCheckNeeded = false
		}
	}
	if finderName == nil {
		def.Id("buf").Dot("WriteByte").Call(LitRune(')'))
	}

	def.Id("data").Op("=").Id("buf").Dot("String").Call()
	def.Return()
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
