package restli

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/dave/jennifer/jen"
)

// Namespace field delimiter for Java
const NamespaceSep = "."

// Restli fields
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
)

// Restli types
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

type SnapshotParser struct {
	DestinationPackage string
	GeneratedTypes     []*GeneratedType
}

func (p *SnapshotParser) GenerateTypes(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Panicln(err)
	}

	snapshot := struct {
		Models []interface{} `json:"models"`
		Schema struct {
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
			Doc       string `json:"doc"`
			Path      string `json:"path"`
		} `json:"schema"`
	}{}

	decoder := json.NewDecoder(file)
	decoder.UseNumber()

	if err = decoder.Decode(&snapshot); err != nil {
		log.Panicln(err)
	}

	for _, m := range snapshot.Models {
		if this, isMap := m.(map[string]interface{}); isMap {
			p.generateType(this)
		}
	}
}

func (p *SnapshotParser) generateType(this map[string]interface{}) {
	namespace := getStringField(this, Namespace)
	if namespace == "" {
		log.Panicln("top-level defs must specify a namespace")
	}

	switch this[Type] {
	case Typeref:
		p.typeForTyperefField(jen.Empty(), this, namespace)
	case Enum:
		p.typeForEnumField(jen.Empty(), this, namespace)
	case Record:
		p.typeForRecordField(jen.Empty(), this, namespace)
	case Fixed:
		f := p.newGeneratedType(getStringField(this, Name), namespace)
		f.Definition.Type().Id(f.Name).Add(p.newFixedType(this).GoType())
	default:
		log.Panicln("illegal type in top-level definition", this[Type])
	}
}

func (p *SnapshotParser) typeForField(field *jen.Statement, this map[string]interface{}, namespace string) {
	if newNamespace := getStringField(this, Namespace); newNamespace != "" {
		namespace = newNamespace
	}
	switch this[Type] {
	case Typeref:
		p.typeForTyperefField(field, this, namespace)
	case Map:
		p.typeForMapField(field, this[Values], namespace)
	case Array:
		p.typeForArrayField(field, this[Items], namespace)
	case Record:
		p.typeForRecordField(field, this, namespace)
	case Enum:
		p.typeForEnumField(field, this, namespace)
	case Fixed:
		field.Add(p.newFixedType(this).GoType())
	default:
		field.Add(p.primitiveOrReference(this[Type].(string), namespace).GoType())
	}
}

func (p *SnapshotParser) typeForTyperefField(field *jen.Statement, this map[string]interface{}, namespace string) {
	tr := p.newGeneratedType(getStringField(this, Name), namespace)

	tr.Definition.Comment(getStringField(this, Doc)).Line()
	tr.Definition.Type().Id(tr.Name).Add(p.primitiveOrReference(getStringField(this, "ref"), namespace).GoType()).Line()

	field.Add(tr.GoType())
}

func (p *SnapshotParser) typeForEnumField(field *jen.Statement, this map[string]interface{}, namespace string) {
	e := p.newGeneratedType(getStringField(this, Name), namespace)

	var consts []jen.Code
	var symbolDocs map[string]interface{}
	if symbolDocsI, hasSymbolDocs := this[SymbolDocs]; hasSymbolDocs {
		symbolDocs = symbolDocsI.(map[string]interface{})
	}
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

	field.Add(e.GoType())
}

func (p *SnapshotParser) typeForUnionField(field *jen.Statement, unionTypes []interface{}, namespace string) {
	var unionFields []jen.Code
	for _, unionType := range unionTypes {
		unionField := jen.Empty()
		switch unionType.(type) {
		case string:
			stringUnionType := unionType.(string)
			t := p.primitiveOrReference(stringUnionType, namespace)
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
			unionField.Id(name).Op("*")
			p.typeForField(unionField, this, namespace)
			unionField.Tag(jsonTag(tag))
		default:
			log.Panicln("illegal type in union", unionType)
		}
		unionFields = append(unionFields, unionField)
	}
	field.Struct(unionFields...)
}

func (p *SnapshotParser) typeForArrayField(field *jen.Statement, items interface{}, namespace string) {
	field.Index()
	switch items.(type) {
	case string:
		field.Add(p.primitiveOrReference(items.(string), namespace).GoType())
	case []interface{}:
		p.typeForUnionField(field, items.([]interface{}), namespace)
	case map[string]interface{}:
		p.typeForField(field, items.(map[string]interface{}), namespace)
	default:
		log.Panicln("illegal items type for array field", items)
	}
}

func (p *SnapshotParser) typeForMapField(field *jen.Statement, values interface{}, namespace string) {
	field.Map(jen.String())
	switch values.(type) {
	case string:
		field.Add(p.primitiveOrReference(values.(string), namespace).GoType())
	case []interface{}:
		p.typeForUnionField(field, values.([]interface{}), namespace)
	case map[string]interface{}:
		p.typeForField(field, values.(map[string]interface{}), namespace)
	default:
		log.Panicln("illegal values type for map field", values)
	}
}

func (p *SnapshotParser) typeForRecordField(field *jen.Statement, this map[string]interface{}, namespace string) {
	r := p.newGeneratedType(getStringField(this, Name), namespace)
	r.Definition.Comment(getStringField(this, Doc)).Line()

	var structFields []jen.Code
	if include, hasInclude := this[Include]; hasInclude {
		for _, i := range include.([]interface{}) {
			switch i.(type) {
			case string:
				structFields = append(structFields, jen.Add(p.newReferenceType(i.(string), namespace).GoType()))
			case map[string]interface{}:
				mapIncludeField := i.(map[string]interface{})
				if mapIncludeField[Type] != Record {
					log.Panicln("include fields can only be references or records")
				}
				f := jen.Empty()
				p.typeForRecordField(f, mapIncludeField, namespace)
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
			p.typeForField(f, field, namespace)
		case []interface{}:
			p.typeForUnionField(f, field[Type].([]interface{}), namespace)
		case map[string]interface{}:
			p.typeForField(f, field[Type].(map[string]interface{}), namespace)
		default:
			log.Panicln("illegal field type", field[Type])
		}
		f.Tag(jsonTag(field[Name].(string)))
		structFields = append(structFields, f)
	}
	r.Definition.Type().Id(r.Name).Struct(structFields...)

	field.Add(r.GoType())
}

func (p *SnapshotParser) primitiveOrReference(ident string, namespace string) RestliType {
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
	case Fixed:
		log.Panicln("cannot specify", Fixed, "types without size")
		return nil
	default:
		return p.newReferenceType(ident, namespace)
	}
}

func (p *SnapshotParser) newFixedType(this map[string]interface{}) *PrimitiveType {
	size, err := this[Size].(json.Number).Int64()
	if err != nil {
		log.Panicln(err)
	}
	return &PrimitiveType{
		Type: Fixed,
		Size: int(size),
	}
}

func (p *SnapshotParser) newReferenceType(name, namespace string) (rt *ReferenceType) {
	rt = &ReferenceType{
		NamespacePrefix: p.DestinationPackage,
	}
	if strings.Count(string(name), NamespaceSep) != 0 {
		lastDot := strings.LastIndex(name, NamespaceSep)
		rt.Name = name[lastDot+1:]
		rt.Namespace = name[:lastDot]
	} else {
		if namespace == "" {
			log.Panicln("namespace cannot be empty")
		} else {
			rt.Name = name
			rt.Namespace = namespace
		}
	}
	return
}

func (p *SnapshotParser) newGeneratedType(name, namespace string) *GeneratedType {
	gt := &GeneratedType{
		ReferenceType: *p.newReferenceType(name, namespace),
	}
	p.GeneratedTypes = append(p.GeneratedTypes, gt)
	return gt
}
