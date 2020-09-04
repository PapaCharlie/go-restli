package types

import (
	"encoding/json"
	"regexp"
	"sort"

	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
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

func (r *Record) InnerTypes() utils.IdentifierSet {
	innerTypes := make(utils.IdentifierSet)
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
	return utils.ExportedIdentifier(f.Name)
}

func (r *Record) SortedFields() (fields []Field) {
	fields = append(fields, r.Fields...)
	sort.Slice(fields, func(i, j int) bool { return fields[i].Name < fields[j].Name })
	return fields
}

func (r *Record) GenerateCode() *Statement {
	return Empty().
		Add(r.GenerateStruct()).Line().Line().
		Add(r.generateDefaultValuesCode()).Line().Line().
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
			for _, f := range r.Fields {
				field := def.Empty()
				utils.AddWordWrappedComment(field, f.Doc).Line()
				field.Id(f.FieldName())

				if f.IsOptionalOrDefault() {
					field.Add(f.Type.PointerType())
				} else {
					field.Add(f.Type.GoType())
				}

				field.Tag(utils.JsonFieldTag(f.Name, f.IsOptionalOrDefault()))
			}
		})
}

func (r *Record) generateDefaultValuesCode() Code {
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
					def.Add(r.field(f)).Op("=").Op("*").Qual(record.PackagePath(), record.defaultValuesConstructor()).Call()
				}
			}
			def.Id(r.Receiver()).Dot(PopulateLocalDefaultValues).Call()
			def.Return()
		}).Line().Line()

	utils.AddFuncOnReceiver(def, r.Receiver(), r.Name, PopulateLocalDefaultValues).Params().BlockFunc(func(def *Group) {
		for _, f := range r.Fields {
			if f.DefaultValue != nil {
				r.setDefaultValue(def, f.FieldName(), *f.DefaultValue, &f.Type)
				def.Line()
			}
		}
	}).Line().Line()

	return def
}

func (r *Record) setDefaultValue(def *Group, name, rawJson string, t *RestliType) {
	def.If(Id(r.Receiver()).Dot(name).Op("==").Nil()).BlockFunc(func(def *Group) {
		accessor := Code(Id(r.Receiver()).Dot(name))
		declareReader := func() {
			def.Var().Err().Error()
			def.Add(Reader).Op(":=").Add(NewJsonReader).Call(Index().Byte().Call(Lit(rawJson)))
		}
		addPanic := func() {
			def.If(Err().Op("!=").Nil()).Block(Qual("log", "Panicln").Call(Lit("Illegal default value"), Err()))
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
				def.Id(r.Receiver()).Dot(name).Op("=").New(t.GoType())
				declareReader()
				def.Add(Reader.Read(*t, Reader, accessor))
				addPanic()
			} else {
				utils.Logger.Panic("Unknown reference type for default value", t.Reference.Resolve())
			}
		case t.Array != nil:
			// If the default value for an array is the empty array, we can leave it as nil since that will behave
			// identically to an empty slice
			if emptyArrayRegex.MatchString(rawJson) {
				return
			}
			def.Add(accessor).Op("=").New(t.GoType())
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
