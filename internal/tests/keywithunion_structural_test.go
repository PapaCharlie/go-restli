package tests

import (
	"context"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/keywithunion/keywithunion"
)

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR testsuite.keywithunion.keywithunion
var _ = keywithunion.Client(new(keyWithUnionClient))

type keyWithUnionClient int

func (k *keyWithUnionClient) Get(*keywithunion.Keywithunion_ComplexKey) (*conflictresolution.LargeRecord, error) {
	panic(nil)
}

func (k *keyWithUnionClient) GetWithContext(context.Context, *keywithunion.Keywithunion_ComplexKey) (*conflictresolution.LargeRecord, error) {
	panic(nil)
}
