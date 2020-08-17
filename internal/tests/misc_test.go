package tests

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/stretchr/testify/require"
)

func TestInclude(t *testing.T) {
	expected := &testsuite.Include{
		Integer: int32(1),
		F1:      4.27,
	}
	testJsonEncoding(t, expected, `{ "integer": 1, "f1": 4.27 }`)
}

// TestDefaults tests that default values are loaded correctly (see
// rest.li-test-suite/client-testsuite/schemas/testsuite/Defaults.pdsc) for the default values used here
func TestDefaults(t *testing.T) {
	five := int32(5)
	d := testsuite.NewDefaultsWithDefaultValues()
	require.Equal(t, int32(1), *d.DefaultInteger)
	require.Equal(t, int64(23), *d.DefaultLong)
	require.Equal(t, float32(52.5), *d.DefaultFloat)
	require.Equal(t, float64(66.5), *d.DefaultDouble)
	require.Equal(t, protocol.Bytes("@ABC"), *d.DefaultBytes)
	require.Equal(t, string("default string"), *d.DefaultString)
	require.Equal(t, conflictresolution.Fruits_APPLE, *d.DefaultEnum)
	require.Equal(t, testsuite.Fixed5{1, 2, 3, 4, 5}, *d.DefaultFixed)
	require.Equal(t, testsuite.PrimitiveField{Integer: 10}, *d.DefaultRecord)
	require.Equal(t, []int32{1, 3, 5}, *d.DefaultArray)
	require.Equal(t, map[string]int32{"a": 1, "b": 2}, *d.DefaultMap)
	require.Equal(t, testsuite.Defaults_DefaultUnion{Int: &five}, *d.DefaultUnion)
}
