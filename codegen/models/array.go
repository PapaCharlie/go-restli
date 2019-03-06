package models

import (
	"encoding/json"
	. "github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

const ArrayModelTypeName = "array"

type ArrayModel struct {
	Items *Model
}

func (a *ArrayModel) UnmarshalJSON(data []byte) error {
	t := &struct {
		Type  string
		Items *Model
	}{}
	if err := json.Unmarshal(data, t); err != nil {
		return err
	}
	if t.Type != ArrayModelTypeName {
		return errors.Errorf("Not an array type: %s", string(data))
	}
	a.Items = t.Items
	return nil
}

func (a *ArrayModel) GoType() *Statement {
	return Index().Add(a.Items.GoType())
}

func (a *ArrayModel) InnerModels() []*Model {
	return []*Model{a.Items}
}
