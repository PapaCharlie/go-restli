package nativerefs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseNativeTyperefs(t *testing.T) {
	require.NoError(t, ParseNativeTyperefs("/Users/pchesnai/code/personal/go-restli/internal/tests/native"))
}
