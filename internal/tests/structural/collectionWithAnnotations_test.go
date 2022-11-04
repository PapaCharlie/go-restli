package structural

import (
	"context"

	"github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated_extras/extras"
	collectionwithannotations "github.com/PapaCharlie/go-restli/v2/internal/tests/testdata/generated_extras/extras/collectionWithAnnotations"
)

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR extras.collectionWithAnnotations
var _ = collectionwithannotations.Client(new(collectionWithAnnotationsClient))

type collectionWithAnnotationsClient struct{}

func (c *collectionWithAnnotationsClient) Create(*extras.MultiplePrimitiveFields) (*collectionwithannotations.CreatedEntity, error) {
	panic(nil)
}

func (c *collectionWithAnnotationsClient) CreateWithContext(context.Context, *extras.MultiplePrimitiveFields) (*collectionwithannotations.CreatedEntity, error) {
	panic(nil)
}

func (c *collectionWithAnnotationsClient) Update(extras.Temperature, *extras.MultiplePrimitiveFields) error {
	panic(nil)
}

func (c *collectionWithAnnotationsClient) UpdateWithContext(context.Context, extras.Temperature, *extras.MultiplePrimitiveFields) error {
	panic(nil)
}

func (c *collectionWithAnnotationsClient) PartialUpdate(extras.Temperature, *extras.MultiplePrimitiveFields_PartialUpdate) error {
	panic(nil)
}

func (c *collectionWithAnnotationsClient) PartialUpdateWithContext(context.Context, extras.Temperature, *extras.MultiplePrimitiveFields_PartialUpdate) error {
	panic(nil)
}
