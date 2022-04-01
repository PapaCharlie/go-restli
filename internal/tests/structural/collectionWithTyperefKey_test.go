package structural

import (
	"context"

	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	collectionwithtyperefkey "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithTyperefKey"
)

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR extras.collectionWithTyperefKey
var _ = collectionwithtyperefkey.Client(new(collectionWithTyperefKeyClient))

type collectionWithTyperefKeyClient struct{}

func (c *collectionWithTyperefKeyClient) BatchCreate([]*extras.SinglePrimitiveField, *collectionwithtyperefkey.BatchCreateParams) ([]*collectionwithtyperefkey.CreatedEntity, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) BatchCreateWithContext(context.Context, []*extras.SinglePrimitiveField, *collectionwithtyperefkey.BatchCreateParams) ([]*collectionwithtyperefkey.CreatedEntity, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) BatchGet([]extras.Temperature, *collectionwithtyperefkey.BatchGetParams) (*collectionwithtyperefkey.BatchEntities, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) BatchGetWithContext(context.Context, []extras.Temperature, *collectionwithtyperefkey.BatchGetParams) (*collectionwithtyperefkey.BatchEntities, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) Create(*extras.SinglePrimitiveField) (*collectionwithtyperefkey.CreatedEntity, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) CreateWithContext(context.Context, *extras.SinglePrimitiveField) (*collectionwithtyperefkey.CreatedEntity, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) Get(extras.Temperature) (*extras.SinglePrimitiveField, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) GetWithContext(context.Context, extras.Temperature) (*extras.SinglePrimitiveField, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) FindBySearch(*collectionwithtyperefkey.FindBySearchParams) (*collectionwithtyperefkey.Elements, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) FindBySearchWithContext(context.Context, *collectionwithtyperefkey.FindBySearchParams) (*collectionwithtyperefkey.Elements, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) OnEntityAction(extras.Temperature, *collectionwithtyperefkey.OnEntityActionParams) error {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) OnEntityActionWithContext(context.Context, extras.Temperature, *collectionwithtyperefkey.OnEntityActionParams) error {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) GetAll(*collectionwithtyperefkey.GetAllParams) (*collectionwithtyperefkey.Elements, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) GetAllWithContext(context.Context, *collectionwithtyperefkey.GetAllParams) (*collectionwithtyperefkey.Elements, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) FindByNoParams() (*collectionwithtyperefkey.Elements, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) FindByNoParamsWithContext(context.Context) (*collectionwithtyperefkey.Elements, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) FindByNoParamsWithPaging(*collectionwithtyperefkey.FindByNoParamsWithPagingParams) (*collectionwithtyperefkey.Elements, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) FindByNoParamsWithPagingWithContext(context.Context, *collectionwithtyperefkey.FindByNoParamsWithPagingParams) (*collectionwithtyperefkey.Elements, error) {
	panic(nil)
}
