package collection

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

type SortedSet[T constraints.Ordered] []T

func (set *SortedSet[T]) Add(t T) {
	i, ok := slices.BinarySearch(*set, t)
	if ok {
		return
	}
	*set = slices.Insert(*set, i, t)
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

func (set *SortedSet[T]) Remove(t T) {
	i, ok := slices.BinarySearch(*set, t)
	if !ok {
		return
	}
	*set = append((*set)[:i], (*set)[i+1:]...)
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
