package models

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
)

func LoadModels(reader io.Reader) ([]*Model, error) {
	snapshot := &struct {
		Models []*Model `json:"models"`
	}{}
	err := json.NewDecoder(reader).Decode(snapshot)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	snapshot.Models = append(snapshot.Models, flattenModels(snapshot.Models)...)
	return snapshot.Models, nil
}

func flattenModels(models []*Model) (innerModels []*Model) {
	for _, m := range models {
		innerModels = append(innerModels, m.InnerModels()...)
	}
	if len(innerModels) > 0 {
		innerModels = append(innerModels, flattenModels(innerModels)...)
	}
	return innerModels
}
