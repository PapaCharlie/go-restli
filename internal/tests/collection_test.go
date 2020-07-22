package tests

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/generated/conflictResolution"
	. "github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/collection"
	colletionSubCollection "github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/collection/subcollection"
	colletionSubSimple "github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/collection/subsimple"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/stretchr/testify/require"
)

func (s *TestServer) CollectionCreate(t *testing.T, c Client) {
	id, err := c.Create(&conflictresolution.Message{
		Message: "test message",
	})
	require.NoError(t, err)
	require.Equal(t, id, int64(1))
}

func (s *TestServer) CollectionCreate500(t *testing.T, c Client) {
	id, err := c.Create(newMessage(3, "internal error test"))
	require.Errorf(t, err, "Did not receive an error from the server (got %+v)", id)
	require.Equal(t, err.(*protocol.RestLiError).Status, 500)
}

func (s *TestServer) CollectionCreateErrorDetails(t *testing.T, c Client) {
	id, err := c.Create(newMessage(3, "error details test"))
	require.Errorf(t, err, "Did not receive an error from the server (got %+v)", id)
	require.Equal(t, err.(*protocol.RestLiError).Status, 400)
}

func (s *TestServer) CollectionGet(t *testing.T, c Client) {
	id := int64(1)
	res, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t, newMessage(id, "test message"), res)
}

func (s *TestServer) CollectionUpdate(t *testing.T, c Client) {
	id := int64(1)
	err := c.Update(id, newMessage(id, "updated message"))
	require.NoError(t, err)
}

func (s *TestServer) CollectionDelete(t *testing.T, c Client) {
	id := int64(1)
	err := c.Delete(id)
	require.NoError(t, err)
}

func (s *TestServer) CollectionGet404(t *testing.T, c Client) {
	m, err := c.Get(2)
	require.Errorf(t, err, "Did not receive an error from the server (got %+v)", m)
	require.Equal(t, 404, err.(*protocol.RestLiError).Status, "Unexpected status code from server")
}

func (s *TestServer) CollectionUpdate400(t *testing.T, c Client) {
	t.Skip("It is impossible to craft the request required using the generated code because it would require a field " +
		"to be deliberately missing. This can be chalked up as a win for the generated code's safety.")
}

func (s *TestServer) CollectionSearchFinder(t *testing.T, c Client) {
	params := &FindBySearchParams{Keyword: "message"}
	expectedMessages := []*conflictresolution.Message{newMessage(1, "test message"), newMessage(2, "another message")}
	res, err := c.FindBySearch(params)
	require.NoError(t, err)
	require.Equal(t, expectedMessages, res)
}

func (s *TestServer) CollectionPartialUpdate(t *testing.T, c Client) {
	id := int64(1)
	patch := new(conflictresolution.Message_PartialUpdate)
	message := "partial updated message"
	patch.Update.Message = &message
	err := c.PartialUpdate(id, patch)
	require.NoError(t, err)
}

func (s *TestServer) SubCollectionOfCollectionGet(t *testing.T, c Client) {
	id := int64(100)
	res, err := colletionSubCollection.NewClient(s.client).Get(1, id)
	require.NoError(t, err)
	require.Equal(t, newMessage(id, "sub collection message"), res)
}

func (s *TestServer) SubSimpleOfCollectionGet(t *testing.T, c Client) {
	res, err := colletionSubSimple.NewClient(s.client).Get(1)
	require.NoError(t, err)
	require.Equal(t, &conflictresolution.Message{Message: "sub simple message"}, res, "Invalid response from server")
}

func newMessage(id int64, message string) *conflictresolution.Message {
	return &conflictresolution.Message{
		Id:      &id,
		Message: message,
	}
}
