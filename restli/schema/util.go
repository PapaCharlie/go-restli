package schema

import (
	"encoding/json"
	"io"
)

func LoadSchema(reader io.Reader) (*Schema, error) {
	schema := &struct {
		Schema *Schema `json:"schema"`
	}{}

	err := json.NewDecoder(reader).Decode(schema)
	if err != nil {
		return nil, err
	} else {
		return schema.Schema, nil
	}
}
