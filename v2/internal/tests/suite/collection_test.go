package suite

import (
	"net/http"
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated/testsuite"
	. "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated/testsuite/collection"
	"github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated/testsuite/collection/subcollection"
	subcollectiontest "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated/testsuite/collection/subcollection_test"
	colletionSubSimple "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated/testsuite/collection/subsimple"
	colletionSubSimpletest "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated/testsuite/collection/subsimple_test"
	. "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated/testsuite/collection_test"
	"github.com/PapaCharlie/go-restli/v2/restli"
	"github.com/PapaCharlie/go-restli/v2/restlidata/generated/com/linkedin/restli/common"
	"github.com/stretchr/testify/require"
)

func (o *Operation) CollectionCreate(t *testing.T, c Client) func(*testing.T) *MockResource {
	message := &conflictresolution.Message{
		Message: "test message",
	}
	returned := &CreatedEntity{
		Id:     1,
		Status: 201,
	}
	id, err := c.Create(message)
	require.NoError(t, err)
	require.Equal(t, returned, id)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockCreate: func(ctx *restli.RequestContext, entity *conflictresolution.Message) (createdEntity *CreatedEntity, err error) {
				require.Equal(t, message, entity)
				return returned, nil
			},
		}
	}
}

func (o *Operation) CollectionCreate500(t *testing.T, c Client) func(*testing.T) *MockResource {
	message := newMessage(3, "internal error test")
	id, err := c.Create(message)
	require.Errorf(t, err, "Did not receive an error from the server (got %+v)", id)
	require.Equal(t, err.(*restli.Error).Response.StatusCode, 500)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockCreate: func(ctx *restli.RequestContext, entity *conflictresolution.Message) (createdEntity *CreatedEntity, err error) {
				require.Equal(t, message, entity)
				return nil, &common.ErrorResponse{
					Status: restli.Int32Pointer(http.StatusInternalServerError),
				}
			},
		}
	}
}

func (o *Operation) CollectionCreateErrorDetails(t *testing.T, c Client) func(*testing.T) *MockResource {
	message := newMessage(3, "error details test")
	id, err := c.Create(message)
	require.Errorf(t, err, "Did not receive an error from the server (got %+v)", id)
	require.Equal(t, err.(*restli.Error).Response.StatusCode, 400)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockCreate: func(ctx *restli.RequestContext, entity *conflictresolution.Message) (createdEntity *CreatedEntity, err error) {
				require.Equal(t, message, entity)
				return nil, &common.ErrorResponse{Status: restli.Int32Pointer(http.StatusBadRequest)}
			},
		}
	}
}

func (o *Operation) CollectionGet(t *testing.T, c Client) func(*testing.T) *MockResource {
	id := int64(1)
	expected := newMessage(id, "test message")
	res, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockGet: func(ctx *restli.RequestContext, collectionId int64) (entity *conflictresolution.Message, err error) {
				require.Equal(t, id, collectionId)
				return expected, nil
			},
		}
	}
}

func (o *Operation) CollectionUpdate(t *testing.T, c Client) func(*testing.T) *MockResource {
	id := int64(1)
	expected := newMessage(id, "updated message")
	err := c.Update(id, expected)
	require.NoError(t, err)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockUpdate: func(ctx *restli.RequestContext, collectionId int64, entity *conflictresolution.Message) (err error) {
				require.Equal(t, id, collectionId)
				require.Equal(t, expected, entity)
				return nil
			},
		}
	}
}

func (o *Operation) CollectionDelete(t *testing.T, c Client) func(*testing.T) *MockResource {
	id := int64(1)
	err := c.Delete(id)
	require.NoError(t, err)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockDelete: func(ctx *restli.RequestContext, collectionId int64) (err error) {
				require.Equal(t, id, collectionId)
				return nil
			},
		}
	}
}

func (o *Operation) CollectionGet404(t *testing.T, c Client) func(*testing.T) *MockResource {
	id := int64(2)
	_, err := c.Get(id)
	require.Error(t, err)
	require.Equal(t, 404, err.(*restli.Error).Response.StatusCode, "Unexpected status code from server")

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockGet: func(ctx *restli.RequestContext, collectionId int64) (entity *conflictresolution.Message, err error) {
				return nil, &common.ErrorResponse{
					Status: restli.Int32Pointer(http.StatusNotFound),
				}
			},
		}
	}
}

func (o *Operation) CollectionUpdate400(t *testing.T, _ Client) func(*testing.T) *MockResource {
	t.Log("It is impossible to craft the request required using the generated code because it would require a field " +
		"to be deliberately missing. This can be chalked up as a win for the generated code's safety.")
	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockUpdate: func(ctx *restli.RequestContext, collectionId int64, entity *conflictresolution.Message) (err error) {
				id := int64(1)
				require.Equal(t, id, collectionId)
				require.Equal(t, &conflictresolution.Message{Id: &id}, entity)
				return nil
			},
		}
	}
}

func (o *Operation) CollectionSearchFinder(t *testing.T, c Client) func(*testing.T) *MockResource {
	params := &FindBySearchParams{Keyword: "message"}
	expectedMessages := &FindBySearchElements{
		Elements: []*conflictresolution.Message{
			newMessage(1, "test message"),
			newMessage(2, "another message"),
		},
		Paging: &common.CollectionMetadata{
			Count: 10,
			Total: restli.Int32Pointer(2),
		},
		Metadata: &testsuite.Optionals{
			OptionalLong:   restli.Int64Pointer(5),
			OptionalString: restli.StringPointer("metadata"),
		},
	}

	res, err := c.FindBySearch(params)
	require.NoError(t, err)
	require.Equal(t, expectedMessages, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockFindBySearch: func(ctx *restli.RequestContext, queryParams *FindBySearchParams) (results *FindBySearchElements, err error) {
				require.Equal(t, params, queryParams)
				return expectedMessages, nil
			},
		}
	}
}

func (o *Operation) CollectionPartialUpdate(t *testing.T, c Client) func(*testing.T) *MockResource {
	id := int64(1)
	patch := &conflictresolution.Message_PartialUpdate{
		Set_Fields: conflictresolution.Message_PartialUpdate_Set_Fields{
			Message: restli.StringPointer("partial updated message"),
		},
	}
	err := c.PartialUpdate(id, patch)
	require.NoError(t, err)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockPartialUpdate: func(ctx *restli.RequestContext, collectionId int64, entity *conflictresolution.Message_PartialUpdate) (err error) {
				require.Equal(t, id, collectionId)
				require.Equal(t, patch, entity)
				return nil
			},
		}
	}
}

func (o *Operation) SubCollectionOfCollectionGet(t *testing.T, c subcollection.Client) func(*testing.T) *subcollectiontest.MockResource {
	id := int64(1)
	subId := int64(100)
	expected := newMessage(subId, "sub collection message")
	res, err := c.Get(id, subId)
	require.NoError(t, err)
	require.Equal(t, expected, res)

	return func(t *testing.T) *subcollectiontest.MockResource {
		return &subcollectiontest.MockResource{
			MockGet: func(ctx *restli.RequestContext, collectionId int64, subcollectionId int64) (entity *conflictresolution.Message, err error) {
				require.Equal(t, id, collectionId)
				require.Equal(t, subId, subcollectionId)
				return expected, nil
			},
		}
	}
}

func (o *Operation) SubSimpleOfCollectionGet(t *testing.T, c colletionSubSimple.Client) func(*testing.T) *colletionSubSimpletest.MockResource {
	id := int64(1)
	expected := &conflictresolution.Message{Message: "sub simple message"}
	res, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t, expected, res, "Invalid response from server")

	return func(t *testing.T) *colletionSubSimpletest.MockResource {
		return &colletionSubSimpletest.MockResource{
			MockGet: func(ctx *restli.RequestContext, collectionId int64) (entity *conflictresolution.Message, err error) {
				require.Equal(t, id, collectionId)
				return expected, nil
			},
		}
	}
}

func newMessage(id int64, message string) *conflictresolution.Message {
	return &conflictresolution.Message{
		Id:      &id,
		Message: message,
	}
}

func (o *Operation) CollectionBatchDelete(t *testing.T, c Client) func(*testing.T) *MockResource {
	expectedKeys := []int64{1, 3}
	expected := &BatchResponse{
		Results: map[int64]*common.BatchEntityUpdateResponse{
			expectedKeys[0]: {
				Status: 204,
			},
			expectedKeys[1]: {
				Status: 404,
			},
		},
	}
	res, err := c.BatchDelete(expectedKeys)
	require.NoError(t, err)
	requiredBatchResponseEquals(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockBatchDelete: func(ctx *restli.RequestContext, keys []int64) (results *BatchResponse, err error) {
				require.Equal(t, expectedKeys, keys)
				return expected, nil
			},
		}
	}
}

func (o *Operation) CollectionBatchGet(t *testing.T, c Client) func(*testing.T) *MockResource {
	expectedKeys := []int64{1, 3}
	expected := &BatchEntities{
		Errors: map[int64]*common.ErrorResponse{},
		Results: map[int64]*conflictresolution.Message{
			expectedKeys[0]: {
				Id:      &expectedKeys[0],
				Message: "test message",
			},
			expectedKeys[1]: {
				Id:      restli.Int64Pointer(2),
				Message: "another message",
			},
		},
		Statuses: map[int64]int{},
	}
	actual, err := c.BatchGet(expectedKeys)
	require.NoError(t, err)
	requiredBatchResponseEquals(t, expected, actual)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockBatchGet: func(ctx *restli.RequestContext, keys []int64) (results *BatchEntities, err error) {
				require.Equal(t, expectedKeys, keys)
				return expected, nil
			},
		}
	}
}

func (o *Operation) CollectionBatchUpdate(t *testing.T, c Client) func(*testing.T) *MockResource {
	keys := []int64{1, 3}
	update := map[int64]*conflictresolution.Message{
		keys[0]: {
			Id:      &keys[0],
			Message: "updated message",
		},
		keys[1]: {
			Id:      &keys[1],
			Message: "inserted message",
		},
	}
	expected := &BatchResponse{
		Results: map[int64]*common.BatchEntityUpdateResponse{
			keys[0]: {
				Status: 204,
			},
			keys[1]: {
				Status: 201,
			},
		},
	}
	actual, err := c.BatchUpdate(update)
	require.NoError(t, err)
	requiredBatchResponseEquals(t, expected, actual)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockBatchUpdate: func(ctx *restli.RequestContext, entities map[int64]*conflictresolution.Message) (results *BatchResponse, err error) {
				require.Equal(t, update, entities)
				return expected, nil
			},
		}
	}
}

func (o *Operation) CollectionBatchUpdateErrors(t *testing.T, _ Client) {
	t.Log("It's impossible to produce the desired update for the same reason CollectionUpdate400 is skipped. Parsing " +
		"batch response errors is tested in SimpleComplexKeyBatchUpdateWithErrors")
}

func (o *Operation) CollectionBatchPartialUpdate(t *testing.T, c Client) func(*testing.T) *MockResource {
	keys := []int64{1, 3}
	partial := map[int64]*conflictresolution.Message_PartialUpdate{
		keys[0]: {
			Set_Fields: conflictresolution.Message_PartialUpdate_Set_Fields{
				Message: restli.StringPointer("partial updated message"),
			},
		},
		keys[1]: {
			Set_Fields: conflictresolution.Message_PartialUpdate_Set_Fields{
				Message: restli.StringPointer("another partial message"),
			},
		},
	}
	expected := &BatchResponse{
		Results: map[int64]*common.BatchEntityUpdateResponse{
			keys[0]: {
				Status: 204,
			},
			keys[1]: {
				Status: 204,
			},
		},
	}
	actual, err := c.BatchPartialUpdate(partial)
	require.NoError(t, err)
	requiredBatchResponseEquals(t, expected, actual)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockBatchPartialUpdate: func(ctx *restli.RequestContext, entities map[int64]*conflictresolution.Message_PartialUpdate) (results *BatchResponse, err error) {
				require.Equal(t, partial, entities)
				return expected, nil
			},
		}
	}
}

func (o *Operation) CollectionBatchCreate(t *testing.T, c Client) func(*testing.T) *MockResource {
	create := []*conflictresolution.Message{
		{
			Message: "test message",
		},
		{
			Message: "another message",
		},
	}
	expected := []*CreatedEntity{
		{
			Location: restli.StringPointer("/collection/1"),
			Id:       1,
			Status:   201,
		},
		{
			Location: restli.StringPointer("/collection/3"),
			Id:       3,
			Status:   201,
		},
	}

	res, err := c.BatchCreate(create)
	require.NoError(t, err)
	require.Equal(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockBatchCreate: func(ctx *restli.RequestContext, entities []*conflictresolution.Message) (createdEntities []*CreatedEntity, err error) {
				require.Equal(t, create, entities)
				return expected, nil
			},
		}
	}
}
