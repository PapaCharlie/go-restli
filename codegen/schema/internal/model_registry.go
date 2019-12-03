package internal

import (
	"log"

	"github.com/pkg/errors"
)

var (
	currentFile      string
	currentNamespace string
)

type PdscModel struct {
	Type ComplexType
	File string
}

func (p *PdscModel) toModel() *Model {
	return &Model{ComplexType: p.Type}
}

type unresolvedModel struct {
	model *Model
	file  string
}

var ModelRegistry = &modelRegistry{
	resolvedTypes:    make(map[Identifier]*PdscModel),
	unresolvedModels: make(map[Identifier][]unresolvedModel),
}

type modelRegistry struct {
	resolvedTypes    map[Identifier]*PdscModel
	unresolvedModels map[Identifier][]unresolvedModel
}

func (reg *modelRegistry) registerComplexType(t ComplexType) {
	id := t.GetIdentifier()
	if id.Namespace == "" {
		log.Panicf("Cannot register %s without a namespace", id.Name)
	}
	if reg.resolvedTypes[id] == nil {
		reg.resolvedTypes[id] = &PdscModel{
			Type: t,
			File: currentFile,
		}
	}
}

func (reg *modelRegistry) addUnresolvedModel(id Identifier, model *Model) {
	reg.unresolvedModels[id] = append(reg.unresolvedModels[id], unresolvedModel{
		model: model,
		file:  currentFile,
	})
}

func (reg *modelRegistry) trimUnneededModels(loadedModels []*Model) {
	loadedModelIdentifiers := make(IdentifierSet)
	for _, m := range loadedModels {
		if m.ComplexType != nil {
			loadedModelIdentifiers.AddAll(DependencyGraph.AllDependencies(m.ComplexType.GetIdentifier(), nil))
		}
	}
	for id := range reg.unresolvedModels {
		if !loadedModelIdentifiers[id] {
			delete(reg.unresolvedModels, id)
		}
	}
}

func (reg *modelRegistry) Resolve(id Identifier) (ComplexType, error) {
	if t, ok := reg.resolvedTypes[id]; ok {
		return t.Type, nil
	} else {
		return nil, errors.Errorf("Unknown model: %s", id)
	}
}

func (reg *modelRegistry) GetModels() (types []*PdscModel) {
	for _, t := range reg.resolvedTypes {
		types = append(types, t)
	}
	return types
}

func (reg *modelRegistry) resolveModels() error {
	for id, models := range reg.unresolvedModels {
		t, err := reg.Resolve(id)
		if err != nil {
			return errors.Wrapf(err, "Could not resolve %s in %s", id, models[0].file)
		}
		for _, um := range models {
			um.model.ComplexType = t
		}
		delete(reg.unresolvedModels, id)
	}
	return nil
}
