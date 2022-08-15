package sets

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"k8s.io/apimachinery/pkg/types"
)

func TestGenericSet_New(t *testing.T) {
	ss := New[types.NamespacedName](
		nsn("test", "test-ns"),
		nsn("test2", "test-ns"),
		nsn("test", "test-ns"),
	)

	want := []types.NamespacedName{
		nsn("test", "test-ns"),
		nsn("test2", "test-ns"),
	}
	if diff := cmp.Diff(want, ss.List()); diff != "" {
		t.Fatalf("failed to create set with items:\n%s", diff)
	}
}

func TestGenericSet_Insert(t *testing.T) {
	ss := New[types.NamespacedName]()

	ss = ss.Insert(nsn("test", "test-ns"))
	ss = ss.Insert(nsn("test2", "test-ns"), nsn("test", "test-ns"))

	want := []types.NamespacedName{
		nsn("test2", "test-ns"),
		nsn("test", "test-ns"),
	}
	if diff := cmp.Diff(want, ss.List(), cmpopts.SortSlices(
		func(x, y types.NamespacedName) bool {
			return strings.Compare(x.String(), y.String()) < 0
		})); diff != "" {
		t.Fatalf("failed to create set with items:\n%s", diff)
	}
}

func TestGenericSet_List(t *testing.T) {
	ss := New[types.NamespacedName]()

	if ss.List() != nil {
		t.Fatal("list did not return nil")
	}
}

func nsn(name, ns string) types.NamespacedName {
	return types.NamespacedName{
		Name:      name,
		Namespace: ns,
	}
}
