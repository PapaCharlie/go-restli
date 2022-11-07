package suite

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated/testsuite"
	"github.com/PapaCharlie/go-restli/v2/restli"
	"github.com/stretchr/testify/require"

	. "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated/testsuite/typerefs/collectionTyperef"
	. "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated/testsuite/typerefs/collectionTyperef_test"
)

func (o *Operation) CollectionTyperefGet(t *testing.T, c Client) func(*testing.T) *MockResource {
	id := testsuite.Url("http://rest.li")
	expected := &conflictresolution.Message{Message: "test message"}
	message, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t, expected, message)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockGet: func(ctx *restli.RequestContext, collectionTyperefId testsuite.Url) (entity *conflictresolution.Message, err error) {
				require.Equal(t, collectionTyperefId, id)
				return expected, nil
			},
		}
	}
}
