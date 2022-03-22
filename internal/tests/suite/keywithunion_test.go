package suite

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/complexkey"
	"github.com/PapaCharlie/go-restli/restli"
	"github.com/stretchr/testify/require"

	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/keywithunion/keywithunion"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/keywithunion/keywithunion_test"
)

func (o *Operation) KeywithunionGet(t *testing.T, c Client) func(*testing.T) *MockResource {
	id := &Keywithunion_ComplexKey{
		Params: newKeyParams("param1", 5),
	}
	id.KeyWithUnion.Union.ComplexKey = &conflictresolution.ComplexKey{
		Part1: "one",
		Part2: 2,
		Part3: conflictresolution.Fruits_APPLE,
	}
	expected := &conflictresolution.LargeRecord{
		Key:     *id.KeyWithUnion.Union.ComplexKey,
		Message: conflictresolution.Message{Message: "test message"},
	}

	res, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockGet: func(ctx *restli.RequestContext, keywithunionId *Keywithunion_ComplexKey) (entity *conflictresolution.LargeRecord, err error) {
				require.Equal(t, id, keywithunionId)
				return expected, nil
			},
		}
	}
}

func newKeyParams(param1 string, param2 int64) *complexkey.KeyParams {
	return &complexkey.KeyParams{
		Param1: param1,
		Param2: &param2,
	}
}
