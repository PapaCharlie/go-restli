package suite

import (
	"testing"

	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/simpleComplexKey"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) SimpleComplexKeyBatchGet(t *testing.T, c Client) {
	keys := []*SimpleComplexKey_ComplexKey{
		{
			SinglePrimitiveField: extras.SinglePrimitiveField{
				String: "1",
			},
			Params: &extras.SingleTyperefField{
				Temp: 1,
			},
		},
		{
			SinglePrimitiveField: extras.SinglePrimitiveField{
				String: "string:with:colons",
			},
			Params: &extras.SingleTyperefField{
				Temp: 42,
			},
		},
	}
	res, err := c.BatchGet(keys)
	require.NoError(t, err)
	expected := make(map[*SimpleComplexKey_ComplexKey]*extras.SinglePrimitiveField)
	for _, k := range keys {
		expected[k] = &k.SinglePrimitiveField
	}
	require.Equal(t, expected, res)
}

func (s *TestServer) SimpleComplexKeyGet(t *testing.T, c Client) {
	expected := &SimpleComplexKey_ComplexKey{
		SinglePrimitiveField: extras.SinglePrimitiveField{String: "string:with:colons"},
	}
	actual, err := c.Get(expected)
	require.NoError(t, err)
	require.Equal(t, &expected.SinglePrimitiveField, actual)
}

func (s *TestServer) SimpleComplexKeyBatchPartialUpdate(t *testing.T, c Client) {
	keys := []*SimpleComplexKey_ComplexKey{
		{
			SinglePrimitiveField: extras.SinglePrimitiveField{
				String: "1",
			},
			Params: &extras.SingleTyperefField{
				Temp: 1,
			},
		},
		{
			SinglePrimitiveField: extras.SinglePrimitiveField{
				String: "string:with:colons",
			},
			Params: &extras.SingleTyperefField{
				Temp: 42,
			},
		},
	}
	msg0 := "partial updated message"
	msg1 := "another partial message"

	actual, err := c.BatchPartialUpdate(map[*SimpleComplexKey_ComplexKey]*extras.SinglePrimitiveField_PartialUpdate{
		keys[0]: {
			Update_Fields: struct {
				String *string
			}{
				String: &msg0,
			},
		},
		keys[1]: {
			Update_Fields: struct {
				String *string
			}{
				String: &msg1,
			},
		},
	})
	require.NoError(t, err)

	expected := map[*SimpleComplexKey_ComplexKey]*protocol.BatchEntityUpdateResponse{
		keys[0]: {Status: 204},
		keys[1]: {Status: 205},
	}
	require.Equal(t, expected, actual)
}
