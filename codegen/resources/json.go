package resources

import (
	"encoding/json"

	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/pkg/errors"
)

type Resource struct {
	Namespace        string                 `json:"namespace"`
	Doc              string                 `json:"doc"`
	SourceFile       string                 `json:"sourceFile"`
	RootResourceName string                 `json:"rootResourceName"`
	ResourceSchema   *types.RestliType      `json:"resourceSchema"`
	Methods          []MethodImplementation `json:"-"`
	ReturnEntity     bool                   `json:"returnEntity"`
	ReadOnlyFields   []string               `json:"readOnlyFields"`
	CreateOnlyFields []string               `json:"createOnlyFields"`
	IsCollection     bool                   `json:"isCollection"`
}

func (r *Resource) UnmarshalJSON(data []byte) (err error) {
	type t Resource
	err = json.Unmarshal(data, (*t)(r))
	if err != nil {
		return err
	}

	var methods struct {
		Methods []*Method `json:"methods"`
	}

	err = json.Unmarshal(data, &methods)
	if err != nil {
		return err
	}

	for _, m := range methods.Methods {
		mI := methodImplementation{
			Resource: r,
			Method:   m,
		}
		var impl MethodImplementation
		switch m.MethodType {
		case REST_METHOD:
			impl = &RestMethod{mI}
		case ACTION:
			impl = &Action{mI}
		case FINDER:
			impl = &Finder{mI}
		default:
			return errors.Errorf("unknown method type: %s", m.MethodType)
		}
		r.Methods = append(r.Methods, impl)
	}

	return nil
}

type MethodType string

const (
	REST_METHOD MethodType = "REST_METHOD"
	ACTION      MethodType = "ACTION"
	FINDER      MethodType = "FINDER"
)

type Method struct {
	MethodType      MethodType        `json:"methodType"`
	Name            string            `json:"name"`
	Doc             string            `json:"doc"`
	Path            string            `json:"path"`
	OnEntity        bool              `json:"onEntity"`
	EntityPathKey   *PathKey          `json:"entityPathKey"`
	PathKeys        []PathKey         `json:"pathKeys"`
	Params          []types.Field     `json:"params"`
	Return          *types.RestliType `json:"return"`
	Metadata        *types.RestliType `json:"metadata"`
	ReturnEntity    bool              `json:"returnEntity"`
	PagingSupported bool              `json:"pagingSupported"`
}

type PathKey struct {
	Name string           `json:"name"`
	Type types.RestliType `json:"type"`
}
