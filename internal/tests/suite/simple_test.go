package suite

import (
	"testing"

	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/simple"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/simpleWithPartialUpdate"
	"github.com/PapaCharlie/go-restli/protocol"
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

func (s *TestServer) SimplePartialUpdate(t *testing.T, c simplewithpartialupdate.Client) {
	err := c.PartialUpdate(&extras.SinglePrimitiveField_PartialUpdate{
		Set_Fields: extras.SinglePrimitiveField_PartialUpdate_Set_Fields{
			String: protocol.StringPointer("updated string"),
		},
	})
	require.NoError(t, err)
}
