package models

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	. "github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

var (
	emptyMapRegex   = regexp.MustCompile("{ *}")
	emptyArrayRegex = regexp.MustCompile("\\[ *]")
)

const RecordTypeModelTypeName = "record"

type RecordModel struct {
	Identifier
	Doc string

	Include []*Model
	Fields  []Field

	populateDefaultValues *Statement
	validateUnionFields   *Statement
}

type Field struct {
	Name     string          `json:"name"`
	Doc      string          `json:"doc"`
	Type     *Model          `json:"type"`
	Optional bool            `json:"optional"`
	Default  json.RawMessage `json:"default"`
}

func (r *RecordModel) UnmarshalJSON(data []byte) error {
	t := &struct {
		Identifier
		typeField
		docField
		Include []*Model `json:"include"`
		Fields  []Field  `json:"fields"`
	}{}
	if err := json.Unmarshal(data, t); err != nil {
		return err
	}
	if t.Type != RecordTypeModelTypeName {
		return errors.Errorf("Not a record type: %s", string(data))
	}
	r.Identifier = t.Identifier
	r.Doc = t.Doc
	r.Include = t.Include
	r.Fields = t.Fields
	return nil
}

func (r *RecordModel) innerModels() (models []*Model) {
	models = append(models, r.Include...)
	for _, f := range r.Fields {
		models = append(models, f.Type)
	}
	return models
}

func (f *Field) IsPointer() bool {
	if f.Optional || f.Default != nil {
		return true
	}
	if u, ok := f.Type.BuiltinType.(*UnionModel); ok {
		return u.IsOptional
	}
	return false
}

func (r *RecordModel) allFields() (allFields []Field) {
	for _, i := range r.Include {
		if rec, ok := i.ComplexType.(*RecordModel); ok {
			allFields = append(allFields, rec.Fields...)
		} else {
			log.Panic("Illegal included type:", i)
		}
	}
	allFields = append(allFields, r.Fields...)
	return allFields
}

func (r *RecordModel) field(fieldName string) *Statement {
	return Id(r.receiver()).Dot(ExportedIdentifier(fieldName))
}

func (r *RecordModel) GenerateCode() (def *Statement) {
	def = Empty()

	AddWordWrappedComment(def, r.Doc).Line()

	def.Type().Id(r.Name).StructFunc(func(def *Group) {
		for _, f := range r.allFields() {
			field := def.Empty()
			AddWordWrappedComment(field, f.Doc).Line()
			field.Id(ExportedIdentifier(f.Name))

			var tag FieldTag
			tag.Json.Name = f.Name
			if f.IsPointer() {
				tag.Json.Optional = true
				field.Add(f.Type.PointerType())
			} else {
				field.Add(f.Type.GoType())
			}

			field.Tag(tag.ToMap())
		}
	}).Line().Line()

	hasDefaultValue := r.generatePopulateDefaultValues(def)
	hasUnionField := r.generateValidateUnionFields(def)

	def.Func().
		Id("New" + r.Name).Params().
		Params(Id(r.receiver()).Op("*").Id(r.Name))
	def.BlockFunc(func(def *Group) {
		def.Id(r.receiver()).Op("=").New(Id(r.Name))
		for _, f := range r.allFields() {
			if record, ok := f.Type.ComplexType.(*RecordModel); ok && !f.IsPointer() {
				def.Add(r.field(f.Name)).Op("=").Op("*").Qual(record.PackagePath(), "New"+record.Name).Call()
			}
		}
		def.Add(r.populateDefaultValues)
		def.Return()
	}).Line().Line()

	if hasDefaultValue || hasUnionField {
		r.jsonSerDe(def)
	}
	r.restLiSerDe(def)

	return def
}

func (r *RecordModel) restLiSerDe(def *Statement) {
	AddRestLiEncode(def, r.receiver(), r.Name, func(def *Group) {
		def.Add(r.populateDefaultValues, r.validateUnionFields)

		def.Var().Id("buf").Qual("strings", "Builder")
		def.Id("buf").Dot("WriteByte").Call(LitRune('('))

		allFields := r.allFields()
		for i, f := range allFields {
			serialize := def.Empty()
			if f.IsPointer() {
				serialize.If(r.field(f.Name).Op("!=").Nil())
			}

			serialize.BlockFunc(func(def *Group) {
				accessor := r.field(f.Name)
				if f.IsPointer() && f.Type.IsBytesOrPrimitive() {
					accessor = Op("*").Add(accessor)
				}

				if i != 0 {
					def.Id("buf").Dot("WriteByte").Call(LitRune(','))
				}

				def.Id("buf").Dot("WriteString").Call(Lit(f.Name + ":"))
				f.Type.restLiWriteToBuf(def, accessor)
			})
			serialize.Line()
		}
		def.Id("buf").Dot("WriteByte").Call(LitRune(')'))

		def.Id("data").Op("=").Id("buf").Dot("String").Call()
		def.Return()
	}).Line().Line()
}

func (r *RecordModel) jsonSerDe(def *Statement) {
	AddMarshalJSON(def, r.receiver(), r.Name, func(def *Group) {
		def.Add(r.populateDefaultValues, r.validateUnionFields)
		def.Type().Id("_t").Id(r.Name)
		def.Return(Qual(EncodingJson, Marshal).Call(Call(Op("*").Id("_t")).Call(Id(r.receiver()))))
	}).Line().Line()

	AddUnmarshalJSON(def, r.receiver(), r.Name, func(def *Group) {
		def.Type().Id("_t").Id(r.Name)
		def.Err().Op("=").Qual(EncodingJson, Unmarshal).Call(Id("data"), Call(Op("*").Id("_t")).Call(Id(r.receiver())))
		IfErrReturn(def).Line()
		def.Add(r.populateDefaultValues, r.validateUnionFields)
		def.Return()
	}).Line().Line()
}

func (r *RecordModel) setDefaultValue(def *Group, name, rawJson string, model *Model) {
	def.If(Id(r.receiver()).Dot(name).Op("==").Nil()).BlockFunc(func(def *Group) {
		// Special case for primitives, instead of parsing them from JSON every time, we can leave them as literals
		if primitive, ok := model.BuiltinType.(*PrimitiveModel); ok {
			def.Id("val").Op(":=").Lit(primitive.getLit(rawJson))
			def.Id(r.receiver()).Dot(name).Op("=").Op("&").Id("val")
			return
		}

		// Empty arrays and maps can be initialized directly, regardless of type
		_, isArray := model.BuiltinType.(*ArrayModel)
		_, isMap := model.BuiltinType.(*MapModel)
		if (isArray && emptyArrayRegex.MatchString(rawJson)) || (isMap && emptyMapRegex.MatchString(rawJson)) {
			def.Id(r.receiver()).Dot(name).Op("=").Make(model.GoType(), Lit(0))
			return
		}

		// Enum values can also be added as literals
		if enum, ok := model.ComplexType.(*EnumModel); ok {
			var v string
			err := json.Unmarshal([]byte(rawJson), &v)
			if err != nil {
				log.Panicln("illegal enum", err)
			}
			def.Id("val").Op(":=").Qual(enum.PackagePath(), enum.SymbolIdentifier(v))
			def.Id(r.receiver()).Dot(name).Op("= &").Id("val")
			return
		}

		if !model.IsMapOrArray() {
			def.Id(r.receiver()).Dot(name).Op("=").New(model.GoType())
		}

		field := Empty()
		if model.IsMapOrArray() {
			field.Op("&")
		}
		field.Id(r.receiver()).Dot(name)

		def.Err().Op(":=").Qual(EncodingJson, Unmarshal).Call(Index().Byte().Call(Lit(rawJson)), field)
		def.If(Err().Op("!=").Nil()).Block(Qual("log", "Panicln").Call(Lit("Illegal default value"), Err()))
	})
}

func (r *RecordModel) generatePopulateDefaultValues(def *Statement) bool {
	r.populateDefaultValues = Empty()

	hasDefault := false
	for _, f := range r.allFields() {
		if f.Default != nil {
			hasDefault = true
			break
		}
	}
	if !hasDefault {
		return false
	}

	AddFuncOnReceiver(def, r.receiver(), r.Name, PopulateDefaultValues).Params().BlockFunc(func(def *Group) {
		for _, f := range r.allFields() {
			if f.Default != nil {
				r.setDefaultValue(def, ExportedIdentifier(f.Name), string(f.Default), f.Type)
				def.Line()
			}
		}
	}).Line().Line()

	r.populateDefaultValues.Id(r.receiver()).Dot(PopulateDefaultValues).Call().Line()
	return true
}

func (r *RecordModel) generateValidateUnionFields(def *Statement) bool {
	r.validateUnionFields = Empty()

	hasUnion := false
	for _, f := range r.allFields() {
		if _, ok := f.Type.BuiltinType.(*UnionModel); ok {
			hasUnion = true
			break
		}
	}
	if !hasUnion {
		return false
	}

	AddFuncOnReceiver(def, r.receiver(), r.Name, ValidateUnionFields).
		Params().
		Params(Err().Error()).
		BlockFunc(func(def *Group) {
			for _, f := range r.allFields() {
				if union, ok := f.Type.BuiltinType.(*UnionModel); ok {
					def.BlockFunc(func(def *Group) {
						if f.IsPointer() {
							def.If(Id(r.receiver()).Dot(ExportedIdentifier(f.Name)).Op("==").Nil()).
								Block(Return(Nil())).Line()
						}

						isSet := "is" + ExportedIdentifier(f.Name) + "Set"
						def.Id(isSet).Op(":=").False().Line()
						errorMessage := fmt.Sprintf("must specify exactly one member of %s.%s",
							r.Name, ExportedIdentifier(f.Name))

						for i, u := range union.Types {
							cond := Id(r.receiver()).Dot(ExportedIdentifier(f.Name)).Dot(u.name()).Op("!=").Nil()
							def.If(cond).BlockFunc(func(def *Group) {
								if i == 0 {
									def.Id(isSet).Op("=").True()
								} else {
									def.If(Op("!").Id(isSet)).BlockFunc(func(def *Group) {
										def.Id(isSet).Op("=").True()
									}).Else().BlockFunc(func(def *Group) {
										def.Err().Op("=").Qual("fmt", "Errorf").Call(Lit(errorMessage))
										def.Return()
									})
								}
							}).Line()
						}
						def.If(Op("!").Id(isSet)).BlockFunc(func(def *Group) {
							def.Err().Op("=").Qual("fmt", "Errorf").Call(Lit(errorMessage))
							def.Return()
						})
					})
				}
			}
			def.Return()
		}).Line().Line()

	r.validateUnionFields.Err().Op("=").Id(r.receiver()).Dot(ValidateUnionFields).Call().Line()
	r.validateUnionFields.If(Err().Op("!=").Nil()).Block(Return()).Line()

	return true
}
