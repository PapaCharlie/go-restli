package types

import (
	"encoding/json"
	"regexp"
	"sort"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

var (
	emptyMapRegex   = regexp.MustCompile("{ *}")
	emptyArrayRegex = regexp.MustCompile("\\[ *]")
)

const RecordShouldUsePointer = utils.Yes

type Record struct {
	NamedType
	Fields []Field
}

func (r *Record) InnerTypes() utils.IdentifierSet {
	innerTypes := make(utils.IdentifierSet)
	for _, f := range r.Fields {
		innerTypes.AddAll(f.Type.InnerTypes())
	}

	return innerTypes
}

func (r *Record) ShouldReference() utils.ShouldUsePointer {
	return RecordShouldUsePointer
}

func (r *Record) PartialUpdateStructName() string {
	return r.Name + PartialUpdate
}

func (r *Record) PartialUpdateDeleteFieldsStructName() string {
	return r.Name + PartialUpdate + "_" + DeleteFields
}

func (r *Record) PartialUpdateSetFieldsStructName() string {
	return r.Name + PartialUpdate + "_" + SetFields
}

func (r *Record) PartialUpdateStruct() *Statement {
	return Qual(r.PackagePath(), r.PartialUpdateStructName())
}

type Field struct {
	Type               RestliType
	Name               string
	Doc                string
	IsOptional         bool
	DefaultValue       *string
	IncludedFrom       *utils.Identifier
	isComplexKeyParams bool
}

func (r *Record) fieldAccessor(f Field) *Statement {
	return fieldAccessor(Id(r.Receiver()), f)
}

func (r *Record) rawFieldAccessor(f Field) *Statement {
	return Id(r.Receiver()).Dot(f.FieldName())
}

func fieldAccessor(receiver Code, f Field) *Statement {
	accessor := Add(receiver)
	if r := f.ResolveRecord(); r != nil {
		accessor.Dot(r.Name)
	}
	return accessor.Dot(f.FieldName())
}

func (f *Field) IsOptionalOrDefault() bool {
	return f.IsOptional || f.DefaultValue != nil
}

func (f *Field) FieldName() string {
	if f.isComplexKeyParams {
		return utils.ComplexKeyParamsField
	} else {
		return utils.ExportedIdentifier(f.Name)
	}
}

func (f *Field) ResolveRecord() *Record {
	if f.IncludedFrom == nil {
		return nil
	} else {
		// This is guaranteed to be a record
		return f.IncludedFrom.Resolve().(*Record)
	}
}

func (r *Record) SortedFields() (fields []Field) {
	fields = append(fields, r.Fields...)
	sortFields(fields)
	return fields
}

func sortFields(fields []Field) {
	sort.Slice(fields, func(i, j int) bool { return fields[i].Name < fields[j].Name })
}

func (r *Record) GenerateCode() *Statement {
	return Empty().
		Add(r.GenerateStruct()).Line().Line().
		Add(r.GeneratePopulateDefaultValues()).Line().Line().
		Add(r.GenerateEquals()).Line().Line().
		Add(r.GenerateComputeHash()).Line().Line().
		Add(r.GenerateMarshalRestLi()).Line().Line().
		Add(r.GenerateUnmarshalRestLi()).Line().Line().
		Add(r.generatePartialUpdateStruct()).Line()
}

func (r *Record) GenerateStruct() *Statement {
	return utils.AddWordWrappedComment(Empty(), r.Doc).Line().
		Type().Id(r.Name).
		StructFunc(func(def *Group) {
			var uniqueIncludedRecords []utils.Identifier
			includedRecords := utils.IdentifierSet{}
			for _, f := range r.Fields {
				if f.IncludedFrom != nil && includedRecords.Add(*f.IncludedFrom) {
					uniqueIncludedRecords = append(uniqueIncludedRecords, *f.IncludedFrom)
				}
			}

			for _, id := range uniqueIncludedRecords {
				def.Add(id.Qual())
			}
			for _, f := range r.Fields {
				if f.IncludedFrom != nil {
					continue
				}
				field := def.Empty()
				utils.AddWordWrappedComment(field, f.Doc).Line()
				field.Id(f.FieldName())

				if f.IsOptionalOrDefault() {
					field.Add(f.Type.PointerType())
				} else {
					field.Add(f.Type.GoType())
				}
			}
		})
}

func (r *Record) GeneratePopulateDefaultValues() Code {
	if !r.hasDefaultValue() {
		return Empty()
	}

	def := Empty()

	def.Commentf("Sanity check %s has no illegal default values", r.defaultValuesConstructor()).Line()
	def.Var().Id("_").Op("=").Id(r.defaultValuesConstructor()).Call().Line().Line()

	def.Func().
		Id(r.defaultValuesConstructor()).Params().
		Params(Id(r.Receiver()).Op("*").Id(r.Name)).
		BlockFunc(func(def *Group) {
			def.Id(r.Receiver()).Op("=").New(Id(r.Name))
			for _, f := range r.Fields {
				if f.Type.Reference == nil {
					continue
				}
				if record := f.Type.Record(); record != nil && !f.IsOptionalOrDefault() && record.hasDefaultValue() {
					def.Add(r.fieldAccessor(f)).Op("=").Op("*").Qual(record.PackagePath(), record.defaultValuesConstructor()).Call()
				}
			}
			def.Id(r.Receiver()).Dot(utils.PopulateLocalDefaultValues).Call()
			def.Return()
		}).Line().Line()

	utils.AddFuncOnReceiver(def, r.Receiver(), r.Name, utils.PopulateLocalDefaultValues, RecordShouldUsePointer).
		Params().
		BlockFunc(func(def *Group) {
			for _, f := range r.Fields {
				if f.DefaultValue != nil {
					r.setDefaultValue(def, r.fieldAccessor(f), *f.DefaultValue, &f.Type)
					def.Line()
				}
			}
		}).Line().Line()

	return def
}

func (r *Record) setDefaultValue(def *Group, accessor Code, rawJson string, t *RestliType) {
	def.If(Add(accessor).Op("==").Nil()).BlockFunc(func(def *Group) {
		addPanic := func() {
			def.If(Err().Op("!=").Nil()).Block(Qual("log", "Panicln").Call(Lit("Illegal default value"), Err()))
		}
		declareReader := func() {
			def.List(Reader, Err()).Op(":=").Add(utils.NewJsonReader).Call(Index().Byte().Call(Lit(rawJson)))
			addPanic()
		}
		switch {
		// Special case for primitives, instead of parsing them from JSON every time, we can leave them as literals
		case t.Primitive != nil:
			def.Id("val").Op(":=").Add(t.Primitive.getLit(rawJson))
			def.Add(accessor).Op("= &").Id("val")
			return
		case t.Reference != nil:
			if enum := t.Enum(); enum != nil {
				var v string
				err := json.Unmarshal([]byte(rawJson), &v)
				if err != nil {
					utils.Logger.Panicln("illegal enum", err)
				}
				if !enum.isValidSymbol(v) {
					utils.Logger.Panicf("illegal enum value %q for %q (not in %q)", v, enum.Identifier, enum.Symbols)
				}
				def.Id("val").Op(":=").Qual(enum.PackagePath(), enum.SymbolIdentifier(v))
				def.Add(accessor).Op("= &").Id("val")
			} else if fixed, ok := t.Reference.Resolve().(*Fixed); ok {
				def.Id("val").Op(":=").Add(fixed.getLit(rawJson))
				def.Add(accessor).Op("= &").Id("val")
			} else if typeref := t.Typeref(); typeref != nil {
				def.Id("val").Op(":=").Add(t.GoType()).Call(typeref.Type.getLit(rawJson))
				def.Add(accessor).Op("= &").Id("val")
			} else if t.Record() != nil || t.StandaloneUnion() != nil {
				def.Add(accessor).Op("=").New(t.GoType())
				declareReader()
				def.Add(Reader.Read(*t, Reader, accessor))
				addPanic()
			} else {
				utils.Logger.Panic("Unknown reference type for default value", t.Reference.Resolve())
			}
		case t.Array != nil:
			def.Add(accessor).Op("=").New(t.GoType())
			if emptyArrayRegex.MatchString(rawJson) {
				return
			}
			declareReader()
			def.Add(Reader.Read(*t, Reader, Op("*").Add(accessor)))
			addPanic()
			return
		case t.Map != nil:
			def.Add(accessor).Op("= &").Add(t.GoType()).Values()
			if emptyMapRegex.MatchString(rawJson) {
				return
			}
			declareReader()
			def.Add(Reader.Read(*t, Reader, Op("*").Add(accessor)))
			addPanic()
			return
		}
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
