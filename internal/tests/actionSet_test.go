package tests

import (
	"testing"

	conflictresolution "github.com/bored-engineer/go-restli/internal/tests/generated/conflictResolution"
	"github.com/bored-engineer/go-restli/internal/tests/generated/testsuite"
	. "github.com/bored-engineer/go-restli/internal/tests/generated/testsuite/actionSet"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) ActionsetEcho(t *testing.T, c Client) {
	input := "Is anybody out there?"
	output, err := c.EchoAction(&EchoActionParams{Input: &input})
	require.NoError(t, err)
	require.Equal(t, &input, output, "Invalid response from server")
}

func (s *TestServer) ActionsetReturnInt(t *testing.T, c Client) {
	res, err := c.ReturnIntAction()
	require.NoError(t, err)
	i := int32(42)
	require.Equal(t, &i, res, "Invalid response from server")
}

func (s *TestServer) ActionsetReturnBool(t *testing.T, c Client) {
	res, err := c.ReturnBoolAction()
	require.NoError(t, err)
	b := true
	require.Equal(t, &b, res, "Invalid response from server")
}

func (s *TestServer) ActionsetEchoMessage(t *testing.T, c Client) {
	msg := "test message"
	message := conflictresolution.Message{Message: &msg}
	res, err := c.EchoMessageAction(&EchoMessageActionParams{Message: &message})
	require.NoError(t, err)
	require.Equal(t, &message, res, "Invalid response from server")
}

func (s *TestServer) ActionsetEchoMessageArray(t *testing.T, c Client) {
	msg1 := "test message"
	msg2 := "another message"
	messageArray := []*conflictresolution.Message{
		{Message: &msg1},
		{Message: &msg2},
	}
	res, err := c.EchoMessageArrayAction(&EchoMessageArrayActionParams{Messages: messageArray})
	require.NoError(t, err)
	require.Equal(t, messageArray, res, "Invalid response from server")
}

func (s *TestServer) ActionsetEchoStringArray(t *testing.T, c Client) {
	stringArray := []string{"string one", "string two"}
	res, err := c.EchoStringArrayAction(&EchoStringArrayActionParams{Strings: stringArray})
	require.NoError(t, err)
	require.Equal(t, stringArray, res, "Invalid response from server")
}

func (s *TestServer) ActionsetEchoStringMap(t *testing.T, c Client) {
	stringMap := map[string]string{
		"one": "string one",
		"two": "string two",
	}
	res, err := c.EchoStringMapAction(&EchoStringMapActionParams{Strings: stringMap})
	require.NoError(t, err)
	require.Equal(t, stringMap, res, "Invalid response from server")
}

func (s *TestServer) ActionsetEchoTyperefUrl(t *testing.T, c Client) {
	var urlTyperef testsuite.Url = "http://rest.li"
	res, err := c.EchoTyperefUrlAction(&EchoTyperefUrlActionParams{UrlTyperef: &urlTyperef})
	require.NoError(t, err)
	require.Equal(t, urlTyperef, *res, "Invalid response from server")
}

func (s *TestServer) ActionsetEchoPrimitiveUnion(t *testing.T, c Client) {
	union := &testsuite.UnionOfPrimitives{}
	//union.InitializePrimitivesUnion()
	union.PrimitivesUnion.Long = new(int64)
	*union.PrimitivesUnion.Long = 100

	res, err := c.EchoPrimitiveUnionAction(&EchoPrimitiveUnionActionParams{PrimitiveUnion: union})
	require.NoError(t, err)
	require.Equal(t, *union, *res, "Invalid response from server")
}

func (s *TestServer) ActionsetEchoComplexTypesUnion(t *testing.T, c Client) {
	union := &testsuite.UnionOfComplexTypes{}
	union.ComplexTypeUnion.Fruits = new(conflictresolution.Fruits)
	*union.ComplexTypeUnion.Fruits = conflictresolution.Fruits_APPLE

	res, err := c.EchoComplexTypesUnionAction(&EchoComplexTypesUnionActionParams{ComplexTypesUnion: union})
	require.NoError(t, err)
	require.Equal(t, *union, *res, "Invalid response from server")
}

func (s *TestServer) ActionsetEmptyResponse(t *testing.T, c Client) {
	msg1 := "test message"
	msg2 := "another message"
	err := c.EmptyResponseAction(&EmptyResponseActionParams{
		Message1: &conflictresolution.Message{Message: &msg1},
		Message2: &conflictresolution.Message{Message: &msg2},
	})
	require.NoError(t, err)
}

func (s *TestServer) ActionsetMultipleInputs(t *testing.T, c Client) {
	optionalString := "optional string"
	str := "string"
	url := testsuite.Url("http://rest.li")
	msg := "test message"
	res, err := c.MultipleInputsAction(&MultipleInputsActionParams{
		String:         &str,
		Message:        &conflictresolution.Message{Message: &msg},
		UrlTyperef:     &url,
		OptionalString: &optionalString,
	})
	require.NoError(t, err)
	require.True(t, *res, "Invalid response from server")
}

func (s *TestServer) ActionsetMultipleInputsNoOptional(t *testing.T, c Client) {
	str := "string"
	url := testsuite.Url("http//rest.li")
	msg := "test message"
	res, err := c.MultipleInputsAction(&MultipleInputsActionParams{
		String:     &str,
		Message:    &conflictresolution.Message{Message: &msg},
		UrlTyperef: &url,
	})
	require.NoError(t, err)
	require.True(t, *res, "Invalid response from server")
}
