package assert

import "fmt"

func NilError(e error) {
	if e != nil {
		panic(fmt.Errorf("Assert NilError Failed: %w", e))
	}
}
