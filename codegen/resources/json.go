package resources

import (
	"encoding/json"

	"github.com/PapaCharlie/go-restli/codegen/types"
	"github.com/pkg/errors"
)

type ResourcePathSegment struct {
	ResourceName string   `json:"resourceName"`
	PathKey      *PathKey `json:"pathKey"`
}

type PathKey struct {
	Name string           `json:"name"`
	Type types.RestliType `json:"type"`
}

func (pk *PathKey) UnmarshalJSON(data []byte) error {
	type _t PathKey
	err := json.Unmarshal(data, (*_t)(pk))
	if err != nil {
		return err
	}

	if !(pk.Type.Primitive != nil || (pk.Type.Reference != nil && !pk.Type.RawRecord)) {
		return errors.Errorf("go-restli: Invalid PathKey, cannot be array, map or raw record.")
	}
	return nil
}

type Resource struct {
	PackageRoot string `json:"-"`

	Namespace            string                 `json:"namespace"`
	Doc                  string                 `json:"doc"`
	SourceFile           string                 `json:"sourceFile"`
	ResourcePathSegments []ResourcePathSegment  `json:"resourcePathSegments"`
	ResourceSchema       *types.RestliType      `json:"resourceSchema"`
	Methods              []MethodImplementation `json:"-"`
	ReadOnlyFields       []string               `json:"readOnlyFields"`
	CreateOnlyFields     []string               `json:"createOnlyFields"`
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
	MethodType   MethodType        `json:"methodType"`
	Name         string            `json:"name"`
	Doc          string            `json:"doc"`
	OnEntity     bool              `json:"onEntity"`
	Params       []types.Field     `json:"params"`
	Return       *types.RestliType `json:"return"`
	Metadata     *types.RestliType `json:"metadata"`
	ReturnEntity bool              `json:"returnEntity"`
}
