package tests

import (
	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/generated/conflictResolution"
	"github.com/stretchr/testify/require"
	"testing"

	. "github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/complexkey"
)

func (s *TestServer) ComplexkeyGet(t *testing.T, c Client) {
	id := &Complexkey_ComplexKey{
		Params: newKeyParams("param1", 5),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "one",
			Part2: 2,
			Part3: conflictresolution.Fruits_APPLE,
		},
	}
	res, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t, &conflictresolution.LargeRecord{
		Message: conflictresolution.Message{Message: "test message"},
		Key:     id.ComplexKey,
	}, res)
}

func (s *TestServer) ComplexkeyUpdate(t *testing.T, c Client) {
	t.SkipNow()
}

func (s *TestServer) ComplexkeyDelete(t *testing.T, c Client) {
	t.SkipNow()
}
