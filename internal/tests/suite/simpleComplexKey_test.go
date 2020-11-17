package suite

import (
	"testing"

	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/simpleComplexKey"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) SimpleComplexKeyBatchGet(t *testing.T, c Client) {
	keys := []*extras.SinglePrimitiveField{
		{Integer: 1},
		{Integer: 2},
	}
	res, err := c.BatchGet(keys)
	require.NoError(t, err)
	expected := make(map[*extras.SinglePrimitiveField]*extras.SinglePrimitiveField)
	for _, k := range keys {
		expected[k] = k
	}
	require.Equal(t, expected, res)
}
