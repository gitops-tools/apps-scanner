package repositories

import (
	"sort"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"

	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
)

// Repository is a summarised version of a list of GitRepository objects.
type Repository struct {
	URL  string
	Refs []RepositoryRef
}

// RepositoryRef indicates which ref a specific GitRepository is tracking.
type RepositoryRef struct {
	types.NamespacedName
	Ref sourcev1.GitRepositoryRef
}

// Parser parses a list of Kustomization objects and extracts information from
// them.
type Parser struct {
	Accessor meta.MetadataAccessor
	// map of URL -> discovered data
	repositories map[string]discoveryRepository
}

// NewParser creates and returns a new Parser ready for use.
func NewParser() *Parser {
	return &Parser{
		Accessor:     meta.NewAccessor(),
		repositories: make(map[string]discoveryRepository),
	}
}

// Add a list of Kustomization objects to be parsed.
func (p *Parser) Add(list *sourcev1.GitRepositoryList) error {
	for _, repo := range list.Items {
		k, ok := p.repositories[repo.Spec.URL]
		if !ok {
			k = discoveryRepository{
				refs: newRepositoryRefSet(),
			}
		}
		k.refs.Insert(RepositoryRef{NamespacedName: namespacedNameFromRepository(repo), Ref: *repo.Spec.Reference})
		p.repositories[repo.Spec.URL] = k
	}
	return nil
}

func (p *Parser) Repositories() []Repository {
	res := []Repository{}
	for k, v := range p.repositories {
		res = append(res, Repository{
			URL:  k,
			Refs: v.refs.List(),
		})
	}
	return res
}

// discoveryRepository is a temporary holding type to simplify identification
// of repositories.
type discoveryRepository struct {
	refs repositoryRefSet
}

func namespacedNameFromRepository(o sourcev1.GitRepository) types.NamespacedName {
	return types.NamespacedName{
		Name:      o.GetName(),
		Namespace: o.GetNamespace(),
	}
}

// repositoryRefSet is a set of RepositoryRefs to simplify  discovery.
type repositoryRefSet map[RepositoryRef]sets.Empty

func newRepositoryRefSet(items ...RepositoryRef) repositoryRefSet {
	ss := repositoryRefSet{}
	return ss.Insert(items...)
}

func (s repositoryRefSet) Insert(items ...RepositoryRef) repositoryRefSet {
	for _, item := range items {
		s[item] = sets.Empty{}
	}
	return s
}

// List returns the contents as a sorted slice.
// WARNING: This is suboptimal as it's stringifying on each comparison, there
// aren't expected to be a huge number of repository refs.
func (s repositoryRefSet) List() []RepositoryRef {
	if len(s) == 0 {
		return nil
	}
	res := []RepositoryRef{}
	for key := range s {
		res = append(res, key)
	}
	sort.Slice(res, func(i, j int) bool { return res[i].String() < res[j].String() })
	return res
}
