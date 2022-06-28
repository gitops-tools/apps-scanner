package visualise

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gitops-tools/apps-scanner/pkg/applications"
)

func TestNewDOT(t *testing.T) {
	g := NewDOT([]applications.Application{makeApplication(), applications.Application{Name: "billing-system"}})

	want := `digraph  {
	
	n2[label="billing-system"];
	n1[label="frontend"];
	n1->n2;
	
}
`
	if diff := cmp.Diff(want, g.String()); diff != "" {
		t.Fatalf("failed to visualise applications: %s\n", diff)
	}
}

func TestNewDOT_with_multi_level(t *testing.T) {
	g := NewDOT([]applications.Application{
		makeApplication(),
		applications.Application{Name: "billing-system"},
		makeApplication(func(a *applications.Application) {
			a.Name = "backend"
			a.Parents = []applications.Application{{Name: "frontend"}}
		})})

	want := `digraph  {
	
	n3[label="backend"];
	n2[label="billing-system"];
	n1[label="frontend"];
	n3->n1;
	n1->n2;
	
}
`
	if diff := cmp.Diff(want, g.String()); diff != "" {
		t.Fatalf("failed to visualise applications: %s\n", diff)
	}
}

func makeApplication(opts ...func(*applications.Application)) applications.Application {
	a := applications.Application{
		Name:           "frontend",
		Instances:      []string{"staging", "production"},
		Components:     []string{"database", "web"},
		Parents:        []applications.Application{{Name: "billing-system"}},
		Kustomizations: []types.NamespacedName{{Name: "repo-main", Namespace: "flux-system"}},
	}
	for _, opt := range opts {
		opt(&a)
	}
	return a
}
