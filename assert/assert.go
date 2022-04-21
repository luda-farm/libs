package assert

import "fmt"

func NilError(e error) {
	if e != nil {
		panic(fmt.Errorf("Assert NilError Failed: %w", e))
	}
}

func Must[V any](value V, e error) V {
	if e != nil {
		panic(fmt.Errorf("Assert Must Failed: %w", e))
	}
	return value
}
