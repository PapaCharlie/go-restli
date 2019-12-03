package internal

import (
	"os"
	"path/filepath"

	"github.com/PapaCharlie/go-restli/codegen"
	"github.com/pkg/errors"
)

func LoadModels() error {
	failedFiles := make(map[string]error)
	err := filepath.Walk(codegen.PdscDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		currentFile = path
		m := new(Model)
		err = codegen.ReadJSONFromFile(path, m)
		if me, ok := errors.Cause(err).(*codegen.MalformedPdscFileError); ok {
			return me
		}
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

	return ModelRegistry.resolveModels()
}

func ResolveCyclicDependencies(loadedModels []*Model) {
	buildDependencyGraph()
	ModelRegistry.trimUnneededModels(loadedModels)
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
