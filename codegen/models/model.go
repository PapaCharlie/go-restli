package models

import (
	"log"
	"path/filepath"
	"strings"

	. "github.com/PapaCharlie/go-restli/codegen"
	. "github.com/dave/jennifer/jen"
)

type Identifier struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func (i *Identifier) GetIdentifier() Identifier {
	iCopy := *i
	return iCopy
}

func (i Identifier) GetQualifiedClasspath() string {
	return i.Namespace + "." + i.Name
}

func (i *Identifier) GoType() *Statement {
	if i.Name == "" {
		log.Panicln("name cannot be empty!", i)
	}
	return Qual(i.PackagePath(), i.Name)
}

func (i *Identifier) PackagePath() string {
	if i.Namespace == "" {
		log.Panicf("%+v has no namespace!", i)
	}
	var p string
	if CyclicModels[*i] {
		p = "conflictResolution"
	} else {
		p = strings.Replace(namespaceEscape.ReplaceAllString(i.Namespace, "${1}_internal${2}"), ".", "/", -1)
	}
	if GetPackagePrefix() != "" {
		p = filepath.Join(GetPackagePrefix(), p)
	}
	return p
}

func (i *Identifier) receiver() string {
	return ReceiverName(i.Name)
}

func (m *Model) GoType() (def *Statement) {
	if m.BuiltinType != nil {
		return m.BuiltinType.GoType()
	} else {
		return m.ComplexType.GoType()
	}
}

type BuiltinType interface {
	GoType() (def *Statement)
	restLiWriteToBuf(def *Group, accessor *Statement)
}

type ComplexType interface {
	GoType() (def *Statement)
	GenerateCode() (def *Statement)
	PackagePath() string
	GetIdentifier() Identifier
}

func GenerateModelCode(m ComplexType) *CodeFile {
	return &CodeFile{
		PackagePath: m.PackagePath(),
		Filename:    m.GetIdentifier().Name,
		Code:        m.GenerateCode(),
	}
}

func (m *Model) IsMapOrArray() bool {
	if _, isMap := m.BuiltinType.(*MapModel); isMap {
		return true
	}
	if _, isArray := m.BuiltinType.(*ArrayModel); isArray {
		return true
	}
	return false
}

func (m *Model) IsBytesOrPrimitive() bool {
	if _, isBytes := m.BuiltinType.(*BytesModel); isBytes {
		return true
	}
	if _, isPrimitive := m.BuiltinType.(*PrimitiveModel); isPrimitive {
		return true
	}
	return false
}

func (m *Model) PointerType() *Statement {
	c := Empty()
	if !m.IsMapOrArray() {
		c.Op("*")
	}
	c.Add(m.GoType())
	return c
}

func (m *Model) restLiWriteToBuf(def *Group, accessor *Statement) {
	if m.BuiltinType != nil {
		m.BuiltinType.restLiWriteToBuf(def, accessor)
	} else {
		def.Var().Id("tmp").String()
		def.List(Id("tmp"), Err()).Op("=").Add(accessor).Dot(RestLiEncode).Call(Id(Codec))
		IfErrReturn(def)
		writeStringToBuf(def, Id("tmp"))
	}
}

func writeStringToBuf(def *Group, s *Statement) *Statement {
	return def.Id("buf").Dot("WriteString").Call(s)
}

func (m *Model) RestLiURLEncodeModel(accessor *Statement) (def *Statement, hasError bool) {
	return m.RestLiEncodeModel(RestLiUrlEncoder, accessor)
}

func (m *Model) RestLiReducedEncodeModel(accessor *Statement) (def *Statement, hasError bool) {
	return m.RestLiEncodeModel(RestLiReducedEncoder, accessor)
}

func (m *Model) RestLiEncodeModel(encoder string, accessor *Statement) (*Statement, bool) {
	encoderRef := Qual(ProtocolPackage, encoder)
	if m.BuiltinType == nil {
		return Add(accessor).Dot(RestLiEncode).Call(encoderRef), true
	}

	if primitive, ok := m.BuiltinType.(*PrimitiveModel); ok {
		return Add(encoderRef).Dot("Encode" + ExportedIdentifier(primitive[0])).Call(accessor), false
	}

	if _, ok := m.BuiltinType.(*BytesModel); ok {
		return Add(encoderRef).Dot("EncodeBytes").Call(accessor), false
	}

	log.Panicln(m, "cannot be url encoded")
	return nil, false
}
