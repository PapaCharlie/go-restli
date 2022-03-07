package types

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func (r *Record) GenerateUnmarshalerFunc() *Statement {
	return AddUnmarshalerFunc(Empty(), r.Receiver(), r.Identifier, RecordShouldUsePointer)
}

func (r *Record) GenerateUnmarshalRestLi() *Statement {
	def := Empty()

	hasRequiredFields := false
	for _, f := range r.Fields {
		if !f.IsOptionalOrDefault() {
			hasRequiredFields = true
			break
		}
	}

	var requiredFields Code
	if hasRequiredFields {
		requiredFields = Code(Id("_" + r.Name + "RequiredFields"))
		def.Var().Add(requiredFields).Op("=").Add(utils.RequiredFields).ValuesFunc(func(def *Group) {
			if len(r.Fields) > 0 {
				for _, f := range r.Fields {
					if !f.IsOptionalOrDefault() {
						def.Line().Lit(f.Name)
					}
				}
				def.Line()
			}
		}).Line().Line()
	} else {
		requiredFields = Nil()
	}

	return AddUnmarshalRestLi(def, r.Receiver(), r.Name, func(_, reader Code, def *Group) {
		r.generateUnmarshaler(def, reader, requiredFields)
	})
}

func (r *Record) generateUnmarshaler(def *Group, reader, requiredFields Code) {
	if len(r.Fields) == 0 {
		def.Return(ReaderUtils.ReadMap(reader, func(reader Code, field Code, def *Group) {
			def.Return(ReaderUtils.Skip(reader))
		}))
		return
	}

	def.Err().Op("=").Add(ReaderUtils.ReadRecord(reader, requiredFields, func(reader, field Code, def *Group) {
		def.Switch(field).BlockFunc(func(def *Group) {
			for _, f := range r.Fields {
				def.Case(Lit(f.Name)).BlockFunc(func(def *Group) {
					accessor := r.fieldAccessor(f)

					if f.IsOptionalOrDefault() {
						def.Add(accessor).Op("=").New(f.Type.GoType())
						if f.Type.Reference == nil {
							accessor = Op("*").Add(accessor)
						}
					}

					def.Add(ReaderUtils.Read(f.Type, reader, accessor))
				})
			}
			def.Default().BlockFunc(func(def *Group) {
				def.Err().Op("=").Add(ReaderUtils.Skip(reader))
			})
		})
		def.Return(Err())
	}))
	def.Add(utils.IfErrReturn(Err())).Line()

	if r.hasDefaultValue() {
		def.Id(r.Receiver()).Dot(utils.PopulateLocalDefaultValues).Call()
	}

	def.Return(Nil())
}

// TODO
// func (r *Record) generateQueryParamsUnmarhsaler(def *Group, finderName *string) {
// 	fields := r.SortedFields()
// 	qIndex := -1
// 	if finderName != nil {
// 		qIndex = sort.Search(len(fields), func(i int) bool { return fields[i].Name >= FinderNameParam })
// 		fields = append(fields[:qIndex], append([]Field{{
// 			Name:       FinderNameParam,
// 			IsOptional: false,
// 			Type:       RestliType{Primitive: &StringPrimitive},
// 		}}, fields[qIndex:]...)...)
// 	}
//
// 	finderNameVar := Id("finderName")
// 	if finderName != nil {
// 		def.Var().Add(finderNameVar).String()
// 		def.Line()
// 	}
//
// }
