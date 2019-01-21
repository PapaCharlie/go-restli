package schema

import (
	"encoding/json"
	"io"
)

func LoadSchema(reader io.Reader) (*Resource, error) {
	schema := &struct {
		Schema *Resource `json:"schema"`
	}{}

	err := json.NewDecoder(reader).Decode(schema)
	if err != nil {
		return nil, err
	} else {
		return schema.Schema, nil
	}
}
