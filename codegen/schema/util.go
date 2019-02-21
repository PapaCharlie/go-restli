package schema

import (
	"encoding/json"
	"io"
)

func LoadResources(reader io.Reader) ([]*Resource, error) {
	resources := &struct {
		Resources map[string]*Resource `json:"resources"`
	}{}

	err := json.NewDecoder(reader).Decode(resources)
	if err != nil {
		return nil, err
	} else {
		var r []*Resource
		for _, v := range resources.Resources {
			r = append(r, v)
		}
		return r, nil
	}
}
