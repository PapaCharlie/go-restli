package suite

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collectionReturnEntity"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) CollectionReturnEntityCreate(t *testing.T, c Client) {
	e, err := c.Create(&conflictresolution.Message{
		Message: "test message",
	})
	require.NoError(t, err)
	require.Equal(t, &protocol.CreatedAndReturnedEntity[int64, *conflictresolution.Message]{
		CreatedEntity: protocol.CreatedEntity[int64]{
			Id:     1,
			Status: 201,
		},
		Entity: &conflictresolution.Message{
			Message: "test message",
		},
	}, e)
}
