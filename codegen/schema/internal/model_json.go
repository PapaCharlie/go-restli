package internal

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type typeField struct {
	Type string `json:"type"`
}

type docField struct {
	Doc string `json:"doc"`
}

type WrongTypeError struct {
	Expected, Actual string
}

func (w *WrongTypeError) Error() string {
	return fmt.Sprintf("models: Incorrect type, expected %s got %s", w.Expected, w.Actual)
}

type Model struct {
	BuiltinType BuiltinType
	ComplexType ComplexType
}

func (m *Model) String() string {
	var t string
	var s interface{}
	if m.BuiltinType != nil {
		t = "BuiltinType"
		s = m.BuiltinType
	}

	if m.ComplexType != nil {
		t = "ComplexType"
		s = m.ComplexType
	}

	return fmt.Sprintf("Model{%s: %+v}", t, s)
}

func (m *Model) innerModels() []*Model {
	type hasInnerModels interface {
		innerModels() []*Model
	}
	if im, ok := m.ComplexType.(hasInnerModels); ok {
		return im.innerModels()
	}
	if im, ok := m.BuiltinType.(hasInnerModels); ok {
		return im.innerModels()
	}
	return nil
}

func (m *Model) UnmarshalJSON(data []byte) (err error) {
	defer func() {
		if err != nil {
			m.register()

		}
	}()

	model := &struct {
		Namespace string `json:"namespace"`
		Type      json.RawMessage
		Aliases   []string `json:"aliases"`
	}{}

	if err = json.Unmarshal(data, model); err != nil {
		originalErr := err
		if strings.Contains(err.Error(), "cannot unmarshal array") {
			union := &UnionModel{}
			if err = json.Unmarshal(data, union); err == nil {
				m.BuiltinType = union
				return nil
			} else {
				err = errors.Wrapf(err, "could not deserialize union: %v, (%s)", originalErr, string(data))
				return err
			}
		} else {
			var unmarshalErrors []error

			var bytes BytesModel
			if err = json.Unmarshal(data, &bytes); err == nil {
				m.BuiltinType = &bytes
				return nil
			} else {
				unmarshalErrors = append(unmarshalErrors, err)
			}

			var primitive PrimitiveModel
			if err = json.Unmarshal(data, &primitive); err == nil {
				m.BuiltinType = &primitive
				return nil
			} else {
				unmarshalErrors = append(unmarshalErrors, err)
			}

			var reference ModelReference
			if err = json.Unmarshal(data, &reference); err == nil {
				m.ComplexType, err = reference.Resolve()
				return err
			} else {
				unmarshalErrors = append(unmarshalErrors, err)
			}

			err = errors.Errorf("illegal model type: %v, %+v (%s)", originalErr, unmarshalErrors, string(data))
			return err
		}
	}

	if model.Namespace != "" {
		oldNamespace := currentNamespace
		defer func() {
			currentNamespace = oldNamespace
		}()
		currentNamespace = model.Namespace
	}

	if len(model.Aliases) > 0 {
		defer func() {
			for _, alias := range model.Aliases {
				registerComplexType(m.ComplexType.CopyWithAlias(alias))
			}
		}()
	}

	var modelType string
	if err = json.Unmarshal(model.Type, &modelType); err != nil {
		err = errors.Wrapf(err, "type must either be a string or union (%s)", model.Type)
		return err
	}

	switch modelType {
	case RecordModelTypeName:
		recordType := &RecordModel{}
		if err = errors.WithStack(json.Unmarshal(data, recordType)); err == nil {
			m.ComplexType = recordType
			return nil
		} else {
			return err
		}
	case EnumModelTypeName:
		enumType := &EnumModel{}
		if err = errors.WithStack(json.Unmarshal(data, enumType)); err == nil {
			m.ComplexType = enumType
			return nil
		} else {
			return err
		}
	case FixedModelTypeName:
		fixedType := &FixedModel{}
		if err = errors.WithStack(json.Unmarshal(data, fixedType)); err == nil {
			m.ComplexType = fixedType
			return nil
		} else {
			return err
		}
	case MapModelTypeName:
		mapType := &MapModel{}
		if err = errors.WithStack(json.Unmarshal(data, mapType)); err == nil {
			m.BuiltinType = mapType
			return nil
		} else {
			return err
		}
	case ArrayModelTypeName:
		arrayType := &ArrayModel{}
		if err = errors.WithStack(json.Unmarshal(data, arrayType)); err == nil {
			m.BuiltinType = arrayType
			return nil
		} else {
			return err
		}
	case TyperefModelTypeName:
		typerefType := &TyperefModel{}
		if err = errors.WithStack(json.Unmarshal(data, typerefType)); err == nil {
			m.ComplexType = typerefType
			return nil
		} else {
			return err
		}
	case BytesModelTypeName:
		m.BuiltinType = &BytesModel{}
		return nil
	}

	var primitiveType PrimitiveModel
	if err = json.Unmarshal(model.Type, &primitiveType); err == nil {
		m.BuiltinType = &primitiveType
		return nil
	}

	var referenceType ModelReference
	if err = json.Unmarshal(model.Type, &referenceType); err == nil {
		m.ComplexType, err = referenceType.Resolve()
		return err
	}

	err = errors.Errorf("could not deserialize %v into %v", string(data), m)
	return err
}

func (m *Model) register() {
	if m.ComplexType != nil {
		registerComplexType(m.ComplexType)
	}

	for _, child := range m.innerModels() {
		child.register()
	}
}
