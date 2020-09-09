package structural

import (
	"context"

	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	collectionwithtyperefkey "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithTyperefKey"
)

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR extras.collectionWithTyperefKey
var _ = collectionwithtyperefkey.Client(new(collectionWithTyperefKeyClient))

type collectionWithTyperefKeyClient struct{}

func (c *collectionWithTyperefKeyClient) Create(*extras.SinglePrimitiveField) (extras.Temperature, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) CreateWithContext(context.Context, *extras.SinglePrimitiveField) (extras.Temperature, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) BatchGet([]extras.Temperature, *collectionwithtyperefkey.BatchGetParams) (map[extras.Temperature]*extras.SinglePrimitiveField, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) BatchGetWithContext(context.Context, []extras.Temperature, *collectionwithtyperefkey.BatchGetParams) (map[extras.Temperature]*extras.SinglePrimitiveField, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) Get(extras.Temperature) (*extras.SinglePrimitiveField, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) GetWithContext(context.Context, extras.Temperature) (*extras.SinglePrimitiveField, error) {
	panic(nil)
}
