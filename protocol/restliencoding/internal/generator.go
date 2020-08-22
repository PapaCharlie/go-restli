package main

import (
	"github.com/PapaCharlie/go-restli/internal/codegen"
	. "github.com/dave/jennifer/jen"
)

var types = []struct {
	Name string
	Type *Statement
}{
	{Name: "Int32", Type: Int32()},
	{Name: "Int64", Type: Int64()},
	{Name: "Float32", Type: Float32()},
	{Name: "Float64", Type: Float64()},
	{Name: "Bool", Type: Bool()},
	{Name: "String", Type: String()},
	{Name: "Bytes", Type: Index().Byte()},
}

func main() {
	def := NewFile("restliencoding")

	e := Id("e")
	encoder := Add(e).Dot("encoder")
	fieldName := Id("fieldName")
	fieldValue := Id("fieldValue")

	for _, t := range types {
		funcDef := func(suffix string, valueType *Statement, block func(def *Group)) {
			def.Func().
				Params(Add(e).Op("*").Id("Encoder")).
				Id(t.Name+suffix).
				Params(Add(fieldName).String(), Add(fieldValue).Add(valueType)).
				BlockFunc(func(def *Group) {
					def.Add(e).Dot("writeField").Call(fieldName)
					block(def)
				}).
				Line().Line()
		}

		// func (e *Encoder) BoolField(fieldName string, fieldValue bool) {
		// 	e.writeField(fieldName)
		// 	e.Bool(fieldValue)
		// }
		funcDef("Field", t.Type, func(def *Group) {
			def.Add(e).Dot(t.Name).Call(fieldValue)
		})

		// func (e *Encoder) BoolMapField(fieldName string, fieldValue map[string]bool) {
		// 	e.writeField(fieldName)
		// 	e.encoder.WriteMapStart()
		// 	first := true
		// 	for k, v := range fieldValue {
		// 		if first {
		// 			first = false
		// 		} else {
		// 			e.encoder.WriteMapEntryDelimiter()
		// 		}
		// 		e.encoder.WriteMapKey(k)
		// 		e.encoder.WriteMapKeyDelimiter()
		// 		e.encoder.Bool(v)
		// 	}
		// 	e.encoder.WriteMapEnd()
		// }
		funcDef("MapField", Map(String()).Add(t.Type), func(def *Group) {
			first := Id("first")
			k, v := Id("k"), Id("v")

			def.Add(encoder).Dot("WriteMapStart").Call()
			def.Add(first).Op(":=").True()
			def.For(List(k, v).Op(":=").Range().Add(fieldValue)).BlockFunc(func(def *Group) {
				def.If(first).
					Block(Add(first).Op("=").False()).
					Else().
					Block(Add(encoder).Dot("WriteMapEntryDelimiter").Call())
				def.Add(encoder).Dot("WriteMapKey").Call(k)
				def.Add(encoder).Dot("WriteMapKeyDelimiter").Call()
				def.Add(e).Dot(t.Name).Call(v)
			})
			def.Add(encoder).Dot("WriteMapEnd").Call()
		})

		// func (e *Encoder) BoolArrayField(fieldName string, fieldValue []bool) {
		// 	e.writeField(fieldName)
		// 	e.encoder.WriteArrayStart()
		// 	for i, v := range fieldValue {
		// 		if i > 0 {
		// 			e.encoder.WriteArrayItemDelimiter()
		// 		}
		// 		e.Bool(v)
		// 	}
		// 	e.encoder.WriteArrayEnd()
		// }
		funcDef("ArrayField", Index().Add(t.Type), func(def *Group) {

			def.Add(encoder).Dot("WriteArrayStart").Call()
			index, item := Id("index"), Id("item")
			def.For(List(index, item).Op(":=").Range().Add(fieldValue)).BlockFunc(func(def *Group) {
				def.If(Add(index).Op(">").Lit(0)).Block(
					Add(encoder).Dot("WriteArrayItemDelimiter").Call(),
				)
				def.Add(e).Dot(t.Name).Call(item)
			})
			def.Add(encoder).Dot("WriteArrayEnd").Call()
		})
	}

	err := codegen.WriteJenFile("interface_primitives.go", def)
	if err != nil {
		panic(err)
	}
}
