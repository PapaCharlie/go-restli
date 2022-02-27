package types

import (
	"fmt"
	"strings"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

const (
	PartialUpdate = "_PartialUpdate"
	DeleteFields  = "Delete_Fields"
	SetFields     = "Set_Fields"
)

func (r *Record) generatePartialUpdateStruct() *Statement {
	delimiter := strings.Repeat("=", 80)
	def := Commentf("%s\nPARTIAL UPDATE STRUCTS\n%s", delimiter, delimiter).Line().Line()

	comment := fmt.Sprintf(
		"%s is used to represent a partial update on %s. Toggling the value of a field\n"+
			"in Delete represents selecting it for deletion in a partial update, while\n"+
			"setting the value of a field in Update represents setting that field in the\n"+
			"current struct. Other fields in this struct represent record fields that can\n"+
			"themselves be partially updated.",
		r.PartialUpdateStructName(), r.Name)

	if len(r.Fields) == 0 {
		utils.AddWordWrappedComment(def, comment).Line()
		def.Type().Id(r.PartialUpdateStructName()).Struct().Line()

		return AddMarshalRestLi(def, r.Receiver(), r.PartialUpdateStructName(), RecordShouldUsePointer, func(def *Group) {
			def.Return(Writer.WriteMap(Writer, func(keyWriter Code, def *Group) {
				def.Return(Nil())
			}))
		})
	}

	fields := r.SortedFields()

	deletableFieldsStruct := r.generatePartialUpdateDeleteFieldsStruct()
	if deletableFieldsStruct != nil {
		def.Add(deletableFieldsStruct).Line()
	}
	def.Add(r.generatePartialUpdateSetFieldsStruct())

	// Generate the struct
	utils.AddWordWrappedComment(def, comment).Line()

	def.Type().Id(r.PartialUpdateStructName()).StructFunc(func(def *Group) {
		if deletableFieldsStruct != nil {
			def.Id(DeleteFields).Id(r.PartialUpdateDeleteFieldsStructName())
		}

		def.Id(SetFields).Id(r.PartialUpdateSetFieldsStructName())

		for _, f := range r.Fields {
			if record := f.Type.Record(); record != nil {
				def.Id(f.FieldName()).Op("*").Add(record.PartialUpdateStruct())
			}
		}
	}).Line().Line()

	deleteAccessor := Code(Id(r.Receiver()).Dot(DeleteFields))
	setAccessor := Code(Id(r.Receiver()).Dot(SetFields))
	deleteAccessorF := func(f Field) *Statement { return Add(deleteAccessor).Dot(f.FieldName()) }
	setAccessorF := func(f Field) *Statement { return Add(setAccessor).Dot(f.FieldName()) }

	AddMarshalRestLi(def, r.Receiver(), r.PartialUpdateStructName(), RecordShouldUsePointer, func(def *Group) {
		checker := Id("checker")
		def.Add(checker).Op(":=").Qual(utils.ProtocolPackage, "PartialUpdateFieldChecker").Values(Dict{
			Id("RecordType"): Lit(r.Identifier.String()),
		})

		def.Return(Writer.WriteMap(Writer, func(keyWriter Code, def *Group) {
			for _, f := range r.Fields {
				var d Code
				if f.IsOptionalOrDefault() {
					d = deleteAccessorF(f)
				} else {
					d = False()
				}

				var p Code
				if f.Type.Record() != nil {
					p = r.rawFieldAccessor(f).Op("!=").Nil()
				} else {
					p = False()
				}

				def.If(Err().Op("=").Add(checker).Dot("CheckField").Call(
					Writer,
					Lit(f.Name),
					d,
					setAccessorF(f).Op("!=").Nil(),
					p,
				), Err().Op("!=").Nil()).Block(Return(Err()))
			}
			def.Line()

			if deletableFieldsStruct != nil {
				def.If(Add(checker).Dot("HasDeletes")).BlockFunc(func(def *Group) {
					def.Err().Op("=").Add(deleteAccessor).Dot(utils.MarshalRestLi).Call(Add(keyWriter).Call(Lit("$delete")))
					def.Add(utils.IfErrReturn(Err()))
				}).Line()
			}

			def.If(Add(checker).Dot("HasSets")).BlockFunc(func(def *Group) {
				def.Err().Op("=").Add(setAccessor).Dot(utils.MarshalRestLi).Call(Add(keyWriter).Call(Lit("$set")))
				def.Add(utils.IfErrReturn(Err()))
			}).Line()

			for _, f := range fields {
				if f.Type.Record() != nil {
					def.Line()
					accessor := r.rawFieldAccessor(f)
					def.If(Add(accessor).Op("!=").Nil()).BlockFunc(func(def *Group) {
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

func (r *Record) generatePartialUpdateDeleteFieldsStruct() *Statement {
	var deletableFields []Field
	for _, f := range r.Fields {
		if f.IsOptionalOrDefault() {
			deletableFields = append(deletableFields, f)
		}
	}

	if len(deletableFields) == 0 {
		return nil
	}

	fields := append([]Field(nil), deletableFields...)
	sortFields(fields)

	def := Empty()
	structName := r.PartialUpdateStructName() + "_" + DeleteFields
	def.Type().Id(structName).StructFunc(func(def *Group) {
		for _, f := range deletableFields {
			def.Id(f.FieldName()).Bool()
		}
	}).Line().Line()

	return AddMarshalRestLi(def, r.Receiver(), structName, RecordShouldUsePointer, func(def *Group) {
		def.Return(Writer.WriteArray(Writer, func(itemWriter Code, def *Group) {
			deleteType := RestliType{Primitive: &StringPrimitive}

			for _, f := range fields {
				def.If(r.rawFieldAccessor(f)).Block(
					Writer.Write(deleteType, Add(itemWriter).Call(), Lit(f.Name)),
				)
			}
			def.Add(Return(Nil()))
		}))
	})
}

func (r *Record) generatePartialUpdateSetFieldsStruct() *Statement {
	def := Empty()

	fields := r.SortedFields()

	def.Type().Id(r.PartialUpdateSetFieldsStructName()).StructFunc(func(def *Group) {
		for _, f := range r.Fields {
			def.Id(f.FieldName()).Add(Op("*").Add(f.Type.GoType()))
		}
	}).Line().Line()

	return AddMarshalRestLi(def, r.Receiver(), r.PartialUpdateSetFieldsStructName(), RecordShouldUsePointer, func(def *Group) {
		def.Return(Writer.WriteMap(Writer, func(keyWriter Code, def *Group) {
			for _, f := range fields {
				accessor := r.rawFieldAccessor(f)
				def.If(Add(accessor).Op("!=").Nil()).BlockFunc(func(def *Group) {
					if f.Type.Reference == nil {
						accessor = Op("*").Add(accessor)
					}
					def.Add(Writer.Write(f.Type, Add(keyWriter).Call(Lit(f.Name)), accessor, Err()))
				})
				def.Line()
			}
			def.Return(Nil())
		}))
	}).Line()
}
