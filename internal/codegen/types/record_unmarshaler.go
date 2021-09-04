package types

import (
	"github.com/PapaCharlie/go-restli/internal/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func AddUnmarshalRestli(def *Statement, receiver, typeName string, f func(def *Group)) *Statement {
	utils.AddFuncOnReceiver(def, receiver, typeName, utils.UnmarshalRestLi).
		Params(Add(Reader).Add(ReaderQual)).
		Params(Err().Error()).
		BlockFunc(f).
		Line().Line()

	data := Id("data")
	utils.AddFuncOnReceiver(def, receiver, typeName, "UnmarshalJSON").
		Params(Add(data).Index().Byte()).
		Params(Error()).
		BlockFunc(func(def *Group) {
			def.Add(Reader).Op(":=").Add(utils.NewJsonReader).Call(data)
			def.Return(Id(receiver).Dot(utils.UnmarshalRestLi).Call(Reader))
		}).Line().Line()

	utils.AddFuncOnReceiver(def, receiver, typeName, "UnmarshalProtobuf").
		Params(Id("data").Index().Byte()).
		Params(Error()).
		BlockFunc(func(def *Group) {
			def.Add(Reader).Op(":=").Add(utils.NewProtobufReader).Call(data)
			def.Return(Id(receiver).Dot(utils.UnmarshalRestLi).Call(Reader))
		}).Line().Line()

	return def
}

func (r *Record) GenerateUnmarshalRestLi() *Statement {
	return AddUnmarshalRestli(Empty(), r.Receiver(), r.Name, func(def *Group) {
		r.generateUnmarshaler(def)
	})
}

func (r *Record) generateUnmarshaler(def *Group) {
	fields := r.SortedFields()

	if len(fields) == 0 {
		def.Return(Reader.ReadMap(Reader, func(reader Code, field Code, def *Group) {
			def.Return(Reader.Skip(reader))
		}))
		return
	}

	atInputStart := Id("atInputStart")
	def.Add(atInputStart).Op(":=").Add(Reader).Dot("AtInputStart").Call()

	requiredFieldsRemaining := Id("requiredFieldsRemaining")
	def.Add(requiredFieldsRemaining).Op(":=").Map(String()).Struct().Values(DictFunc(func(dict Dict) {
		for _, f := range fields {
			if !f.IsOptionalOrDefault() {
				dict[Lit(f.Name)] = Values()
			}
		}
	})).Line()

	def.Err().Op("=").Add(Reader.ReadMap(Reader, func(reader, field Code, def *Group) {
		def.Switch(field).BlockFunc(func(def *Group) {
			for _, f := range fields {
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
		def.Add(utils.IfErrReturn(Err()))
		def.Delete(requiredFieldsRemaining, field)
		def.Return(Nil())
	})).Line()

	def.Add(utils.IfErrReturn(Err())).Line()
	def.Add(Reader).Dot("RecordMissingRequiredFields").Call(requiredFieldsRemaining).Line()

	if r.hasDefaultValue() {
		def.Id(r.Receiver()).Dot(utils.PopulateLocalDefaultValues).Call()
	}

	def.If(atInputStart).Block(
		Return(Add(Reader).Dot("CheckMissingFields").Call()),
	).Else().Block(
		Return(Nil()),
	)
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
