package fnv1a

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddMap(t *testing.T) {
	// Generate two identical maps
	m1 := map[string]int32{}
	m2 := map[string]int32{}
	for i := 0; i < 10; i++ {
		k := strconv.Itoa(i)
		v := rand.Int31()
		m1[k] = v
		m2[k] = v
	}

	// Run the test multiple times to ensure iteration order does not change the hash value
	for i := 0; i < 10; i++ {
		// Hash the maps and compare the values
		h1 := NewHash()
		AddMap(h1, m1, Hash.AddInt32)
		h2 := NewHash()
		AddMap(h2, m2, Hash.AddInt32)

		require.True(t, h1.Equals(h2), "Expected %s but got %s", h1, h2)
	}
}
