package sets

import (
	"k8s.io/apimachinery/pkg/util/sets"
)

// Set is a generic Set implementation with some basic Set methods.
type Set[T comparable] map[T]sets.Empty

// New returns a new Set of the given items.
func New[T comparable](items ...T) Set[T] {
	ss := Set[T](make(map[T]sets.Empty))
	return ss.Insert(items...)
}

// Insert inserts items into the set and returns an updated Set.
func (s Set[T]) Insert(items ...T) Set[T] {
	for _, item := range items {
		s[item] = sets.Empty{}
	}

	return s
}

// Has returns true if item is contained in the set.
func (s Set[T]) Has(item T) bool {
	_, contained := s[item]

	return contained
}

// Delete removes the item, and returns an updated Set.
func (s Set[T]) Delete(items ...T) Set[T] {
	for _, item := range items {
		delete(s, item)
	}

	return s
}

// List returns a slice with all the items.
//
// These items are not sorted.
func (s Set[T]) List() []T {
	if len(s) == 0 {
		return nil
	}

	res := make([]T, 0, len(s))
	for key := range s {
		res = append(res, key)
	}

	return res
}

// HasAll returns true if all items are contained in the set.
func (s Set[T]) HasAll(items ...T) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}

	return true
}

// HasAny returns true if any of the items are contained in the set.
func (s Set[T]) HasAny(items ...T) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}

	return false
}

// Difference returns a set of objects that are not in s2
// For example:
// s1 = {a1, a2, a3}
// s2 = {a1, a2, a4, a5}
// s1.Difference(s2) = {a3}
// s2.Difference(s1) = {a4, a5}
func (s Set[T]) Difference(s2 Set[T]) Set[T] {
	result := New[T]()
	for key := range s {
		if !s2.Has(key) {
			result.Insert(key)
		}
	}

	return result
}

// Union returns a new set which includes items in either s1 or s2.
// For example:
// s1 = {a1, a2}
// s2 = {a3, a4}
// s1.Union(s2) = {a1, a2, a3, a4}
// s2.Union(s1) = {a1, a2, a3, a4}
func (s1 Set[T]) Union(s2 Set[T]) Set[T] {
	result := New[T]()
	for key := range s1 {
		result.Insert(key)
	}
	for key := range s2 {
		result.Insert(key)
	}

	return result
}

// Intersection returns a new set which includes the item in BOTH s and s2
// For example:
// s = {a1, a2}
// s2 = {a2, a3}
// s.Intersection(s2) = {a2}
func (s Set[T]) Intersection(s2 Set[T]) Set[T] {
	var walk, other Set[T]
	result := New[T]()
	if s.Len() < s2.Len() {
		walk = s
		other = s2
	} else {
		walk = s2
		other = s
	}
	for key := range walk {
		if other.Has(key) {
			result.Insert(key)
		}
	}

	return result
}

// Len returns the size of the set.
func (s Set[T]) Len() int {
	return len(s)
}

// Equal returns true if s is equal (as a set) to s2.
//
// Two sets are equal if their membership is identical.
func (s Set[T]) Equal(s2 Set[T]) bool {
	return len(s) == len(s2) && s.IsSuperset(s2)
}

// IsSuperset returns true if s is a superset of s2.
func (s Set[T]) IsSuperset(s2 Set[T]) bool {
	for item := range s2 {
		if !s.Has(item) {
			return false
		}
	}
	return true
}
