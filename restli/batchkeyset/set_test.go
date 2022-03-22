package batchkeyset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenericBatchKeySetAddTwice(t *testing.T) {
	set := NewBytesKeySet()
	k := []byte{1, 2, 3, 4}
	kCopy := append([]byte(nil), k...)
	require.NoError(t, set.AddKey(k))
	require.Error(t, set.AddKey(kCopy))
}

func TestPrimitiveBatchKeySetAddTwice(t *testing.T) {
	set := NewPrimitiveKeySet[int64]()
	k := int64(42)
	require.NoError(t, set.AddKey(k))
	require.Error(t, set.AddKey(k))
}
