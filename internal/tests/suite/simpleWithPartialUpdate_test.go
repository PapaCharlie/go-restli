package suite

import (
	"testing"

	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/simpleWithPartialUpdate"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/simpleWithPartialUpdate_test"
	"github.com/PapaCharlie/go-restli/restli"
	"github.com/stretchr/testify/require"
)

func (o *Operation) SimplePartialUpdate(t *testing.T, c Client) func(*testing.T) *MockResource {
	update := &extras.SinglePrimitiveField_PartialUpdate{
		Set_Fields: extras.SinglePrimitiveField_PartialUpdate_Set_Fields{
			String: restli.StringPointer("updated string"),
		},
	}
	err := c.PartialUpdate(update)
	require.NoError(t, err)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockPartialUpdate: func(ctx *restli.RequestContext, entity *extras.SinglePrimitiveField_PartialUpdate) (err error) {
				require.Equal(t, update, entity)
				return nil
			},
		}
	}
}
