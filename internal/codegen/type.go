package codegen

import (
	"encoding/json"
	"log"

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
	default:
		return errors.Errorf("go-restli: RestliType declares no underlying type! (%s)", string(data))
	}
}

func (t *RestliType) InnerTypes() IdentifierSet {
	switch {
	case t.Primitive != nil:
		return nil
	case t.Reference != nil:
		return NewIdentifierSet(*t.Reference)
	case t.Array != nil:
		return t.Array.InnerTypes()
	case t.Map != nil:
		return t.Map.InnerTypes()
	default:
		log.Panicf("Illegal restli type: %+v", t)
		return nil
	}
}

func (t *RestliType) GoType() *Statement {
	switch {
	case t.Primitive != nil:
		return t.Primitive.GoType()
	case t.Reference != nil:
		return t.Reference.Qual()
	case t.Array != nil:
		return Index().Add(t.Array.ReferencedType())
	case t.Map != nil:
		return Map(String()).Add(t.Map.ReferencedType())
	default:
		log.Panicf("Illegal restli type: %+v", t)
		return nil
	}
}

func (t *RestliType) UnderlyingPrimitive() *PrimitiveType {
	switch {
	case t.Primitive != nil:
		return t.Primitive
	case t.Reference != nil:
		if typeref, ok := t.Reference.Resolve().(*Typeref); ok {
			return typeref.Type
		}
	}
	return nil
}

func (t *RestliType) UnderlyingPrimitiveZeroValueLit() *Statement {
	if t.Primitive != nil {
		return t.Primitive.zeroValueLit()
	} else {
		return Add(t.GoType()).Call(t.UnderlyingPrimitive().zeroValueLit())
	}
}

func (t *RestliType) ShouldReference() bool {
	switch {
	case t.UnderlyingPrimitive() != nil:
		// No need to reference primitive types or typerefs, makes it more convenient to call methods
		return false
	case t.IsMapOrArray():
		// Maps and arrays are already reference types, no need to take the pointer
		return false
	case t.Enum() != nil:
		return false
	}
	return true
}

func (t *RestliType) IsReferenceEncodable() bool {
	return t.Reference != nil && !t.ShouldReference()
}

func (t *RestliType) ReferencedType() *Statement {
	if t.ShouldReference() {
		return t.PointerType()
	} else {
		return t.GoType()
	}
}

func (t *RestliType) ZeroValueReference() *Statement {
	if p := t.UnderlyingPrimitive(); p != nil {
		return p.zeroValueLit()
	} else {
		return Nil()
	}
}

func (t *RestliType) IsMapOrArray() bool {
	return t.Array != nil || t.Map != nil || (t.UnderlyingPrimitive() != nil && t.UnderlyingPrimitive().IsBytes())
}

func (t *RestliType) Enum() *Enum {
	if t.Reference == nil {
		return nil
	}

	enum, _ := t.Reference.Resolve().(*Enum)
	return enum
}

func (t *RestliType) Record() *Record {
	if t.Reference == nil {
		return nil
	}

	record, _ := t.Reference.Resolve().(*Record)
	return record
}

func (t *RestliType) PointerType() *Statement {
	return Op("*").Add(t.GoType())
}

func (t *RestliType) WriteToBuf(def *Group, accessor *Statement, encoderAccessor *Statement, returnOnError ...Code) {
	switch {
	case t.Primitive != nil:
		writeStringToBuf(def, t.Primitive.encode(encoderAccessor, accessor))
	case t.Reference != nil:
		def.Err().Op("=").Add(accessor).Dot(RestLiEncode).Call(encoderAccessor, Id("buf"))
		IfErrReturn(def, append(append([]Code(nil), returnOnError...), Err())...)
	case t.Array != nil:
		writeArrayToBuf(def, accessor, t.Array, encoderAccessor, returnOnError...)
	case t.Map != nil:
		writeMapToBuf(def, accessor, t.Map, encoderAccessor, returnOnError...)
	default:
		log.Panicf("Illegal restli type: %+v", t)
	}
}

type GoRestliSpec struct {
	DataTypes []struct {
		Enum            *Enum            `json:"enum"`
		Fixed           *Fixed           `json:"fixed"`
		Record          *Record          `json:"record"`
		ComplexKey      *ComplexKey      `json:"complexKey"`
		StandaloneUnion *StandaloneUnion `json:"standaloneUnion"`
		Typeref         *Typeref         `json:"typeref"`
	} `json:"dataTypes"`
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
		case t.ComplexKey != nil:
			complexType = t.ComplexKey
		case t.StandaloneUnion != nil:
			complexType = t.StandaloneUnion
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
