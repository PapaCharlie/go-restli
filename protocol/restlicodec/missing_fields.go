package restlicodec

import (
	"fmt"
	"sort"
	"strings"
)

type deserializationScopeSegment struct {
	segment string
	isArray bool
}

type missingFieldsTracker struct {
	currentScope []deserializationScopeSegment
	fields       []string
}

func (t *missingFieldsTracker) enterMapScope(key string) {
	t.currentScope = append(t.currentScope, deserializationScopeSegment{
		segment: key,
		isArray: false,
	})
}

func (t *missingFieldsTracker) enterArrayScope(index int) {
	t.currentScope = append(t.currentScope, deserializationScopeSegment{
		segment: fmt.Sprintf("[%d]", index),
		isArray: true,
	})
}

func (t *missingFieldsTracker) exitScope() {
	t.currentScope = t.currentScope[:len(t.currentScope)-1]
}

func (t *missingFieldsTracker) scopeString() string {
	var buf strings.Builder
	for i, p := range t.currentScope {
		if i > 0 && !p.isArray {
			buf.WriteByte('.')
		}
		buf.WriteString(p.segment)
	}
	return buf.String()
}

func (t *missingFieldsTracker) wrapDeserializationError(err error) error {
	switch err.(type) {
	case *DeserializationError:
		return err
	case nil:
		return nil
	default:
		return &DeserializationError{
			Scope: t.scopeString(),
			Err:   err,
		}
	}
}

func (t *missingFieldsTracker) RecordMissingRequiredFields(fields map[string]struct{}) {
	scope := t.scopeString()
	if len(scope) > 0 {
		scope += "."
	}
	for f := range fields {
		t.fields = append(t.fields, scope+f)
	}
}

func (t *missingFieldsTracker) CheckMissingFields() error {
	if len(t.fields) != 0 {
		sort.Strings(t.fields)
		return &MissingRequiredFieldsError{Fields: t.fields}
	} else {
		return nil
	}
}

type MissingRequiredFieldsError struct {
	Fields []string
}

func (m *MissingRequiredFieldsError) Error() string {
	return fmt.Sprintf("go-restli: Missing required fields %s", m.Fields)
}
