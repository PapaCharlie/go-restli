package suite

import (
	"testing"

	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithAnnotations"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/stretchr/testify/require"
)

var multiplePrimitiveFields = &extras.MultiplePrimitiveFields{
	Field1: "one",
	Field2: "two",
	Field3: "three",
}

func (s *TestServer) CollectionWithAnnotationsPartialUpdate(t *testing.T, c Client) {
	update := &extras.MultiplePrimitiveFields_PartialUpdate{
		Set_Fields: extras.MultiplePrimitiveFields_PartialUpdate_Set_Fields{
			Field3: protocol.StringPointer("trois"),
		},
	}
	require.NoError(t, c.PartialUpdate(1, update))
}

func (s *TestServer) CollectionWithAnnotationsCreate(t *testing.T, c Client) {
	_, err := c.Create(multiplePrimitiveFields)
	require.NoError(t, err)
}

func (s *TestServer) CollectionWithAnnotationsUpdate(t *testing.T, c Client) {
	update := *multiplePrimitiveFields
	update.Field3 = "trois"
	require.NoError(t, c.Update(1, &update))
}
