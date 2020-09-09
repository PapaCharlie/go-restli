package suite

import (
	"testing"

	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithTyperefKey"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) CollectionWithTyperefKeyBatchGetWithParams(t *testing.T, c Client) {
	keys := []extras.Temperature{1, 3}
	res, err := c.BatchGet(keys, &BatchGetParams{Test: "foo"})
	require.NoError(t, err)
	expected := make(map[extras.Temperature]*extras.SinglePrimitiveField)
	for _, k := range keys {
		expected[k] = &extras.SinglePrimitiveField{Integer: int32(k)}
	}
	require.Equal(t, expected, res)
}

func (s *TestServer) CollectionWithTyperefKeyGet(t *testing.T, c Client) {
	var id extras.Temperature = 42
	res, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t, &extras.SinglePrimitiveField{Integer: 42}, res)
}
