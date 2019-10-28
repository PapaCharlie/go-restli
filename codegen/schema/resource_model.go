package schema

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/codegen/models"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/pkg/errors"
)

func (r *Resource) UnmarshalJSON(data []byte) error {
	type t Resource
	if err := json.Unmarshal(data, (*t)(r)); err != nil {
		return err
	}

	if r.Schema != nil && r.Schema.Model.ComplexType == nil && !r.Schema.Model.IsBytesOrPrimitive() {
		return errors.Errorf("Failed to deserialize Resource model (can only be primitive or reference type), got: %+v",
			r.Schema.Model.BuiltinType)
	}

	if r.Association != nil {
		r.Association.Namespace = r.Namespace
	}

	return nil
}

type ResourceModel struct {
	*models.Model
}

func (r *ResourceModel) UnmarshalJSON(data []byte) error {
	r.Model = new(models.Model)

	var primitive models.PrimitiveModel
	if err := json.Unmarshal(data, &primitive); err == nil {
		r.Model.BuiltinType = &primitive
		return nil
	}

	var bytes models.BytesModel
	if err := json.Unmarshal(data, &bytes); err == nil {
		r.Model.BuiltinType = &bytes
		return nil
	}

	var ref models.ModelReference
	if err := json.Unmarshal(data, &ref); err == nil {
		if t := ref.Resolve(); t == nil {
			return errors.Errorf("Unresolved reference %+v", ref)
		} else {
			r.Model.ComplexType = t
			return nil
		}
	}

	var unescapedType string
	_ = json.Unmarshal(data, &unescapedType)

	if err := json.Unmarshal([]byte(unescapedType), r.Model); err != nil {
		return errors.Errorf("Failed to deserialize Resource model from %s: %+v", unescapedType, err)
	}

	return nil
}

type parameter struct {
	Name, Doc string
	Type      *ResourceModel
	Optional  bool
	Default   *string
}

func (p parameter) toField() (f models.Field) {
	f.Name = p.Name
	f.Doc = p.Doc
	f.Type = p.Type.Model
	f.Optional = p.Optional
	if p.Default != nil {
		// Special case for string default values: the @Optional annotation doesn't escape the string, it puts it as a
		// literal, therefore we need to escape it before passing it in. Maps and lists are represented as `{}` and
		// `[]` respectively, so no escaping there, and numeric values don't need to be escaped in JSON.
		if primitive, ok := p.Type.BuiltinType.(*models.PrimitiveModel); ok && *primitive == models.StringPrimitive {
			raw, _ := json.Marshal(*p.Default)
			f.Default = json.RawMessage(raw)
		} else {
			f.Default = json.RawMessage(*p.Default)
		}
	}
	return f
}

func (e *Endpoint) UnmarshalJSON(data []byte) error {
	t := &struct {
		Name, Doc  string
		Parameters []*parameter
		Returns    *ResourceModel
	}{}

	err := json.Unmarshal(data, t)
	if err != nil {
		return err
	}

	e.Name = t.Name
	e.Doc = t.Doc
	e.Returns = t.Returns
	for _, p := range t.Parameters {
		e.Fields = append(e.Fields, p.toField())
	}

	return nil
}

func (m *Method) UnmarshalJSON(data []byte) error {
	t := &struct {
		Method          string
		Doc             string
		Parameters      []*parameter
		PagingSupported bool
	}{}

	err := json.Unmarshal(data, t)
	if err != nil {
		return err
	}

	if name, ok := protocol.RestLiMethodNameMapping[strings.ToLower(t.Method)]; ok {
		m.Method = name
		m.Name = string(name)
	} else {
		log.Panicln("Unknown method", t.Method)
	}
	m.Doc = t.Doc
	m.PagingSupported = t.PagingSupported

	for _, p := range t.Parameters {
		m.Fields = append(m.Fields, p.toField())
	}

	return nil
}

func (a *Action) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &a.Endpoint)
	if err != nil {
		return err
	}

	a.ActionName = a.Endpoint.Name
	a.StructName = codegen.ExportedIdentifier(a.Name + "ActionParams")
	a.Endpoint.Name = a.StructName

	return nil
}

func (f *Finder) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &f.Endpoint)
	if err != nil {
		return err
	}

	f.FinderName = f.Endpoint.Name
	f.StructName = codegen.ExportedIdentifier(FindBy + codegen.ExportedIdentifier(f.Name) + "Params")
	f.Endpoint.Name = f.StructName

	p := &struct {
		PagingSupported bool
	}{}
	err = json.Unmarshal(data, p)
	if err != nil {
		return err
	}
	f.PagingSupported = p.PagingSupported

	if f.PagingSupported {
		f.Fields = append(f.Fields, models.Field{
			Name:     "start",
			Doc:      "PagingContext parameter",
			Type:     &models.Model{BuiltinType: &models.IntPrimitive},
			Optional: true,
		})
		f.Fields = append(f.Fields, models.Field{
			Name:     "count",
			Doc:      "PagingContext parameter",
			Type:     &models.Model{BuiltinType: &models.IntPrimitive},
			Optional: true,
		})
	}

	return nil
}

func (i *Identifier) EncodedVariableName() string {
	return i.Name + "Str"
}
