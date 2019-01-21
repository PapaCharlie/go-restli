package models

import (
	"encoding/json"
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

const ArrayType = "array"

type Array struct {
	Type  string
	Items *Model
}

func (a *Array) UnmarshalJSON(data []byte) error {
	type t Array
	if err := json.Unmarshal(data, (*t)(a)); err != nil {
		return err
	}
	if a.Type != ArrayType {
		return errors.Errorf("Not an array type: %s", string(data))
	}
	return nil
}

func (a *Array) GoType(packagePrefix string) *jen.Statement {
	return jen.Index().Add(a.Items.GoType(packagePrefix))
}

func (a *Array) InnerModels() []*Model {
	return []*Model{a.Items}
}
