package std

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// New
// ---------------------------------------------------------------------------------------------------------------------

func Test_NewSyncMap(t *testing.T) {
	m := NewSyncMap[string, int]()

	assert.NotNil(t, m)
	assert.Equal(t, 0, m.Len())
}

// ---------------------------------------------------------------------------------------------------------------------
// Set
// ---------------------------------------------------------------------------------------------------------------------

func Test_SyncMap_Set(t *testing.T) {
	t.Run("set new key", func(t *testing.T) {
		m := NewSyncMap[string, int]()

		m.Set("key1", 100)

		value, exists := m.Get("key1")
		assert.True(t, exists)
		assert.Equal(t, 100, value)
		assert.Equal(t, 1, m.Len())
	})

	t.Run("overwrite existing key", func(t *testing.T) {
		m := NewSyncMap[string, int]()

		m.Set("key1", 100)
		m.Set("key1", 200)

		value, exists := m.Get("key1")
		assert.True(t, exists)
		assert.Equal(t, 200, value)
		assert.Equal(t, 1, m.Len())
	})

	t.Run("thread-safe", func(t *testing.T) {
		m := NewSyncMap[int, int]()
		const goroutines = 10
		const iterations = 1000

		var wg sync.WaitGroup
		wg.Add(goroutines)
		for g := 0; g < goroutines; g++ {
			go func(g int) {
				defer wg.Done()
				for i := 0; i < iterations; i++ {
					key := g*iterations + i
					m.Set(key, key*2)
				}
			}(g)
		}
		wg.Wait()

		for k := 0; k < goroutines*iterations; k++ {
			value, ok := m.Get(k)
			assert.True(t, ok)
			assert.Equal(t, k*2, value)
		}
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// Modify
// ---------------------------------------------------------------------------------------------------------------------

func Test_SyncMap_Modify(t *testing.T) {
	t.Run("modify existing value", func(t *testing.T) {
		m := NewSyncMap[string, int]()
		m.Set("num", 5)

		m.Modify("num", func(v int) int { return v * 2 })

		value, _ := m.Get("num")
		assert.Equal(t, 10, value)
	})

	t.Run("modify non‑existing value (zero value)", func(t *testing.T) {
		m := NewSyncMap[string, int]()

		m.Modify("missing", func(v int) int { return v + 10 })

		value, _ := m.Get("missing")
		assert.Equal(t, 10, value)
		assert.Equal(t, 1, m.Len())
	})

	t.Run("thread safe", func(t *testing.T) {
		m := NewSyncMap[int, int]()
		for i := 0; i < 100; i++ {
			m.Set(i, i)
		}

		const goroutines = 10
		const iterations = 1000

		var wg sync.WaitGroup
		wg.Add(goroutines)
		for g := 0; g < goroutines; g++ {
			go func() {
				defer wg.Done()
				for i := 0; i < iterations; i++ {
					key := i % 100
					m.Modify(key, func(v int) int { return v + 1 })
				}
			}()
		}
		wg.Wait()

		for i := 0; i < 100; i++ {
			value, _ := m.Get(i)
			expected := i + iterations*goroutines/100
			assert.Equal(t, expected, value)
		}
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------------------------------------------------

func Test_SyncMap_Delete(t *testing.T) {
	t.Run("delete existing key", func(t *testing.T) {
		m := NewSyncMap[string, bool]()
		m.Set("flag", true)

		m.Delete("flag")

		_, exists := m.Get("flag")
		assert.False(t, exists)
		assert.Equal(t, 0, m.Len())
	})

	t.Run("delete non‑existing key (no panic)", func(t *testing.T) {
		m := NewSyncMap[string, bool]()

		assert.NotPanics(t, func() {
			m.Delete("nonexistent")
		})

		assert.Equal(t, 0, m.Len())
	})

	t.Run("thread safe", func(t *testing.T) {
		const elems = 10_000
		const goroutines = 10
		m := NewSyncMap[int, string]()
		for i := 0; i < elems; i++ {
			m.Set(i, "value")
		}

		var wg sync.WaitGroup
		wg.Add(goroutines)
		for g := 0; g < goroutines; g++ {
			go func(g int) {
				defer wg.Done()
				for i := g; i < elems; i += goroutines {
					m.Delete(i)
				}
			}(g)
		}
		wg.Wait()

		assert.Equal(t, 0, m.Len())
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// Delete and Get
// ---------------------------------------------------------------------------------------------------------------------

func Test_SyncMap_DeleteAndGet(t *testing.T) {
	t.Run("get and delete existing key", func(t *testing.T) {
		m := NewSyncMap[int, string]()
		m.Set(42, "answer")

		value, ok := m.DeleteAndGet(42)

		assert.True(t, ok)
		assert.Equal(t, "answer", value)

		zeroVal, exists := m.Get(42)
		assert.Equal(t, "", zeroVal)
		assert.False(t, exists)
		assert.Equal(t, 0, m.Len())
	})

	t.Run("try get and delete non‑existing key", func(t *testing.T) {
		m := NewSyncMap[int, string]()

		value, ok := m.DeleteAndGet(999)

		assert.False(t, ok)
		assert.Equal(t, "", value)
		assert.Equal(t, 0, m.Len())
	})

	t.Run("thread safe", func(t *testing.T) {
		const goroutines = 10
		const elems = 10_000
		m := NewSyncMap[int, int]()
		for i := 0; i < elems; i++ {
			m.Set(i, i*i)
		}
		ch := make(chan [2]int, elems)
		restored := NewSyncMap[int, int]()

		wg := new(sync.WaitGroup)
		wg.Add(goroutines)
		for g := 0; g < goroutines; g++ {
			go func(g int) {
				defer wg.Done()
				for i := g; i < elems; i += goroutines {
					v, ok := m.DeleteAndGet(i)
					assert.True(t, ok)
					assert.Equal(t, v, i*i)

					ch <- [2]int{i, v}

					v, ok = m.Get(i)
					assert.False(t, ok)
					assert.Equal(t, 0, v)
				}
			}(g)
		}
		wg.Wait()
		close(ch)

		for e := range ch {
			restored.Set(e[0], e[1])
		}

		assert.Equal(t, 0, m.Len())
		assert.Equal(t, elems, restored.Len())
		for i := 0; i < elems; i++ {
			v, ok := restored.Get(i)
			assert.True(t, ok)
			assert.Equal(t, v, i*i)
		}
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// Truncate
// ---------------------------------------------------------------------------------------------------------------------

func Test_SyncMap_Truncate(t *testing.T) {
	t.Run("truncate non‑empty map", func(t *testing.T) {
		m := NewSyncMap[string, int]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Truncate()

		keys := m.Keys()
		assert.Len(t, keys, 0)
		assert.Equal(t, 0, m.Len())
	})

	t.Run("truncate already empty map", func(t *testing.T) {
		m := NewSyncMap[string, int]()
		m.Truncate()

		keys := m.Keys()
		assert.Len(t, keys, 0)
		assert.Equal(t, 0, m.Len())
	})

	t.Run("thread safe", func(t *testing.T) {
		m := NewSyncMap[string, int]()
		m.Set("a", 1)

		// --

		// don't see a better way to check this,
		// because it is unable to check difference between multiple executions of `Truncate` method.
		sleepTime := 10 * time.Millisecond

		var time1, time2 time.Time
		go func() {
			m.Modify("a", func(v int) int {
				time1 = time.Now()
				time.Sleep(sleepTime)
				return v
			})
		}()
		time.Sleep(sleepTime / 2) // be sure goroutine is started
		m.Truncate()              // here we must be locked
		time2 = time.Now()        // time2 = time1 + sleepTime + ...tiny execution time.

		// --

		diff := time2.Sub(time1) - sleepTime
		assert.Equal(t, 0, m.Len())
		assert.True(t, 0 <= diff && diff < time.Millisecond)
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// Len
// ---------------------------------------------------------------------------------------------------------------------

func Test_SyncMap_Len(t *testing.T) {
	t.Run("length of empty map", func(t *testing.T) {
		m := NewSyncMap[string, string]()
		assert.Equal(t, 0, m.Len())
	})

	t.Run("length after adding keys", func(t *testing.T) {
		m := NewSyncMap[int, bool]()
		assert.Equal(t, 0, m.Len())

		m.Set(1, true)
		assert.Equal(t, 1, m.Len())

		m.Set(2, false)
		assert.Equal(t, 2, m.Len())
	})

	t.Run("length after deletion", func(t *testing.T) {
		m := NewSyncMap[string, int]()
		m.Set("x", 10)
		m.Set("y", 20)
		assert.Equal(t, 2, m.Len())

		m.Delete("x")
		assert.Equal(t, 1, m.Len())
	})

	t.Run("thread safe", func(t *testing.T) {
		m := NewSyncMap[string, int]()
		m.Set("a", 1)

		// --

		// don't see a better way to check this,
		// because it is unable to check difference between multiple executions of `Len` method.
		sleepTime := 10 * time.Millisecond

		var time1, time2 time.Time
		go func() {
			m.Modify("a", func(v int) int {
				time1 = time.Now()
				time.Sleep(sleepTime)
				return v
			})
		}()
		time.Sleep(sleepTime / 2) // be sure goroutine is started
		_ = m.Len()               // here we must be locked
		time2 = time.Now()        // time2 = time1 + sleepTime + ...tiny execution time.

		// --

		diff := time2.Sub(time1) - sleepTime
		assert.True(t, 0 <= diff && diff < time.Millisecond)
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------------------------------------------------

func Test_SyncMap_Get(t *testing.T) {
	t.Run("get existing key", func(t *testing.T) {
		m := NewSyncMap[string, float64]()
		m.Set("pi", 3.14)
		value, ok := m.Get("pi")

		assert.True(t, ok)
		assert.Equal(t, 3.14, value)
	})

	t.Run("get non‑existing key", func(t *testing.T) {
		m := NewSyncMap[string, float64]()
		value, ok := m.Get("unknown")

		var zeroValue float64
		assert.False(t, ok)
		assert.Equal(t, zeroValue, value)
	})

	t.Run("thread safe", func(t *testing.T) {
		m := NewSyncMap[string, int]()
		m.Set("a", 1)

		// --

		// don't see a better way to check this,
		// because it is unable to check difference between multiple executions of `Get` method.
		sleepTime := 10 * time.Millisecond

		var time1, time2 time.Time
		go func() {
			m.Modify("a", func(v int) int {
				time1 = time.Now()
				time.Sleep(sleepTime)
				return v
			})
		}()
		time.Sleep(sleepTime / 2) // be sure goroutine is started
		_, _ = m.Get("a")         // here we must be locked
		time2 = time.Now()        // time2 = time1 + sleepTime + ...tiny execution time.

		// --

		diff := time2.Sub(time1) - sleepTime
		assert.True(t, 0 <= diff && diff < time.Millisecond)
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// Keys
// ---------------------------------------------------------------------------------------------------------------------

func Test_SyncMap_Keys(t *testing.T) {
	t.Run("get keys from empty map", func(t *testing.T) {
		m := NewSyncMap[string, struct{}]()
		keys := m.Keys()
		assert.Empty(t, keys)
	})

	t.Run("get keys from non‑empty map", func(t *testing.T) {
		m := NewSyncMap[string, int]()
		m.Set("first", 1)
		m.Set("second", 2)
		m.Set("third", 3)
		m.Delete("third")

		expectedKeys := map[string]bool{}
		for _, key := range m.Keys() {
			expectedKeys[key] = true
		}

		assert.Equal(t, 2, len(expectedKeys))
		assert.True(t, expectedKeys["first"])
		assert.True(t, expectedKeys["second"])
	})

	t.Run("thread safe", func(t *testing.T) {
		m := NewSyncMap[string, int]()
		m.Set("a", 1)

		// --

		// don't see a better way to check this,
		// because it is unable to check difference between multiple executions of `Keys` method.
		sleepTime := 10 * time.Millisecond

		var time1, time2 time.Time
		go func() {
			m.Modify("a", func(v int) int {
				time1 = time.Now()
				time.Sleep(sleepTime)
				return v
			})
		}()
		time.Sleep(sleepTime / 2) // be sure goroutine is started
		_ = m.Keys()              // here we must be locked
		time2 = time.Now()        // time2 = time1 + sleepTime + ...tiny execution time.

		// --

		diff := time2.Sub(time1) - sleepTime
		assert.True(t, 0 <= diff && diff < time.Millisecond)
	})
}

// ---------------------------------------------------------------------------------------------------------------------
