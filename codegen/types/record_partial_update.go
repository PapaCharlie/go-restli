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

	isEmpty := len(r.Fields) == 0
	fields := r.SortedFields()

	deletableFieldsStruct := r.generatePartialUpdateDeleteFieldsStruct()
	if deletableFieldsStruct != nil {
		def.Add(deletableFieldsStruct).Line()
	}
	def.Add(r.generatePartialUpdateSetFieldsStruct())

	// Generate the struct
	utils.AddWordWrappedComment(def, comment).Line()

	def.Type().Id(r.PartialUpdateStructName()).StructFunc(func(def *Group) {
		if !isEmpty {
			if deletableFieldsStruct != nil {
				def.Id(DeleteFields).Id(r.PartialUpdateDeleteFieldsStructName())
			}

			def.Id(SetFields).Id(r.PartialUpdateSetFieldsStructName())

			for _, f := range r.Fields {
				if record := f.Type.Record(); record != nil {
					def.Id(f.FieldName()).Op("*").Add(record.PartialUpdateStruct())
				}
			}
		}
	}).Line().Line()

	deleteAccessor := Code(Id(r.Receiver()).Dot(DeleteFields))
	setAccessor := Code(Id(r.Receiver()).Dot(SetFields))
	deleteAccessorF := func(f Field) *Statement { return Add(deleteAccessor).Dot(f.FieldName()) }
	setAccessorF := func(f Field) *Statement { return Add(setAccessor).Dot(f.FieldName()) }

	checker := Code(Id("checker"))
	checkAllFields := func(def *Group, keyChecker Code) {
		def.Add(checker).Op(":=").Qual(utils.RestLiPatchPackage, "PartialUpdateFieldChecker").Values(Dict{
			Id("RecordType"): Lit(r.Identifier.String()),
		})
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
				keyChecker,
				Lit(f.Name),
				d,
				setAccessorF(f).Op("!=").Nil(),
				p,
			), Err().Op("!=").Nil()).Block(Return(Err()))
		}
	}

	patch := Qual(utils.RestLiPatchPackage, "PatchField")

	const marshalPatch = "MarshalRestLiPatch"
	utils.AddFuncOnReceiver(def, r.Receiver(), r.PartialUpdateStructName(), marshalPatch, utils.Yes).
		Params(WriterParam).
		Params(Err().Error()).
		BlockFunc(func(def *Group) {
			def.Return(Writer.WriteMap(Writer, func(keyWriter Code, def *Group) {
				if isEmpty {
					def.Return(Nil())
					return
				}
				checkAllFields(def, Writer)
				marshal := func(def *Group, source Code, name string) {
					def.Err().Op("=").Add(source).Dot(utils.MarshalRestLi).Call(Add(keyWriter).Call(Lit(name)))
					def.Add(utils.IfErrReturn(Err()))
				}
				if deletableFieldsStruct != nil {
					def.If(Add(checker).Dot("HasDeletes")).BlockFunc(func(def *Group) {
						marshal(def, deleteAccessor, "$delete")
					}).Line()
				}

				def.If(Add(checker).Dot("HasSets")).BlockFunc(func(def *Group) {
					marshal(def, setAccessor, "$set")
				}).Line()

				for _, f := range fields {
					if f.Type.Record() != nil {
						def.Line()
						accessor := r.rawFieldAccessor(f)
						def.If(Add(accessor).Op("!=").Nil()).BlockFunc(func(def *Group) {
							def.Err().Op("=").Add(accessor).Dot(marshalPatch).Call(Add(keyWriter).Call(Lit(f.Name)))
							def.Add(utils.IfErrReturn(Err()))
						})
					}
				}
				def.Line()

				def.Return(Nil())
			}))
		}).Line().Line()
	AddMarshalRestLi(def, r.Receiver(), r.PartialUpdateStructName(), RecordShouldUsePointer, func(def *Group) {
		def.Return(Writer.WriteMap(Writer, func(keyWriter Code, def *Group) {
			def.Return(Id(r.Receiver()).Dot(marshalPatch).Call(Add(keyWriter).Call(patch).Dot("SetScope").Call()))
		}))
	})

	const unmarshalPatch = "UnmarshalRestLiPatch"
	utils.AddFuncOnReceiver(def, r.Receiver(), r.PartialUpdateStructName(), unmarshalPatch, utils.Yes).
		Params(ReaderParam).
		Params(Err().Error()).
		BlockFunc(func(def *Group) {
			if isEmpty {
				def.Return(Add(Reader.ReadMap(Reader, func(reader, key Code, def *Group) {
					def.Return(Reader.Skip(reader))
				})))
				return
			}

			def.Err().Op("=").Add(Reader.ReadMap(Reader, func(reader, key Code, def *Group) {
				def.Switch(key).BlockFunc(func(def *Group) {
					unmarshal := func(def *Group, accessor Code) {
						def.Err().Op("=").Add(accessor).Dot(utils.UnmarshalRestLi).Call(reader)
					}
					if deletableFieldsStruct != nil {
						def.Case(Lit("$delete")).BlockFunc(func(def *Group) {
							unmarshal(def, deleteAccessor)
						})
					}
					def.Case(Lit("$set")).BlockFunc(func(def *Group) {
						unmarshal(def, setAccessor)
					})

					for _, f := range fields {
						if rec := f.Type.Record(); rec != nil {
							accessor := r.rawFieldAccessor(f)
							def.Case(Lit(f.Name)).BlockFunc(func(def *Group) {
								def.Add(accessor).Op("=").New(rec.PartialUpdateStruct())
								def.Err().Op("=").Add(accessor).Dot(unmarshalPatch).Call(reader)
							})
						}
					}

					def.Default().BlockFunc(func(def *Group) {
						def.Err().Op("=").Add(Reader.Skip(reader))
					})
				})
				def.Return(Err())
			}))
			def.Add(utils.IfErrReturn(Err()))

			checkAllFields(def, Reader)

			def.Return(Nil())
		}).Line().Line()

	AddUnmarshalRestli(def, r.Receiver(), r.PartialUpdateStructName(), RecordShouldUsePointer, func(def *Group) {
		def.Return(Add(Reader.ReadRecord(Reader, Qual(utils.RestLiPatchPackage, "RequiredPatchRecordFields"), func(reader, key Code, def *Group) {
			def.If(Add(key).Op("==").Add(patch)).Block(
				Return(Id(r.Receiver()).Dot(unmarshalPatch).Call(reader)),
			).Else().Block(
				Return(Reader.Skip(reader)),
			)
		})))
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
	def.Type().Id(r.PartialUpdateDeleteFieldsStructName()).StructFunc(func(def *Group) {
		for _, f := range deletableFields {
			def.Id(f.FieldName()).Bool()
		}
	}).Line().Line()

	deleteType := RestliType{Primitive: &StringPrimitive}
	AddMarshalRestLi(def, r.Receiver(), r.PartialUpdateDeleteFieldsStructName(), RecordShouldUsePointer, func(def *Group) {
		def.Return(Writer.WriteArray(Writer, func(itemWriter Code, def *Group) {
			for _, f := range fields {
				def.If(r.rawFieldAccessor(f)).Block(
					Writer.Write(deleteType, Add(itemWriter).Call(), Lit(f.Name)),
				)
			}
			def.Add(Return(Nil()))
		}))
	})

	AddUnmarshalRestli(def, r.Receiver(), r.PartialUpdateDeleteFieldsStructName(), RecordShouldUsePointer, func(def *Group) {
		field := Id("field")
		def.Var().Add(field).String()

		def.Return(Reader.ReadArray(Reader, func(itemReader Code, def *Group) {
			def.Add(Reader.Read(deleteType, itemReader, field))
			def.Add(utils.IfErrReturn(Err())).Line()

			def.Switch(field).BlockFunc(func(def *Group) {
				for _, f := range r.Fields {
					if !f.IsOptionalOrDefault() {
						continue
					}
					def.Case(Lit(f.Name)).Block(r.rawFieldAccessor(f).Op("=").True())
				}
			})
			def.Add(Return(Nil()))
		}))
	})

	return def
}

func (r *Record) generatePartialUpdateSetFieldsStruct() *Statement {
	setRecord := &Record{
		NamedType: NamedType{
			Identifier: utils.Identifier{
				Namespace: r.Namespace,
				Name:      r.PartialUpdateSetFieldsStructName(),
			},
			SourceFile: r.SourceFile,
			Doc:        "",
		},
	}

	for _, f := range r.Fields {
		setRecord.Fields = append(setRecord.Fields, Field{
			Type:       f.Type,
			Name:       f.Name,
			Doc:        f.Name,
			IsOptional: true,
		})
	}

	return Empty().
		Add(setRecord.GenerateStruct()).Line().Line().
		Add(setRecord.GenerateMarshalRestLi()).Line().Line().
		Add(setRecord.GenerateUnmarshalRestLi()).Line().Line()
}
