package structural

import (
	"context"

	"github.com/PapaCharlie/go-restli/internal/tests/native/testsuite"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	collectionwithtyperefkey "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithTyperefKey"
	"github.com/PapaCharlie/go-restli/protocol"
)

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR extras.collectionWithTyperefKey
var _ = collectionwithtyperefkey.Client(new(collectionWithTyperefKeyClient))

type collectionWithTyperefKeyClient struct{}

func (c *collectionWithTyperefKeyClient) Create(*extras.SinglePrimitiveField) (*protocol.CreatedEntity[testsuite.Temperature], error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) CreateWithContext(context.Context, *extras.SinglePrimitiveField) (*protocol.CreatedEntity[testsuite.Temperature], error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) BatchGet([]testsuite.Temperature, *collectionwithtyperefkey.BatchGetParams) (map[testsuite.Temperature]*extras.SinglePrimitiveField, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) BatchGetWithContext(context.Context, []testsuite.Temperature, *collectionwithtyperefkey.BatchGetParams) (map[testsuite.Temperature]*extras.SinglePrimitiveField, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) Get(testsuite.Temperature) (*extras.SinglePrimitiveField, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) GetWithContext(context.Context, testsuite.Temperature) (*extras.SinglePrimitiveField, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) FindBySearch(*collectionwithtyperefkey.FindBySearchParams) (*protocol.FinderResults[*extras.SinglePrimitiveField], error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) FindBySearchWithContext(context.Context, *collectionwithtyperefkey.FindBySearchParams) (*protocol.FinderResults[*extras.SinglePrimitiveField], error) {
	panic(nil)
}
