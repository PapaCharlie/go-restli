package tests

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/stretchr/testify/require"

	. "github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/typerefs/collectionTyperef"
)

func (s *TestServer) CollectionTyperefGet(t *testing.T, c Client) {
	id := testsuite.Url("http://rest.li")
	message, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t, &conflictresolution.Message{Message: "test message"}, message)
}

func TestDecode(t *testing.T) {
	var u testsuite.Url
	expected := "asd"
	require.NoError(t, u.RestLiDecode(protocol.RestLiReducedEncoder, "asd"))
	require.Equal(t, string(u), expected)
}
