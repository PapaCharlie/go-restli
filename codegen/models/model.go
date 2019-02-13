package models

import (
	"encoding/json"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
	. "go-restli/codegen"
	"log"
	"strings"
)

type NameAndDoc struct {
	Name string `json:"name"`
	Doc  string `json:"doc"`
}

type Ns struct {
	Namespace string `json:"namespace"`
}

func (n *Ns) PackagePath(packagePrefix string) string {
	if n.Namespace == "" {
		panic("no namespace for package!")
	}
	p := strings.Replace(n.Namespace, ".", "/", -1)
	p = strings.Replace(p, "/internal/", "/internal_/", -1)
	if packagePrefix != "" {
		p = packagePrefix + "/" + p
	}
	return p
}

func (m *Model) Qual(packagePrefix string) *jen.Statement {
	if m.Name == "" {
		log.Panicln("name cannot be empty!", m)
	}
	return jen.Qual(m.PackagePath(packagePrefix), m.Name)
}

type ModelCodeGenerator interface {
	GenerateCode(packagePrefix string, previousNamespace string) (packagePath string, typeName string, def *jen.Statement)
}

type Model struct {
	Ns
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

	isTopLevel bool
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
		model = string(*m.Primitive)
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

func (m *Model) GenerateModelCode(packagePrefix string, sourceFilename string) (f *CodeFile) {
	f = &CodeFile{
		SourceFilename: sourceFilename,
	}
	if m.Namespace != "" {
		f.PackagePath = m.PackagePath(packagePrefix)
	}

	if m.Enum != nil {
		f.Code = m.Enum.generateCode(packagePrefix)
		f.Filename = m.Name
	}

	if m.Record != nil {
		f.Code = m.Record.generateCode(packagePrefix)
		f.Filename = m.Name
	}

	if m.Typeref != nil {
		f.Code = m.Typeref.generateCode(packagePrefix)
		f.Filename = m.Name
	}

	if m.Union != nil {
		f.Code = m.Union.generateCode(packagePrefix)
		f.Filename = m.Name
	}

	if m.isTopLevel && m.Fixed != nil {
		f.Code = m.Fixed.generateCode()
		f.Filename = m.Name
	}

	if f.Code == nil {
		return nil
	}

	if f.Code != nil && (f.PackagePath == "" || f.Filename == "") {
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
	if m.Typeref != nil {
		models = m.Typeref.InnerModels()
	}

	for _, innerModel := range models {
		if innerModel.Namespace == "" {
			innerModel.Namespace = m.Namespace
		}
	}

	return
}

func (m *Model) GoType(packagePrefix string) *jen.Statement {
	// Arrays and maps have special notation
	if m.Array != nil {
		return m.Array.GoType(packagePrefix)
	}
	if m.Map != nil {
		return m.Map.GoType(packagePrefix)
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
		return m.Union.GoType(packagePrefix)
	}

	if m.Reference != nil {
		return m.Reference.GoType(packagePrefix, m.Namespace)
	}

	// All of the following are type references
	if m.Enum != nil || m.Record != nil || m.Typeref != nil {
		if m.Namespace == "" {
			log.Panicln(m.Name, "has no namespace!")
		} else {
			return m.Qual(packagePrefix)
		}
	}

	panic("all fields nil")
}

func (m *Model) UnmarshalJSON(data []byte) error {
	model := &struct {
		Ns
		NameAndDoc
		Type json.RawMessage
	}{}

	if err := json.Unmarshal(data, model); err != nil {
		var unmarshalErrors []error

		var primitive Primitive
		if err := json.Unmarshal(data, &primitive); err == nil {
			m.Primitive = &primitive
			return nil
		} else {
			unmarshalErrors = append(unmarshalErrors, err)
		}

		var reference Reference
		if err := json.Unmarshal(data, &reference); err == nil {
			m.Reference = &reference
			m.Namespace = reference.Namespace
			m.Name = reference.Name
			return nil
		} else {
			unmarshalErrors = append(unmarshalErrors, err)
		}

		union := &Union{}
		if err := json.Unmarshal(data, union); err == nil {
			m.Union = union
			return nil
		} else {
			unmarshalErrors = append(unmarshalErrors, err)
		}

		return errors.Errorf("illegal model type: %v, %v, (%s)", unmarshalErrors, err, string(data))
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
			return errors.WithStack(err)
		}
	case EnumType:
		enumType := &Enum{}
		if err := json.Unmarshal(data, enumType); err == nil {
			m.Enum = enumType
			return nil
		} else {
			return errors.WithStack(err)
		}
	case FixedType:
		fixedType := &Fixed{}
		if err := json.Unmarshal(data, fixedType); err == nil {
			m.Fixed = fixedType
			return nil
		} else {
			return errors.WithStack(err)
		}
	case MapType:
		mapType := &Map{}
		if err := json.Unmarshal(data, mapType); err == nil {
			m.Map = mapType
			return nil
		} else {
			return errors.WithStack(err)
		}
	case ArrayType:
		arrayType := &Array{}
		if err := json.Unmarshal(data, arrayType); err == nil {
			m.Array = arrayType
			return nil
		} else {
			return errors.WithStack(err)
		}
	case TyperefType:
		typerefType := &Typeref{}
		if err := json.Unmarshal(data, typerefType); err == nil {
			m.Typeref = typerefType
			return nil
		} else {
			return errors.WithStack(err)
		}
	}

	var primitiveType Primitive
	if err := json.Unmarshal(model.Type, &primitiveType); err == nil {
		m.Primitive = &primitiveType
		return nil
	}

	var referenceType Reference
	if err := json.Unmarshal(model.Type, &referenceType); err == nil {
		m.Reference = &referenceType
		//if referenceType.Namespace != "" {
		//	m.Namespace = referenceType.Namespace
		//}
		//m.Name = referenceType.Name
		return nil
	}

	return errors.Errorf("could not deserialize %v into %v", string(data), m)
}
