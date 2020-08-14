package tests

import (
	"testing"

	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite"
)

func TestInclude(t *testing.T) {
	expected := &testsuite.Include{
		Integer: int32(1),
		F1:      4.27,
	}
	testJsonEncoding(t, expected, `{ "integer": 1, "f1": 4.27 }`)
}
