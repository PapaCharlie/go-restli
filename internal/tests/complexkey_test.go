package tests

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/generated/conflictResolution"
	"github.com/stretchr/testify/require"

	. "github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite/complexkey"
)

func (s *TestServer) ComplexkeyGet(t *testing.T, c Client) {
	id := &Complexkey_ComplexKey{
		Params: newKeyParams("param1", 5),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "one",
			Part2: 2,
			Part3: conflictresolution.Fruits_APPLE,
		},
	}
	res, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t, &conflictresolution.LargeRecord{
		Message: conflictresolution.Message{Message: "test message"},
		Key:     id.ComplexKey,
	}, res)
}

func (s *TestServer) ComplexkeyUpdate(t *testing.T, c Client) {
	key := conflictresolution.ComplexKey{
		Part1: "one",
		Part2: 2,
		Part3: conflictresolution.Fruits_APPLE,
	}
	id := &Complexkey_ComplexKey{
		Params:     newKeyParams("param1", 5),
		ComplexKey: key,
	}
	record := &conflictresolution.LargeRecord{
		Key: key,
		Message: conflictresolution.Message{
			Message: "updated message",
		},
	}
	err := c.Update(id, record)
	require.NoError(t, err)
}

func (s *TestServer) ComplexkeyDelete(t *testing.T, c Client) {
	id := &Complexkey_ComplexKey{
		Params: newKeyParams("param1", 5),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "one",
			Part2: 2,
			Part3: conflictresolution.Fruits_APPLE,
		},
	}
	err := c.Delete(id)
	require.NoError(t, err)
}

func (s *TestServer) ComplexkeyCreate(t *testing.T, c Client) {
	_, err := c.Create(&conflictresolution.LargeRecord{
		Key: conflictresolution.ComplexKey{
			Part1: "one",
			Part2: 2,
			Part3: conflictresolution.Fruits_APPLE,
		},
		Message: conflictresolution.Message{
			Message: "test message",
		},
	})
	require.NoError(t, err)
}

func (s *TestServer) ComplexkeyPartialUpdate(t *testing.T, c Client) {
	id := &Complexkey_ComplexKey{
		Params: newKeyParams("param1", 5),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "one",
			Part2: 2,
			Part3: conflictresolution.Fruits_APPLE,
		},
	}
	keyPatch := new(conflictresolution.ComplexKey_PartialUpdate)
	newPart1 := "newpart1"
	keyPatch.Update.Part1 = &newPart1

	err := c.PartialUpdate(id, &conflictresolution.LargeRecord_PartialUpdate{Key: keyPatch})
	require.NoError(t, err)
}
