package gfcontainer

import "sync"

type SafeMap[K comparable, V any] struct {
	m sync.Map
}

func (m *SafeMap[K, V]) Load(key K) (val V, ok bool) {
	v, ok := m.m.Load(key)
	if !ok {
		return val, false
	}
	ret, ok := v.(V)
	if !ok {
		return val, false
	}
	return ret, ok
}

func (m *SafeMap[K, V]) Store(key K, val V) {
	m.m.Store(key, val)
}

func (m *SafeMap[K, V]) LoadOrStore(key K, val V) (actual V, loaded bool) {
	v, ok := m.m.LoadOrStore(key, val)
	if !ok {
		return val, false
	}
	ret, ok := v.(V)
	if !ok {
		return val, false
	}
	return ret, ok
}

func (m *SafeMap[K, V]) Range(f func(key K, val V) bool) {
	m.m.Range(func(k, v any) bool {
		key, ok := k.(K)
		if !ok {
			return false
		}
		val, ok := v.(V)
		if !ok {
			return false
		}
		return f(key, val)
	})
}

func (m *SafeMap[K, V]) Swap(key K, val V) (previous V, loaded bool) {
	v, ok := m.m.Swap(key, val)
	if !ok {
		return val, false
	}
	ret, ok := v.(V)
	if !ok {
		return val, false
	}
	return ret, ok
}

func (m *SafeMap[K, V]) CompareAndSwap(key K, old, new V) bool {
	return m.m.CompareAndSwap(key, old, new)
}

func (m *SafeMap[K, V]) Delete(key K) {
	m.m.Delete(key)
}

func (m *SafeMap[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return m.m.CompareAndDelete(key, old)
}

func (m *SafeMap[K, V]) LoadAndDelete(key K) (val V, loaded bool) {
	v, ok := m.m.LoadAndDelete(key)
	if !ok {
		return val, false
	}
	ret, ok := v.(V)
	if !ok {
		return val, false
	}
	return ret, ok
}
