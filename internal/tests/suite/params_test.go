package suite

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/params"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/params_test"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/stretchr/testify/require"
)

func (o *Operation) ParamsGetWithQueryparams(t *testing.T, c Client) func(*testing.T) *MockResource {
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
	expected := newMessage(1, "test message")

	m, err := c.Get(long, params)
	require.NoError(t, err)
	require.Equal(t, expected, m)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockGet: func(ctx *protocol.RequestContext, paramsId int64, queryParams *GetParams) (entity *conflictresolution.Message, err error) {
				require.Equal(t, long, paramsId)
				require.Equal(t, queryParams, params)
				return expected, nil
			},
		}
	}
}
