package z

import (
	"fmt"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func FormatCents[I constraints.Integer](i I) string {
	if i < 0 {
		return fmt.Sprintf("%d.%02d", i/100, -i%100)
	} else {
		return fmt.Sprintf("%d.%02d", i/100, i%100)
	}
}

func Must(e error) {
	if e != nil {
		panic(e)
	}
}

func MustChain[V any](v V, e error) V {
	Must(e)
	return v
}

func SortedKeys[M map[K]V, K constraints.Ordered, V any](m M) []K {
	keys := maps.Keys(m)
	slices.Sort(keys)
	return keys
}
