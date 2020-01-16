package codegen

import (
	"encoding/json"

	. "github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

type NamedType struct {
	Identifier
	SourceFile string
	Doc        string
}

func (t *NamedType) GetSourceFile() string {
	return t.SourceFile
}

type RestliType struct {
	Primitive *PrimitiveType
	Reference *Identifier
	Array     *RestliType
	Map       *RestliType
	Union     *UnionType
}

func (t *RestliType) UnmarshalJSON(data []byte) error {
	type _t RestliType
	err := json.Unmarshal(data, (*_t)(t))
	if err != nil {
		return err
	}

	switch {
	case t.Primitive != nil:
		return nil
	case t.Reference != nil:
		return nil
	case t.Array != nil:
		return nil
	case t.Map != nil:
		return nil
	case t.Union != nil:
		return nil
	default:
		return errors.Errorf("go-restli: RestliType declares no underlying type! (%s)", string(data))
	}
}

func (t *RestliType) InnerTypes() IdentifierSet {
	switch {
	case t.Primitive != nil:
		return nil
	case t.Reference != nil:
		innerTypes := make(IdentifierSet)
		innerTypes.Add(*t.Reference)
		return innerTypes
	case t.Array != nil:
		return t.Array.InnerTypes()
	case t.Map != nil:
		return t.Map.InnerTypes()
	default:
		return t.Union.InnerModels()
	}
}

func (t *RestliType) GoType() *Statement {
	switch {
	case t.Primitive != nil:
		return t.Primitive.GoType()
	case t.Reference != nil:
		return Qual(t.Reference.PackagePath(), t.Reference.Name)
	case t.Array != nil:
		return Index().Add(t.Array.GoType())
	case t.Map != nil:
		return Map(String()).Add(t.Map.GoType())
	default:
		return t.Union.GoType()
	}
}

func (t *RestliType) PointerType() *Statement {
	return Op("*").Add(t.GoType())
}

func (t *RestliType) WriteToBuf(def *Group, accessor *Statement) {
	switch {
	case t.Primitive != nil:
		writeStringToBuf(def, t.Primitive.encode(accessor))
	case t.Reference != nil:
		def.Var().Id("tmp").String()
		def.List(Id("tmp"), Err()).Op("=").Add(accessor).Dot(RestLiEncode).Call(Id(Codec))
		IfErrReturn(def)
		writeStringToBuf(def, Id("tmp"))
	case t.Array != nil:
		writeStringToBuf(def, Lit("List("))

		def.For(List(Id("idx"), Id("val")).Op(":=").Range().Add(accessor)).BlockFunc(func(def *Group) {
			def.If(Id("idx").Op("!=").Lit(0)).Block(Id("buf").Dot("WriteByte").Call(LitRune(','))).Line()
			t.Array.WriteToBuf(def, Id("val"))
		})

		def.Id("buf").Dot("WriteByte").Call(LitRune(')'))
		return
	case t.Map != nil:
		def.Id("buf").Dot("WriteByte").Call(LitRune('('))

		def.Id("idx").Op(":=").Lit(0)
		def.For(List(Id("key"), Id("val")).Op(":=").Range().Add(accessor)).BlockFunc(func(def *Group) {
			def.If(Id("idx").Op("!=").Lit(0)).Block(Id("buf").Dot("WriteByte").Call(LitRune(','))).Line()
			def.Id("idx").Op("++")
			writeStringToBuf(def, Id(Codec).Dot("EncodeString").Call(Id("key")))
			def.Id("buf").Dot("WriteByte").Call(LitRune(':'))
			t.Map.WriteToBuf(def, Id("val"))
		})

		def.Id("buf").Dot("WriteByte").Call(LitRune(')'))
		return
	default:
		label := "end" + canonicalizeAccessor(accessor)

		for _, m := range *t.Union {
			def.If(Add(accessor).Dot(m.name()).Op("!=").Nil()).BlockFunc(func(def *Group) {
				writeStringToBuf(def, Lit("("+m.Alias+":"))
				fieldAccessor := Add(accessor).Dot(m.name())
				if m.Type.Reference == nil {
					fieldAccessor = Op("*").Add(fieldAccessor)
				}
				m.Type.WriteToBuf(def, fieldAccessor)
				def.Id("buf").Dot("WriteByte").Call(LitRune(')'))
				def.Goto().Id(label)
			}).Line()
		}

		def.Id(label).Op(":")
	}
}

type GoRestliSpec struct {
	DataTypes []struct {
		Enum    *Enum
		Fixed   *Fixed
		Record  *Record
		Typeref *Typeref
	}
	Resources []Resource
}

func (s *GoRestliSpec) UnmarshalJSON(data []byte) error {
	type t GoRestliSpec
	err := json.Unmarshal(data, (*t)(s))
	if err != nil {
		return err
	}

	for _, t := range s.DataTypes {
		var complexType ComplexType
		switch {
		case t.Enum != nil:
			complexType = t.Enum
		case t.Fixed != nil:
			complexType = t.Fixed
		case t.Record != nil:
			complexType = t.Record
		case t.Typeref != nil:
			complexType = t.Typeref
		default:
			return errors.New("go-restli: Must declare at least one underlying type")
		}
		TypeRegistry.Register(complexType)
	}

	TypeRegistry.FlagCyclicDependencies()
	return nil
}

func (s *GoRestliSpec) GenerateClientCode() (codeFiles []*CodeFile) {
	for _, r := range s.Resources {
		codeFiles = append(codeFiles, r.GenerateCode()...)
	}
	return codeFiles
}
