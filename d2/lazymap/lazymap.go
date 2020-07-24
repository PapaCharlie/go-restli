package lazymap

import "sync"

type LazySyncMap sync.Map

type inFlightValue struct {
	wg sync.WaitGroup
	v  interface{}
}

func (m *LazySyncMap) LoadOrStore(key interface{}, f func() interface{}) interface{} {
	value := new(inFlightValue)
	value.wg.Add(1)

	if s, loaded := (*sync.Map)(m).LoadOrStore(key, value); loaded {
		if v, ok := s.(*inFlightValue); ok {
			v.wg.Wait()
			return v.v
		} else {
			return s
		}
	}

	value.v = f()
	(*sync.Map)(m).Store(key, value.v)
	value.wg.Done()
	return value.v
}

func (m *LazySyncMap) Load(key interface{}) (interface{}, bool) {
	s, ok := (*sync.Map)(m).Load(key)
	if !ok {
		return nil, false
	}

	if v, ok := s.(*inFlightValue); ok {
		v.wg.Wait()
		return v.v, true
	} else {
		return s, true
	}
}

func (m *LazySyncMap) Store(key interface{}, value interface{}) {
	stored := false
	m.LoadOrStore(key, func() interface{} {
		stored = true
		return value
	})

	if !stored {
		(*sync.Map)(m).Store(key, value)
	}
}
