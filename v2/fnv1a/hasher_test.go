package fnv1a

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddMap(t *testing.T) {
	// Generate two identical maps
	m1 := randomMap()
	m2 := map[string]int32{}
	for k, v := range m1 {
		m2[k] = v
	}
	m3 := randomMap()

	// Run the test multiple times to ensure iteration order does not change the hash value
	for i := 0; i < 10; i++ {
		// Hash the maps and compare the values
		compareMaps(t, m1, m2, true)
		compareMaps(t, m1, m3, false)
		compareMaps(t, m2, m3, false)
	}
}

func randomMap() map[string]int32 {
	m := map[string]int32{}
	for i := 0; i < 10; i++ {
		m[strconv.Itoa(i)] = rand.Int31()
	}
	return m
}

func compareMaps(t *testing.T, m1, m2 map[string]int32, equals bool) {
	h1 := NewHash()
	AddMap(h1, m1, Hash.AddInt32)
	h2 := NewHash()
	AddMap(h2, m2, Hash.AddInt32)

	if equals {
		require.Equal(t, m1, m2)
		require.True(t, h1.Equals(h2))
	} else {
		require.NotEqual(t, m1, m2)
		require.False(t, h1.Equals(h2))
	}
}
