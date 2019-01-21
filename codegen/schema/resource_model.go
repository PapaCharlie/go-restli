package schema

import (
	"encoding/json"
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
	"go-restli/codegen/models"
	"log"
)

type ResourceModel struct {
	Primitive *models.Primitive
	Reference *models.Reference
	Array     *Array
	Map       *Map
}

type Array struct {
	Type  string
	Items *ResourceModel
}

func (a *Array) UnmarshalJSON(data []byte) error {
	type t Array
	if err := json.Unmarshal(data, (*t)(a)); err != nil {
		return err
	}
	if a.Type != models.ArrayType {
		return errors.Errorf("Not an array type: %s", string(data))
	}
	return nil
}

type Map struct {
	Type   string
	Values *ResourceModel
}

func (m *Map) UnmarshalJSON(data []byte) error {
	type t Map
	if err := json.Unmarshal(data, (*t)(m)); err != nil {
		return err
	}
	if m.Type != models.MapType {
		return errors.Errorf("Not a map type: %s", string(data))
	}
	return nil
}

func (t *ResourceModel) UnmarshalJSON(data []byte) error {
	var unmarshallErrors []error

	var primitive models.Primitive
	if err := json.Unmarshal(data, &primitive); err == nil {
		t.Primitive = &primitive
		return nil
	} else {
		unmarshallErrors = append(unmarshallErrors, err)
	}

	var reference models.Reference
	if err := json.Unmarshal(data, &reference); err == nil {
		t.Reference = &reference
		return nil
	} else {
		unmarshallErrors = append(unmarshallErrors, err)
	}

	var unescapedType string
	if err := json.Unmarshal(data, &unescapedType); err != nil {
		return errors.Wrap(err, "Could not deserialize type")
	}

	var array Array
	if err := json.Unmarshal([]byte(unescapedType), &array); err == nil {
		t.Array = &array
		return nil
	} else {
		unmarshallErrors = append(unmarshallErrors, err)
	}

	var _map Map
	if err := json.Unmarshal([]byte(unescapedType), &_map); err == nil {
		t.Map = &_map
		return nil
	} else {
		unmarshallErrors = append(unmarshallErrors, err)
	}

	return errors.Errorf("Failed to deserialize Resource model (can only be primitive, array, map or reference type): %v",
		unmarshallErrors)
}

func (t *ResourceModel) GoType(packagePrefix string) *jen.Statement {
	if t.Primitive != nil {
		return t.Primitive.GoType()
	}
	if t.Reference != nil {
		return jen.Qual(t.Reference.PackagePath(packagePrefix), t.Reference.Name)
	}
	if t.Array != nil {
		return jen.Index().Add(t.Array.Items.GoType(packagePrefix))
	}
	if t.Map != nil {
		return jen.Map(jen.String()).Add(t.Map.Values.GoType(packagePrefix))
	}
	log.Panicln("All models nil!", t)
	return nil
}
