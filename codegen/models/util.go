package models

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/pkg/errors"
)

var namespaceEscape = regexp.MustCompile("([/.])_?internal([/.]?)")
var ModelCache = make(map[Identifier]ComplexType)

func LoadModels(reader io.Reader) ([]ComplexType, error) {
	spec := &struct {
		Models map[string]*Model `json:"models"`
	}{}

	err := ReadJSON(reader, spec)
	if err != nil {
		return nil, err
	}

	for _, m := range spec.Models {
		m.sanityCheck(nil)
	}

	return getRegisteredModels(), err
}

func LoadSnapshotModels(reader io.Reader) ([]ComplexType, error) {
	snapshot := &struct {
		Models []*Model `json:"models"`
	}{}

	err := ReadJSON(reader, snapshot)
	if err != nil {
		return nil, err
	}

	for _, m := range snapshot.Models {
		m.sanityCheck(nil)
	}

	return getRegisteredModels(), nil
}

func ReadJSON(reader io.Reader, s interface{}) error {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.WithStack(err)
	}

	err = json.Unmarshal(bytes, s)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func getRegisteredModels() (models []ComplexType) {
	for _, m := range ModelCache {
		models = append(models, m)
	}
	return models
}

func (m *Model) sanityCheck(parentModels []*Model) {
	if m.ref != nil {
		log.Panicln(parentModels, m)
	}

	if m.ComplexType != nil && m.BuiltinType != nil {
		log.Panicln(m)
	}

	if primitive, ok := m.BuiltinType.(*PrimitiveModel); ok && *primitive == NullPrimitive {
		log.Panicln(m)
	}

	parentModels = append(append([]*Model(nil), parentModels...), m)
	for _, im := range m.innerModels() {
		im.sanityCheck(parentModels)
	}
}
