package lazymap

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLazySyncMap_LoadOrStore(t *testing.T) {
	var m LazySyncMap
	v := m.LoadOrStore("", func() interface{} {
		return 1
	})
	require.Equal(t, 1, v)
}

func TestLazySyncMap_SubsequentLoadOrStore(t *testing.T) {
	var m LazySyncMap
	m.LoadOrStore("", func() interface{} {
		return 1
	})

	v := m.LoadOrStore("", func() interface{} {
		t.Fail()
		return nil
	})

	require.Equal(t, v, 1)
}

func TestLazySyncMap_CompetingWrites(t *testing.T) {
	var m LazySyncMap

	inLambda := new(sync.WaitGroup)
	inLambda.Add(1)

	go func() {
		m.LoadOrStore("", func() interface{} {
			inLambda.Done()
			return 1
		})
	}()

	inLambda.Wait()
	v := m.LoadOrStore("", func() interface{} {
		t.Fail()
		return nil
	})

	require.Equal(t, v, 1)
}

func TestLazySyncMap_Store(t *testing.T) {
	var m LazySyncMap

	inLambda := new(sync.WaitGroup)
	inLambda.Add(1)

	go func() {
		m.LoadOrStore("", func() interface{} {
			inLambda.Done()
			return 1
		})
	}()

	inLambda.Wait()
	m.Store("", 2)

	loaded, ok := m.Load("")
	require.True(t, ok)
	require.Equal(t, 2, loaded)
}
