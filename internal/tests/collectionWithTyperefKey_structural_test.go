package tests

import (
	"context"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/generated/conflictResolution"
	collectionwithtyperefkey "github.com/PapaCharlie/go-restli/internal/tests/generated/extras/collectionWithTyperefKey"
	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite"
)

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR extras.collectionWithTyperefKey
var _ = collectionwithtyperefkey.Client(new(collectionWithTyperefKeyClient))

type collectionWithTyperefKeyClient int

func (c *collectionWithTyperefKeyClient) Create(*conflictresolution.Message) (testsuite.Time, error) {
	panic(nil)
}

func (c *collectionWithTyperefKeyClient) CreateWithContext(context.Context, *conflictresolution.Message) (testsuite.Time, error) {
	panic(nil)
}
