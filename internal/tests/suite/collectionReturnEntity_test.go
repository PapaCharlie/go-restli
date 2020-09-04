package suite

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collectionReturnEntity"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) CollectionReturnEntityCreate(t *testing.T, c Client) {
	id, m, err := c.Create(&conflictresolution.Message{
		Message: "test message",
	})
	require.NoError(t, err)
	require.Equal(t, id, int64(1))
	require.Equal(t, m, &conflictresolution.Message{
		Message: "test message",
	})
}
