package tests

import (
	"context"

	collectionwithtyperefkey "github.com/PapaCharlie/go-restli/internal/tests/generated/extras/collectionWithTyperefKey"
	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite"
)

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR extras.collectionWithTyperefKey
var _ = collectionwithtyperefkey.Client(new(collectionWithTyperefKeyClient))

type collectionWithTyperefKeyClient struct{}

func (c *collectionWithTyperefKeyClient) Create(*testsuite.PrimitiveField) (testsuite.Time, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) CreateWithContext(context.Context, *testsuite.PrimitiveField) (testsuite.Time, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) BatchGet([]testsuite.Time, *collectionwithtyperefkey.BatchGetParams) (map[testsuite.Time]*testsuite.PrimitiveField, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) BatchGetWithContext(context.Context, []testsuite.Time, *collectionwithtyperefkey.BatchGetParams) (map[testsuite.Time]*testsuite.PrimitiveField, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) Get(testsuite.Time) (*testsuite.PrimitiveField, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) GetWithContext(context.Context, testsuite.Time) (*testsuite.PrimitiveField, error) {
	panic(nil)
}
