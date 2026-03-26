package std

import "sync"

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type SyncMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		mu: sync.RWMutex{},
		m:  make(map[K]V),
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Set
//
// (Re-)Associates value `v` with the key `k`.
func (m *SyncMap[K, V]) Set(k K, v V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.m[k] = v
}

// Modify
//
// (Re-)Associates key `k` with value `v` modified by the `fn(v)`.
// If key `k` does not exist, `fn` gets zero-value argument.
func (m *SyncMap[K, V]) Modify(k K, fn func(v V) V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	v, _ := m.m[k]

	m.m[k] = fn(v)
}

// --

// Delete
//
// Removes association with the key `k`.
// Does nothing, if key `k` does not exist.
func (m *SyncMap[K, V]) Delete(k K) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.m, k)
}

// DeleteAndGet
//
// Removes association with the key `k` and then returns associated value like the Get() method does.
func (m *SyncMap[K, V]) DeleteAndGet(k K) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	v, ok := m.m[k]
	if ok {
		delete(m.m, k)
	}

	return v, ok
}

// Truncate
//
// Truncates the map.
func (m *SyncMap[K, V]) Truncate() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.m = make(map[K]V)
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (m *SyncMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.m)
}

func (m *SyncMap[K, V]) Get(k K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.m[k]

	return v, ok
}

func (m *SyncMap[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]K, len(m.m))

	i := 0
	for k := range m.m {
		keys[i] = k
		i++
	}

	return keys
}

// ---------------------------------------------------------------------------------------------------------------------
