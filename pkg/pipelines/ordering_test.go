package pipelines

import (
	"testing"

	"github.com/gitops-tools/apps-scanner/test"
	"github.com/google/go-cmp/cmp"
)

func TestOrderEnvironments(t *testing.T) {
	environmentTests := []struct {
		name         string
		environments []environment
		want         []string
	}{
		{
			name:         "no environments",
			environments: []environment{},
			want:         []string{},
		},
		{
			name:         "single environment",
			environments: []environment{{name: "first"}},
			want:         []string{"first"},
		},
		{
			name:         "two environments",
			environments: []environment{{name: "first"}, {name: "second", after: "first"}},
			want:         []string{"first", "second"},
		},
		{
			name:         "three environments",
			environments: []environment{{name: "first"}, {name: "second", after: "first"}, {name: "third", after: "second"}},
			want:         []string{"first", "second", "third"},
		},
		{
			name:         "three environments, two depending on one",
			environments: []environment{{name: "first"}, {name: "second", after: "first"}, {name: "third", after: "first"}},
			want:         []string{"first", "second", "third"},
		},
	}

	for _, tt := range environmentTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OrderEnvironments(tt.environments)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("failed to order environments:\n%s", diff)
			}
		})
	}
}

func TestOrderEnvironments_errors(t *testing.T) {
	environmentTests := []struct {
		name         string
		environments []environment
		wantErr      string
	}{
		{
			// TODO: maybe improve on this error message?
			name:         "duplicate environments",
			environments: []environment{{name: "first"}, {name: "first"}},
			wantErr:      "'first' is already known",
		},
		{
			name:         "missing after environment",
			environments: []environment{{name: "first"}, {name: "third", after: "second"}},
			wantErr:      "reference to unknown environment \"second\"",
		},
	}

	for _, tt := range environmentTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := OrderEnvironments(tt.environments)
			test.AssertErrorMatch(t, tt.wantErr, err)
		})
	}
}
