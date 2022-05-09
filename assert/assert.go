package assert

import "fmt"

// Panics on error.
func Must(e error) {
	if e != nil {
		panic(fmt.Errorf("assert.Must failed: %w", e))
	}
}

// Panics on error. Returns "value" if no error.
func MustChain[V any](value V, e error) V {
	if e != nil {
		panic(fmt.Errorf("assert.MustChain failed: %w", e))
	}
	return value
}
