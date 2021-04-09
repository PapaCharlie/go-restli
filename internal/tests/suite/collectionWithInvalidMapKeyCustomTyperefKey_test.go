package suite

import (
	"testing"

	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithInvalidMapKeyCustomTyperefKey"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) CollectionWithInvalidMapKeyCustomTyperefKeyBatchGet(t *testing.T, c Client) {
	k1, k2 := extras.BigDecimal2(1), extras.BigDecimal2(3)
	res, err := c.BatchGet([]extras.BigDecimal2{k1, k2})
	require.NoError(t, err)
	require.Equal(t, map[extras.BigDecimal2]*extras.SinglePrimitiveField{
		k1: {Integer: 1},
		k2: {Integer: 3},
	}, res)
}
