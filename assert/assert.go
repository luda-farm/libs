package assert

import "fmt"

// Panics on error.
func NilError(e error) {
	if e != nil {
		panic(fmt.Errorf("Assert NilError Failed: %w", e))
	}
}

// Panics on error. Returns "value" if no error.
func Must[V any](value V, e error) V {
	if e != nil {
		panic(fmt.Errorf("Assert Must Failed: %w", e))
	}
	return value
}
