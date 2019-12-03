package internal

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type ModelReference Identifier

var validReferenceType = regexp.MustCompile("^[a-zA-Z_]([a-zA-Z_0-9.])*$")

var illegalReferenceTypes = map[string]bool{
	ArrayModelTypeName:   true,
	MapModelTypeName:     true,
	BytesModelTypeName:   true,
	EnumModelTypeName:    true,
	FixedModelTypeName:   true,
	TyperefModelTypeName: true,
	RecordModelTypeName:  true,
}

func (r *ModelReference) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return errors.WithStack(err)
	}

	// ensure not a primitive
	var p PrimitiveModel
	if err := json.Unmarshal(data, &p); err == nil {
		return errors.New("Reference types cannot be primitives")
	}

	// ensure not bytes
	var b BytesModel
	if err := json.Unmarshal(data, &b); err == nil {
		return errors.New("Reference types cannot be \"bytes\"")
	}

	// sanity check: ensure the data type is neither map nor array
	if illegalReferenceTypes[name] {
		return errors.Errorf("Cannot be in %v, got: %s", illegalReferenceTypes, string(data))
	}

	if !validReferenceType.MatchString(name) {
		return errors.Errorf("Illegal reference type: |%s|", name)
	}

	lastDot := strings.LastIndex(name, ".")
	r.Name = name[lastDot+1:]
	if lastDot != -1 {
		r.Namespace = name[:lastDot]
	}
	if r.Namespace == "" {
		r.Namespace = currentNamespace
	}

	return nil
}

func (r *ModelReference) resolveOrRegisterPending(m *Model) bool {
	if r.Namespace == "" || r.Name == "" {
		log.Panicf("Unresolvable reference in %s: %+v", currentFile, r)
	}
	id := Identifier(*r)
	if t, ok := ModelRegistry.resolvedTypes[id]; ok {
		m.ComplexType = t.Type
		return true
	} else {
		ModelRegistry.addUnresolvedModel(id, m)
		return false
	}
}

func (r *ModelReference) Resolve() (ComplexType, error) {
	return ModelRegistry.Resolve(Identifier(*r))
}
