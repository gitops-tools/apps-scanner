package repositories

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"k8s.io/apimachinery/pkg/types"

	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
)

func TestParser(t *testing.T) {
	discoverTests := []struct {
		name  string
		items [][]sourcev1.GitRepository
		want  []Repository
	}{
		{
			name:  "empty list",
			items: [][]sourcev1.GitRepository{},
			want:  []Repository{},
		},
		{
			name: "single repository",
			items: [][]sourcev1.GitRepository{
				{
					makeGitRepository(withURL("git@github.com:demo/demo-repo.git"), named("test", "test-ns"), branch("main")),
				},
			},
			want: []Repository{
				Repository{
					URL: "git@github.com:demo/demo-repo.git",
					Refs: []RepositoryRef{
						{NamespacedName: types.NamespacedName{Name: "test", Namespace: "test-ns"}, Ref: sourcev1.GitRepositoryRef{Branch: "main"}},
					},
				},
			},
		},
		{
			name: "multiple refs for a repository",
			items: [][]sourcev1.GitRepository{
				{
					makeGitRepository(withURL("git@github.com:demo/demo-repo.git"), named("test1", "test-ns"), branch("main")),
					makeGitRepository(withURL("git@github.com:demo/demo-repo.git"), named("test2", "test-ns"), branch("production")),
				},
			},
			want: []Repository{
				{
					URL: "git@github.com:demo/demo-repo.git",
					Refs: []RepositoryRef{
						{
							NamespacedName: types.NamespacedName{Name: "test1", Namespace: "test-ns"},
							Ref:            sourcev1.GitRepositoryRef{Branch: "main"},
						},
						{
							NamespacedName: types.NamespacedName{Name: "test2", Namespace: "test-ns"},
							Ref:            sourcev1.GitRepositoryRef{Branch: "production"},
						},
					},
				},
			},
		},
		{
			name: "multiple repositories",
			items: [][]sourcev1.GitRepository{
				{
					makeGitRepository(withURL("git@github.com:demo/demo-repo1.git"), named("test1", "test-ns"), branch("main")),
					makeGitRepository(withURL("git@github.com:demo/demo-repo2.git"), named("test2", "test-ns"), branch("main")),
				},
			},
			want: []Repository{
				{
					URL: "git@github.com:demo/demo-repo1.git",
					Refs: []RepositoryRef{
						{
							NamespacedName: types.NamespacedName{Name: "test1", Namespace: "test-ns"},
							Ref:            sourcev1.GitRepositoryRef{Branch: "main"},
						},
					},
				},
				{
					URL: "git@github.com:demo/demo-repo2.git",
					Refs: []RepositoryRef{
						{
							NamespacedName: types.NamespacedName{Name: "test2", Namespace: "test-ns"},
							Ref:            sourcev1.GitRepositoryRef{Branch: "main"},
						},
					},
				},
			},
		},
	}
	for _, tt := range discoverTests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			for _, v := range tt.items {
				repos := &sourcev1.GitRepositoryList{
					Items: v,
				}

				err := p.Add(repos)
				if err != nil {
					t.Fatal(err)
				}
			}
			kustomizations := p.Repositories()
			if diff := cmp.Diff(tt.want, kustomizations, sortOpts()...); diff != "" {
				t.Fatalf("failed discovery:\n%s", diff)
			}
		})
	}
}

func branch(b string) func(*sourcev1.GitRepository) {
	return func(o *sourcev1.GitRepository) {
		o.Spec.Reference = &sourcev1.GitRepositoryRef{
			Branch: b,
		}
	}
}

func named(name, namespace string) func(*sourcev1.GitRepository) {
	return func(o *sourcev1.GitRepository) {
		o.ObjectMeta.Name = name
		o.ObjectMeta.Namespace = namespace
	}
}

func withURL(u string) func(*sourcev1.GitRepository) {
	return func(o *sourcev1.GitRepository) {
		o.Spec.URL = u
	}
}

func makeGitRepository(opts ...func(*sourcev1.GitRepository)) sourcev1.GitRepository {
	p := sourcev1.GitRepository{}
	for _, o := range opts {
		o(&p)
	}
	return p
}

func sortOpts() []cmp.Option {
	return []cmp.Option{
		cmpopts.SortSlices(
			func(x, y string) bool {
				return strings.Compare(x, y) < 0
			}),
		cmpopts.SortSlices(
			func(x, y RepositoryRef) bool {
				return strings.Compare(x.NamespacedName.String(), y.NamespacedName.String()) < 0
			}),
		cmpopts.SortSlices(
			func(x, y Repository) bool {
				return strings.Compare(x.URL, y.URL) < 0
			}),
	}
}
