package types

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func AddUnmarshalRestli(def *Statement, receiver string, id utils.Identifier, pointer utils.ShouldUsePointer, f func(def *Group)) *Statement {
	target := Id(receiver)
	if pointer.ShouldUsePointer() {
		target = target.Op("*")
	}

	def.Func().Id(id.UnmarshalerFuncName()).
		Params(Add(Reader).Add(ReaderQual)).
		Params(target.Add(id.Qual()), Err().Error()).
		BlockFunc(func(def *Group) {
			if pointer.ShouldUsePointer() {
				def.Add(Id(receiver).Op("=").New(id.Qual()))
			}
			def.Add(Reader.Read(RestliType{Reference: &id}, Reader, Id(receiver)))
			def.Return(Id(receiver), Err())
		}).Line().Line()

	utils.AddFuncOnReceiver(def, receiver, id.Name, utils.UnmarshalRestLi, utils.Yes).
		Params(Add(Reader).Add(ReaderQual)).
		Params(Err().Error()).
		BlockFunc(f).
		Line().Line()

	data := Id("data")
	utils.AddFuncOnReceiver(def, receiver, id.Name, "UnmarshalJSON", utils.Yes).
		Params(Add(data).Index().Byte()).
		Params(Error()).
		BlockFunc(func(def *Group) {
			def.Add(Reader).Op(":=").Add(utils.NewJsonReader).Call(data)
			def.Return(Id(receiver).Dot(utils.UnmarshalRestLi).Call(Reader))
		}).Line().Line()

	return def
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

	return AddUnmarshalRestli(def, r.Receiver(), r.Identifier, RecordShouldUsePointer, func(def *Group) {
		r.generateUnmarshaler(def, requiredFields)
	})
}

func (r *Record) generateUnmarshaler(def *Group, requiredFields Code) {
	if len(r.Fields) == 0 {
		def.Return(Reader.ReadMap(Reader, func(reader Code, field Code, def *Group) {
			def.Return(Reader.Skip(reader))
		}))
		return
	}

	def.Err().Op("=").Add(Reader.ReadRecord(Reader, requiredFields, func(reader, field Code, def *Group) {
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

					def.Add(Reader.Read(f.Type, reader, accessor))
				})
			}
			def.Default().BlockFunc(func(def *Group) {
				def.Err().Op("=").Add(Reader.Skip(reader))
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
