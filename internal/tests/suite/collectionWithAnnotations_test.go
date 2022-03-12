package suite

import (
	"net/http"
	"testing"

	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithAnnotations"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithAnnotations_test"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/stretchr/testify/require"
)

var multiplePrimitiveFields = &extras.MultiplePrimitiveFields{
	Field1: "one",
	Field2: "two",
	Field3: "three",
}

func (o *Operation) CollectionWithAnnotationsPartialUpdate(t *testing.T, c Client) func(t *testing.T) *MockResource {
	id := extras.Temperature(1)
	update := &extras.MultiplePrimitiveFields_PartialUpdate{
		Set_Fields: extras.MultiplePrimitiveFields_PartialUpdate_Set_Fields{
			Field3: protocol.StringPointer("trois"),
		},
	}
	require.NoError(t, c.PartialUpdate(id, update))

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockPartialUpdate: func(ctx *protocol.RequestContext, key extras.Temperature, entity *extras.MultiplePrimitiveFields_PartialUpdate) (err error) {
				require.Equal(t, id, key)
				require.Equal(t, update, entity)
				return nil
			},
		}
	}
}

func (o *Operation) CollectionWithAnnotationsCreate(t *testing.T, c Client) func(t *testing.T) *MockResource {
	expected := &CreatedEntity{
		Id:     1,
		Status: http.StatusCreated,
	}
	actual, err := c.Create(multiplePrimitiveFields)
	require.NoError(t, err)
	require.Equal(t, expected, actual)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockCreate: func(ctx *protocol.RequestContext, entity *extras.MultiplePrimitiveFields) (createdEntity *CreatedEntity, err error) {
				require.Equal(t, &extras.MultiplePrimitiveFields{
					Field1: "", // field1 is a read-only field so it's not supposed to be populated
					Field2: multiplePrimitiveFields.Field2,
					Field3: multiplePrimitiveFields.Field3,
				}, entity)
				return expected, nil
			},
		}
	}
}

func (o *Operation) CollectionWithAnnotationsUpdate(t *testing.T, _ Client) func(t *testing.T) *MockResource {
	deliberateSkip(t, "Skipped because testing for field exclusion is done elsewhere")
	return func(t *testing.T) *MockResource {
		deliberateSkip(t, "Skipped because go-restli's server automatically rejects incoming update requests with excluded fields")
		return nil
	}
}
