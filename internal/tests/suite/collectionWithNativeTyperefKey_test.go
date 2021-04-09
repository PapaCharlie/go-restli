package suite

import (
	"testing"
	"time"

	. "github.com/PapaCharlie/go-restli/internal/tests/native"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithNativeTyperefKey"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) CollectionWithNativeTyperefKeyEcho(t *testing.T, c Client) {
	ts := NewTime(42)
	params := &EchoActionParams{Input: ts}
	res, err := c.EchoAction(params)
	require.NoError(t, err)
	require.Equal(t, ts, res)
}

func (s *TestServer) CollectionWithNativeTyperefKeyFind(t *testing.T, c Client) {
	res, err := c.FindBySearch(&FindBySearchParams{Keyword: "test"})
	require.NoError(t, err)
	expected := extras.NewRecordWithTimeWithDefaultValues()
	expected.Time = NewTime(42)
	require.Equal(t, res, []*extras.RecordWithTime{expected})
}

func (s *TestServer) CollectionWithNativeTyperefKeyBatchGet(t *testing.T, c Client) {
	t1, t2 := NewTime(1), NewTime(3)
	res, err := c.BatchGet([]time.Time{t1, t2})
	require.NoError(t, err)

	expected1 := extras.NewRecordWithTimeWithDefaultValues()
	expected1.Time = t1
	optionalTime := NewTime(2)
	expected1.OptionalTime = &optionalTime

	expected2 := extras.NewRecordWithTimeWithDefaultValues()
	expected2.Time = t2

	require.Equal(t, map[time.Time]*extras.RecordWithTime{t1: expected1, t2: expected2}, res)
}

func (s *TestServer) CollectionWithNativeTyperefKeyGet(t *testing.T, c Client) {
	ts := NewTime(42)
	res, err := c.Get(ts)
	require.NoError(t, err)

	expected := extras.NewRecordWithTimeWithDefaultValues()
	expected.Time = ts

	require.Equal(t, expected, res)
}
