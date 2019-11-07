package main

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/tests/generated/conflictResolution"
	. "github.com/PapaCharlie/go-restli/tests/generated/testsuite/actionSet"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) ActionsetEcho(t *testing.T) {
	c := Client{RestLiClient: s.client}
	input := "Is anybody out there?"
	output, err := c.EchoAction(&EchoActionParams{Input: input})
	require.NoError(t, err)
	require.Equal(t, &input, output, "Invalid response from server")
}

func (s *TestServer) ActionsetReturnInt(t *testing.T) {
	c := Client{RestLiClient: s.client}

	res, err := c.ReturnIntAction()
	require.NoError(t, err)
	i := int32(42)
	require.Equal(t, &i, res, "Invalid response from server")
}

func (s *TestServer) ActionsetReturnBool(t *testing.T) {
	c := Client{RestLiClient: s.client}

	res, err := c.ReturnBoolAction()
	require.NoError(t, err)
	b := true
	require.Equal(t, &b, res, "Invalid response from server")
}

func (s *TestServer) ActionsetEchoMessage(t *testing.T) {
	c := Client{RestLiClient: s.client}

	message := conflictresolution.Message{Message: "test message"}
	res, err := c.EchoMessageAction(&EchoMessageActionParams{Message: message})
	require.NoError(t, err)
	require.Equal(t, &message, res, "Invalid response from server")
}
