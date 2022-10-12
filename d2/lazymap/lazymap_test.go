package lazymap

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLazySyncMap_LoadOrStore(t *testing.T) {
	var m LazySyncMap[string, int]
	v, _ := m.LazyLoad("", func() (int, error) {
		return 1, nil
	})
	require.Equal(t, 1, v)
}

func TestLazySyncMap_SubsequentLoadOrStore(t *testing.T) {
	var m LazySyncMap[string, int]
	_, _ = m.LazyLoad("", func() (int, error) {
		return 1, nil
	})

	v, _ := m.LazyLoad("", func() (int, error) {
		t.Fail()
		return 0, nil
	})

	require.Equal(t, v, 1)
}

func TestLazySyncMap_CompetingWrites(t *testing.T) {
	var m LazySyncMap[string, int]

	inLambda := new(sync.WaitGroup)
	inLambda.Add(1)

	go func() {
		_, _ = m.LazyLoad("", func() (int, error) {
			inLambda.Done()
			return 1, nil
		})
	}()

	inLambda.Wait()
	v, _ := m.LazyLoad("", func() (int, error) {
		t.Fail()
		return 0, nil
	})

	require.Equal(t, v, 1)
}
