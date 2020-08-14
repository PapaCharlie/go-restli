package tests

import (
	"context"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/complexkey"
	"github.com/PapaCharlie/go-restli/protocol"
)

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR testsuite.typerefs.collectionTyperef
var _ = complexkey.Client(new(complexKeyClient))

type complexKeyClient int

func (c *complexKeyClient) Create(*conflictresolution.LargeRecord) (protocol.RawComplexKey, error) {
	panic(nil)
}

func (c *complexKeyClient) CreateWithContext(context.Context, *conflictresolution.LargeRecord) (protocol.RawComplexKey, error) {
	panic(nil)
}

func (c *complexKeyClient) Get(*complexkey.Complexkey_ComplexKey) (*conflictresolution.LargeRecord, error) {
	panic(nil)
}

func (c *complexKeyClient) GetWithContext(context.Context, *complexkey.Complexkey_ComplexKey) (*conflictresolution.LargeRecord, error) {
	panic(nil)
}

func (c *complexKeyClient) Update(*complexkey.Complexkey_ComplexKey, *conflictresolution.LargeRecord) error {
	panic(nil)
}

func (c *complexKeyClient) UpdateWithContext(context.Context, *complexkey.Complexkey_ComplexKey, *conflictresolution.LargeRecord) error {
	panic(nil)
}

func (c *complexKeyClient) PartialUpdate(*complexkey.Complexkey_ComplexKey, *conflictresolution.LargeRecord_PartialUpdate) error {
	panic(nil)
}

func (c *complexKeyClient) PartialUpdateWithContext(context.Context, *complexkey.Complexkey_ComplexKey, *conflictresolution.LargeRecord_PartialUpdate) error {
	panic(nil)
}

func (c *complexKeyClient) Delete(*complexkey.Complexkey_ComplexKey) error {
	panic(nil)
}

func (c *complexKeyClient) DeleteWithContext(context.Context, *complexkey.Complexkey_ComplexKey) error {
	panic(nil)
}
