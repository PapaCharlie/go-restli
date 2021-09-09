package suite

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/stretchr/testify/require"

	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/complexkey"
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
	expectedKey := conflictresolution.ComplexKey{
		Part1: "one",
		Part2: 2,
		Part3: conflictresolution.Fruits_APPLE,
	}
	_, err := c.Create(&conflictresolution.LargeRecord{
		Key: expectedKey,
		Message: conflictresolution.Message{
			Message: "test message",
		},
	})
	require.IsType(t, new(protocol.CreateResponseHasNoEntityHeaderError), err)
	// TODO: Merge https://github.com/linkedin/rest.li-test-suite/pull/6 and actually test the contents of the key
	// require.Equal(t, expectedKey, actualKey.ComplexKey)
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
	keyPatch.Update_Fields.Part1 = &newPart1

	err := c.PartialUpdate(id, &conflictresolution.LargeRecord_PartialUpdate{Key: keyPatch})
	require.NoError(t, err)
}

func (s *TestServer) ComplexkeyBatchGet(t *testing.T, c Client) {
	k1 := &Complexkey_ComplexKey{
		Params: newKeyParams("param1", 5),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "one",
			Part2: 2,
			Part3: conflictresolution.Fruits_APPLE,
		},
	}
	k2 := &Complexkey_ComplexKey{
		Params: newKeyParams("param2", 11),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "two",
			Part2: 7,
			Part3: conflictresolution.Fruits_ORANGE,
		},
	}
	res, err := c.BatchGet([]*Complexkey_ComplexKey{k1, k2})
	require.NoError(t, err)
	require.Equal(t, map[*Complexkey_ComplexKey]*conflictresolution.LargeRecord{
		k1: {
			Key: k1.ComplexKey,
			Message: conflictresolution.Message{
				Message: "test message",
			},
		},
		k2: {
			Key: k2.ComplexKey,
			Message: conflictresolution.Message{
				Message: "test message",
			},
		},
	}, res)
}

const specialChars = `!*'();:@&=+$,/?#[].~`

var specialCharsKey = &Complexkey_ComplexKey{
	Params: newKeyParams("param"+specialChars, 5),
	ComplexKey: conflictresolution.ComplexKey{
		Part1: "key" + specialChars,
		Part2: 2,
		Part3: conflictresolution.Fruits_APPLE,
	},
}

func (s *TestServer) ComplexkeyGetWithSpecialChars(t *testing.T, c Client) {
	_, err := c.Get(specialCharsKey)
	require.NoError(t, err)
}

func (s *TestServer) ComplexkeyBatchGetWithSpecialChars(t *testing.T, c Client) {
	k := &Complexkey_ComplexKey{
		Params: newKeyParams("param2", 11),
		ComplexKey: conflictresolution.ComplexKey{
			Part1: "two",
			Part2: 7,
			Part3: conflictresolution.Fruits_ORANGE,
		},
	}
	res, err := c.BatchGet([]*Complexkey_ComplexKey{specialCharsKey, k})
	require.NoError(t, err)
	_, ok := res[specialCharsKey]
	require.True(t, ok)
	_, ok = res[k]
	require.True(t, ok)
}

func (s *TestServer) ComplexkeyGetProjection(t *testing.T, c Client) {
	t.SkipNow()
}

func (s *TestServer) ComplexkeyBatchCreate(t *testing.T, c Client) {
	t.SkipNow()
}

func (s *TestServer) ComplexkeyBatchUpdate(t *testing.T, c Client) {
	t.SkipNow()
}

func (s *TestServer) ComplexkeyBatchDelete(t *testing.T, c Client) {
	t.SkipNow()
}
