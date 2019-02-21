package models

import (
	"encoding/json"
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

const ArrayType = "array"

type Array struct {
	//Type  string
	Items *Model
}

func (a *Array) UnmarshalJSON(data []byte) error {
	t := &struct {
		Type  string
		Items *Model
	}{}
	if err := json.Unmarshal(data, t); err != nil {
		return err
	}
	if t.Type != ArrayType {
		return errors.Errorf("Not an array type: %s", string(data))
	}
	a.Items = t.Items
	return nil
}

func (a *Array) GoType() *jen.Statement {
	return jen.Index().Add(a.Items.GoType())
}

func (a *Array) InnerModels() []*Model {
	return []*Model{a.Items}
}

//func (a *Array) GetLit(rawJson string) interface{} {
//
//}
