package structural

import (
	"context"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/complexkey"
)

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR testsuite.typerefs.collectionTyperef
var _ = complexkey.Client(new(complexKeyClient))

type complexKeyClient struct{}

func (c *complexKeyClient) Create(*conflictresolution.LargeRecord) (*complexkey.Complexkey_ComplexKey, error) {
	panic(nil)
}

func (c *complexKeyClient) CreateWithContext(context.Context, *conflictresolution.LargeRecord) (*complexkey.Complexkey_ComplexKey, error) {
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

func (c *complexKeyClient) BatchGet([]*complexkey.Complexkey_ComplexKey) (map[*complexkey.Complexkey_ComplexKey]*conflictresolution.LargeRecord, error) {
	panic(nil)
}

func (c *complexKeyClient) BatchGetWithContext(context.Context, []*complexkey.Complexkey_ComplexKey) (map[*complexkey.Complexkey_ComplexKey]*conflictresolution.LargeRecord, error) {
	panic(nil)
}
