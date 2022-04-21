package collection

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func AscendingKeys[M map[K]V, K constraints.Ordered, V any](m M) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func DescendingKeys[M map[K]V, K constraints.Ordered, V any](m M) []K {
	return Reverse(AscendingKeys(m))
}

func Reverse[S []V, V any](s S) S {
	i, j := 0, len(s)-1
	for i < j {
		tmp := s[i]
		s[i] = s[j]
		s[j] = tmp
		i++
		j--
	}
	return s
}
