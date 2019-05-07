package schema

import (
	"github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/codegen/models"
	. "github.com/dave/jennifer/jen"
	"io"
)

func LoadResources(reader io.Reader) ([]*Resource, error) {
	resources := &struct {
		Resources map[string]*Resource `json:"resources"`
	}{}

	err := models.ReadJSON(reader, resources)
	if err != nil {
		return nil, err
	}

	removeSubResourcesFromTopLevel(resources.Resources, nil)
	var r []*Resource
	for _, v := range resources.Resources {
		r = append(r, v)
	}
	return r, nil
}

func removeSubResourcesFromTopLevel(resources map[string]*Resource, res *Resource) {
	if res == nil {
		for _, v := range resources {
			if e := v.getEntity(); e != nil {
				for _, sr := range e.Subresources {
					removeSubResourcesFromTopLevel(resources, sr)
				}
			}
		}
	} else {
		fullResourceName := res.Name
		if res.Namespace != "" {
			fullResourceName = res.Namespace + "." + fullResourceName
		}
		if r, ok := resources[fullResourceName]; ok && r != nil {
			delete(resources, fullResourceName)
		}
	}
}

func LoadSnapshotResource(reader io.Reader) ([]*Resource, error) {
	schema := &struct {
		Schema *Resource `json:"schema"`
	}{}

	err := models.ReadJSON(reader, schema)
	if err != nil {
		return nil, err
	}
	return []*Resource{schema.Schema}, nil
}

func addClientFunc(def *Statement, funcName string) *Statement {
	return codegen.AddFuncOnReceiver(def, ClientReceiver, Client, funcName)
}
