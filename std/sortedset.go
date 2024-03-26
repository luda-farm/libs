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
func (set *SortedSet[T]) AddAll(t []T) bool {
	modified := false
	for _, element := range t {
		if set.Add(element) {
			modified = true
		}
	}
	return modified
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

func (set SortedSet[T]) Complement(other SortedSet[T]) SortedSet[T] {
	var complement SortedSet[T]
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
	var intersection SortedSet[T]
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

func (set SortedSet[T]) Union(other SortedSet[T]) SortedSet[T] {
	if len(set) < len(other) {
		set, other = other, set
	}
	union := SliceToSortedSet(set)
	for _, t := range other {
		union.Add(t)
	}
	return union
}

func SliceToSortedSet[T constraints.Ordered](s []T) SortedSet[T] {
	var set SortedSet[T]
	set.AddAll(s)
	return set
}
