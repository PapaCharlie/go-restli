package internal

import (
	"encoding/json"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PapaCharlie/go-restli/codegen"
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

	if !validReferenceType.Match([]byte(name)) {
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

func (r *ModelReference) Resolve() (t ComplexType) {
	if r.Namespace == "" || r.Name == "" {
		log.Panicf("Unresolvable reference in %s: %+v", currentFile, r)
	}
	if m, ok := ModelRegistry[Identifier(*r)]; ok {
		return m.Type
	}

	oldCurrentFile := currentFile
	currentFile = filepath.Join(append(append([]string{codegen.PdscDirectory}, strings.Split(r.Namespace, ".")...), r.Name+".pdsc")...)
	m := new(Model)
	err := codegen.ReadJSONFromFile(currentFile, m)
	if err != nil {
		log.Panicln("Could not deserialize", currentFile, err)
	}
	if m.ComplexType == nil {
		log.Panicf("PDSC model loaded from %s does not define a ComplexType: %+v", currentFile, m)
	}
	ModelRegistry[m.ComplexType.GetIdentifier()] = &PdscModel{
		Type: m.ComplexType,
		File: currentFile,
	}
	currentFile = oldCurrentFile
	return m.ComplexType
}
