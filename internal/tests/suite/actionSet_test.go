package suite

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/actionSet"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/actionSet_test"
	"github.com/PapaCharlie/go-restli/restli"
	"github.com/stretchr/testify/require"
)

func (o *Operation) ActionsetEcho(t *testing.T, c Client) func(*testing.T) *MockResource {
	input := "Is anybody out there?"
	output, err := c.EchoAction(&EchoActionParams{Input: input})
	require.NoError(t, err)
	require.Equal(t, input, output, "Invalid response from server")

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockEchoAction: func(ctx *restli.RequestContext, actionParams *EchoActionParams) (actionResult string, err error) {
				require.Equal(t, input, actionParams.Input)
				return actionParams.Input, nil
			},
		}
	}
}

func (o *Operation) ActionsetReturnInt(t *testing.T, c Client) func(*testing.T) *MockResource {
	expected := int32(42)
	res, err := c.ReturnIntAction()
	require.NoError(t, err)
	require.Equal(t, expected, res, "Invalid response from server")

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockReturnIntAction: func(ctx *restli.RequestContext) (actionResult int32, err error) {
				return 42, nil
			},
		}
	}
}

func (o *Operation) ActionsetReturnBool(t *testing.T, c Client) func(*testing.T) *MockResource {
	expected := true
	res, err := c.ReturnBoolAction()
	require.NoError(t, err)
	require.Equal(t, expected, res, "Invalid response from server")

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockReturnBoolAction: func(ctx *restli.RequestContext) (actionResult bool, err error) {
				return expected, nil
			},
		}
	}
}

func (o *Operation) ActionsetEchoMessage(t *testing.T, c Client) func(*testing.T) *MockResource {
	message := conflictresolution.Message{Message: "test message"}
	res, err := c.EchoMessageAction(&EchoMessageActionParams{Message: message})
	require.NoError(t, err)
	require.Equal(t, &message, res, "Invalid response from server")

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockEchoMessageAction: func(ctx *restli.RequestContext, actionParams *EchoMessageActionParams) (actionResult *conflictresolution.Message, err error) {
				require.Equal(t, message, actionParams.Message)
				return &message, nil
			},
		}
	}
}

func (o *Operation) ActionsetEchoMessageArray(t *testing.T, c Client) func(*testing.T) *MockResource {
	messageArray := []*conflictresolution.Message{
		{Message: "test message"},
		{Message: "another message"},
	}
	res, err := c.EchoMessageArrayAction(&EchoMessageArrayActionParams{Messages: messageArray})
	require.NoError(t, err)
	require.Equal(t, messageArray, res, "Invalid response from server")

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockEchoMessageArrayAction: func(ctx *restli.RequestContext, actionParams *EchoMessageArrayActionParams) (actionResult []*conflictresolution.Message, err error) {
				require.Equal(t, messageArray, actionParams.Messages)
				return messageArray, nil
			},
		}
	}
}

func (o *Operation) ActionsetEchoStringArray(t *testing.T, c Client) func(*testing.T) *MockResource {
	stringArray := []string{"string one", "string two"}
	res, err := c.EchoStringArrayAction(&EchoStringArrayActionParams{Strings: stringArray})
	require.NoError(t, err)
	require.Equal(t, stringArray, res, "Invalid response from server")

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockEchoStringArrayAction: func(ctx *restli.RequestContext, actionParams *EchoStringArrayActionParams) (actionResult []string, err error) {
				require.Equal(t, stringArray, actionParams.Strings)
				return stringArray, nil
			},
		}
	}
}

func (o *Operation) ActionsetEchoStringMap(t *testing.T, c Client) func(*testing.T) *MockResource {
	stringMap := map[string]string{
		"one": "string one",
		"two": "string two",
	}
	res, err := c.EchoStringMapAction(&EchoStringMapActionParams{Strings: stringMap})
	require.NoError(t, err)
	require.Equal(t, stringMap, res, "Invalid response from server")

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockEchoStringMapAction: func(ctx *restli.RequestContext, actionParams *EchoStringMapActionParams) (actionResult map[string]string, err error) {
				require.Equal(t, stringMap, actionParams.Strings)
				return stringMap, nil
			},
		}
	}
}

func (o *Operation) ActionsetEchoTyperefUrl(t *testing.T, c Client) func(*testing.T) *MockResource {
	urlTyperef := testsuite.Url("http://rest.li")
	res, err := c.EchoTyperefUrlAction(&EchoTyperefUrlActionParams{UrlTyperef: urlTyperef})
	require.NoError(t, err)
	require.Equal(t, urlTyperef, res, "Invalid response from server")

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockEchoTyperefUrlAction: func(_ *restli.RequestContext, actionParams *EchoTyperefUrlActionParams) (actionResult testsuite.Url, err error) {
				require.Equal(t, urlTyperef, actionParams.UrlTyperef)
				return urlTyperef, nil
			},
		}
	}
}

func (o *Operation) ActionsetEchoPrimitiveUnion(t *testing.T, c Client) func(*testing.T) *MockResource {
	union := &testsuite.UnionOfPrimitives{
		PrimitivesUnion: testsuite.UnionOfPrimitives_PrimitivesUnion{Long: restli.Int64Pointer(100)},
	}
	res, err := c.EchoPrimitiveUnionAction(&EchoPrimitiveUnionActionParams{PrimitiveUnion: *union})
	require.NoError(t, err)
	require.Equal(t, union, res, "Invalid response from server")

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockEchoPrimitiveUnionAction: func(ctx *restli.RequestContext, actionParams *EchoPrimitiveUnionActionParams) (actionResult *testsuite.UnionOfPrimitives, err error) {
				require.Equal(t, *union, actionParams.PrimitiveUnion)
				return union, nil
			},
		}
	}
}

func (o *Operation) ActionsetEchoComplexTypesUnion(t *testing.T, c Client) func(*testing.T) *MockResource {
	union := &testsuite.UnionOfComplexTypes{
		ComplexTypeUnion: testsuite.UnionOfComplexTypes_ComplexTypeUnion{
			Fruits: conflictresolution.Fruits_APPLE.Pointer(),
		},
	}
	res, err := c.EchoComplexTypesUnionAction(&EchoComplexTypesUnionActionParams{ComplexTypesUnion: *union})
	require.NoError(t, err)
	require.Equal(t, union, res, "Invalid response from server")

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockEchoComplexTypesUnionAction: func(_ *restli.RequestContext, actionParams *EchoComplexTypesUnionActionParams) (actionResult *testsuite.UnionOfComplexTypes, err error) {
				require.Equal(t, *union, actionParams.ComplexTypesUnion)
				return union, nil
			},
		}
	}
}

func (o *Operation) ActionsetEmptyResponse(t *testing.T, c Client) func(*testing.T) *MockResource {
	params := &EmptyResponseActionParams{
		Message1: conflictresolution.Message{Message: "test message"},
		Message2: conflictresolution.Message{Message: "another message"},
	}
	require.NoError(t, c.EmptyResponseAction(params))

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockEmptyResponseAction: func(_ *restli.RequestContext, actionParams *EmptyResponseActionParams) (err error) {
				require.Equal(t, params, actionParams)
				return nil
			},
		}
	}
}

func (o *Operation) ActionsetMultipleInputs(t *testing.T, c Client) func(*testing.T) *MockResource {
	params := &MultipleInputsActionParams{
		String:         "string",
		Message:        conflictresolution.Message{Message: "test message"},
		UrlTyperef:     "http://rest.li",
		OptionalString: restli.StringPointer("optional string"),
	}
	res, err := c.MultipleInputsAction(params)
	require.NoError(t, err)
	require.True(t, res, "Invalid response from server")

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockMultipleInputsAction: func(ctx *restli.RequestContext, actionParams *MultipleInputsActionParams) (actionResult bool, err error) {
				require.Equal(t, params, actionParams)
				return true, nil
			},
		}
	}
}

func (o *Operation) ActionsetMultipleInputsNoOptional(t *testing.T, c Client) func(*testing.T) *MockResource {
	params := &MultipleInputsActionParams{
		String:     "string",
		Message:    conflictresolution.Message{Message: "test message"},
		UrlTyperef: "http//rest.li",
	}
	res, err := c.MultipleInputsAction(params)
	require.NoError(t, err)
	require.True(t, res, "Invalid response from server")

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockMultipleInputsAction: func(_ *restli.RequestContext, actionParams *MultipleInputsActionParams) (actionResult bool, err error) {
				require.Equal(t, params, actionParams)
				return true, nil
			},
		}
	}
}
