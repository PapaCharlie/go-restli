package schema

import (
	"log"

	"github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/codegen/schema/internal"
	"github.com/pkg/errors"
)

var loadedModels []*internal.Model

func LoadRestSpecs(restSpecs []string) (resources []*Resource, types []*internal.PdscModel, err error) {
	err = internal.LoadModels()
	if err != nil {
		return nil, nil, err
	}

	for _, f := range restSpecs {
		log.Println(f)
		r := &Resource{file: f}
		err = errors.Wrapf(codegen.ReadJSONFromFile(f, r), "Failed to read restspec from %s", f)
		if err != nil {
			return nil, nil, err
		}
		resources = append(resources, r)
	}

	internal.ResolveCyclicDependencies(loadedModels)

	types = internal.ModelRegistry.GetModels()

	return resources, types, nil
}
