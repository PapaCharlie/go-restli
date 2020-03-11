package tests

import (
	"testing"

	conflictresolution "github.com/bored-engineer/go-restli/internal/tests/generated/conflictResolution"
	. "github.com/bored-engineer/go-restli/internal/tests/generated/testsuite/simple"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) SimpleGet(t *testing.T, c Client) {
	res, err := c.Get()
	require.NoError(t, err)
	msg := "test message"
	require.Equal(t, &msg, res.Message, "Invalid response from server")
}

func (s *TestServer) SimpleUpdate(t *testing.T, c Client) {
	msg := "updated message"
	err := c.Update(&conflictresolution.Message{Message: &msg})
	require.NoError(t, err)
}

func (s *TestServer) SimpleDelete(t *testing.T, c Client) {
	err := c.Delete()
	require.NoError(t, err)
}
