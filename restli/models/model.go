package models

import (
	"encoding/json"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
	"log"
	"strings"
)

type NameAndDoc struct {
	Name string `json:"name"`
	Doc  string `json:"doc"`
}

type namespace struct {
	Namespace string `json:"namespace"`
}

func (n *namespace) PackagePath(destinationPackage string) string {
	if n.Namespace == "" {
		panic("no namespace for package!")
	}
	p := strings.Replace(n.Namespace, ".", "/", -1)
	if destinationPackage != "" {
		p = destinationPackage + "/" + p
	}
	return p
}

func (m *Model) Qual(destinationPackage string) *jen.Statement {
	if m.Name == "" {
		log.Panicln("name cannot be empty!", m)
	}
	return jen.Qual(m.PackagePath(destinationPackage), m.Name)
}

type ModelCodeGenerator interface {
	GenerateCode(destinationPackage string, previousNamespace string) (packagePath string, typeName string, def *jen.Statement)
}

type Model struct {
	namespace
	NameAndDoc

	*Array
	*Enum
	*Fixed
	*Map
	*Primitive
	*Record
	*Reference
	*Typeref
	*Union
}

func (m *Model) String() string {
	var modelType string
	var model interface{}

	if m.Array != nil {
		modelType = "Array"
		model = m.Array
	}
	if m.Enum != nil {
		modelType = "Enum"
		model = m.Enum
	}
	if m.Fixed != nil {
		modelType = "Fixed"
		model = m.Fixed
	}
	if m.Map != nil {
		modelType = "Map"
		model = m.Map
	}
	if m.Primitive != nil {
		modelType = "Primitive"
		model = m.Primitive
	}
	if m.Record != nil {
		modelType = "Record"
		model = m.Record
	}
	if m.Reference != nil {
		modelType = "Reference"
		model = m.Reference
	}
	if m.Typeref != nil {
		modelType = "Typeref"
		model = m.Typeref
	}
	if m.Union != nil {
		modelType = "Union"
		model = m.Union
	}

	return fmt.Sprintf("Model{{Name: %s, Namespace: %s, Doc: %s}, %s: %s", m.Name, m.Namespace, m.Doc, modelType, model)
}

func (m *Model) generateCode(destinationPackage string) (def *jen.Statement, packagePath string, typeName string) {
	if m.Enum != nil {
		def = m.Enum.generateCode(destinationPackage)
		typeName = m.Name
	}

	if m.Record != nil {
		def = m.Record.generateCode(destinationPackage)
		typeName = m.Name
	}

	if m.Typeref != nil {
		def = m.Typeref.generateCode(destinationPackage)
		typeName = m.Name
	}

	if m.Union != nil {
		def = m.Union.generateCode(destinationPackage)
		typeName = m.Name
	}

	if m.Namespace != "" {
		packagePath = m.PackagePath(destinationPackage)
	}

	if def != nil && (packagePath == "" || typeName == "") {
		log.Panicf("code generators must have a namespace and name: %+v", m)
	}

	return
}

func (m *Model) InnerModels() (models []*Model) {
	if m.Array != nil {
		models = m.Array.InnerModels()
	}
	if m.Map != nil {
		models = m.Map.InnerModels()
	}
	if m.Record != nil {
		models = m.Record.InnerModels()
	}
	if m.Union != nil {
		models = m.Union.InnerModels()
	}

	for _, innerModel := range models {
		if innerModel.Namespace == "" {
			innerModel.Namespace = m.Namespace
		}
	}

	return
}

func (m *Model) GoType(destinationPackage string) *jen.Statement {
	// Arrays and maps have special notation
	if m.Array != nil {
		return m.Array.GoType(destinationPackage)
	}
	if m.Map != nil {
		return m.Map.GoType(destinationPackage)
	}

	// "Fixed" is an alias for [n]byte
	if m.Fixed != nil {
		return m.Fixed.GoType()
	}

	// primitives don't need to be imported
	if m.Primitive != nil {
		return m.Primitive.GoType()
	}

	if m.Union != nil {
		return m.Union.GoType(destinationPackage)
	}

	// All of the following are type references
	if m.Enum != nil || m.Record != nil || m.Typeref != nil || m.Reference != nil {
		if m.Namespace == "" {
			log.Panicln(m.Name, "has no namespace!")
		} else {
			return m.Qual(destinationPackage)
		}
	}

	panic("all fields nil")
}

func (m *Model) UnmarshalJSON(data []byte) error {
	model := &struct {
		namespace
		NameAndDoc
		Type json.RawMessage
	}{}

	if err := json.Unmarshal(data, model); err != nil {
		var subErr error

		var primitive Primitive
		if subErr = json.Unmarshal(data, &primitive); subErr == nil {
			m.Primitive = &primitive
			return nil
		}

		var reference Reference
		if subErr = json.Unmarshal(data, &reference); subErr == nil {
			m.Reference = &reference
			m.Namespace = reference.Namespace
			m.Name = reference.Name
			return nil
		}

		union := &Union{}
		if subErr = json.Unmarshal(data, union); subErr == nil {
			m.Union = union
			return nil
		}

		return errors.Wrapf(subErr, "illegal model type (original error: %v)", err)
	}

	m.Namespace = model.Namespace
	m.Name = model.Name
	m.Doc = model.Doc

	var modelType string
	if err := json.Unmarshal(model.Type, &modelType); err != nil {
		return errors.Wrap(err, "type must either be a string or union")
	}

	switch modelType {
	case RecordType:
		recordType := &Record{}
		if err := json.Unmarshal(data, recordType); err == nil {
			m.Record = recordType
			return nil
		} else {
			return errors.Cause(err)
		}
	case EnumType:
		enumType := &Enum{}
		if err := json.Unmarshal(data, enumType); err == nil {
			m.Enum = enumType
			return nil
		} else {
			return errors.Cause(err)
		}
	case FixedType:
		fixedType := &Fixed{}
		if err := json.Unmarshal(data, fixedType); err == nil {
			m.Fixed = fixedType
			return nil
		} else {
			return errors.Cause(err)
		}
	case MapType:
		mapType := &Map{}
		if err := json.Unmarshal(data, mapType); err == nil {
			m.Map = mapType
			return nil
		} else {
			return errors.Cause(err)
		}
	case ArrayType:
		arrayType := &Array{}
		if err := json.Unmarshal(data, arrayType); err == nil {
			m.Array = arrayType
			return nil
		} else {
			return errors.Cause(err)
		}
	case TyperefType:
		typerefType := &Typeref{}
		if err := json.Unmarshal(data, typerefType); err == nil {
			m.Typeref = typerefType
			return nil
		} else {
			return errors.Cause(err)
		}
	default:
		var primitiveType Primitive
		if err := json.Unmarshal(model.Type, &primitiveType); err == nil {
			m.Primitive = &primitiveType
			return nil
		} else {
			return errors.Cause(err)
		}
	}

	var referenceType Reference
	if err := json.Unmarshal(model.Type, &referenceType); err == nil {
		m.Reference = &referenceType
		return nil
	}

	return errors.Errorf("could not deserialize %v into %v", string(data), m)
}
