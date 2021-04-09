package suite

import (
	"testing"

	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithCustomTyperefKey"
	"github.com/stretchr/testify/require"
)

func newBigDecimal(t *testing.T, s string) *extras.BigDecimal {
	d := new(extras.BigDecimal)
	require.NoError(t, d.UnmarshalRaw(s))
	return d
}

func (s *TestServer) CollectionWithCustomTyperefKeyEcho(t *testing.T, c Client) {
	d := newBigDecimal(t, "42.27")
	params := &EchoActionParams{Input: *d}
	res, err := c.EchoAction(params)
	require.NoError(t, err)
	require.Equal(t, d, res)
}

func (s *TestServer) CollectionWithCustomTyperefKeyFind(t *testing.T, c Client) {
	res, err := c.FindBySearch(&FindBySearchParams{Keyword: "test"})
	require.NoError(t, err)
	expected := extras.NewRecordWithBigDecimalWithDefaultValues()
	expected.BigDecimal = *newBigDecimal(t, "42")
	require.Equal(t, []*extras.RecordWithBigDecimal{expected}, res)
}

func (s *TestServer) CollectionWithCustomTyperefKeyBatchGet(t *testing.T, c Client) {
	k1, k2 := newBigDecimal(t, "1"), newBigDecimal(t, "3")
	res, err := c.BatchGet([]*extras.BigDecimal{k1, k2})
	require.NoError(t, err)

	expected1 := extras.NewRecordWithBigDecimalWithDefaultValues()
	expected1.BigDecimal = *k1
	expected1.OptionalBigDecimal = newBigDecimal(t, "27.27")

	expected2 := extras.NewRecordWithBigDecimalWithDefaultValues()
	expected2.BigDecimal = *k2

	require.Equal(t, map[*extras.BigDecimal]*extras.RecordWithBigDecimal{
		k1: expected1,
		k2: expected2,
	}, res)
}

func (s *TestServer) CollectionWithCustomTyperefKeyGet(t *testing.T, c Client) {
	d := newBigDecimal(t, "42")
	res, err := c.Get(d)
	require.NoError(t, err)

	expected := extras.NewRecordWithBigDecimalWithDefaultValues()
	expected.BigDecimal = *d

	require.Equal(t, expected, res)
}
