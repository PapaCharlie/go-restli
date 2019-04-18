package schema

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/codegen/models"
	. "github.com/dave/jennifer/jen"
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

func (t *ResourceModel) restLiURLEncode(accessor *Statement) (hasError bool, def *Statement) {
	return t.restLiEncode(codegen.RestLiUrlEncoder, accessor)
}

func (t *ResourceModel) restLiReducedEncode(accessor *Statement) (hasError bool, def *Statement) {
	return t.restLiEncode(codegen.RestLiReducedEncoder, accessor)
}

func (t *ResourceModel) restLiEncode(encoder string, accessor *Statement) (hasError bool, def *Statement) {
	def = Empty()
	encoderRef := Qual(codegen.ProtocolPackage, encoder)
	if t.Primitive != nil {
		def.Add(encoderRef).Dot("Encode" + codegen.ExportedIdentifier(t.Primitive[0])).Call(accessor)
		hasError = false
		return hasError, def
	}

	if t.Bytes != nil {
		def.Add(encoderRef).Dot("EncodeBytes").Call(accessor)
		hasError = false
		return hasError, def
	}

	if t.Typeref != nil || t.Enum != nil || t.Record != nil || t.Fixed != nil {
		def.Add(accessor).Dot(codegen.RestLiEncode).Call(encoderRef)
		hasError = true
		return hasError, def
	}

	log.Panicln(t, "cannot be url encoded")
	return
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
		f.Default = json.RawMessage(*p.Default)
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

	m.Method = t.Method
	if name, ok := RestliMethodNameMapping[strings.ToLower(t.Method)]; ok {
		m.Name = name
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
