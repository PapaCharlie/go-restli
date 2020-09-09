package types

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func (r *Record) generatePartialUpdateStruct() *Statement {
	def := Empty()

	const (
		DeleteField = "Delete"
		UpdateField = "Update"
	)

	// Generate the struct
	utils.AddWordWrappedComment(def, fmt.Sprintf(
		"%s is used to represent a partial update on %s. Toggling the value of a field\n"+
			"in Delete represents selecting it for deletion in a partial update, while\n"+
			"setting the value of a field in Update represents setting that field in the\n"+
			"current struct. Other fields in this struct represent record fields that can\n"+
			"themselves be partially updated.",
		r.PartialUpdateStructName(), r.Name,
	)).Line()

	var deletableFields []Field
	for _, f := range r.Fields {
		if f.IsOptionalOrDefault() {
			deletableFields = append(deletableFields, f)
		}
	}

	def.Type().Id(r.PartialUpdateStructName()).StructFunc(func(def *Group) {
		if len(deletableFields) > 0 {
			def.Id(DeleteField).StructFunc(func(def *Group) {
				for _, f := range deletableFields {
					def.Id(f.FieldName()).Bool()
				}
			})
		}
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

	AddMarshalRestLi(def, r.Receiver(), r.PartialUpdateStructName(), func(def *Group) {
		fields := r.SortedFields()

		deleteAccessor := func(f Field) *Statement { return Id(r.Receiver()).Dot(DeleteField).Dot(f.FieldName()) }
		updateAccessor := func(f Field) *Statement { return Id(r.Receiver()).Dot(UpdateField).Dot(f.FieldName()) }
		noUpdateAndDelete := func(f Field) string {
			return fmt.Sprintf("cannot both delete and update %q of %q", f.Name, r.Name)
		}
		noUpdateAndSet := func(f Field) string {
			return fmt.Sprintf("cannot both set and update %q of %q", f.Name, r.Name)
		}
		createOrReadOnlyFieldError := func(op string, f Field) string {
			return fmt.Sprintf("cannot %s read-only or create-only field %q of %q", op, f.Name, r.Name)
		}
		newError := func(message string) Code {
			return Qual("errors", "New").Call(Lit(message))
		}

		if len(fields) == 0 {
			def.Return(Writer.WriteMap(Writer, func(keyWriter Code, def *Group) {
				def.Return(Nil())
			}))
			return
		}

		def.Return(Writer.WriteMap(Writer, func(keyWriter Code, def *Group) {
			if len(deletableFields) > 0 {
				var hasDeletes *Statement
				for _, f := range deletableFields {
					if hasDeletes == nil {
						hasDeletes = deleteAccessor(f)
					} else {
						hasDeletes.Op("||").Line().Add(deleteAccessor(f))
					}
				}

				def.If(hasDeletes).BlockFunc(func(def *Group) {
					deleteType := RestliType{Primitive: &StringPrimitive}
					def.Err().Op("=").Add(Writer.WriteArray(Add(keyWriter).Call(Lit("$delete")), func(itemWriter Code, def *Group) {
						for _, f := range deletableFields {
							accessor := deleteAccessor(f)
							def.If(accessor).BlockFunc(func(def *Group) {
								def.If(Writer.IsKeyExcluded(Writer, Lit(f.Name))).Block(
									Return(newError(createOrReadOnlyFieldError("delete", f))),
								)
								def.Add(Writer.Write(deleteType, Add(itemWriter).Call(), Lit(f.Name)))
							}).Line()
						}
						def.Add(Return(Nil()))
					}))
					def.Add(utils.IfErrReturn(Err()))
				})
				def.Line()
			}

			hasSets := updateAccessor(fields[0]).Op("!=").Nil()
			for i := 1; i < len(fields); i++ {
				hasSets.Op("||").Line().Add(updateAccessor(fields[i])).Op("!=").Nil()
			}

			def.If(hasSets).BlockFunc(func(def *Group) {
				def.Err().Op("=").Add(Writer.WriteMap(Add(keyWriter).Call(Lit("$set")), func(keyWriter Code, def *Group) {
					for _, f := range fields {
						accessor := updateAccessor(f)
						def.If(Add(accessor).Op("!=").Nil()).BlockFunc(func(def *Group) {
							if f.IsOptionalOrDefault() {
								def.If(deleteAccessor(f)).Block(
									Return(newError(noUpdateAndDelete(f))),
								)
							}
							def.If(Writer.IsKeyExcluded(Writer, Lit(f.Name))).Block(
								Return(newError(createOrReadOnlyFieldError("update", f))),
							)
							if f.Type.Reference == nil {
								accessor = Op("*").Add(accessor)
							}
							def.Add(Writer.Write(f.Type, Add(keyWriter).Call(Lit(f.Name)), accessor, Err()))
						})
						def.Line()
					}
					def.Return(Nil())
				}))
				def.Add(utils.IfErrReturn(Err()))
			})

			for _, f := range fields {
				if f.Type.Record() != nil {
					def.Line()
					accessor := Id(r.Receiver()).Dot(f.FieldName())
					def.If(Add(accessor).Op("!=").Nil()).BlockFunc(func(def *Group) {
						if f.IsOptionalOrDefault() {
							def.If(Add(deleteAccessor(f))).Block(
								Return(newError(noUpdateAndDelete(f))),
							)
						}
						def.If(Add(updateAccessor(f).Op("!=").Nil())).Block(
							Return(newError(noUpdateAndSet(f))),
						)
						def.If(Writer.IsKeyExcluded(Writer, Lit(f.Name))).Block(
							Return(newError(createOrReadOnlyFieldError("update", f))),
						)
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
