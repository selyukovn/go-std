package std

import (
	"sync"
)

// Wg
//
// Helper function to simplify usage of the sync.WaitGroup.
//
// Panics in case of empty arguments.
// Panics, if `fns` contains nil-element.
func Wg(fns ...func()) {
	if len(fns) == 0 {
		panic("`fns` must have at least one function")
	}

	for _, fn := range fns {
		if fn == nil {
			panic("`fns` must have a non-nil functions")
		}
	}

	// --

	wg := new(sync.WaitGroup)

	wg.Add(len(fns))

	for _, fn := range fns {
		go func(fn func()) {
			defer wg.Done()
			fn()
		}(fn)
	}

	wg.Wait()
}

// WgRange
//
// Helper function to simplify usage of the sync.WaitGroup.
//
// Makes a slice of functions and calls Wg.
//
// Panics in case of empty arguments.
// Panics in case of negative `n`.
// Panics in case of inner call of the Wg panics.
func WgRange(n int, fnProvideGoFn func(i int) func()) {
	if n <= 0 {
		panic("`n` must be greater than zero")
	}

	if fnProvideGoFn == nil {
		panic("`fnProvideGoFn` must have a non-nil function")
	}

	fns := make([]func(), n)

	for i := 0; i < n; i++ {
		fns[i] = fnProvideGoFn(i)
	}

	Wg(fns...)
}
