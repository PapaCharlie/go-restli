package lazymap

import (
	"sync"
)

type LazySyncMap[K comparable, V any] sync.Map

type inFlightValue[V any] struct {
	wg  sync.WaitGroup
	v   V
	err error
}

func (m *LazySyncMap[K, V]) LazyLoad(key K, f func() (V, error)) (V, error) {
	value := new(inFlightValue[V])
	value.wg.Add(1)
	defer value.wg.Done()

	if s, loaded := (*sync.Map)(m).LoadOrStore(key, value); loaded {
		if v, ok := s.(*inFlightValue[V]); ok {
			v.wg.Wait()
			return v.v, v.err
		} else {
			return s.(V), nil
		}
	}

	value.v, value.err = f()
	if value.err != nil {
		(*sync.Map)(m).Delete(key)
	} else {
		(*sync.Map)(m).Store(key, value.v)
	}
	return value.v, value.err
}
