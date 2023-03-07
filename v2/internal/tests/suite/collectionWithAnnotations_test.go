package suite

import (
	"net/http"
	"testing"

	"github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated_extras/extras"
	. "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated_extras/extras/collectionWithAnnotations"
	. "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated_extras/extras/collectionWithAnnotations_test"
	"github.com/PapaCharlie/go-restli/v2/restli"
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
			Field3: restli.StringPointer("trois"),
		},
	}
	require.NoError(t, c.PartialUpdate(id, update))

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockPartialUpdate: func(ctx *restli.RequestContext, key extras.Temperature, entity *extras.MultiplePrimitiveFields_PartialUpdate) (err error) {
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
			MockCreate: func(ctx *restli.RequestContext, entity *extras.MultiplePrimitiveFields) (createdEntity *CreatedEntity, err error) {
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

// func (o *Operation) CollectionWithAnnotationsBatchCreate(t *testing.T, c Client) func(t *testing.T) *MockResource {
// 	expected := &CreatedEntity{
// 		Id:     1,
// 		Status: http.StatusCreated,
// 	}
// 	actual, err := c.BatchCreate(multiplePrimitiveFields)
// 	require.NoError(t, err)
// 	require.Equal(t, expected, actual)
//
// 	return func(t *testing.T) *MockResource {
// 		return &MockResource{
// 			MockCreate: func(ctx *restli.RequestContext, entity *extras.MultiplePrimitiveFields) (createdEntity *CreatedEntity, err error) {
// 				require.Equal(t, &extras.MultiplePrimitiveFields{
// 					Field1: "", // field1 is a read-only field so it's not supposed to be populated
// 					Field2: multiplePrimitiveFields.Field2,
// 					Field3: multiplePrimitiveFields.Field3,
// 				}, entity)
// 				return expected, nil
// 			},
// 		}
// 	}
// }

func (o *Operation) CollectionWithAnnotationsUpdate(t *testing.T, c Client) func(t *testing.T) *MockResource {
	one := extras.Temperature(1)
	err := c.Update(one, multiplePrimitiveFields)
	require.NoError(t, err)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockUpdate: func(ctx *restli.RequestContext, key extras.Temperature, entity *extras.MultiplePrimitiveFields) (err error) {
				require.Equal(t, one, key)
				require.Equal(t, &extras.MultiplePrimitiveFields{
					Field1: "", // field1 is a read-only field so it's not supposed to be populated
					Field2: "", // field2 is create-only, so it's not supposed to be populated here either
					Field3: multiplePrimitiveFields.Field3,
				}, entity)
				return nil
			},
		}
	}
}
