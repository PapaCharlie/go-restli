package schema

import (
	"log"

	"github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/codegen/schema/internal"
	"github.com/pkg/errors"
)

func LoadRestSpecs(restSpecs []string) (resources []*Resource, types []*internal.PdscModel, err error) {
	var f string
	defer func() {
		if r := recover(); r != nil {
			if f != "" {
				log.Panicf("Failed to read %s: %+v", f, r)
			} else {
				log.Panicln(r)
			}
		}
	}()

	for _, f = range restSpecs {
		log.Println(f)
		r := &Resource{file: f}
		err = codegen.ReadJSONFromFile(f, r)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
		resources = append(resources, r)
	}
	f = ""

	internal.ResolveCyclicDependencies()

	for _, t := range internal.ModelRegistry {
		types = append(types, t)
	}

	return resources, types, nil
}
