package tests

import (
	"context"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/simple"
)

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR testsuite.simple
var _ = simple.Client(new(simpleClient))

type simpleClient struct{}

func (s *simpleClient) Get() (*conflictresolution.Message, error) {
	panic(nil)
}

func (s *simpleClient) GetWithContext(context.Context) (*conflictresolution.Message, error) {
	panic(nil)
}

func (s *simpleClient) Update(*conflictresolution.Message) error {
	panic(nil)
}

func (s *simpleClient) UpdateWithContext(context.Context, *conflictresolution.Message) error {
	panic(nil)
}

func (s *simpleClient) Delete() error {
	panic(nil)
}

func (s *simpleClient) DeleteWithContext(context.Context) error {
	panic(nil)
}
