package suite

import (
	"testing"

	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/simpleComplexKey"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) SimpleComplexKeyBatchGet(t *testing.T, c Client) {
	keys := []*extras.SinglePrimitiveField{
		{String: "1"},
		{String: "string:with:colons"},
	}
	res, err := c.BatchGet(keys)
	require.NoError(t, err)
	expected := make(map[*extras.SinglePrimitiveField]*extras.SinglePrimitiveField)
	for _, k := range keys {
		expected[k] = k
	}
	require.Equal(t, expected, res)
}

func (s *TestServer) SimpleComplexKeyGet(t *testing.T, c Client) {
	expected := &extras.SinglePrimitiveField{String: "string:with:colons"}
	actual, err := c.Get(expected)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
