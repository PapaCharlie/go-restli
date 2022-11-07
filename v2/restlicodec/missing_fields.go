package restlicodec

import (
	"fmt"
	"sort"
	"strings"
)

type KeyChecker interface {
	IsKeyExcluded(string) bool
}

type deserializationScopeSegment struct {
	segment string
	isArray bool
}

type missingFieldsTracker struct {
	excludedFields PathSpec
	scopeToIgnore  int
	currentScope   []deserializationScopeSegment
	missingFields  []string
}

func newMissingFieldsTracker(excludedFields PathSpec, leadingScopeToIgnore int) missingFieldsTracker {
	return missingFieldsTracker{
		excludedFields: excludedFields,
		scopeToIgnore:  leadingScopeToIgnore,
	}
}

type ExcludedFieldError string

func (e ExcludedFieldError) Error() string {
	return fmt.Sprintf("go-restli: Received read-only or create-only field at: %s", string(e))
}

func (t *missingFieldsTracker) IsKeyExcluded(key string) bool {
	excluded := t.enterMapScope(key) != nil
	t.exitScope()
	return excluded
}

func (t *missingFieldsTracker) enterMapScope(key string) error {
	t.currentScope = append(t.currentScope, deserializationScopeSegment{
		segment: key,
		isArray: false,
	})
	if len(t.currentScope) <= t.scopeToIgnore {
		return nil
	}
	excluded := genericMatches(t.excludedFields, t.currentScope[t.scopeToIgnore:], func(t deserializationScopeSegment) string {
		if t.isArray {
			return WildCard
		} else {
			return t.segment
		}
	})
	if excluded {
		return ExcludedFieldError(t.scopeString())
	} else {
		return nil
	}

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

func (t *missingFieldsTracker) recordMissingRequiredFields(fields map[string]struct{}) {
	scope := t.scopeString()
	if len(scope) > 0 {
		scope += "."
	}
	for f := range fields {
		// Required fields that are excluded (e.g. read-only fields during a Create call) should not cause a
		// MissingRequiredFieldsError
		if t.IsKeyExcluded(f) {
			continue
		}
		t.missingFields = append(t.missingFields, scope+f)
	}
}

func (t *missingFieldsTracker) checkMissingFields() error {
	if len(t.missingFields) != 0 {
		sort.Strings(t.missingFields)
		return &MissingRequiredFieldsError{Fields: t.missingFields}
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
