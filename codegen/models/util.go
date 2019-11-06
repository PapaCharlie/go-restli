package models

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var namespaceEscape = regexp.MustCompile("([/.])_?internal([/.]?)")
var (
	modelRegistry = make(map[Identifier]ComplexType)
	CyclicModels  = make(map[Identifier]bool)
)

type SnapshotModels []*Model

func (models *SnapshotModels) UnmarshalJSON(data []byte) error {
	var modelList []*Model
	err := json.Unmarshal(data, &modelList)
	if err != nil {
		return errors.WithStack(err)
	}
	*models = modelList

	for _, m := range *models {
		m.sanityCheck(nil)
		m.resolveCyclicReferences()
	}

	return nil
}

func GetRegisteredModels() (models []ComplexType) {
	for _, m := range modelRegistry {
		models = append(models, m)
	}
	return models
}

func (m *Model) sanityCheck(parentModels []*Model) {
	if m.ComplexType == nil && m.BuiltinType == nil {
		if m.ref != nil && m.ref.Resolve() != nil {
			log.Panicln("Model has unresolved known ref:", m)
		} else {
			log.Panicln("Model defines neither ComplexType nor BuiltinType:", m)
		}
	}

	if m.ref != nil {
		var identifiers []string
		for _, im := range parentModels {
			if im.ComplexType != nil {
				identifiers = append(identifiers, im.ComplexType.GetIdentifier().GetQualifiedClasspath())
			}
		}
		log.Panicf("Model has a ref! (%s) %s", strings.Join(identifiers, " -> "), m)
	}

	if primitive, ok := m.BuiltinType.(*PrimitiveModel); ok && *primitive == NullPrimitive {
		log.Panicln("Model defines NullPrimitive type", m)
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
