package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyRecord_UnmarshalRestLi(t *testing.T) {
	require.True(t, IsEmptyRecord[EmptyRecord](EmptyRecord{}))
	require.False(t, IsEmptyRecord[struct{}](struct{}{}))
}
