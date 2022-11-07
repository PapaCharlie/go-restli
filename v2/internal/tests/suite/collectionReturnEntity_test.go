package suite

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated/conflictResolution"
	. "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated/testsuite/collectionReturnEntity"
	. "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated/testsuite/collectionReturnEntity_test"
	"github.com/PapaCharlie/go-restli/v2/restli"
	"github.com/stretchr/testify/require"
)

func (o *Operation) CollectionReturnEntityCreate(t *testing.T, c Client) func(t *testing.T) *MockResource {
	key := &conflictresolution.Message{
		Message: "test message",
	}
	expected := &CreatedAndReturnedEntity{
		CreatedEntity: CreatedEntity{
			Id:     1,
			Status: 201,
		},
		Entity: &conflictresolution.Message{
			Message: "test message",
		},
	}
	e, err := c.Create(key)
	require.NoError(t, err)
	require.Equal(t, expected, e)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockCreate: func(ctx *restli.RequestContext, entity *conflictresolution.Message) (createdEntity *CreatedAndReturnedEntity, err error) {
				require.Equal(t, key, entity)
				return expected, nil
			},
		}
	}
}
