package suite

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/params"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) ParamsGetWithQueryparams(t *testing.T, c Client) {
	long := int64(100)
	apple := conflictresolution.Fruits_APPLE
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
	m, err := c.Get(100, params)
	require.NoError(t, err)
	require.Equal(t, newMessage(1, "test message"), m)
}
