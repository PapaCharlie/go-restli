package suite

import (
	"testing"

	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/simpleComplexKey"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/simpleComplexKey_test"
	"github.com/PapaCharlie/go-restli/restli"
	"github.com/PapaCharlie/go-restli/restlidata/generated/com/linkedin/restli/common"
	"github.com/stretchr/testify/require"
)

func (o *Operation) SimpleComplexKeyBatchCreate(t *testing.T, c Client) func(*testing.T) *MockResource {
	create := []*extras.SinglePrimitiveField{
		{
			String: "1",
		},
		{
			String: "string:with:colons",
		},
	}
	res, err := c.BatchCreate(create)
	require.NoError(t, err)
	expected := []*CreatedAndReturnedEntity{
		{
			CreatedEntity: CreatedEntity{
				Id:     &SimpleComplexKey_ComplexKey{SinglePrimitiveField: *create[0]},
				Status: 201,
			},
			Entity: create[0],
		},
		{
			CreatedEntity: CreatedEntity{
				Id:     &SimpleComplexKey_ComplexKey{SinglePrimitiveField: *create[1]},
				Status: 202,
			},
			Entity: create[1],
		},
	}
	require.Equal(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockBatchCreate: func(ctx *restli.RequestContext, entities []*extras.SinglePrimitiveField) (createdEntities []*CreatedAndReturnedEntity, err error) {
				require.Equal(t, create, entities)
				return expected, nil
			},
		}
	}
}

func (o *Operation) SimpleComplexKeyBatchGet(t *testing.T, c Client) func(*testing.T) *MockResource {
	ids := []*SimpleComplexKey_ComplexKey{
		{
			SinglePrimitiveField: extras.SinglePrimitiveField{
				String: "1",
			},
			Params: &extras.SingleTyperefField{
				Temp: 1,
			},
		},
		{
			SinglePrimitiveField: extras.SinglePrimitiveField{
				String: "string:with:colons",
			},
			Params: &extras.SingleTyperefField{
				Temp: 42,
			},
		},
	}
	res, err := c.BatchGet(ids)
	require.NoError(t, err)
	expected := new(BatchEntities)
	for _, k := range ids {
		expected.AddResult(k, &k.SinglePrimitiveField)
	}
	requiredBatchResponseEquals(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockBatchGet: func(ctx *restli.RequestContext, keys []*SimpleComplexKey_ComplexKey) (results *BatchEntities, err error) {
				require.Equal(t, ids, keys)

				results = new(BatchEntities)
				for _, k := range ids {
					kCopy := &SimpleComplexKey_ComplexKey{SinglePrimitiveField: k.SinglePrimitiveField}
					results.AddResult(kCopy, &kCopy.SinglePrimitiveField)
				}

				return results, nil
			},
		}
	}
}

func (o *Operation) SimpleComplexKeyGet(t *testing.T, c Client) func(*testing.T) *MockResource {
	expected := &SimpleComplexKey_ComplexKey{
		SinglePrimitiveField: extras.SinglePrimitiveField{String: "string:with:colons"},
	}
	actual, err := c.Get(expected)
	require.NoError(t, err)
	require.Equal(t, &expected.SinglePrimitiveField, actual)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockGet: func(ctx *restli.RequestContext, key *SimpleComplexKey_ComplexKey) (entity *extras.SinglePrimitiveField, err error) {
				require.Equal(t, expected, key)
				return &expected.SinglePrimitiveField, nil
			},
		}
	}
}

func (o *Operation) SimpleComplexKeyBatchUpdateWithErrors(t *testing.T, c Client) func(*testing.T) *MockResource {
	keys := []*SimpleComplexKey_ComplexKey{
		{
			SinglePrimitiveField: extras.SinglePrimitiveField{
				String: "a",
			},
		},
		{
			SinglePrimitiveField: extras.SinglePrimitiveField{
				String: "b",
			},
		},
		{
			SinglePrimitiveField: extras.SinglePrimitiveField{
				String: "c",
			},
		},
	}
	ids := map[*SimpleComplexKey_ComplexKey]*extras.SinglePrimitiveField{
		keys[0]: &keys[0].SinglePrimitiveField,
		keys[1]: &keys[1].SinglePrimitiveField,
		keys[2]: &keys[2].SinglePrimitiveField,
	}
	expected := &BatchResponse{
		Results: map[*SimpleComplexKey_ComplexKey]*common.BatchEntityUpdateResponse{
			keys[1]: {Status: 204},
		},
		Errors: map[*SimpleComplexKey_ComplexKey]*common.ErrorResponse{
			keys[0]: {
				Status:         restli.Int32Pointer(400),
				Message:        restli.StringPointer("message"),
				ExceptionClass: restli.StringPointer("com.linkedin.restli.server.RestLiServiceException"),
				StackTrace:     restli.StringPointer("trace"),
			},
			keys[2]: {
				Status:         restli.Int32Pointer(500),
				Message:        nil,
				ExceptionClass: restli.StringPointer("com.linkedin.restli.server.RestLiServiceException"),
				StackTrace:     restli.StringPointer("trace"),
			},
		},
	}

	actual, err := c.BatchUpdate(ids)
	require.NoError(t, err)
	requiredBatchResponseEquals(t, expected, actual)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockBatchUpdate: func(ctx *restli.RequestContext, entities map[*SimpleComplexKey_ComplexKey]*extras.SinglePrimitiveField) (results *BatchResponse, err error) {
				requireComplexKeyMapEquals(t, ids, entities)
				return expected, nil
			},
		}
	}
}

func (o *Operation) SimpleComplexKeyBatchPartialUpdate(t *testing.T, c Client) func(*testing.T) *MockResource {
	keys := []*SimpleComplexKey_ComplexKey{
		{
			SinglePrimitiveField: extras.SinglePrimitiveField{
				String: "1",
			},
			Params: &extras.SingleTyperefField{
				Temp: 1,
			},
		},
		{
			SinglePrimitiveField: extras.SinglePrimitiveField{
				String: "string:with:colons",
			},
			Params: &extras.SingleTyperefField{
				Temp: 42,
			},
		},
	}
	partialUpdate := map[*SimpleComplexKey_ComplexKey]*extras.SinglePrimitiveField_PartialUpdate{
		keys[0]: {
			Set_Fields: extras.SinglePrimitiveField_PartialUpdate_Set_Fields{
				String: restli.StringPointer("partial updated message"),
			},
		},
		keys[1]: {
			Set_Fields: extras.SinglePrimitiveField_PartialUpdate_Set_Fields{
				String: restli.StringPointer("another partial message"),
			},
		},
	}

	expected := &BatchResponse{
		Results: map[*SimpleComplexKey_ComplexKey]*common.BatchEntityUpdateResponse{
			keys[0]: {Status: 204},
			keys[1]: {Status: 205},
		},
	}

	params := &BatchPartialUpdateParams{Param: 42}

	actual, err := c.BatchPartialUpdate(partialUpdate, params)
	require.NoError(t, err)

	requiredBatchResponseEquals(t, expected, actual)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockBatchPartialUpdate: func(ctx *restli.RequestContext, entities map[*SimpleComplexKey_ComplexKey]*extras.SinglePrimitiveField_PartialUpdate, queryParams *BatchPartialUpdateParams) (results *BatchResponse, err error) {
				require.Equal(t, params, queryParams)
				requireComplexKeyMapEquals(t, partialUpdate, entities)
				return expected, nil
			},
		}
	}
}
