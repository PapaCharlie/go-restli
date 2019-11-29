package internal

import (
	"encoding/json"

	. "github.com/dave/jennifer/jen"
)

const ArrayModelTypeName = "array"

type ArrayModel struct {
	Items *Model
}

func (a *ArrayModel) UnmarshalJSON(data []byte) error {
	t := &struct {
		typeField
		Items *Model `json:"items"`
	}{}
	if err := json.Unmarshal(data, t); err != nil {
		return err
	}
	if t.Type != ArrayModelTypeName {
		return &WrongTypeError{Expected: ArrayModelTypeName, Actual: t.Type}
	}
	a.Items = t.Items
	return nil
}

func (a *ArrayModel) GoType() *Statement {
	return Index().Add(a.Items.GoType())
}

func (a *ArrayModel) innerModels() []*Model {
	return []*Model{a.Items}
}

func (a *ArrayModel) restLiWriteToBuf(def *Group, accessor *Statement) {
	writeStringToBuf(def, Lit("List("))

	def.For(List(Id("idx"), Id("val")).Op(":=").Range().Add(accessor)).BlockFunc(func(def *Group) {
		def.If(Id("idx").Op("!=").Lit(0)).Block(Id("buf").Dot("WriteByte").Call(LitRune(','))).Line()
		a.Items.restLiWriteToBuf(def, Id("val"))
	})

	def.Id("buf").Dot("WriteByte").Call(LitRune(')'))
	return
}
