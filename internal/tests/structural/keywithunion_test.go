package structural

import (
	"context"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/keywithunion/keywithunion"
)

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR testsuite.keywithunion.keywithunion
var _ = keywithunion.Client(new(keyWithUnionClient))

type keyWithUnionClient struct{}

func (k *keyWithUnionClient) Get(*keywithunion.Keywithunion_ComplexKey) (*conflictresolution.LargeRecord, error) {
	panic(nil)
}

func (k *keyWithUnionClient) GetWithContext(context.Context, *keywithunion.Keywithunion_ComplexKey) (*conflictresolution.LargeRecord, error) {
	panic(nil)
}
