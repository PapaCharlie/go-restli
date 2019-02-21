package models

import (
	"encoding/json"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type Reference struct {
	Ns
	NameAndDoc
}

var ValidReferenceType = regexp.MustCompile("^[a-zA-Z_]([a-zA-Z_0-9.])*$")

func (r *Reference) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return errors.WithStack(err)
	}

	// ensure not a primitive
	var p Primitive
	if err := json.Unmarshal(data, &p); err == nil {
		return errors.New("Reference types cannot be primitives")
	}

	if ! ValidReferenceType.Match([]byte(name)) {
		return errors.Errorf("Illegal reference type: |%s|", name)
	}

	lastDot := strings.LastIndex(name, ".")
	r.Name = name[lastDot+1:]
	if lastDot != -1 {
		r.Namespace = name[:lastDot]
	}

	return nil
}

func (r *Reference) GetRegisteredModel() *Model {
	return GetRegisteredModel(Ns{r.Namespace}, r.Name)
}
