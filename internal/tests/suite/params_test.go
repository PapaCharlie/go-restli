package suite

import (
	"testing"

	nativetestsuite "github.com/PapaCharlie/go-restli/internal/tests/native/testsuite"
	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/params"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) ParamsGetWithQueryparams(t *testing.T, c Client) {
	long := int64(100)
	apple := nativetestsuite.Fruits_APPLE
	params := &GetParams{
		Int:         3,
		String:      "string",
		Long:        9223372036854775807,
		StringArray: []string{"string one", "string two"},
		MessageArray: []*conflictresolution.Message{
			{
				Message: "test message",
			},
			{
				Message: "another message",
			},
		},
		StringMap: map[string]string{
			"one": "string one",
			"two": "string two",
		},
		PrimitiveUnion: testsuite.UnionOfPrimitives{
			PrimitivesUnion: testsuite.UnionOfPrimitives_PrimitivesUnion{
				Long: &long,
			},
		},
		ComplexTypesUnion: testsuite.UnionOfComplexTypes{
			ComplexTypeUnion: testsuite.UnionOfComplexTypes_ComplexTypeUnion{
				Fruits: &apple,
			},
		},
		UrlTyperef: "http://rest.li",
	}

	var m *conflictresolution.Message
	var err error

	// The test server expects queries to match exactly (unlike json POST bodies which are matched structurally) because
	// it also serves to test the encoding of certain complex characters. The problem is that because one of the query
	// parameters is a map and map iteration ordering is nondeterministic by design, sometimes the queries don't match.
	// In this case, the server responds with a specific HTTP code that can be used to simply retry the request and
	// query encoding a number of times in an effort to reduce this test's flakiness.
	for i := 0; i < 10; i++ {
		m, err = c.Get(100, params)
		if e, ok := err.(*protocol.RestLiError); ok && e.Response.StatusCode == MismatchedQueriesStatus {
			continue
		} else {
			break
		}
	}
	require.NoError(t, err)
	require.Equal(t, newMessage(1, "test message"), m)
}
