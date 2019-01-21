package models

import (
	"encoding/json"
	"testing"
)

func TestRecord_UnmarshalJSON(t *testing.T) {
	r := &Record{}
	err := json.Unmarshal([]byte(`{"name": "Foo", "namespace": "io.papacharlie"}`), r)
	if err != nil {
		t.Fatalf("failed to unmarshal record: %v", err)
	}
	t.Logf("%v", r)
}

