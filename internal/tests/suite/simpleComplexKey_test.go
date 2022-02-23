package suite

import (
	"testing"

	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/simpleComplexKey"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/PapaCharlie/go-restli/protocol/stdstructs"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) SimpleComplexKeyBatchCreate(t *testing.T, c Client) {
	entities := []*extras.SinglePrimitiveField{
		{
			String: "1",
		},
		{
			String: "string:with:colons",
		},
	}
	res, err := c.BatchCreate(entities)
	require.NoError(t, err)
	expected := []*protocol.CreatedAndReturnedEntity[*SimpleComplexKey_ComplexKey, *extras.SinglePrimitiveField]{
		{
			CreatedEntity: protocol.CreatedEntity[*SimpleComplexKey_ComplexKey]{
				Id:     &SimpleComplexKey_ComplexKey{SinglePrimitiveField: *entities[0]},
				Status: 201,
			},
			Entity: entities[0],
		},
		{
			CreatedEntity: protocol.CreatedEntity[*SimpleComplexKey_ComplexKey]{
				Id:     &SimpleComplexKey_ComplexKey{SinglePrimitiveField: *entities[1]},
				Status: 202,
			},
			Entity: entities[1],
		},
	}
	require.Equal(t, expected, res)
}

func (s *TestServer) SimpleComplexKeyBatchGet(t *testing.T, c Client) {
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
	res, err := c.BatchGet(keys)
	require.NoError(t, err)
	expected := make(map[*SimpleComplexKey_ComplexKey]*extras.SinglePrimitiveField)
	for _, k := range keys {
		expected[k] = &k.SinglePrimitiveField
	}
	require.Equal(t, expected, res)
}

func (s *TestServer) SimpleComplexKeyGet(t *testing.T, c Client) {
	expected := &SimpleComplexKey_ComplexKey{
		SinglePrimitiveField: extras.SinglePrimitiveField{String: "string:with:colons"},
	}
	actual, err := c.Get(expected)
	require.NoError(t, err)
	require.Equal(t, &expected.SinglePrimitiveField, actual)
}

func (s *TestServer) SimpleComplexKeyBatchUpdateWithErrors(t *testing.T, c Client) {
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

	actual, err := c.BatchUpdate(map[*SimpleComplexKey_ComplexKey]*extras.SinglePrimitiveField{
		keys[0]: &keys[0].SinglePrimitiveField,
		keys[1]: &keys[1].SinglePrimitiveField,
		keys[2]: &keys[2].SinglePrimitiveField,
	})
	require.Equal(t, protocol.BatchRequestResponseError[*SimpleComplexKey_ComplexKey]{
		keys[0]: &stdstructs.ErrorResponse{
			Status:         protocol.Int32Pointer(400),
			Message:        protocol.StringPointer("message"),
			ExceptionClass: protocol.StringPointer("com.linkedin.restli.server.RestLiServiceException"),
			StackTrace:     protocol.StringPointer("trace"),
		},
		keys[2]: &stdstructs.ErrorResponse{
			Status:         protocol.Int32Pointer(500),
			Message:        nil,
			ExceptionClass: protocol.StringPointer("com.linkedin.restli.server.RestLiServiceException"),
			StackTrace:     protocol.StringPointer("trace"),
		},
	}, err)

	expected := map[*SimpleComplexKey_ComplexKey]*protocol.BatchEntityUpdateResponse{
		keys[1]: {Status: 204},
	}
	require.Equal(t, expected, actual)
}

func (s *TestServer) SimpleComplexKeyBatchPartialUpdate(t *testing.T, c Client) {
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

	actual, err := c.BatchPartialUpdate(map[*SimpleComplexKey_ComplexKey]*extras.SinglePrimitiveField_PartialUpdate{
		keys[0]: {
			Set_Fields: extras.SinglePrimitiveField_PartialUpdate_Set_Fields{
				String: protocol.StringPointer("partial updated message"),
			},
		},
		keys[1]: {
			Set_Fields: extras.SinglePrimitiveField_PartialUpdate_Set_Fields{
				String: protocol.StringPointer("another partial message"),
			},
		},
	})
	require.NoError(t, err)

	expected := map[*SimpleComplexKey_ComplexKey]*protocol.BatchEntityUpdateResponse{
		keys[0]: {Status: 204},
		keys[1]: {Status: 205},
	}
	require.Equal(t, expected, actual)
}
