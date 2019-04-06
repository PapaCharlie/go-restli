package models

import (
	"encoding/json"
	"fmt"
	"log"

	. "github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
)

const RecordTypeModelTypeName = "record"

type RecordModel struct {
	NameAndDoc
	Include []*Model
	Fields  []Field

	populateDefaultValues *Statement
	validateUnionFields   *Statement
}

type Field struct {
	NameAndDoc
	Type     *Model          `json:"type"`
	Optional bool            `json:"optional"`
	Default  json.RawMessage `json:"default"`
}

func (f *Field) isPointer() bool {
	return f.Optional || f.Default != nil
}

func (r *RecordModel) allFields() (allFields []Field) {
	for _, i := range r.Include {
		if i.Record == nil {
			log.Panic("Illegal included type:", i)
		}
		allFields = append(allFields, i.Record.Fields...)
	}
	allFields = append(allFields, r.Fields...)
	return
}

func (r *RecordModel) receiver() string {
	return ReceiverName(r.Name)
}

func (r *RecordModel) InnerModels() (models []*Model) {
	models = append(models, r.Include...)
	for _, f := range r.Fields {
		models = append(models, f.Type)
	}
	return
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
			if f.isPointer() {
				tag.Json.Optional = true
				field.Add(f.Type.PointerType())
			} else if f.Type.Union != nil {
				field.Add(f.Type.Union.GoType())
			} else {
				field.Add(f.Type.GoType())
			}

			if f.Type.Union != nil {
				tag.RestLi.Union = true
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
			if f.Type.Record != nil && !f.isPointer() {
				def.Add(r.field(f.Name)).Op("=").Op("*").Qual(f.Type.PackagePath(), "New"+f.Type.Record.Name).Call()
			}
		}
		def.Add(r.populateDefaultValues)
		def.Return()
	}).Line().Line()

	if hasDefaultValue || hasUnionField {
		r.jsonSerDe(def)
	}
	r.restLiSerDe(def)

	return
}

func (r *RecordModel) restLiSerDe(def *Statement) {
	AddRestLiEncode(def, r.receiver(), r.Name, func(def *Group) {
		def.Add(r.populateDefaultValues, r.validateUnionFields)

		def.Var().Id("buf").Qual("strings", "Builder")
		def.Id("buf").Dot("WriteByte").Call(LitRune('('))

		allFields := r.allFields()
		for i, f := range allFields {
			serialize := def.Empty()
			if f.isPointer() {
				serialize.If(r.field(f.Name).Op("!=").Nil())
			}

			serialize.BlockFunc(func(def *Group) {
				accessor := r.field(f.Name)
				if f.isPointer() && (f.Type.Primitive != nil || f.Type.Bytes != nil) {
					accessor = Op("*").Add(accessor)
				}

				if i != 0 {
					def.Id("buf").Dot("WriteByte").Call(LitRune(','))
				}

				def.Id("buf").Dot("WriteString").Call(Lit(f.Name + ":"))

				if f.Type.Union != nil {
					isSet := "is" + ExportedIdentifier(f.Name) + "Set"
					def.Id(isSet).Op(":=").False().Line()
					errorMessage := fmt.Sprintf("must specify exactly one member of %s.%s", r.Name, ExportedIdentifier(f.Name))

					for j, u := range f.Type.Union.Types {
						unionFieldAccessor := Id(r.receiver()).Dot(ExportedIdentifier(f.Name)).Dot(u.name())
						def.If(Id(r.receiver()).Dot(ExportedIdentifier(f.Name)).Dot(u.name()).Op("!=").Nil()).BlockFunc(func(def *Group) {
							if j == 0 {
								def.Id(isSet).Op("=").True()
								writeToBuf(def, Lit("("+u.alias()+":"))
								u.Model.writeToBuf(def, unionFieldAccessor)
								def.Id("buf").Dot("WriteByte").Call(LitRune(')'))
							} else {
								def.If(Id(isSet)).BlockFunc(func(def *Group) {
									def.Err().Op("=").Qual("fmt", "Errorf").Call(Lit(errorMessage))
									def.Return()
								}).Else().BlockFunc(func(def *Group) {
									def.Id(isSet).Op("=").True()
									writeToBuf(def, Lit("("+u.alias()+":"))
									u.Model.writeToBuf(def, unionFieldAccessor)
									def.Id("buf").Dot("WriteByte").Call(LitRune(')'))
								})
							}
						}).Line()
					}
					def.If(Op("!").Id(isSet)).BlockFunc(func(def *Group) {
						def.Err().Op("=").Qual("fmt", "Errorf").Call(Lit(errorMessage))
						def.Return()
					})
				} else {
					f.Type.writeToBuf(def, accessor)
				}

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
		if model.Primitive != nil {
			def.Id("val").Op(":=").Lit(model.Primitive.GetLit(rawJson))
			def.Id(r.receiver()).Dot(name).Op("=").Op("&").Id("val")
			return
		}

		// Empty arrays and maps can be initialized directly, regardless of type
		if (model.Array != nil && rawJson == "[]") || (model.Map != nil && rawJson == "{}") {
			def.Id(r.receiver()).Dot(name).Op("=").Make(model.GoType(), Lit(0))
			return
		}

		// Enum values can also be added as literals
		if model.Enum != nil {
			var v string
			err := json.Unmarshal([]byte(rawJson), &v)
			if err != nil {
				log.Panicln("illegal enum", err)
			}
			def.Id("val").Op(":=").Qual(model.PackagePath(), model.Enum.SymbolIdentifier(v))
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
		if f.Type.Union != nil {
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
				if f.Type.Union != nil {
					isSet := "is" + ExportedIdentifier(f.Name) + "Set"
					def.Id(isSet).Op(":=").False().Line()
					errorMessage := fmt.Sprintf("must specify exactly one member of %s.%s", r.Name, ExportedIdentifier(f.Name))

					for i, u := range f.Type.Union.Types {
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
				}
			}
			def.Return()
		}).Line().Line()

	r.validateUnionFields.Err().Op("=").Id(r.receiver()).Dot(ValidateUnionFields).Call().Line()
	r.validateUnionFields.If(Err().Op("!=").Nil()).Block(Return())

	return true
}

func (r *RecordModel) field(fieldName string) *Statement {
	return Id(r.receiver()).Dot(ExportedIdentifier(fieldName))
}
