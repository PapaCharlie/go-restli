package models

import (
	"encoding/json"
	"fmt"
	. "github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
	. "go-restli/codegen"
	"log"
	"path/filepath"
	"strings"
)

type NameAndDoc struct {
	Name string `json:"name"`
	Doc  string `json:"doc"`
}

type Ns struct {
	Namespace string `json:"namespace"`
}

func (n *Ns) PackagePath() string {
	if n.Namespace == "" {
		panic("no namespace for package!")
	}
	p := strings.Replace(n.Namespace, ".", "/", -1)
	if GetPackagePrefix() != "" {
		p = filepath.Join(GetPackagePrefix(), p)
	}
	return p
}

func (m *Model) Qual() *Statement {
	if m.Name == "" {
		log.Panicln("name cannot be empty!", m)
	}
	return Qual(m.PackagePath(), m.Name)
}

type ModelCodeGenerator interface {
	GenerateCode(packagePrefix string, previousNamespace string) (packagePath string, typeName string, def *Statement)
}

type Model struct {
	Ns
	NameAndDoc

	Array     *ArrayModel
	Enum      *EnumModel
	Fixed     *FixedModel
	Map       *MapModel
	Bytes     *BytesModel
	Primitive *PrimitiveModel
	Record    *RecordModel
	Reference *ModelReference
	Typeref   *TyperefModel
	Union     *UnionModel
}

func (m *Model) String() string {
	var modelTypeName string
	var model interface{}

	if m.Array != nil {
		modelTypeName = ArrayModelTypeName
		model = m.Array
	}
	if m.Enum != nil {
		modelTypeName = EnumModelTypeName
		model = m.Enum
	}
	if m.Fixed != nil {
		modelTypeName = FixedModelTypeName
		model = m.Fixed
	}
	if m.Map != nil {
		modelTypeName = MapModelTypeName
		model = m.Map
	}
	if m.Bytes != nil {
		modelTypeName = "bytes"
		model = "Bytes"
	}
	if m.Primitive != nil {
		modelTypeName = "primitive"
		model = m.Primitive
	}
	if m.Record != nil {
		modelTypeName = RecordTypeModelTypeName
		model = m.Record
	}
	if m.Reference != nil {
		modelTypeName = "reference"
		model = m.Reference
	}
	if m.Typeref != nil {
		modelTypeName = TyperefModelTypeName
		model = m.Typeref
	}
	if m.Union != nil {
		modelTypeName = "union"
		model = m.Union
	}
	if modelTypeName == "" {
		log.Panicln("all fields nil", m.Ns, m.NameAndDoc)
	}
	modelTypeName = strings.ToUpper(modelTypeName[:1]) + modelTypeName[1:]

	return fmt.Sprintf("Model{{Name: %s, Namespace: %s, Doc: %s}, %s: %s}", m.Name, m.Namespace, m.Doc, modelTypeName, model)
}

func (m *Model) GenerateModelCode(sourceFilename string) (f *CodeFile) {
	f = &CodeFile{
		SourceFilename: sourceFilename,
	}
	if m.Namespace != "" {
		f.PackagePath = m.PackagePath()
	}

	if m.Enum != nil {
		f.Code = m.Enum.generateCode()
		f.Filename = m.Name
	}

	if m.Record != nil {
		f.Code = m.Record.GenerateCode()
		f.Filename = m.Name
	}

	if m.Typeref != nil {
		f.Code = m.Typeref.generateCode()
		f.Filename = m.Name
	}

	if m.Fixed != nil {
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
			innerModel.Namespace = escapeNamespace(m.Namespace)
		}
	}

	return
}

func (m *Model) GoType() *Statement {
	// Arrays and maps have special notation
	if m.Array != nil {
		return m.Array.GoType()
	}
	if m.Map != nil {
		return m.Map.GoType()
	}

	// primitives don't need to be imported
	if m.Primitive != nil {
		return m.Primitive.GoType()
	}

	if m.Bytes != nil {
		return m.Bytes.GoType()
	}

	if m.Union != nil {
		return m.Union.GoType()
	}

	if m.Reference != nil {
		log.Panicln("ModelReference type not replaced", m)
	}

	// All of the following are type references
	if m.Enum != nil || m.Record != nil || m.Typeref != nil || m.Fixed != nil {
		if m.Namespace == "" {
			log.Panicln(m.Name, "has no namespace!")
		} else {
			return m.Qual()
		}
	}

	log.Panicln("all fields nil", m)
	return nil
}

func (m *Model) PointerType() *Statement {
	c := Empty()
	if !m.IsMapOrArray() {
		c.Op("*")
	}
	c.Add(m.GoType())
	return c
}

func (m *Model) IsMapOrArray() bool {
	return m.Array != nil || m.Map != nil
}

func (m *Model) UnmarshalJSON(data []byte) (error) {
	defer func() {
		if m.Reference != nil && m.Reference.Namespace != "" {
			if rm := m.Reference.GetRegisteredModel(); rm != nil {
				*m = *rm
			}
		}
	}()

	model := &struct {
		Ns
		NameAndDoc
		Type json.RawMessage
	}{}

	if err := json.Unmarshal(data, model); err != nil {
		var unmarshalErrors []error
		var subErr error

		var bytes BytesModel
		if subErr = json.Unmarshal(data, &bytes); subErr == nil {
			m.Bytes = &bytes
			return nil
		} else {
			unmarshalErrors = append(unmarshalErrors, subErr)
		}

		var primitive PrimitiveModel
		if subErr = json.Unmarshal(data, &primitive); subErr == nil {
			m.Primitive = &primitive
			return nil
		} else {
			unmarshalErrors = append(unmarshalErrors, subErr)
		}

		var reference ModelReference
		if subErr = json.Unmarshal(data, &reference); subErr == nil {
			m.Reference = &reference
			m.Namespace = escapeNamespace(reference.Namespace)
			m.Name = reference.Name
			return nil
		} else {
			unmarshalErrors = append(unmarshalErrors, subErr)
		}

		union := &UnionModel{}
		if subErr = json.Unmarshal(data, union); subErr == nil {
			m.Union = union
			return nil
		} else {
			unmarshalErrors = append(unmarshalErrors, subErr)
		}

		return errors.Errorf("illegal model type: %v, %v, (%s)", unmarshalErrors, err, string(data))
	}

	m.Namespace = escapeNamespace(model.Namespace)
	m.Name = model.Name
	m.Doc = model.Doc

	var modelType string
	if err := json.Unmarshal(model.Type, &modelType); err != nil {
		return errors.Wrap(err, "type must either be a string or union")
	}

	switch modelType {
	case RecordTypeModelTypeName:
		recordType := &RecordModel{}
		if err := json.Unmarshal(data, recordType); err == nil {
			m.Record = recordType
			return nil
		} else {
			return errors.WithStack(err)
		}
	case EnumModelTypeName:
		enumType := &EnumModel{}
		if err := json.Unmarshal(data, enumType); err == nil {
			m.Enum = enumType
			return nil
		} else {
			return errors.WithStack(err)
		}
	case FixedModelTypeName:
		fixedType := &FixedModel{}
		if err := json.Unmarshal(data, fixedType); err == nil {
			m.Fixed = fixedType
			return nil
		} else {
			return errors.WithStack(err)
		}
	case MapModelTypeName:
		mapType := &MapModel{}
		if err := json.Unmarshal(data, mapType); err == nil {
			m.Map = mapType
			return nil
		} else {
			return errors.WithStack(err)
		}
	case ArrayModelTypeName:
		arrayType := &ArrayModel{}
		if err := json.Unmarshal(data, arrayType); err == nil {
			m.Array = arrayType
			return nil
		} else {
			return errors.WithStack(err)
		}
	case TyperefModelTypeName:
		typerefType := &TyperefModel{}
		if err := json.Unmarshal(data, typerefType); err == nil {
			if m.Name == "IPAddress" {
				log.Println(typerefType)
			}

			if typerefType.Ref.Primitive == nil && typerefType.Ref.Bytes == nil {
				return errors.Errorf("illegal typeref is not a reference to a primitive or \"bytes\": %+v", typerefType)
			}

			m.Typeref = typerefType
			return nil
		} else {
			return errors.WithStack(err)
		}
	case BytesModelTypeName:
		m.Bytes = &BytesModel{}
		return nil
	}

	var primitiveType PrimitiveModel
	if err := json.Unmarshal(model.Type, &primitiveType); err == nil {
		m.Primitive = &primitiveType
		return nil
	}

	var referenceType ModelReference
	if err := json.Unmarshal(model.Type, &referenceType); err == nil {
		m.Reference = &referenceType
		//if referenceType.Namespace != "" {
		//	m.Namespace = escapeNamespace(referenceType.Namespace)
		//}
		//m.Name = referenceType.Name
		return nil
	}

	return errors.Errorf("could not deserialize %v into %v", string(data), m)
}
