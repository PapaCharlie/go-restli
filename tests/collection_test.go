package main

import (
	"testing"

	"github.com/PapaCharlie/go-restli/protocol"
	conflictresolution "github.com/PapaCharlie/go-restli/tests/generated/conflictResolution"
	. "github.com/PapaCharlie/go-restli/tests/generated/testsuite/collection"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) CollectionGet(t *testing.T, c *Client) {
	id := int64(1)
	res, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t, &conflictresolution.Message{Id: &id, Message: "test message"}, res, "Invalid response from server")
}

func (s *TestServer) CollectionUpdate(t *testing.T, c *Client) {
	id := int64(1)
	m := &conflictresolution.Message{
		Id:      &id,
		Message: "updated message",
	}
	err := c.Update(id, m)
	require.NoError(t, err)
}

func (s *TestServer) CollectionDelete(t *testing.T, c *Client) {
	id := int64(1)
	err := c.Delete(id)
	require.NoError(t, err)
}

func (s *TestServer) CollectionGet404(t *testing.T, c *Client) {
	m, err := c.Get(2)
	require.Errorf(t, err, "Did not receive an error from the server (got %+v)", m)
	require.Equal(t, 404, err.(*protocol.RestLiError).Status, "Unexpected status code from server")
}

func (s *TestServer) CollectionUpdate400(t *testing.T, c *Client) {
	id := int64(1)
	err := c.Update(id, &conflictresolution.Message{Id: &id, Message: "updated message"})
	require.Error(t, err, "Did not receive an error from the server")
	require.Equal(t, 400, err.(*protocol.RestLiError).Status, "Unexpected status code from server")
}
