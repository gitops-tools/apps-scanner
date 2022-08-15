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
