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

type ResourceModel struct {
	models.Model
}

func (t *ResourceModel) UnmarshalJSON(data []byte) error {
	var unmarshallErrors []error

	var primitive models.PrimitiveModel
	if err := json.Unmarshal(data, &primitive); err == nil {
		t.Primitive = &primitive
		return nil
	} else {
		unmarshallErrors = append(unmarshallErrors, err)
	}

	var reference models.ModelReference
	if err := json.Unmarshal(data, &reference); err == nil {
		if reference.Namespace == "" {
			return errors.Wrapf(err, "%s was provided with no namespace", reference)
		}
		if m := reference.GetRegisteredModel(); m != nil {
			t.Model = *m
		} else {
			return errors.Errorf("unknown type: %+v", reference)
		}
		return nil
	} else {
		unmarshallErrors = append(unmarshallErrors, err)
	}

	var unescapedType string
	if err := json.Unmarshal(data, &unescapedType); err != nil {
		return errors.Wrap(err, "Could not deserialize type")
	}

	if err := json.Unmarshal([]byte(unescapedType), &t.Model); err == nil {
		return nil
	} else {
		unmarshallErrors = append(unmarshallErrors, err)
	}

	return errors.Errorf("Failed to deserialize Resource model (can only be primitive, array, map or reference type): %v",
		unmarshallErrors)
}

type parameter struct {
	models.NameAndDoc
	Type     ResourceModel
	Optional bool
	Default  *string
}

func (p parameter) toField() (f models.Field) {
	f.NameAndDoc = p.NameAndDoc
	f.Type = &p.Type.Model
	f.Optional = p.Optional
	if p.Default != nil {
		// Special case for string default values: the @Optional annotation doesn't escape the string, it puts it as a
		// literal, therefore we need to escape it before passing it in. Maps and lists are represented as `{}` and
		// `[]` respectively, so no escaping there, and numeric values don't need to be escaped in JSON.
		if p.Type.Model.Primitive != nil && *p.Type.Model.Primitive == models.StringPrimitive {
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
		models.NameAndDoc
		Parameters []parameter
		Returns    *ResourceModel
	}{}

	err := json.Unmarshal(data, t)
	if err != nil {
		return err
	}

	e.NameAndDoc = t.NameAndDoc
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
		Parameters      []parameter
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
			NameAndDoc: models.NameAndDoc{Name: "start", Doc: "PagingContext parameter"},
			Type:       &models.Model{Primitive: &models.IntPrimitive},
			Optional:   true,
		})
		f.Fields = append(f.Fields, models.Field{
			NameAndDoc: models.NameAndDoc{Name: "count", Doc: "PagingContext parameter"},
			Type:       &models.Model{Primitive: &models.IntPrimitive},
			Optional:   true,
		})
	}

	return nil
}

func (i *Identifier) EncodedVariableName() string {
	return i.Name + "Str"
}
