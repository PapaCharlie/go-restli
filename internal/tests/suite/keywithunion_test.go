package suite

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/complexkey"
	"github.com/stretchr/testify/require"

	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/keywithunion/keywithunion"
)

func (s *TestServer) KeywithunionGet(t *testing.T, c Client) {
	id := &Keywithunion_ComplexKey{
		Params: newKeyParams("param1", 5),
	}
	id.KeyWithUnion.Union.ComplexKey = &conflictresolution.ComplexKey{
		Part1: "one",
		Part2: 2,
		Part3: conflictresolution.Fruits_APPLE,
	}

	res, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t,
		&conflictresolution.LargeRecord{
			Key:     *id.KeyWithUnion.Union.ComplexKey,
			Message: conflictresolution.Message{Message: "test message"},
		},
		res)
}

func newKeyParams(param1 string, param2 int64) *complexkey.KeyParams {
	return &complexkey.KeyParams{
		Param1: param1,
		Param2: &param2,
	}
}
