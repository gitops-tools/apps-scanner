package sets

import (
	"sort"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
)

// NamespacedNames is a set of types.NamespacedName to simplify Kustomization
// detection
type NamespacedNames map[types.NamespacedName]sets.Empty

// NewNamespacedNames creates and returns a new set of NamespacedNames.
func NewNamespacedNames(items ...types.NamespacedName) NamespacedNames {
	ss := NamespacedNames{}
	return ss.Insert(items...)
}

func (s NamespacedNames) Insert(items ...types.NamespacedName) NamespacedNames {
	for _, item := range items {
		s[item] = sets.Empty{}
	}
	return s
}

// List returns the contents as a sorted slice.
// WARNING: This is suboptimal as it's stringifying on each comparison, there
// aren't expected to be a huge number of NamespacedNames.
func (s NamespacedNames) List() []types.NamespacedName {
	if len(s) == 0 {
		return nil
	}
	res := []types.NamespacedName{}
	for key := range s {
		res = append(res, key)
	}
	sort.Slice(res, func(i, j int) bool { return res[i].String() < res[j].String() })
	return res
}
