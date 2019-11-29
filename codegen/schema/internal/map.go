package internal

import (
	"encoding/json"

	. "github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
)

const MapModelTypeName = "map"

type MapModel struct {
	Values *Model
}

func (m *MapModel) UnmarshalJSON(data []byte) error {
	t := &struct {
		typeField
		Values *Model `json:"values"`
	}{}
	if err := json.Unmarshal(data, t); err != nil {
		return err
	}
	if t.Type != MapModelTypeName {
		return &WrongTypeError{Expected: MapModelTypeName, Actual: t.Type}
	}
	m.Values = t.Values
	return nil
}

func (m *MapModel) GoType() *Statement {
	return Map(String()).Add(m.Values.GoType())
}

func (m *MapModel) innerModels() []*Model {
	return []*Model{m.Values}
}

func (m *MapModel) restLiWriteToBuf(def *Group, accessor *Statement) {
	def.Id("buf").Dot("WriteByte").Call(LitRune('('))

	def.Id("idx").Op(":=").Lit(0)
	def.For(List(Id("key"), Id("val")).Op(":=").Range().Add(accessor)).BlockFunc(func(def *Group) {
		def.If(Id("idx").Op("!=").Lit(0)).Block(Id("buf").Dot("WriteByte").Call(LitRune(','))).Line()
		def.Id("idx").Op("++")
		writeStringToBuf(def, Id(Codec).Dot("EncodeString").Call(Id("key")))
		def.Id("buf").Dot("WriteByte").Call(LitRune(':'))
		m.Values.restLiWriteToBuf(def, Id("val"))
	})

	def.Id("buf").Dot("WriteByte").Call(LitRune(')'))
	return
}
