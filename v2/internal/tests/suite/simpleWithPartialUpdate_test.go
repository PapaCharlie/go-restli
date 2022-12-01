package suite

import (
	"strings"
	"testing"

	"github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated_extras/extras"
	. "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated_extras/extras/simpleWithPartialUpdate"
	. "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated_extras/extras/simpleWithPartialUpdate_test"
	"github.com/PapaCharlie/go-restli/v2/restli"
	"github.com/stretchr/testify/require"
)

func (o *Operation) SimplePartialUpdate(t *testing.T, c Client) func(*testing.T) *MockResource {
	update := &extras.SinglePrimitiveField_PartialUpdate{
		Set_Fields: extras.SinglePrimitiveField_PartialUpdate_Set_Fields{
			String: restli.StringPointer("updated string"),
		},
	}
	params := &PartialUpdateParams{Param: "42"}
	e, err := c.PartialUpdate(update, params)
	require.Equal(t, &extras.SinglePrimitiveField{String: *update.Set_Fields.String}, e)
	require.NoError(t, err)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockPartialUpdate: func(ctx *restli.RequestContext, entity *extras.SinglePrimitiveField_PartialUpdate, queryParams *PartialUpdateParams) (updatedEntity *extras.SinglePrimitiveField, err error) {
				require.Equal(t, update, entity)
				require.Equal(t, params, queryParams)
				return e, nil
			},
		}
	}
}

func (o *Operation) SimplePartialUpdateWithTunnelling(t *testing.T, c Client) func(*testing.T) *MockResource {
	c = o.newClient(t, false, 100).Interface().(Client)

	update := &extras.SinglePrimitiveField_PartialUpdate{
		Set_Fields: extras.SinglePrimitiveField_PartialUpdate_Set_Fields{
			String: restli.StringPointer("updated string"),
		},
	}
	params := &PartialUpdateParams{Param: strings.Repeat("a", 200)}
	e, err := c.PartialUpdate(update, params)
	require.Equal(t, &extras.SinglePrimitiveField{String: *update.Set_Fields.String}, e)
	require.NoError(t, err)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockPartialUpdate: func(ctx *restli.RequestContext, entity *extras.SinglePrimitiveField_PartialUpdate, queryParams *PartialUpdateParams) (updatedEntity *extras.SinglePrimitiveField, err error) {
				require.Equal(t, update, entity)
				require.Equal(t, params, queryParams)
				return e, nil
			},
		}
	}
}
