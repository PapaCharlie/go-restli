package suite

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated_extras/extras"
	. "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated_extras/extras/collectionWithTyperefKey"
	. "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated_extras/extras/collectionWithTyperefKey_test"
	"github.com/PapaCharlie/go-restli/v2/restli"
	"github.com/PapaCharlie/go-restli/v2/restlicodec"
	"github.com/PapaCharlie/go-restli/v2/restlidata"
	"github.com/PapaCharlie/go-restli/v2/restlidata/generated/com/linkedin/restli/common"
	"github.com/stretchr/testify/require"
)

func (o *Operation) CollectionWithTyperefKeyBatchCreateWithParams(t *testing.T, c Client) func(*testing.T) *MockResource {
	create := []*extras.SinglePrimitiveField{
		{String: "1"},
	}
	params := &BatchCreateParams{Test: "foo"}
	res, err := c.BatchCreate(create, params)
	require.NoError(t, err)
	expected := []*CreatedEntity{
		{
			Id:       1,
			Status:   http.StatusCreated,
			Location: restli.StringPointer("/collectionWithTyperefKey/1"),
		},
	}
	require.Equal(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockBatchCreate: func(ctx *restli.RequestContext, entities []*extras.SinglePrimitiveField, queryParams *BatchCreateParams) (createdEntities []*CreatedEntity, err error) {
				require.Equal(t, create, entities)
				require.Equal(t, params, queryParams)
				entity := &CreatedEntity{Id: 1}
				require.NoError(t, restli.SetLocation(ctx, entity))
				return []*CreatedEntity{entity}, nil
			},
		}
	}
}

func (o *Operation) CollectionWithTyperefKeyBatchGetWithParams(t *testing.T, c Client) func(*testing.T) *MockResource {
	expectedKeys := []extras.Temperature{1, 3}
	expectedEntities := &BatchEntities{
		Results: map[extras.Temperature]*extras.SinglePrimitiveField{},
	}
	for _, k := range expectedKeys {
		expectedEntities.Results[k] = &extras.SinglePrimitiveField{String: fmt.Sprint(k)}
	}
	expectedParams := &BatchGetParams{Test: "foo"}
	res, err := c.BatchGet(expectedKeys, expectedParams)
	require.NoError(t, err)
	requiredBatchResponseEquals(t, expectedEntities, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockBatchGet: func(ctx *restli.RequestContext, keys []extras.Temperature, queryParams *BatchGetParams) (results *BatchEntities, err error) {
				require.Equal(t, expectedKeys, keys)
				require.Equal(t, expectedParams, queryParams)
				return expectedEntities, nil
			},
		}
	}
}

func (o *Operation) CollectionWithTyperefKeyGet(t *testing.T, c Client) func(*testing.T) *MockResource {
	k := extras.Temperature(42)
	expected := &extras.SinglePrimitiveField{String: "42"}
	actual, err := c.Get(k)
	require.NoError(t, err)
	require.Equal(t, expected, actual)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockGet: func(ctx *restli.RequestContext, key extras.Temperature) (entity *extras.SinglePrimitiveField, err error) {
				require.Equal(t, k, key)
				return expected, nil
			},
		}
	}
}

func (o *Operation) CollectionWithTyperefKeyGetIncompleteResponse(t *testing.T, c Client) func(*testing.T) *MockResource {
	id := extras.Temperature(42)
	_, err := c.Get(id)
	require.Error(t, err)
	require.IsType(t, new(restlicodec.MissingRequiredFieldsError), err)

	c = o.newClient(t, false, 0).Interface().(Client)
	res, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t, &extras.SinglePrimitiveField{}, res)

	return func(t *testing.T) *MockResource {
		deliberateSkip(t, "Cannot return an incomplete response")
		return nil
	}
}

func (o *Operation) CollectionWithTyperefKeyFindWithPagingContext(t *testing.T, c Client) func(*testing.T) *MockResource {
	params := &FindBySearchParams{
		PagingContext: restlidata.NewPagingContext(0, 10),
		Keyword:       "test",
	}
	expected := &Elements{
		Paging: &common.CollectionMetadata{Total: restli.Int32Pointer(42)},
	}
	results, err := c.FindBySearch(params)
	require.NoError(t, err)
	require.Equal(t, expected, results)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockFindBySearch: func(ctx *restli.RequestContext, queryParams *FindBySearchParams) (results *Elements, err error) {
				require.Equal(t, params, queryParams)
				return expected, nil
			},
		}
	}
}

func (o *Operation) CollectionWithTyperefKeyFindWithPagingContextNoTotal(t *testing.T, c Client) func(*testing.T) *MockResource {
	params := &FindBySearchParams{
		PagingContext: restlidata.NewPagingContext(0, 10),
		Keyword:       "test",
	}
	expected := &Elements{}
	results, err := c.FindBySearch(params)
	require.NoError(t, err)
	require.Equal(t, expected, results)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockFindBySearch: func(ctx *restli.RequestContext, queryParams *FindBySearchParams) (results *Elements, err error) {
				require.Equal(t, params, queryParams)
				return expected, nil
			},
		}
	}
}

func (o *Operation) CollectionWithTyperefKeyActionOnEntity(t *testing.T, c Client) func(*testing.T) *MockResource {
	id := extras.Temperature(64)
	params := &OnEntityActionParams{Input: "Action on entity"}
	require.NoError(t, c.OnEntityAction(id, params))

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockOnEntityAction: func(ctx *restli.RequestContext, key extras.Temperature, actionParams *OnEntityActionParams) (err error) {
				require.Equal(t, id, key)
				require.Equal(t, params, actionParams)

				return nil
			},
		}
	}
}

func (o *Operation) CollectionWithTyperefKeyGetAll(t *testing.T, c Client) func(*testing.T) *MockResource {
	params := &GetAllParams{
		PagingContext: restlidata.NewPagingContext(5, 10),
	}
	res, err := c.GetAll(params)
	require.NoError(t, err)
	expected := new(Elements)
	require.Equal(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockGetAll: func(ctx *restli.RequestContext, queryParams *GetAllParams) (results *Elements, err error) {
				require.Equal(t, params, queryParams)
				return expected, nil
			},
		}
	}
}

func (o *Operation) CollectionWithTyperefKeyFinderNoParams(t *testing.T, c Client) func(*testing.T) *MockResource {
	res, err := c.FindByNoParams()
	require.NoError(t, err)
	expected := new(Elements)
	require.Equal(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockFindByNoParams: func(ctx *restli.RequestContext) (results *Elements, err error) {
				return expected, nil
			},
		}
	}
}

func (o *Operation) CollectionWithTyperefKeyFinderNoParamsWithPaging(t *testing.T, c Client) func(*testing.T) *MockResource {
	params := &FindByNoParamsWithPagingParams{
		PagingContext: restlidata.NewPagingContext(5, 10),
	}
	res, err := c.FindByNoParamsWithPaging(params)
	require.NoError(t, err)
	expected := new(Elements)
	require.Equal(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockFindByNoParamsWithPaging: func(ctx *restli.RequestContext, queryParams *FindByNoParamsWithPagingParams) (results *Elements, err error) {
				require.Equal(t, params, queryParams)
				return expected, nil
			},
		}
	}
}

func (o *Operation) CollectionWithTyperefKeyCreateWithParams(t *testing.T, c Client) func(*testing.T) *MockResource {
	create := &extras.SinglePrimitiveField{String: "1"}
	params := &CreateParams{Test: "foo"}
	res, err := c.Create(create, params)
	require.NoError(t, err)
	expected := &CreatedEntity{
		Id:     1,
		Status: http.StatusCreated,
	}
	require.Equal(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockCreate: func(ctx *restli.RequestContext, entity *extras.SinglePrimitiveField, queryParams *CreateParams) (createdEntity *CreatedEntity, err error) {
				require.Equal(t, create, entity)
				require.Equal(t, params, queryParams)
				return &CreatedEntity{Id: 1}, nil
			},
		}
	}
}
