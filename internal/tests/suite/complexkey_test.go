package suite

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/complexkey"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/complexkey_test"
	"github.com/PapaCharlie/go-restli/restli"
	"github.com/PapaCharlie/go-restli/restlidata"
	"github.com/stretchr/testify/require"
)

func (o *Operation) ComplexkeyGet(t *testing.T, c Client) func(*testing.T) *MockResource {
	id := &Complexkey_ComplexKey{
		Params: newKeyParams("param1", 5),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "one",
			Part2: 2,
			Part3: conflictresolution.Fruits_APPLE,
		},
	}
	expected := &conflictresolution.LargeRecord{
		Message: conflictresolution.Message{Message: "test message"},
		Key:     id.ComplexKey,
	}
	res, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockGet: func(ctx *restli.RequestContext, complexkeyId *Complexkey_ComplexKey) (entity *conflictresolution.LargeRecord, err error) {
				require.Equal(t, id, complexkeyId)
				return expected, nil
			},
		}
	}
}

func (o *Operation) ComplexkeyUpdate(t *testing.T, c Client) func(*testing.T) *MockResource {
	key := conflictresolution.ComplexKey{
		Part1: "one",
		Part2: 2,
		Part3: conflictresolution.Fruits_APPLE,
	}
	id := &Complexkey_ComplexKey{
		Params:     newKeyParams("param1", 5),
		ComplexKey: key,
	}
	record := &conflictresolution.LargeRecord{
		Key: key,
		Message: conflictresolution.Message{
			Message: "updated message",
		},
	}
	err := c.Update(id, record)
	require.NoError(t, err)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockUpdate: func(ctx *restli.RequestContext, complexkeyId *Complexkey_ComplexKey, entity *conflictresolution.LargeRecord) (err error) {
				require.Equal(t, id, complexkeyId)
				require.Equal(t, record, entity)
				return nil
			},
		}
	}
}

func (o *Operation) ComplexkeyDelete(t *testing.T, c Client) func(*testing.T) *MockResource {
	id := &Complexkey_ComplexKey{
		Params: newKeyParams("param1", 5),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "one",
			Part2: 2,
			Part3: conflictresolution.Fruits_APPLE,
		},
	}
	err := c.Delete(id)
	require.NoError(t, err)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockDelete: func(ctx *restli.RequestContext, complexkeyId *Complexkey_ComplexKey) (err error) {
				require.Equal(t, id, complexkeyId)
				return nil
			},
		}
	}
}

func (o *Operation) ComplexkeyCreate(t *testing.T, c Client) func(*testing.T) *MockResource {
	expectedKey := conflictresolution.ComplexKey{
		Part1: "one",
		Part2: 2,
		Part3: conflictresolution.Fruits_APPLE,
	}
	create := &conflictresolution.LargeRecord{
		Key: expectedKey,
		Message: conflictresolution.Message{
			Message: "test message",
		},
	}
	_, err := c.Create(create)
	require.IsType(t, new(restli.CreateResponseHasNoEntityHeaderError), err)
	// TODO: Merge https://github.com/linkedin/rest.li-test-suite/pull/6 and actually test the contents of the key
	// require.Equal(t, expectedKey, actualKey.ComplexKey)
	return func(t *testing.T) *MockResource {
		// TODO: ^ see above
		t.SkipNow()
		return nil
		// return &MockResource{
		// 	MockCreate: func(ctx *restli.RequestContext, entity *conflictresolution.LargeRecord) (createdEntity *restli.CreatedEntity[*Complexkey_ComplexKey], err error) {
		// 		require.Equal(t, create, entity)
		// 		return &restli.CreatedEntity[*Complexkey_ComplexKey]{Id: &Complexkey_ComplexKey{ComplexKey: expectedKey}}, nil
		// 	},
		// }
	}
}

func (o *Operation) ComplexkeyPartialUpdate(t *testing.T, c Client) func(*testing.T) *MockResource {
	id := &Complexkey_ComplexKey{
		Params: newKeyParams("param1", 5),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "one",
			Part2: 2,
			Part3: conflictresolution.Fruits_APPLE,
		},
	}
	keyPatch := &conflictresolution.ComplexKey_PartialUpdate{
		Set_Fields: conflictresolution.ComplexKey_PartialUpdate_Set_Fields{
			Part1: restli.StringPointer("newpart1"),
		},
	}
	update := &conflictresolution.LargeRecord_PartialUpdate{Key: keyPatch}
	err := c.PartialUpdate(id, update)
	require.NoError(t, err)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockPartialUpdate: func(ctx *restli.RequestContext, complexkeyId *Complexkey_ComplexKey, entity *conflictresolution.LargeRecord_PartialUpdate) (err error) {
				require.Equal(t, id, complexkeyId)
				require.Equal(t, entity, update)
				return nil
			},
		}
	}
}

func (o *Operation) ComplexkeyBatchDelete(t *testing.T, c Client) func(*testing.T) *MockResource {
	k1 := &Complexkey_ComplexKey{
		Params: newKeyParams("param1", 5),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "one",
			Part2: 2,
			Part3: conflictresolution.Fruits_APPLE,
		},
	}
	k2 := &Complexkey_ComplexKey{
		Params: newKeyParams("param2", 11),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "two",
			Part2: 7,
			Part3: conflictresolution.Fruits_ORANGE,
		},
	}
	delete := []*Complexkey_ComplexKey{k1, k2}
	res, err := c.BatchDelete(delete)
	require.NoError(t, err)

	expected := &BatchResponse{
		Results: map[*Complexkey_ComplexKey]*restlidata.BatchEntityUpdateResponse{
			k1: {
				Status: 204,
			},
			k2: {
				Status: 204,
			},
		},
	}
	requiredBatchResponseEquals(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockBatchDelete: func(ctx *restli.RequestContext, keys []*Complexkey_ComplexKey) (results *BatchResponse, err error) {
				require.Equal(t, delete, keys)
				results = &BatchResponse{
					Results: map[*Complexkey_ComplexKey]*restlidata.BatchEntityUpdateResponse{},
				}
				for _, k := range keys {
					// Explicitly don't set the status to check if the default status kicks in
					results.Results[&Complexkey_ComplexKey{ComplexKey: k.ComplexKey}] = new(restlidata.BatchEntityUpdateResponse)
				}
				return results, nil
			},
		}
	}
}

func (o *Operation) ComplexkeyBatchGet(t *testing.T, c Client) func(*testing.T) *MockResource {
	k1 := &Complexkey_ComplexKey{
		Params: newKeyParams("param1", 5),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "one",
			Part2: 2,
			Part3: conflictresolution.Fruits_APPLE,
		},
	}
	k2 := &Complexkey_ComplexKey{
		Params: newKeyParams("param2", 11),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "two",
			Part2: 7,
			Part3: conflictresolution.Fruits_ORANGE,
		},
	}
	get := []*Complexkey_ComplexKey{k1, k2}
	res, err := c.BatchGet(get)
	require.NoError(t, err)

	expected := &BatchEntities{
		Results: map[*Complexkey_ComplexKey]*conflictresolution.LargeRecord{
			k1: {
				Key: k1.ComplexKey,
				Message: conflictresolution.Message{
					Message: "test message",
				},
			},
			k2: {
				Key: k2.ComplexKey,
				Message: conflictresolution.Message{
					Message: "test message",
				},
			},
		},
	}
	requiredBatchResponseEquals(t, expected, res)
	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockBatchGet: func(ctx *restli.RequestContext, keys []*Complexkey_ComplexKey) (results *BatchEntities, err error) {
				require.Equal(t, get, keys)
				results = new(BatchEntities)
				for k, v := range expected.Results {
					results.AddResult(&Complexkey_ComplexKey{ComplexKey: k.ComplexKey}, v)
				}
				return results, nil
			},
		}
	}
}

const specialChars = `!*'();:@&=+$,/?#[].~`

var specialCharsKey = &Complexkey_ComplexKey{
	Params: newKeyParams("param"+specialChars, 5),
	ComplexKey: conflictresolution.ComplexKey{
		Part1: "key" + specialChars,
		Part2: 2,
		Part3: conflictresolution.Fruits_APPLE,
	},
}

func (o *Operation) ComplexkeyGetWithSpecialChars(t *testing.T, c Client) func(*testing.T) *MockResource {
	expected := &conflictresolution.LargeRecord{
		Key: specialCharsKey.ComplexKey,
		Message: conflictresolution.Message{
			Message: "test message",
		},
	}
	res, err := c.Get(specialCharsKey)
	require.NoError(t, err)
	require.Equal(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockGet: func(ctx *restli.RequestContext, complexkeyId *Complexkey_ComplexKey) (entity *conflictresolution.LargeRecord, err error) {
				require.Equal(t, specialCharsKey, complexkeyId)
				return expected, nil
			},
		}
	}
}

func (o *Operation) ComplexkeyBatchGetWithSpecialChars(t *testing.T, c Client) func(*testing.T) *MockResource {
	k := &Complexkey_ComplexKey{
		Params: newKeyParams("param2", 11),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "two",
			Part2: 7,
			Part3: conflictresolution.Fruits_ORANGE,
		},
	}
	res, err := c.BatchGet([]*Complexkey_ComplexKey{specialCharsKey, k})
	expected := &BatchEntities{
		Results: map[*Complexkey_ComplexKey]*conflictresolution.LargeRecord{
			specialCharsKey: {
				Key: specialCharsKey.ComplexKey,
				Message: conflictresolution.Message{
					Message: "test message",
				},
			},
			k: {
				Key: k.ComplexKey,
				Message: conflictresolution.Message{
					Message: "test message",
				},
			},
		},
	}
	require.NoError(t, err)
	requiredBatchResponseEquals(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockBatchGet: func(ctx *restli.RequestContext, keys []*Complexkey_ComplexKey) (results *BatchEntities, err error) {
				require.Equal(t, []*Complexkey_ComplexKey{specialCharsKey, k}, keys)
				results = new(BatchEntities)
				for k, v := range expected.Results {
					results.AddResult(&Complexkey_ComplexKey{ComplexKey: k.ComplexKey}, v)
				}
				return results, nil
			},
		}
	}
}

func (o *Operation) ComplexkeyBatchUpdate(t *testing.T, c Client) func(*testing.T) *MockResource {
	k1 := &Complexkey_ComplexKey{
		Params: newKeyParams("param1", 5),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "one",
			Part2: 2,
			Part3: conflictresolution.Fruits_APPLE,
		},
	}
	k2 := &Complexkey_ComplexKey{
		Params: newKeyParams("param2", 11),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "two",
			Part2: 7,
			Part3: conflictresolution.Fruits_ORANGE,
		},
	}
	updates := map[*Complexkey_ComplexKey]*conflictresolution.LargeRecord{
		k1: {
			Key: k1.ComplexKey,
			Message: conflictresolution.Message{
				Message: "updated message",
			},
		},
		k2: {
			Key: k1.ComplexKey,
			Message: conflictresolution.Message{
				Message: "another updated message",
			},
		},
	}
	res, err := c.BatchUpdate(updates)
	require.NoError(t, err)
	requiredBatchResponseEquals(t, &BatchResponse{
		Results: map[*Complexkey_ComplexKey]*restlidata.BatchEntityUpdateResponse{
			k1: {
				Status: 204,
			},
			k2: {
				Status: 204,
			},
		},
	}, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockBatchUpdate: func(ctx *restli.RequestContext, entities map[*Complexkey_ComplexKey]*conflictresolution.LargeRecord) (results *BatchResponse, err error) {
				requireComplexKeyMapEquals(t, updates, entities)
				return &BatchResponse{
					Results: map[*Complexkey_ComplexKey]*restlidata.BatchEntityUpdateResponse{
						&Complexkey_ComplexKey{ComplexKey: k1.ComplexKey}: {},
						&Complexkey_ComplexKey{ComplexKey: k2.ComplexKey}: {},
					},
				}, nil
			},
		}
	}
}

func (o *Operation) ComplexkeyBatchCreate(t *testing.T, c Client) func(*testing.T) *MockResource {
	create := []*conflictresolution.LargeRecord{
		{
			Key: conflictresolution.ComplexKey{
				Part1: "one",
				Part2: 2,
				Part3: conflictresolution.Fruits_APPLE,
			},
			Message: conflictresolution.Message{
				Message: "test message",
			},
		},
		{
			Key: conflictresolution.ComplexKey{
				Part1: "two",
				Part2: 7,
				Part3: conflictresolution.Fruits_ORANGE,
			},
			Message: conflictresolution.Message{
				Message: "another message",
			},
		},
	}
	res, err := c.BatchCreate(create)
	require.NoError(t, err)
	require.Equal(t, []*CreatedEntity{
		{
			Status: 201,
		},
		{
			Status: 201,
		},
	}, res)

	return func(t *testing.T) *MockResource {
		deliberateSkip(t, "Cannot return empty key from batch create (merge https://github.com/linkedin/rest.li-test-suite/pull/6 to fix)")
		return nil
	}
}
