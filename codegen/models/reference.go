package models

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type ModelReference Identifier

var validReferenceType = regexp.MustCompile("^[a-zA-Z_]([a-zA-Z_0-9.])*$")

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
	if name == ArrayModelTypeName || name == MapModelTypeName {
		return errors.Errorf("Cannot be an array or map (%s)", string(data))
	}

	if ! validReferenceType.Match([]byte(name)) {
		return errors.Errorf("Illegal reference type: |%s|", name)
	}

	lastDot := strings.LastIndex(name, ".")
	r.Name = name[lastDot+1:]
	if lastDot != -1 {
		r.Namespace = name[:lastDot]
	}

	return nil
}

func (r *ModelReference) Resolve() ComplexType {
	return ModelCache[Identifier(*r)]
}
