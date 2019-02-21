package models

import (
	"encoding/json"
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

const MapType = "map"

type Map struct {
	Type   string
	Values *Model
}

func (m *Map) UnmarshalJSON(data []byte) error {
	type t Map
	if err := json.Unmarshal(data, (*t)(m)); err != nil {
		return err
	}
	if m.Type != MapType {
		return errors.Errorf("Not a map type: %s", string(data))
	}
	return nil
}

func (m *Map) GoType() *jen.Statement {
	return jen.Map(jen.String()).Add(m.Values.GoType())
}

func (m *Map) InnerModels() []*Model {
	return []*Model{m.Values}
}
