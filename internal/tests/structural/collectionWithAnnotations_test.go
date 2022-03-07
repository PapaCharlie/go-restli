package structural

import (
	"context"

	"github.com/PapaCharlie/go-restli/internal/tests/native/testsuite"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	collectionwithannotations "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithAnnotations"
	"github.com/PapaCharlie/go-restli/protocol"
)

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR extras.collectionWithAnnotations
var _ = collectionwithannotations.Client(new(collectionWithAnnotationsClient))

type collectionWithAnnotationsClient struct{}

func (c *collectionWithAnnotationsClient) Create(*extras.MultiplePrimitiveFields) (*protocol.CreatedEntity[testsuite.Temperature], error) {
	panic(nil)
}

func (c *collectionWithAnnotationsClient) CreateWithContext(context.Context, *extras.MultiplePrimitiveFields) (*protocol.CreatedEntity[testsuite.Temperature], error) {
	panic(nil)
}

func (c *collectionWithAnnotationsClient) Update(testsuite.Temperature, *extras.MultiplePrimitiveFields) error {
	panic(nil)
}

func (c *collectionWithAnnotationsClient) UpdateWithContext(context.Context, testsuite.Temperature, *extras.MultiplePrimitiveFields) error {
	panic(nil)
}

func (c *collectionWithAnnotationsClient) PartialUpdate(testsuite.Temperature, *extras.MultiplePrimitiveFields_PartialUpdate) error {
	panic(nil)
}

func (c *collectionWithAnnotationsClient) PartialUpdateWithContext(context.Context, testsuite.Temperature, *extras.MultiplePrimitiveFields_PartialUpdate) error {
	panic(nil)
}
