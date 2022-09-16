package syncutils

import "sync"

type RwMap[K comparable, V any] interface {
	Load(k K) (V, bool)
	Get(k K) V
	Store(k K, v V)
	Delete(k K)
	Exists(k K) bool
	Len() int
	Range(f func(k K, v V))
	Copy() map[K]V
	CompareDelete(k K, f func(v V) bool)
}

func NewRwMap[K comparable, V any]() RwMap[K, V] {
	return NewRwMapWithLength[K, V](0)
}

func NewRwMapWithLength[K comparable, V any](length int) RwMap[K, V] {
	return &rwMap[K, V]{
		store: make(map[K]V, length),
	}
}

type rwMap[K comparable, V any] struct {
	mut   sync.RWMutex
	store map[K]V
}

func (m *rwMap[K, V]) Load(k K) (V, bool) {
	m.mut.RLock()
	defer m.mut.RUnlock()

	v, ok := m.store[k]
	return v, ok
}

func (m *rwMap[K, V]) Get(k K) V {
	v, _ := m.Load(k)
	return v
}

func (m *rwMap[K, V]) Store(k K, v V) {
	m.mut.Lock()
	defer m.mut.Unlock()

	m.store[k] = v
}

func (m *rwMap[K, V]) Delete(k K) {
	m.mut.Lock()
	defer m.mut.Unlock()

	delete(m.store, k)
}

func (m *rwMap[K, V]) Exists(k K) bool {
	_, ok := m.Load(k)
	return ok
}

func (m *rwMap[K, V]) Len() int {
	m.mut.RLock()
	defer m.mut.RUnlock()

	return len(m.store)
}

func (m *rwMap[K, V]) Range(f func(k K, v V)) {
	m.mut.RLock()
	defer m.mut.RUnlock()

	for k, v := range m.store {
		f(k, v)
	}
}

func (m *rwMap[K, V]) Copy() map[K]V {
	m.mut.RLock()
	defer m.mut.RUnlock()

	cpy := make(map[K]V, len(m.store))
	for k, v := range m.store {
		cpy[k] = v
	}

	return cpy
}

func (m *rwMap[K, V]) CompareDelete(k K, f func(v V) bool) {
	m.mut.Lock()
	defer m.mut.Unlock()

	if f(m.store[k]) {
		delete(m.store, k)
	}
}

func (m *rwMap[K, V]) Swap(k K, f func(v V) V) {
	m.mut.Lock()
	defer m.mut.Unlock()

	m.store[k] = f(m.store[k])
}
