package applications

import (
	"fmt"
	"sort"

	"github.com/gitops-tools/apps-scanner/pkg/sets"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

const (
	// AppLabel is the Kubernetes recommended label to indicate that a component
	// is part of an application.
	AppLabel = "app.kubernetes.io/part-of"

	partOfLabel    = AppLabel
	instanceLabel  = "app.kubernetes.io/instance"
	nameLabel      = "app.kubernetes.io/name"
	componentLabel = "app.kubernetes.io/component"

	kustomizationName      = "kustomize.toolkit.fluxcd.io/name"
	kustomizationNamespace = "kustomize.toolkit.fluxcd.io/namespace"
)

// Application represents a discovered deployment group.
type Application struct {
	Name           string
	Instances      []string
	Components     []string
	Parents        []Application
	Kustomizations []types.NamespacedName
}

// Parser parses the labels and annotations on runtime Objects and extracts apps
// from the labels.
type Parser struct {
	Accessor meta.MetadataAccessor
	apps     map[string]discoveryApplication
}

// NewParser creates and returns a new Parser ready for use.
func NewParser() *Parser {
	return &Parser{
		Accessor: meta.NewAccessor(),
		apps:     make(map[string]discoveryApplication),
	}
}

// Add a list of objects to the parser.
//
// The list should be a List type, e.g. PodList, DeploymentList etc.
func (p *Parser) Add(list runtime.Object) error {
	return meta.EachListItem(list, func(obj runtime.Object) error {
		l, err := p.Accessor.Labels(obj)
		if err != nil {
			return fmt.Errorf("failed to get labels from %v: %w", obj, err)
		}
		appName := l[nameLabel]
		if appName == "" {
			return nil
		}
		a, ok := p.apps[appName]
		if !ok {
			a = discoveryApplication{
				name:           appName,
				instances:      sets.New[string](),
				parents:        sets.New[string](),
				components:     sets.New[string](),
				kustomizations: sets.New[types.NamespacedName](),
			}
		}
		// TODO: this should check for the presence of these labels!
		a.instances.Insert(l[instanceLabel])
		a.components.Insert(l[componentLabel])
		if v := l[partOfLabel]; v != "" {
			a.parents.Insert(l[partOfLabel])
		}
		if nn := kustomizationRefFromLabels(l); nn != nil {
			a.kustomizations.Insert(*nn)
		}
		p.apps[appName] = a
		return nil
	})
}

// Applications returns the Applications that were discovered during the parsing
// process.
func (p *Parser) Applications() []Application {
	apps := map[string]Application{}
	for _, v := range p.apps {
		app, ok := apps[v.name]
		if !ok {
			app = Application{
				Name: v.name,
			}
		}

		app.Instances = v.instances.List()
		app.Components = v.components.List()
		app.Kustomizations = v.kustomizations.List()

		// Scan the parents of the app, link the child to the parent Application
		// object.
		//
		// If the parent is not known at this time, add it to the list of known
		// Applications.
		// TODO: simplify this to a two-pass check?

		for _, p := range v.parents.List() {
			if a, ok := apps[p]; ok {
				app.Parents = append(app.Parents, a)
			} else {
				newApp := Application{Name: p}
				apps[p] = newApp
				app.Parents = append(app.Parents, newApp)
			}
		}
		apps[app.Name] = app
	}

	res := []Application{}
	for _, v := range apps {
		res = append(res, v)
	}
	// Sorting to ensure that the tests are stable
	sort.Slice(res, func(i, j int) bool { return res[i].Name < res[j].Name })
	return res

}

// discoveryApplication is a temporary holding type to simplify identification
// of services/environments/kustomizations.
type discoveryApplication struct {
	name           string
	instances      sets.Set[string]
	parents        sets.Set[string]
	components     sets.Set[string]
	kustomizations sets.Set[types.NamespacedName]
}

func findApplication(name string, apps []Application) *Application {
	for _, v := range apps {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

func kustomizationRefFromLabels(m map[string]string) *types.NamespacedName {
	name, ok := m[kustomizationName]
	if !ok {
		return nil
	}
	ns, ok := m[kustomizationNamespace]
	if !ok {
		return nil
	}
	return &types.NamespacedName{Name: name, Namespace: ns}
}
