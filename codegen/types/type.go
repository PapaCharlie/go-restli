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
	case t.Array != nil:
		return t.Array.InnerTypes()
	case t.Map != nil:
		return t.Map.InnerTypes()
	case t.NativeTyperef != nil:
		return utils.IdentifierSet{}
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
	case t.NativeTyperef != nil:
		return t.NativeTyperef.GoType()
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
			// Assume custom typerefs aren't primitive types under the hood (though this can be wired in later if
			// needed via the TyperefCodeProvider)
			if typeref.codeProvider != nil {
				return nil
			} else {
				return typeref.Type
			}
		}
	}
	return nil
}

func (t *RestliType) ShouldReference() bool {
	switch {
	case t.UnderlyingPrimitive() != nil:
		// No need to reference primitive types, makes it more convenient to call methods
		return false
	case t.Typeref() != nil:
		// Only reference typerefs that are explictly requested to be referenced
		tr := t.Typeref()
		return tr.codeProvider != nil && tr.codeProvider.ShouldReference()
	case t.IsMapOrArray():
		// Maps and arrays are already reference types, no need to take the pointer
		return false
	case t.Enum() != nil:
		return false
	case t.NativeTyperef != nil:
		// For now, it's assumed native typerefs deserialize to types that don't mind being copied
		return false
	}
	return true
}

func (t *RestliType) ReferencedType() Code {
	if t.ShouldReference() {
		return t.PointerType()
	} else {
		return t.GoType()
	}
}

func (t *RestliType) ZeroValueReference() Code {
	if tr := t.Typeref(); tr != nil {
		if tr.codeProvider != nil {
			if tr.codeProvider.ShouldReference() {
				return Nil()
			} else {
				return tr.codeProvider.ZeroValue()
			}
		} else {
			return Add(t.GoType()).Call(t.UnderlyingPrimitive().ZeroValueLit())
		}
	} else if p := t.UnderlyingPrimitive(); p != nil {
		return p.ZeroValueLit()
	} else if e := t.Enum(); e != nil {
		return e.zeroValueLit()
	} else if t.NativeTyperef != nil {
		return t.NativeTyperef.ZeroValue().Call()
	} else {
		return Nil()
	}
}

func (t *RestliType) IsMapOrArray() bool {
	return t.Array != nil || t.Map != nil
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
