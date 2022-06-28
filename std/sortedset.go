package std

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

type SortedSet[T constraints.Ordered] []T

// returns whether the set was modified
func (set *SortedSet[T]) Add(t T) bool {
	i, ok := slices.BinarySearch(*set, t)
	if ok {
		return false
	}
	*set = slices.Insert(*set, i, t)
	return true
}

func (set SortedSet[T]) Complement(other SortedSet[T]) SortedSet[T] {
	complement := SortedSet[T]{}
	for i, j := 0, 0; i < len(set) && j < len(other); {
		if set[i] < other[j] {
			complement = append(complement, set[i])
			i++
		} else if other[j] < set[i] {
			j++
		} else {
			i++
			j++
		}
	}
	return complement
}

func (set SortedSet[T]) Contains(t T) bool {
	_, ok := slices.BinarySearch(set, t)
	return ok
}

func (set SortedSet[T]) Intersection(other SortedSet[T]) SortedSet[T] {
	intersection := SortedSet[T]{}
	for i, j := 0, 0; i < len(set) && j < len(other); {
		if set[i] < other[j] {
			i++
		} else if other[j] < set[i] {
			j++
		} else {
			intersection = append(intersection, set[i])
			i++
			j++
		}
	}
	return intersection
}

// returns whether the set was modified
func (set *SortedSet[T]) Remove(t T) bool {
	i, ok := slices.BinarySearch(*set, t)
	if !ok {
		return false
	}
	*set = append((*set)[:i], (*set)[i+1:]...)
	return true
}

func (set SortedSet[T]) Union(other SortedSet[T]) SortedSet[T] {
	if len(set) < len(other) {
		set, other = other, set
	}
	union := append(SortedSet[T]{}, set...)
	for _, t := range other {
		union.Add(t)
	}
	return union
}
