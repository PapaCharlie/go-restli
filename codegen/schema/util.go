package schema

import (
	"fmt"
	"github.com/PapaCharlie/go-restli/codegen/models"
	"io"
	"strings"

	"github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
)

func LoadResources(reader io.Reader) ([]*Resource, error) {
	resources := &struct {
		Resources map[string]*Resource `json:"resources"`
	}{}

	err := models.ReadJSON(reader, resources)
	if err != nil {
		return nil, err
	} else {
		removeSubResourcesFromTopLevel(resources.Resources, nil)
		var r []*Resource
		for _, v := range resources.Resources {
			r = append(r, v)
		}
		return r, nil
	}
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

func AddClientFunc(def *Statement, funcName string) *Statement {
	return codegen.AddFuncOnReceiver(def, ClientReceiver, Client, funcName)
}

func addEntityParams(def *Group, resources []*Resource) {
	for _, r := range resources {
		if id := r.getIdentifier(); id != nil {
			def.Id(id.Name).Add(id.Type.GoType())
		}
	}
}

func buildQueryPath(resources []*Resource, rawQueryPath string) string {
	for _, r := range resources {
		if id := r.getIdentifier(); id != nil {
			rawQueryPath = strings.Replace(rawQueryPath, fmt.Sprintf("{%s}", id.Name), "%s", 1)
		}
	}
	return rawQueryPath
}

func encodeEntitySegments(def *Group, resources []*Resource) {
	for _, r := range resources {
		if id := r.getIdentifier(); id != nil {
			hasError, assignment := id.Type.restLiURLEncode(Id(id.Name))
			if hasError {
				def.List(Id(id.Name+"Str"), Err()).Op(":=").Add(assignment)
				codegen.IfErrReturn(def)
			} else {
				def.Id(id.Name + "Str").Op(":=").Add(assignment)
			}
		}
	}
}

//
//func k(def *Statement, funcName string, queryPath string, resources []*Resource, extraParams func(def *Group), extraReturnParams func(def *Group), urlSuffix string, body func(def *Group)) *Statement {
//	return codegen.AddFuncOnReceiver(def, ClientReceiver, Client, funcName).
//		ParamsFunc(func(def *Group) {
//			for _, r := range resources {
//				if id := r.getIdentifier(); id != nil {
//					def.Id(id.Name).Add(id.Type.GoType())
//					queryPath = strings.Replace(queryPath, fmt.Sprintf("{%s}", id.Name), "%s", 1)
//				}
//			}
//			extraParams(def)
//		}).
//		ParamsFunc(func(def *Group) {
//			extraReturnParams(def)
//			def.Err().Error()
//		}).
//		BlockFunc(func(def *Group) {
//			def.Id(Url).Op(":=").Id(ClientReceiver).Dot(HostnameClientField).Op("+").Qual("fmt", "Sprintf").
//				CallFunc(func(def *Group) {
//					def.Lit(queryPath + urlSuffix)
//					for _, r := range resources {
//						if id := r.getIdentifier(); id != nil {
//							def.Id(id.Name + "Str")
//						}
//					}
//				})
//		})
//}
