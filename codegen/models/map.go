package models

import (
	"encoding/json"
	. "github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

const MapModelTypeName = "map"

type MapModel struct {
	Type   string
	Values *Model
}

func (m *MapModel) UnmarshalJSON(data []byte) error {
	type t MapModel
	if err := json.Unmarshal(data, (*t)(m)); err != nil {
		return err
	}
	if m.Type != MapModelTypeName {
		return errors.Errorf("Not a map type: %s", string(data))
	}
	return nil
}

func (m *MapModel) GoType() *Statement {
	return Map(String()).Add(m.Values.GoType())
}

func (m *MapModel) InnerModels() []*Model {
	return []*Model{m.Values}
}
