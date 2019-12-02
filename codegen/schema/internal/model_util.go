package internal

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/PapaCharlie/go-restli/codegen"
	"github.com/pkg/errors"
)

var (
	namespaceEscape  = regexp.MustCompile("([/.])_?internal([/.]?)")
	currentFile      string
	currentNamespace string
)

var (
	ModelRegistry = make(map[Identifier]*PdscModel)

	DependencyGraph Graph
)

func registerComplexType(t ComplexType) {
	id := t.GetIdentifier()
	if id.Namespace != "" && ModelRegistry[id] == nil {
		ModelRegistry[id] = &PdscModel{
			Type: t,
			File: currentFile,
		}
	}
}

func LoadModels() error {
	failedFiles := make(map[string]error)
	err := filepath.Walk(codegen.PdscDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		m := new(Model)
		err = codegen.ReadJSONFromFile(path, m)
		if err != nil {
			failedFiles[path] = err
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}

	failedFilesLength := len(failedFiles)
	for len(failedFiles) != 0 {
		for f := range failedFiles {
			m := new(Model)
			err = codegen.ReadJSONFromFile(f, m)
			if err != nil {
				failedFiles[f] = err
			} else {
				delete(failedFiles, f)
			}
		}

		if len(failedFiles) == failedFilesLength {
			return errors.Errorf("Failed to deserialize the following files: %+v", failedFiles)
		}
	}
	return nil
}

func trimUnneededModels(models []*Model) {
	loadedModels := make(IdentifierSet)
	for _, m := range models {
		if m.ComplexType != nil {
			loadedModels.AddAll(DependencyGraph.AllDependencies(m.ComplexType.GetIdentifier(), nil))
		}
	}
	for id := range ModelRegistry {
		if !loadedModels[id] {
			delete(ModelRegistry, id)
		}
	}
}

func ResolveCyclicDependencies(loadedModels []*Model) {
	buildDependencyGraph()
	trimUnneededModels(loadedModels)
	flagCyclicDependencies()
}

func (m *Model) flattenInnerModels() (children IdentifierSet) {
	children = make(IdentifierSet)
	for _, im := range m.innerModels() {
		if im.ComplexType != nil {
			children.Add(im.ComplexType.GetIdentifier())
		} else {
			children.AddAll(im.flattenInnerModels())
		}
	}
	return children
}
