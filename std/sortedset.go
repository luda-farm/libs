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

// returns whether the set was modified
func (set *SortedSet[T]) AddAll(slice []T) bool {
	originalLength := len(*set)
	other := NewSortedSetFromSlice(slice)
	*set = set.Union(other)
	return len(*set) > originalLength
}

// returns whether the set was modified
func (set *SortedSet[T]) Remove(t T) bool {
	i, ok := slices.BinarySearch(*set, t)
	if !ok {
		return false
	}
	*set = slices.Delete(*set, i, i+1)
	return true
}

func (set SortedSet[T]) Contains(t T) bool {
	_, ok := slices.BinarySearch(set, t)
	return ok
}

func (set SortedSet[T]) Difference(other SortedSet[T]) SortedSet[T] {
	difference := make(SortedSet[T], 0, len(set))

	i, j := 0, 0
	for i < len(set) && j < len(other) {
		switch {
		case set[i] < other[j]:
			difference = append(difference, set[i])
			i++
		case other[j] < set[i]:
			j++
		default:
			i++
			j++
		}
	}

	if i < len(set) {
		difference = append(difference, set[i:]...)
	}

	return difference
}

func (set SortedSet[T]) Intersection(other SortedSet[T]) SortedSet[T] {
	capacity := len(set)
	if len(other) < len(set) {
		capacity = len(other)
	}

	intersection := make(SortedSet[T], 0, capacity)

	for i, j := 0, 0; i < len(set) && j < len(other); {
		switch {
		case set[i] < other[j]:
			i++
		case other[j] < set[i]:
			j++
		default:
			intersection = append(intersection, set[i])
			i++
			j++
		}
	}

	return intersection
}

func (set SortedSet[T]) Union(other SortedSet[T]) SortedSet[T] {
	union := make(SortedSet[T], 0, len(set)+len(other))

	i, j := 0, 0
	for i < len(set) && j < len(other) {
		switch {
		case set[i] < other[j]:
			union = append(union, set[i])
			i++
		case other[j] < set[i]:
			union = append(union, other[j])
			j++
		default:
			union = append(union, set[i])
			i++
			j++
		}
	}

	switch {
	case i < len(set):
		union = append(union, set[i:]...)
	case j < len(other):
		union = append(union, other[j:]...)
	}

	return union
}

func NewSortedSetFromSlice[T constraints.Ordered](slice []T) SortedSet[T] {
	set := make(SortedSet[T], len(slice))
	copy(set, slice)
	slices.Sort(set)
	return slices.Compact(set)
}
