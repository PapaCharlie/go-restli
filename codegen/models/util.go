package models

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var namespaceEscape = regexp.MustCompile("([/.])_?internal([/.]?)")
var (
	ModelCache   = make(map[Identifier]ComplexType)
	CyclicModels = make(map[Identifier]bool)
)

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
		m.resolveCyclicReferences()
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
		m.resolveCyclicReferences()
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

func (m *Model) resolveCyclicReferences() {
	for {
		if m.ComplexType == nil {
			continue
		}
		modelChain := m.traverseDependencyGraph(nil, nil)
		if len(modelChain) > 0 {
			if modelChain[0].Name == modelChain[len(modelChain)-1].Name {
				log.Fatalf("%s depends on itself!", modelChain[0])
			} else {
				var identifiers []string
				for _, id := range modelChain {
					identifiers = append(identifiers, id.GetQualifiedClasspath())
				}

				log.Println("Detected cyclic dependency:", strings.Join(identifiers, " -> "))
			}
		} else {
			break
		}
	}

	dependsOnCyclicModel := false
	allDependencies := m.allDependencies(nil)
	for dep := range allDependencies {
		if CyclicModels[dep.GetIdentifier()] {
			dependsOnCyclicModel = true
			break
		}
	}
	if dependsOnCyclicModel {
		for dep := range allDependencies {
			CyclicModels[dep.GetIdentifier()] = true
		}
	}
}

func (m *Model) traverseDependencyGraph(path []Identifier, visitedModels map[Identifier]bool) []Identifier {
	if path == nil && m.ComplexType != nil {
		path = []Identifier{m.ComplexType.GetIdentifier()}
	}
	if visitedModels == nil {
		visitedModels = map[Identifier]bool{}
	}

	for _, im := range m.innerModels() {
		innerPath := append([]Identifier(nil), path...)
		if im.ComplexType != nil && len(path) > 0 {
			startingModelId := path[0]
			previousModelId := path[len(path)-1]
			innerModelId := im.ComplexType.GetIdentifier()

			innerPath = append(innerPath, innerModelId)

			if visitedModels[innerModelId] || CyclicModels[innerModelId] {
				continue
			}

			if innerModelId.Namespace == startingModelId.Namespace && previousModelId.Namespace != innerModelId.Namespace {
				for _, id := range innerPath {
					CyclicModels[id] = true
				}
				return innerPath
			} else {
				visitedModels[innerModelId] = true
			}
		}

		if modelChain := im.traverseDependencyGraph(innerPath, visitedModels); len(modelChain) > 0 {
			return modelChain
		}
	}

	return nil
}

func (m *Model) allDependencies(types map[ComplexType]bool) map[ComplexType]bool {
	if types == nil {
		types = make(map[ComplexType]bool)
	}
	if m.ComplexType != nil {
		types[m.ComplexType] = true
	}
	for _, im := range m.innerModels() {
		if im.ComplexType != nil {
			if types[im.ComplexType] {
				break
			} else {
				types[im.ComplexType] = true
			}
		}
		for k, v := range im.allDependencies(types) {
			types[k] = v
		}
	}
	return types
}
