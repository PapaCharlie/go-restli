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

func (t *RestliType) Nil() *Statement {
	switch {
	case t.Primitive != nil:
		return t.Primitive.Nil()
	default:
		return Nil()
	}
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
		return Index().Add(t.Array.ReferencedType())
	case t.Map != nil:
		return Map(String()).Add(t.Map.ReferencedType())
	default:
		return t.Union.GoType()
	}
}

func (t *RestliType) ReferencedType() *Statement {
	switch {
	case t.Primitive != nil:
		// No need to reference primitive types, makes it more convenient to call methods
		return t.GoType()
	case t.Reference != nil:
		// If the typeref is backed by a primitive, then don't take the reference either
		if ref, ok := t.Reference.Resolve().(*Typeref); ok && ref.isPrimitive() {
			return t.GoType()
		}
	case t.Union != nil:
		// Union types are structs of references, we don't need to add another layer
		return t.GoType()
	}
	return t.PointerType()
}

func (t *RestliType) IsMapOrArray() bool {
	return t.Array != nil || t.Map != nil || (t.Primitive != nil && t.Primitive.IsBytes())
}

func (t *RestliType) PointerType() *Statement {
	if t.IsMapOrArray() {
		// Never use pointers to maps or arrays since they are already reference types. We can just use them as-is
		return t.GoType()
	} else {
		return Op("*").Add(t.GoType())
	}
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
		def.Switch().BlockFunc(func(def *Group) {
			for _, m := range *t.Union {
				def.Case(Add(accessor).Dot(m.name()).Op("!=").Nil()).BlockFunc(func(def *Group) {
					writeStringToBuf(def, Lit("("+m.Alias+":"))
					fieldAccessor := Add(accessor).Dot(m.name())
					if !(m.Type.Reference != nil || m.Type.IsMapOrArray()) {
						fieldAccessor = Op("*").Add(fieldAccessor)
					}
					m.Type.WriteToBuf(def, fieldAccessor)
					def.Id("buf").Dot("WriteByte").Call(LitRune(')'))
				})
			}
		})
	}
}

type GoRestliSpec struct {
	DataTypes []struct {
		Enum       *Enum
		Fixed      *Fixed
		Record     *Record
		Typeref    *Typeref
		ComplexKey *ComplexKey
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
		case t.ComplexKey != nil:
			complexType = t.ComplexKey
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
