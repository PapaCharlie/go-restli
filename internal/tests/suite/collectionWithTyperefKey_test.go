package suite

import (
	"testing"

	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithTyperefKey"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
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
	res, err := c.Get(42)
	require.NoError(t, err)
	require.Equal(t, &extras.SinglePrimitiveField{Integer: 42}, res)
}

func (s *TestServer) CollectionWithTyperefKeyGetIncompleteResponse(t *testing.T, c Client) {
	oldValue := s.client.StrictResponseDeserialization

	s.client.StrictResponseDeserialization = false
	res, err := c.Get(42)
	require.NoError(t, err)
	require.Equal(t, &extras.SinglePrimitiveField{}, res)
	s.client.StrictResponseDeserialization = true
	_, err = c.Get(42)
	require.Error(t, err)
	require.IsType(t, new(restlicodec.MissingRequiredFieldsError), err)

	s.client.StrictResponseDeserialization = oldValue
}

func (s *TestServer) CollectionWithTyperefKeyFindWithPagingContext(t *testing.T, c Client) {
	_, total, err := c.FindBySearch(&FindBySearchParams{
		PagingContext: protocol.NewPagingContext(0, 10),
		Keyword:       "test",
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	require.Equal(t, 42, *total)
}

func (s *TestServer) CollectionWithTyperefKeyFindWithPagingContextNoTotal(t *testing.T, c Client) {
	_, total, err := c.FindBySearch(&FindBySearchParams{
		PagingContext: protocol.NewPagingContext(0, 10),
		Keyword:       "test",
	})
	require.NoError(t, err)
	require.Nil(t, total)
}

func TestEmbeddedPagingContext(t *testing.T) {
	var start, count int32
	start = 10
	count = 20
	tests := []struct {
		name     string
		params   FindBySearchParams
		expected string
	}{
		{
			name: "empty context",
			params: FindBySearchParams{
				Keyword: "foo",
			},
			expected: "keyword=foo&q=search",
		},
		{
			name: "start only",
			params: FindBySearchParams{
				PagingContext: protocol.PagingContext{
					Start: &start,
				},
				Keyword: "foo",
			},
			expected: "keyword=foo&q=search&start=10",
		},
		{
			name: "count only",
			params: FindBySearchParams{
				PagingContext: protocol.PagingContext{
					Count: &count,
				},
				Keyword: "foo",
			},
			expected: "count=20&keyword=foo&q=search",
		},
		{
			name: "full context",
			params: FindBySearchParams{
				PagingContext: protocol.PagingContext{
					Start: &start,
					Count: &count,
				},
				Keyword: "foo",
			},
			expected: "count=20&keyword=foo&q=search&start=10",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.params.EncodeQueryParams()
			require.NoError(t, err)
			require.Equal(t, test.expected, actual)
		})
	}
}
