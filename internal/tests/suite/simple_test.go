package suite

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/simple"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/simple_test"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/stretchr/testify/require"
)

func (o *Operation) SimpleGet(t *testing.T, c Client) func(*testing.T) *MockResource {
	expected := &conflictresolution.Message{Message: "test message"}
	res, err := c.Get()
	require.NoError(t, err)
	require.Equal(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockGet: func(ctx *protocol.RequestContext) (entity *conflictresolution.Message, err error) {
				return expected, nil
			},
		}
	}
}

func (o *Operation) SimpleUpdate(t *testing.T, c Client) func(*testing.T) *MockResource {
	expected := &conflictresolution.Message{Message: "updated message"}
	err := c.Update(expected)
	require.NoError(t, err)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockUpdate: func(ctx *protocol.RequestContext, entity *conflictresolution.Message) (err error) {
				require.Equal(t, expected, entity)
				return nil
			},
		}
	}
}

func (o *Operation) SimpleDelete(t *testing.T, c Client) func(*testing.T) *MockResource {
	err := c.Delete()
	require.NoError(t, err)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockDelete: func(ctx *protocol.RequestContext) (err error) {
				return nil
			},
		}
	}
}
