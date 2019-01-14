package restli

import (
	"encoding/json"
	"github.com/dave/jennifer/jen"
	"log"
	"os"
	"strings"
)

const (
	Name       = "name"
	Namespace  = "namespace"
	Type       = "type"
	Doc        = "doc"
	Include    = "include"
	Fields     = "fields"
	Default    = "default"
	Optional   = "optional"
	Size       = "size"
	Symbols    = "symbols"
	SymbolDocs = "symbolDocs"
	Alias      = "alias"

	NamespaceSep = "."
)

const (
	Int     = "int"
	Long    = "long"
	Float   = "float"
	Double  = "double"
	Boolean = "boolean"
	String  = "string"
	Bytes   = "bytes"
	Fixed   = "fixed"

	Record  = "record"
	Typeref = "typeref"
	Enum    = "enum"
	Array   = "array"
	Items   = "items"
	Map     = "map"
	Values  = "values"
)

type Generator struct {
	DestinationPackage            string
	GeneratedTypesNamespacePrefix string
	TyperefsNamespacePrefix       string
	GeneratedTypes                []GeneratedType

	knownTypeRefs map[string]bool
}

func NewGenerator(destinationPackage string, generatedTypesNamespacePrefix string, typerefsNamespacePrefix string) *Generator {
	return &Generator{
		DestinationPackage:            destinationPackage,
		GeneratedTypesNamespacePrefix: NsJoin(destinationPackage, generatedTypesNamespacePrefix),
		TyperefsNamespacePrefix:       NsJoin(destinationPackage, typerefsNamespacePrefix),
		GeneratedTypes:                nil,

		knownTypeRefs: make(map[string]bool),
	}
}

func (t *Generator) DecodeSnapshotModels(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Panicln(err)
	}

	snapshot := struct {
		Models []interface{} `json:"models"`
	}{}
	if err = json.NewDecoder(file).Decode(&snapshot); err != nil {
		log.Panicln(err)
	}

	for _, m := range snapshot.Models {
		if this, isMap := m.(map[string]interface{}); isMap {
			t.generateType(this)
		}
	}
}

func (t *Generator) generateType(this map[string]interface{}) {
	namespace := getStringField(this, Namespace)
	if namespace == "" {
		log.Panicln("top-level defs must specify a namespace")
	}

	switch this[Type] {
	case Typeref:
		t.registerTyperef(this)
	case Enum:
		t.typeForEnumField(jen.Empty(), this, namespace)
	case Record:
		t.typeForRecordField(jen.Empty(), this, namespace)
	default:
		log.Panicln("illegal type in top-level definition", this[Type])
	}
}

func (t *Generator) typeForField(field *jen.Statement, this map[string]interface{}, namespace string) {
	switch this[Type] {
	case Typeref:
		t.registerTyperef(this)
	case Map:
		t.typeForMapField(field, this[Values], namespace)
	case Array:
		t.typeForArrayField(field, this[Items], namespace)
	case Record:
		t.typeForRecordField(field, this, namespace)
	case Enum:
		t.typeForEnumField(field, this, namespace)
	default:
		field.Add(t.primitiveOrReference(this[Type].(string), namespace).GoType())
	}
}

func (t *Generator) registerTyperef(this map[string]interface{}) {
	t.knownTypeRefs[NsJoin(this[Namespace].(string), this[Name].(string))] = true
}

func (t *Generator) registerGeneratedType(generatedTypes ...GeneratedType) {
	t.GeneratedTypes = append(t.GeneratedTypes, generatedTypes...)
}

func (t *Generator) typeForEnumField(field *jen.Statement, this map[string]interface{}, namespace string) {
	e := GeneratedType{
		NamespacePrefix: t.GeneratedTypesNamespacePrefix,
		Namespace:       namespace,
		Name:            getStringField(this, Name),
	}
	var consts []jen.Code
	symbolDocs := this[SymbolDocs].(map[string]interface{})
	for i, constNameI := range this[Symbols].([]interface{}) {
		constName := constNameI.(string)
		doc := getStringField(symbolDocs, constName)
		def := jen.Id(constName)
		if i == 0 {
			def.Op("=").Id(e.Name).Call(jen.Iota())
		}
		consts = append(consts, def.Comment(doc))
	}
	e.Definition.Comment(getStringField(this, Doc)).Line()
	e.Definition.Type().Id(e.Name).String().Line()
	e.Definition.Const().Defs(consts...)
	t.registerGeneratedType(e)

	field.Add(e.GoType())
}

func (t *Generator) typeForUnionField(field *jen.Statement, unionTypes []interface{}, namespace string) {
	var unionFields []jen.Code
	for _, unionType := range unionTypes {
		unionField := jen.Empty()
		switch unionType.(type) {
		case string:
			stringUnionType := unionType.(string)
			t := t.primitiveOrReference(stringUnionType, namespace)
			unionField.Id(t.UnionFieldName()).Op("*").Add(t.GoType())
			switch t.(type) {
			case *PrimitiveType:
				unionField.Tag(jsonTag(t.(*PrimitiveType).Type))
			case *ReferenceType:
				rt := t.(*ReferenceType)
				unionField.Tag(jsonTag(NsJoin(rt.Namespace, rt.Name)))
			}
		case map[string]interface{}:
			this := unionType.(map[string]interface{})
			alias := getStringField(this, Alias)
			var tag string
			var name string
			if alias == "" {
				name = getStringField(this, Name)
				if name == "" {
					tag = getStringField(this, Type)
					name = capitalizeFirstLetter(tag)
				} else {
					newNamespace := getStringField(this, Namespace)
					if newNamespace == "" {
						newNamespace = namespace
					}
					tag = NsJoin(newNamespace, name)
				}
			}
			unionField.Id(name)
			t.typeForField(unionField, this, namespace)
			unionField.Tag(jsonTag(tag))
		default:
			log.Panicln("illegal type in union", unionType)
		}
		unionFields = append(unionFields, unionField)
	}
	field.Struct(unionFields...)
}

func (t *Generator) typeForArrayField(field *jen.Statement, items interface{}, namespace string) {
	field.Index()
	switch items.(type) {
	case string:
		field.Add(t.primitiveOrReference(items.(string), namespace).GoType())
	case []interface{}:
		t.typeForUnionField(field, items.([]interface{}), namespace)
	case map[string]interface{}:
		t.typeForField(field, items.(map[string]interface{}), namespace)
	default:
		log.Panicln("illegal items type for array field", items)
	}
}

func (t *Generator) typeForMapField(field *jen.Statement, values interface{}, namespace string) {
	field.Map(jen.String())
	switch values.(type) {
	case string:
		field.Add(t.primitiveOrReference(values.(string), namespace).GoType())
	case []interface{}:
		t.typeForUnionField(field, values.([]interface{}), namespace)
	case map[string]interface{}:
		t.typeForField(field, values.(map[string]interface{}), namespace)
	default:
		log.Panicln("illegal values type for map field", values)
	}
}

func (t *Generator) typeForRecordField(field *jen.Statement, this map[string]interface{}, namespace string) {
	r := GeneratedType{
		NamespacePrefix: t.GeneratedTypesNamespacePrefix,
		Namespace:       namespace,
		Name:            getStringField(this, Name),
	}
	r.Definition.Comment(getStringField(this, Doc)).Line()

	var structFields []jen.Code
	if include, hasInclude := this[Include]; hasInclude {
		for _, i := range include.([]interface{}) {
			switch i.(type) {
			case string:
				structFields = append(structFields, jen.Add(t.NewReferenceType(i.(string), namespace).GoType()))
			case map[string]interface{}:
				mapIncludeField := i.(map[string]interface{})
				if mapIncludeField[Type] != Record {
					log.Panicln("include fields can only be references or records")
				}
				f := jen.Empty()
				t.typeForRecordField(f, mapIncludeField, namespace)
				structFields = append(structFields, f)
			default:
				log.Panicln("include fields can only be references or records")
			}
		}
	}
	for _, field := range this[Fields].([]interface{}) {
		field := field.(map[string]interface{})
		f := jen.Id(getFieldName(field))
		switch field[Type].(type) {
		case string:
			t.typeForField(f, field, namespace)
		case []interface{}:
			t.typeForUnionField(f, field[Type].([]interface{}), namespace)
		case map[string]interface{}:
			t.typeForField(f, field[Type].(map[string]interface{}), namespace)
		default:
			log.Panicln("illegal field type", field[Type])
		}
		f.Tag(jsonTag(field[Name].(string)))
		structFields = append(structFields, f)
	}
	r.Definition.Type().Id(r.Name).Struct(structFields...)
	t.registerGeneratedType(r)
	field.Add(r.GoType())
}

func (t *Generator) primitiveOrReference(ident string, namespace string) RestliType {
	switch ident {
	case Int:
		fallthrough
	case Long:
		fallthrough
	case Float:
		fallthrough
	case Double:
		fallthrough
	case Boolean:
		fallthrough
	case String:
		fallthrough
	case Bytes:
		return &PrimitiveType{Type: ident}
	default:
		return t.NewReferenceType(ident, namespace)
	}
}

func (t *Generator) NewReferenceType(name, namespace string) (rt *ReferenceType) {
	if strings.Count(string(name), NamespaceSep) != 0 {
		lastDot := strings.LastIndex(name, NamespaceSep)
		rt = &ReferenceType{
			Name:      name[lastDot+1:],
			Namespace: name[:lastDot],
		}
	} else {
		if namespace == "" {
			log.Panicln("namespace cannot be empty")
		} else {
			rt = &ReferenceType{
				Name:      name,
				Namespace: namespace,
			}
		}
	}
	if t.knownTypeRefs[NsJoin(rt.Namespace, rt.Name)] {
		rt.NamespacePrefix = t.TyperefsNamespacePrefix
	} else {
		rt.NamespacePrefix = t.GeneratedTypesNamespacePrefix
	}
	return
}
