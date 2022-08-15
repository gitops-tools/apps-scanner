package sets

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"k8s.io/apimachinery/pkg/types"
)

func TestSet_New(t *testing.T) {
	ss := New[types.NamespacedName](
		nsn("test"),
		nsn("test2"),
		nsn("test"),
	)

	want := []types.NamespacedName{
		nsn("test"),
		nsn("test2"),
	}
	if diff := cmp.Diff(want, ss.List(), sortNamespaces()); diff != "" {
		t.Fatalf("failed to create set with items:\n%s", diff)
	}
}

func TestSet_Len_when_empty(t *testing.T) {
	ss := New[types.NamespacedName]()

	if l := ss.Len(); l != 0 {
		t.Fatalf("Len() got %v, want 0", l)
	}
}

func TestSet_List_when_empty(t *testing.T) {
	ss := New[types.NamespacedName]()

	if ss.List() != nil {
		t.Fatal("list did not return nil")
	}
}

func TestSet_Insert(t *testing.T) {
	ss := New[types.NamespacedName]()

	ss.Insert(nsn("test"))
	ss.Insert(nsn("test2"), nsn("test"))

	want := []types.NamespacedName{
		nsn("test2"),
		nsn("test"),
	}
	if diff := cmp.Diff(want, ss.List(), sortNamespaces()); diff != "" {
		t.Fatalf("failed to create set with items:\n%s", diff)
	}
	if l := ss.Len(); l != 2 {
		t.Fatalf("Len() got %v, want 2", l)
	}
}

func TestSet_Delete(t *testing.T) {
	ss := New[types.NamespacedName](
		nsn("test"),
		nsn("test2"),
	)
	ss.Delete(nsn("test"))

	want := []types.NamespacedName{
		nsn("test2"),
	}
	if diff := cmp.Diff(want, ss.List()); diff != "" {
		t.Fatalf("failed to create set with items:\n%s", diff)
	}
}

func TestSet_Has(t *testing.T) {
	ss := New[types.NamespacedName](
		nsn("test"),
	)

	if ss.Has(nsn("test2")) {
		t.Errorf("Has(%v) got true", nsn("test2"))
	}
	if !ss.Has(nsn("test")) {
		t.Errorf("Has(%v) got false", nsn("test"))
	}
}

func TestSet_HasAll(t *testing.T) {
	ss := New[types.NamespacedName](
		nsn("test"),
		nsn("test2"),
	)
	hasAllTests := []struct {
		all    []types.NamespacedName
		hasAll bool
	}{
		{
			all: []types.NamespacedName{
				nsn("test"),
				nsn("test2"),
			},
			hasAll: true,
		},
		{
			all: []types.NamespacedName{
				nsn("test"),
				nsn("test2"),
				nsn("test3"),
			},
			hasAll: false,
		},
		{
			all: []types.NamespacedName{
				nsn("test"),
			},
			hasAll: true,
		},
	}

	for i, tt := range hasAllTests {
		t.Run(fmt.Sprintf("hasAll_%d", i), func(t *testing.T) {
			if h := ss.HasAll(tt.all...); h != tt.hasAll {
				t.Fatalf("ss.HasAll(%v) got %v, want %v", tt.all, h, tt.hasAll)
			}
		})
	}
}

func TestSet_IsSuperset(t *testing.T) {
	ss := New[types.NamespacedName](
		nsn("test"),
		nsn("test2"),
	)
	supersetTests := []struct {
		set      []types.NamespacedName
		superset bool
	}{
		{
			set: []types.NamespacedName{
				nsn("test"),
				nsn("test2"),
			},
			superset: true,
		},
		{
			set: []types.NamespacedName{
				nsn("test"),
				nsn("test2"),
				nsn("test3"),
			},
			superset: false,
		},
		{
			set: []types.NamespacedName{
				nsn("test"),
			},
			superset: true,
		},
	}

	for i, tt := range supersetTests {
		t.Run(fmt.Sprintf("superSet_%d", i), func(t *testing.T) {
			if h := ss.IsSuperset(New(tt.set...)); h != tt.superset {
				t.Fatalf("ss.IsSuperset(%v) got %v, want %v", tt.set, h, tt.superset)
			}
		})
	}
}

func TestSet_HasAny(t *testing.T) {
	ss := New[types.NamespacedName](
		nsn("test"),
		nsn("test2"),
	)
	hasAnyTests := []struct {
		anyNames []types.NamespacedName
		hasAny   bool
	}{
		{
			anyNames: []types.NamespacedName{
				nsn("test"),
				nsn("test2"),
			},
			hasAny: true,
		},
		{
			anyNames: []types.NamespacedName{
				nsn("test3"),
				nsn("test4"),
			},
			hasAny: false,
		},
		{
			anyNames: []types.NamespacedName{
				nsn("test"),
			},
			hasAny: true,
		},
	}

	for i, tt := range hasAnyTests {
		t.Run(fmt.Sprintf("hasAny_%d", i), func(t *testing.T) {
			if h := ss.HasAny(tt.anyNames...); h != tt.hasAny {
				t.Fatalf("ss.HasAny(%v) got %v, want %v", tt.anyNames, h, tt.hasAny)
			}
		})
	}
}

func TestSet_Difference(t *testing.T) {
	a := New[types.NamespacedName](
		nsn("test1"),
		nsn("test2"),
		nsn("test3"),
	)
	b := New[types.NamespacedName](
		nsn("test1"),
		nsn("test2"),
		nsn("test4"),
		nsn("test5"),
	)

	c := a.Difference(b)
	d := b.Difference(a)
	if len(c) != 1 {
		t.Errorf("Expected len=1: %d", len(c))
	}
	if !c.Has(nsn("test3")) {
		t.Errorf("Unexpected contents: %#v", c.List())
	}
	if len(d) != 2 {
		t.Errorf("Expected len=2: %d", len(d))
	}
	if !d.Has(nsn("test4")) || !d.Has(nsn("test5")) {
		t.Errorf("Unexpected contents: %#v", d.List())
	}
}

func TestSet_Union(t *testing.T) {
	tests := []struct {
		s1   Set[types.NamespacedName]
		s2   Set[types.NamespacedName]
		want Set[types.NamespacedName]
	}{
		{
			s1:   New[types.NamespacedName](nsn("test1"), nsn("test2"), nsn("test3"), nsn("test4")),
			s2:   New[types.NamespacedName](nsn("test3"), nsn("test4"), nsn("test5"), nsn("test6")),
			want: New[types.NamespacedName](nsn("test1"), nsn("test2"), nsn("test3"), nsn("test4"), nsn("test5"), nsn("test6")),
		},
		{
			s1:   New[types.NamespacedName](nsn("test1"), nsn("test2"), nsn("test3"), nsn("test4")),
			s2:   New[types.NamespacedName](),
			want: New[types.NamespacedName](nsn("test1"), nsn("test2"), nsn("test3"), nsn("test4")),
		},
		{
			s1:   New[types.NamespacedName](),
			s2:   New[types.NamespacedName](nsn("test1"), nsn("test2"), nsn("test3"), nsn("test4")),
			want: New[types.NamespacedName](nsn("test1"), nsn("test2"), nsn("test3"), nsn("test4")),
		},
		{
			s1:   New[types.NamespacedName](),
			s2:   New[types.NamespacedName](),
			want: New[types.NamespacedName](),
		},
	}

	for _, test := range tests {
		union := test.s1.Union(test.s2)
		if union.Len() != test.want.Len() {
			t.Errorf("Expected union.Len()=%d but got %d", test.want.Len(), union.Len())
		}

		if !union.Equal(test.want) {
			t.Errorf("Expected union.Equal(expected) but not true.  union:%v want:%v", union.List(), test.want.List())
		}
	}
}

func TestSet_Intersection(t *testing.T) {
	tests := []struct {
		s1   Set[types.NamespacedName]
		s2   Set[types.NamespacedName]
		want Set[types.NamespacedName]
	}{
		{
			New[types.NamespacedName](nsn("test1"), nsn("test2"), nsn("test3"), nsn("test4")),
			New[types.NamespacedName](nsn("test3"), nsn("test4"), nsn("test5"), nsn("test6")),
			New[types.NamespacedName](nsn("test3"), nsn("test4")),
		},
		{
			New[types.NamespacedName](nsn("test1"), nsn("test2"), nsn("test3"), nsn("test4")),
			New[types.NamespacedName](nsn("test1"), nsn("test2"), nsn("test3"), nsn("test4")),
			New[types.NamespacedName](nsn("test1"), nsn("test2"), nsn("test3"), nsn("test4")),
		},
		{
			New[types.NamespacedName](nsn("test1"), nsn("test2"), nsn("test3"), nsn("test4")),
			New[types.NamespacedName](),
			New[types.NamespacedName](),
		},
		{
			New[types.NamespacedName](),
			New[types.NamespacedName](nsn("test1"), nsn("test2"), nsn("test3"), nsn("test4")),
			New[types.NamespacedName](),
		},
		{
			New[types.NamespacedName](),
			New[types.NamespacedName](),
			New[types.NamespacedName](),
		},
	}

	for _, test := range tests {
		intersection := test.s1.Intersection(test.s2)
		if intersection.Len() != test.want.Len() {
			t.Errorf("Expected intersection.Len()=%d but got %d", test.want.Len(), intersection.Len())
		}

		if !intersection.Equal(test.want) {
			t.Errorf("Expected intersection.Equal(want) but not true.  intersection:%v want:%v", intersection.List(), test.want.List())
		}
	}
}

func nsn(name string) types.NamespacedName {
	return types.NamespacedName{
		Name:      name,
		Namespace: "test-ns",
	}
}

func sortNamespaces() cmp.Option {
	return cmpopts.SortSlices(
		func(x, y types.NamespacedName) bool {
			return strings.Compare(x.String(), y.String()) < 0
		})
}
