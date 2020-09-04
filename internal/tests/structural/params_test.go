package structural

import (
	"context"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/params"
)

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR testsuite.params
var _ = params.Client(new(paramsClient))

type paramsClient struct{}

func (p *paramsClient) Get(int64, *params.GetParams) (*conflictresolution.Message, error) {
	panic(nil)
}

func (p *paramsClient) GetWithContext(context.Context, int64, *params.GetParams) (*conflictresolution.Message, error) {
	panic(nil)
}
