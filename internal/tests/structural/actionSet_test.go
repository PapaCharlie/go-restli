package structural

import (
	"context"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite"
	actionset "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/actionSet"
)

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR testsuite.actionSet
var _ = actionset.Client(new(actionSetClient))

type actionSetClient struct{}

func (a *actionSetClient) EchoAction(*actionset.EchoActionParams) (string, error) {
	panic(nil)
}

func (a *actionSetClient) EchoActionWithContext(context.Context, *actionset.EchoActionParams) (string, error) {
	panic(nil)
}

func (a *actionSetClient) EchoComplexTypesUnionAction(*actionset.EchoComplexTypesUnionActionParams) (*testsuite.UnionOfComplexTypes, error) {
	panic(nil)
}

func (a *actionSetClient) EchoComplexTypesUnionActionWithContext(context.Context, *actionset.EchoComplexTypesUnionActionParams) (*testsuite.UnionOfComplexTypes, error) {
	panic(nil)
}

func (a *actionSetClient) EchoMessageAction(*actionset.EchoMessageActionParams) (*conflictresolution.Message, error) {
	panic(nil)
}

func (a *actionSetClient) EchoMessageActionWithContext(context.Context, *actionset.EchoMessageActionParams) (*conflictresolution.Message, error) {
	panic(nil)
}

func (a *actionSetClient) EchoMessageArrayAction(*actionset.EchoMessageArrayActionParams) ([]*conflictresolution.Message, error) {
	panic(nil)
}

func (a *actionSetClient) EchoMessageArrayActionWithContext(context.Context, *actionset.EchoMessageArrayActionParams) ([]*conflictresolution.Message, error) {
	panic(nil)
}

func (a *actionSetClient) EchoPrimitiveUnionAction(*actionset.EchoPrimitiveUnionActionParams) (*testsuite.UnionOfPrimitives, error) {
	panic(nil)
}

func (a *actionSetClient) EchoPrimitiveUnionActionWithContext(context.Context, *actionset.EchoPrimitiveUnionActionParams) (*testsuite.UnionOfPrimitives, error) {
	panic(nil)
}

func (a *actionSetClient) EchoStringArrayAction(*actionset.EchoStringArrayActionParams) ([]string, error) {
	panic(nil)
}

func (a *actionSetClient) EchoStringArrayActionWithContext(context.Context, *actionset.EchoStringArrayActionParams) ([]string, error) {
	panic(nil)
}

func (a *actionSetClient) EchoStringMapAction(*actionset.EchoStringMapActionParams) (map[string]string, error) {
	panic(nil)
}

func (a *actionSetClient) EchoStringMapActionWithContext(context.Context, *actionset.EchoStringMapActionParams) (map[string]string, error) {
	panic(nil)
}

func (a *actionSetClient) EchoTyperefUrlAction(*actionset.EchoTyperefUrlActionParams) (testsuite.Url, error) {
	panic(nil)
}

func (a *actionSetClient) EchoTyperefUrlActionWithContext(context.Context, *actionset.EchoTyperefUrlActionParams) (testsuite.Url, error) {
	panic(nil)
}

func (a *actionSetClient) EmptyResponseAction(*actionset.EmptyResponseActionParams) error {
	panic(nil)
}

func (a *actionSetClient) EmptyResponseActionWithContext(context.Context, *actionset.EmptyResponseActionParams) error {
	panic(nil)
}

func (a *actionSetClient) MultipleInputsAction(*actionset.MultipleInputsActionParams) (bool, error) {
	panic(nil)
}

func (a *actionSetClient) MultipleInputsActionWithContext(context.Context, *actionset.MultipleInputsActionParams) (bool, error) {
	panic(nil)
}

func (a *actionSetClient) ReturnBoolAction() (bool, error) {
	panic(nil)
}

func (a *actionSetClient) ReturnBoolActionWithContext(context.Context) (bool, error) {
	panic(nil)
}

func (a *actionSetClient) ReturnIntAction() (int32, error) {
	panic(nil)
}

func (a *actionSetClient) ReturnIntActionWithContext(context.Context) (int32, error) {
	panic(nil)
}
