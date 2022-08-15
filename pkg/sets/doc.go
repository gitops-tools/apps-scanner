// This is an implementation of the k8s util sets package using Go 1.18
// generics.
//
// It implements the same API as k8s.io/apimachinery/pkg/util/sets
//
// Instead of creating a sets.StringSet you can use
// sets.New[string]() and the same functionality is available.
package sets
