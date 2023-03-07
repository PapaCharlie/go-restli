package types

import (
	"fmt"
	"strings"

	"github.com/PapaCharlie/go-restli/v2/codegen/utils"
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
			"in Delete_Field represents selecting it for deletion in a partial update, while\n"+
			"setting the value of a field in Set_Fields represents setting that field in the\n"+
			"current struct. Other fields in this struct represent record fields that can\n"+
			"themselves be partially updated.",
		r.PartialUpdateStructName(), r.TypeName())

	fields := r.SortedFields()

	def.Add(r.generatePartialUpdateDeleteFieldsStruct()).Line()
	def.Add(r.generatePartialUpdateSetFieldsStruct()).Line()

	// Generate the struct
	utils.AddWordWrappedComment(def, comment).Line()

	def.Type().Id(r.PartialUpdateStructName()).StructFunc(func(def *Group) {
		for _, ir := range r.includedRecords() {
			def.Qual(ir.PackagePath(), ir.PartialUpdateStructName())
		}
		def.Id(DeleteFields).Id(r.PartialUpdateDeleteFieldsStructName())
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

	const CheckFields = "CheckFields"
	fieldChecker := Code(Id("fieldChecker"))
	fieldCheckerType := Code(Qual(utils.RestLiPatchPackage, "PartialUpdateFieldChecker"))
	keyChecker := Code(Id("keyChecker"))
	utils.AddFuncOnReceiver(def, r.Receiver(), r.PartialUpdateStructName(), CheckFields, RecordShouldUsePointer).
		Params(
			Add(fieldChecker).Add(Op("*").Add(fieldCheckerType)),
			Add(keyChecker).Qual(utils.RestLiCodecPackage, "KeyChecker"),
		).
		Params(Err().Error()).
		BlockFunc(func(def *Group) {
			for _, ir := range r.includedRecords() {
				def.Err().Op("=").Id(r.Receiver()).Dot(ir.PartialUpdateStructName()).Dot(CheckFields).Call(fieldChecker, keyChecker)
				def.Add(utils.IfErrReturn(Err()))
				def.Line()
			}
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

				def.If(Err().Op("=").Add(fieldChecker).Dot("CheckField").Call(
					keyChecker,
					Lit(f.Name),
					d,
					setAccessorF(f).Op("!=").Nil(),
					p,
				), Err().Op("!=").Nil()).Block(Return(Err()))
			}
			def.Return(Nil())
		}).Line().Line()

	checkAllFields := func(def *Group, keyChecker Code) {
		def.Add(fieldChecker).Op(":=&").Add(fieldCheckerType).
			Custom(utils.MultiLineValues, Id("RecordType").Op(":").Lit(r.Identifier.FullName()))
		def.Err().Op("=").Id(r.Receiver()).Dot(CheckFields).Call(fieldChecker, keyChecker)
		def.Add(utils.IfErrReturn(Err()))
	}

	patch := Qual(utils.RestLiPatchPackage, "PatchField")

	const marshalPatch = "MarshalRestLiPatch"
	utils.AddFuncOnReceiver(def, r.Receiver(), r.PartialUpdateStructName(), marshalPatch, RecordShouldUsePointer).
		Params(WriterParam).
		Params(Err().Error()).
		BlockFunc(func(def *Group) {
			def.Return(Writer.WriteMap(Writer, func(keyWriter Code, def *Group) {
				checkAllFields(def, Writer)
				def.If(Add(fieldChecker).Dot("HasDeletes")).BlockFunc(func(def *Group) {
					def.Err().Op("=").Add(Writer.WriteArray(Add(keyWriter).Call(Lit("$delete")), func(itemWriter Code, def *Group) {
						def.Id(r.Receiver()).Dot(MarshalDeleteFields).Call(itemWriter)
						def.Add(Return(Nil()))
					}))
					def.Add(utils.IfErrReturn(Err()))
				}).Line()

				def.If(Add(fieldChecker).Dot("HasSets")).BlockFunc(func(def *Group) {
					def.Err().Op("=").Add(Writer.WriteMap(Add(keyWriter).Call(Lit("$set")), func(keyWriter Code, def *Group) {
						def.Return(Id(r.Receiver()).Dot(MarshalSetFields).Call(keyWriter))
					}))
					def.Add(utils.IfErrReturn(Err()))
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
			def.Err().Op("=").Add(Reader.ReadMap(Reader, func(reader, key Code, def *Group) {
				def.Switch(key).BlockFunc(func(def *Group) {
					def.Case(Lit("$delete")).BlockFunc(func(def *Group) {
						def.Err().Op("=").Add(Reader.ReadArray(Reader, func(itemReader Code, def *Group) {
							def.Var().Add(FieldParamName).String()
							def.Add(Reader.Read(deleteType, itemReader, FieldParamName))
							def.Add(utils.IfErrReturn(Err())).Line()

							def.Err().Op("=").Id(r.Receiver()).Dot(UnmarshalDeleteField).Call(FieldParamName)
							def.If(Err().Op("==").Add(noSuchFieldError)).Block(Err().Op("=").Nil())
							def.Return(Err())
						}))
					})
					def.Case(Lit("$set")).BlockFunc(func(def *Group) {
						def.Err().Op("=").Add(Reader.ReadMap(Reader, func(reader, key Code, def *Group) {
							def.List(Found, Err()).Op(":=").Id(r.Receiver()).Dot(UnmarshalSetField).Call(reader, key)
							def.If(Op("!").Add(Found)).Block(Err().Op("=").Add(Reader.Skip(Reader)))
							def.Return(Err())
						}))
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

const MarshalSetFields = "MarshalSetFields"
const UnmarshalSetField = "UnmarshalSetField"
const MarshalDeleteFields = "MarshalDeleteFields"
const UnmarshalDeleteField = "UnmarshalDeleteField"

var noSuchFieldError = Code(Qual(utils.RestLiPatchPackage, "NoSuchFieldErr"))
var deleteType = RestliType{Primitive: &StringPrimitive}

func (r *Record) generatePartialUpdateDeleteFieldsStruct() *Statement {
	var deletableFields []Field
	for _, f := range r.Fields {
		if f.IsOptionalOrDefault() {
			deletableFields = append(deletableFields, f)
		}
	}

	fields := append([]Field(nil), deletableFields...)
	sortFields(fields)

	def := Empty()
	def.Type().Id(r.PartialUpdateDeleteFieldsStructName()).StructFunc(func(def *Group) {
		for _, ir := range r.includedRecords() {
			def.Qual(ir.PackagePath(), ir.PartialUpdateDeleteFieldsStructName())
		}

		for _, f := range deletableFields {
			def.Id(f.FieldName()).Bool()
		}
	}).Line().Line()

	write := Code(Id("write"))

	utils.AddFuncOnReceiver(def, r.Receiver(), r.PartialUpdateDeleteFieldsStructName(), MarshalDeleteFields, RecordShouldUsePointer).
		Params(Add(write).Func().Params(String())).
		BlockFunc(func(def *Group) {
			for _, f := range fields {
				def.If(Id(r.Receiver()).Dot(f.FieldName())).Block(Add(write).Call(Lit(f.Name)))
			}
		}).
		Line().Line()

	utils.AddFuncOnReceiver(def, r.Receiver(), r.PartialUpdateDeleteFieldsStructName(), UnmarshalDeleteField, RecordShouldUsePointer).
		Params(Add(FieldParamName).String()).
		Params(Err().Error()).
		BlockFunc(func(def *Group) {
			def.Switch(FieldParamName).BlockFunc(func(def *Group) {
				for _, f := range r.Fields {
					if f.IsOptionalOrDefault() {
						def.Case(Lit(f.Name)).Block(
							r.rawFieldAccessor(f).Op("=").True(),
							Return(Nil()),
						)
					} else {
						def.Case(Lit(f.Name)).Block(
							Return(Qual(utils.RestLiPatchPackage, "NewFieldCannotBeDeletedError").
								Call(Lit(f.Name), Lit(r.Identifier.FullName()))),
						)
					}
				}
				def.Default().Block(Return(noSuchFieldError))
			})
		}).
		Line().Line()

	utils.AddFuncOnReceiver(def, r.Receiver(), r.PartialUpdateStructName(), MarshalDeleteFields, RecordShouldUsePointer).
		Params(ItemWriterFunc).
		Params(Err().Error()).
		BlockFunc(func(def *Group) {
			write := Code(Id("write"))
			def.Add(write).Op(":=").Func().Params(Id("name").String()).BlockFunc(func(def *Group) {
				def.Add(Writer.Write(deleteType, Add(ItemWriter).Call(), Id("name")))
			})
			for _, ir := range r.includedRecords() {
				def.Id(r.Receiver()).Dot(ir.PartialUpdateStructName()).Dot(DeleteFields).Dot(MarshalDeleteFields).Call(write)
				def.Line()
			}
			def.Id(r.Receiver()).Dot(DeleteFields).Dot(MarshalDeleteFields).Call(write)
			def.Add(Return(Nil()))
		}).Line().Line()

	utils.AddFuncOnReceiver(def, r.Receiver(), r.PartialUpdateStructName(), UnmarshalDeleteField, RecordShouldUsePointer).
		Params(Add(FieldParamName).String()).
		Params(Err().Error()).
		BlockFunc(func(def *Group) {
			for _, ir := range r.includedRecords() {
				def.Err().Op("=").Id(r.Receiver()).Dot(ir.PartialUpdateStructName()).Dot(DeleteFields).Dot(UnmarshalDeleteField).Call(FieldParamName)
				// Either err is nil in which case the field was successfully set, or it wasn't NoSuchFieldErr in which
				// case bail
				def.If(Err().Op("!=").Add(noSuchFieldError)).Block(Return(Err()))
				def.Line()
			}
			def.Return(Id(r.Receiver()).Dot(DeleteFields).Dot(UnmarshalDeleteField).Call(FieldParamName))
		}).Line().Line()

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

	def := Add(setRecord.GenerateStruct()).Line().Line().
		Add(setRecord.GenerateMarshalFields()).Line().Line().
		Add(setRecord.GenerateUnmarshalField(nil)).Line().Line()

	utils.AddFuncOnReceiver(def, r.Receiver(), r.PartialUpdateStructName(), MarshalSetFields, RecordShouldUsePointer).
		Params(KeyWriterFunc).
		Params(Err().Error()).
		BlockFunc(func(def *Group) {
			for _, ir := range r.includedRecords() {
				def.Err().Op("=").Id(r.Receiver()).Dot(ir.PartialUpdateStructName()).Dot(MarshalSetFields).Call(KeyWriter)
				def.Add(utils.IfErrReturn(Err()))
				def.Line()
			}
			def.Err().Op("=").Id(r.Receiver()).Dot(SetFields).Dot(MarshalFields).Call(KeyWriter)
			def.Return(Err())
		}).Line().Line()

	utils.AddFuncOnReceiver(def, r.Receiver(), r.PartialUpdateStructName(), UnmarshalSetField, RecordShouldUsePointer).
		Params(ReaderParam, Add(FieldParamName).String()).
		Params(Add(Found).Bool(), Err().Error()).
		BlockFunc(func(def *Group) {
			for _, ir := range r.includedRecords() {
				def.List(Found, Err()).Op("=").Id(r.Receiver()).Dot(ir.PartialUpdateStructName()).Dot(UnmarshalSetField).Call(Reader, FieldParamName)
				def.If(Err().Op("!=").Nil().Op("||").Add(Found)).Block(Return(Found, Err()))
				def.Line()
			}
			def.Return(Id(r.Receiver()).Dot(SetFields).Dot(UnmarshalField).Call(Reader, FieldParamName))
		}).Line().Line()

	return def
}

func (r *Record) includedRecords() []*Record {
	included := make([]*Record, len(r.Includes))
	for i, id := range r.Includes {
		included[i] = id.Resolve().(*Record)
	}
	return included
}
