package models

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"regexp"
)

var namespaceEscape = regexp.MustCompile("([/.])_?internal([/.]?)")

func LoadModels(reader io.Reader) ([]*Model, error) {
	snapshot := &struct {
		Models map[string]*Model `json:"models"`
	}{}
	err := json.NewDecoder(reader).Decode(snapshot)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var models []*Model
	for _, m := range snapshot.Models {
		models = append(models, m)
	}

	models = append(models, flattenModels(models)...)
	replaceReferences(models)
	return models, nil
}

func flattenModels(models []*Model) (innerModels []*Model) {
	for _, m := range models {
		m.register()
		for _, im := range m.InnerModels() {
			im.register()
			innerModels = append(innerModels, im)
		}
	}
	if len(innerModels) > 0 {
		innerModels = append(innerModels, flattenModels(innerModels)...)
	}
	return innerModels
}

func replaceReferences(models []*Model) {
	for _, m := range models {
		if m.Reference != nil {
			*m = *GetRegisteredModel(m.Ns, m.Name)
		}
	}
}

func escapeNamespace(namespace string) string {
	return namespaceEscape.ReplaceAllString(namespace, "${1}_internal${2}")
}

var loadedModels = make(map[string]*Model)

func (m *Model) register() {
	if m.Primitive != nil || m.Reference != nil {
		return
	}

	loadedModels[m.PackagePath()+"."+m.Name] = m
}

func GetRegisteredModel(ns Ns, name string) *Model {
	return loadedModels[ns.PackagePath()+"."+name]
}
