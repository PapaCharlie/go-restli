package suite

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collection"
	colletionSubCollection "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collection/subcollection"
	colletionSubSimple "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collection/subsimple"
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
	require.Equal(t, err.(*protocol.RestLiError).Response.StatusCode, 500)
}

func (s *TestServer) CollectionCreateErrorDetails(t *testing.T, c Client) {
	id, err := c.Create(newMessage(3, "error details test"))
	require.Errorf(t, err, "Did not receive an error from the server (got %+v)", id)
	require.Equal(t, err.(*protocol.RestLiError).Response.StatusCode, 400)
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
	require.Equal(t, 404, err.(*protocol.RestLiError).Response.StatusCode, "Unexpected status code from server")
}

func (s *TestServer) CollectionUpdate400(t *testing.T, c Client) {
	t.Log("It is impossible to craft the request required using the generated code because it would require a field " +
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
	patch.Update_Fields.Message = &message
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

func (s *TestServer) CollectionBatchGet(t *testing.T, c Client) {
	keys := []int64{1, 3}
	two := int64(2)
	res, err := c.BatchGet(keys)
	require.NoError(t, err)
	require.Equal(t, map[int64]*conflictresolution.Message{
		keys[0]: {
			Id:      &keys[0],
			Message: "test message",
		},
		keys[1]: {
			Id:      &two,
			Message: "another message",
		},
	}, res)
}

func (s *TestServer) CollectionBatchCreate(t *testing.T, c Client) {
	t.SkipNow()
}

func (s *TestServer) CollectionBatchUpdate(t *testing.T, c Client) {
	t.SkipNow()
}

func (s *TestServer) CollectionBatchPartialUpdate(t *testing.T, c Client) {
	keys := []int64{1, 3}
	msg1 := "partial updated message"
	msg2 := "another partial message"
	res, err := c.BatchPartialUpdate(map[int64]*conflictresolution.Message_PartialUpdate{
		keys[0]: {
			Update_Fields: struct {
				Id      *int64
				Message *string
			}{
				Message: &msg1,
			},
		},
		keys[1]: {
			Update_Fields: struct {
				Id      *int64
				Message *string
			}{
				Message: &msg2,
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, map[int64]*protocol.BatchEntityUpdateResponse{
		keys[0]: {
			Status: 204,
		},
		keys[1]: {
			Status: 204,
		},
	}, res)
}

func (s *TestServer) CollectionBatchDelete(t *testing.T, c Client) {
	t.SkipNow()
}

func (s *TestServer) CollectionGetProjection(t *testing.T, c Client) {
	t.SkipNow()
}

func (s *TestServer) CollectionBatchUpdateErrors(t *testing.T, c Client) {
	t.SkipNow()
}
