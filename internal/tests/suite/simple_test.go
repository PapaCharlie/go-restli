package suite

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/simple"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) SimpleGet(t *testing.T, c Client) {
	res, err := c.Get()
	require.NoError(t, err)
	require.Equal(t, "test message", res.Message, "Invalid response from server")
}

func (s *TestServer) SimpleUpdate(t *testing.T, c Client) {
	err := c.Update(&conflictresolution.Message{Message: "updated message"})
	require.NoError(t, err)
}

func (s *TestServer) SimpleDelete(t *testing.T, c Client) {
	err := c.Delete()
	require.NoError(t, err)
}
