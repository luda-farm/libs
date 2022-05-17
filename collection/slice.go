package collection

import (
	"golang.org/x/exp/slices"
)

// The relative complement of b in a.
func Complement[T comparable](a, b []T) []T {
	for _, e := range b {
		a = Remove(a, e)
	}
	return a
}

// A copy of s with at least one instance of t in it.
func PutUnique[T comparable](s []T, t T) []T {
	if !slices.Contains(s, t) {
		s = append(s, t)
	}
	return s
}

// A copy of s without any instances of t.
func Remove[T comparable](s []T, t T) []T {
	result := []T{}
	for {
		i := slices.Index(s, t)
		if i < 0 {
			break
		}
		result = append(result, s[:i]...)
		s = s[i+1:]
	}
	return s
}
