package tests

import (
	"testing"

	. "github.com/PapaCharlie/go-restli/internal/tests/generated/extras/collectionWithTyperefKey"
	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) CollectionWithTyperefKeyBatchGetWithParams(t *testing.T, c Client) {
	keys := []testsuite.Time{1, 3}
	res, err := c.BatchGet(keys, &BatchGetParams{Test: "foo"})
	require.NoError(t, err)
	expected := make(map[testsuite.Time]*testsuite.PrimitiveField)
	for _, k := range keys {
		expected[k] = &testsuite.PrimitiveField{Integer: int32(k)}
	}
	require.Equal(t, expected, res)
}

func (s *TestServer) CollectionWithTyperefKeyGet(t *testing.T, c Client) {
	var id testsuite.Time = 42
	res, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t, &testsuite.PrimitiveField{Integer: 42}, res)
}
