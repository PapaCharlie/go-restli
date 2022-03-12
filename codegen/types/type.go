package types

import (
	"log"

	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

func (t *NamedType) GetSourceFile() string {
	return t.SourceFile
}

func (t *RestliType) InnerTypes() utils.IdentifierSet {
	switch {
	case t.Primitive != nil:
		return nil
	case t.Reference != nil:
		return utils.NewIdentifierSet(*t.Reference)
	case t.IsMapOrArray():
		innerT, _ := t.InnerMapOrArray()
		return innerT.InnerTypes()
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

func (t *RestliType) ShouldReference() bool {
	switch {
	case t.Primitive != nil:
		// No need to reference primitive types
		return false
	case t.IsMapOrArray():
		// Maps and arrays are already reference types, no need to take the pointer
		return false
	}
	return t.Reference.Resolve().ShouldReference().ShouldUsePointer()
}

func (t *RestliType) ReferencedType() *Statement {
	if t.ShouldReference() {
		return t.PointerType()
	} else {
		return t.GoType()
	}
}

func (t *RestliType) IsMapOrArray() bool {
	return t.Array != nil || t.Map != nil
}

func (t *RestliType) InnerMapOrArray() (innerT RestliType, word string) {
	if t.Array != nil {
		return *t.Array, "Array"
	} else {
		return *t.Map, "Map"
	}
}

func (t *RestliType) Typeref() *Typeref {
	if t.Reference == nil {
		return nil
	}

	typeref, _ := t.Reference.Resolve().(*Typeref)
	return typeref
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

func (t *RestliType) ComplexKey() *ComplexKey {
	if t.Reference == nil {
		return nil
	}

	complexKey, _ := t.Reference.Resolve().(*ComplexKey)
	return complexKey
}

func (t *RestliType) StandaloneUnion() *StandaloneUnion {
	if t.Reference == nil {
		return nil
	}

	union, _ := t.Reference.Resolve().(*StandaloneUnion)
	return union
}

func (t *RestliType) PointerType() *Statement {
	return Op("*").Add(t.GoType())
}
