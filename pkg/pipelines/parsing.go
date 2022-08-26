package pipelines

import (
	"fmt"
	"sort"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	// PipelineNameLabel is a label that indicates a resource is part of a Pipeline.
	PipelineNameLabel = "gitops.pro/pipeline"

	// PipelineEnvironmentabel is a label that indicates which stage a component
	// is from within a Pipeline.
	PipelineEnvironmentLabel = "gitops.pro/pipeline-environment"

	// PipelineEnvironmentAfterLabel is a label that indicates which stage a
	// component follows within a pipeline.
	PipelineEnvironmentAfterLabel = "gitops.pro/pipeline-after"
)

// Parser parses the labels and annotations on runtime Objects and extracts apps
// from the labels.
type Parser struct {
	accessor  meta.MetadataAccessor
	discovery map[string]discoveryPipeline
	Labels    Labels
}

// Labels configures the set of labels to examine resources for.
type Labels struct {
	Pipeline    string
	Environment string
	After       string
}

// WithLabels is a functional option for configuring the Parser with a set of
// labels.
func WithLabels(pipeline, environment, after string) func(*Parser) {
	return func(p *Parser) {
		p.Labels.Pipeline = pipeline
		p.Labels.Environment = environment
		p.Labels.After = after
	}
}

// NewParser creates and returns a new Parser ready for use.
func NewParser(opts ...func(*Parser)) *Parser {
	p := &Parser{
		accessor:  meta.NewAccessor(),
		discovery: make(map[string]discoveryPipeline),
		Labels: Labels{
			Pipeline:    PipelineNameLabel,
			Environment: PipelineEnvironmentLabel,
			After:       PipelineEnvironmentAfterLabel,
		},
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Add accepts a list of objects and records them for parsing with the Pipelines
// method.
func (p *Parser) Add(list []runtime.Object) error {
	for _, obj := range list {
		l, err := p.accessor.Labels(obj)
		if err != nil {
			return fmt.Errorf("failed to get labels from %v: %w", obj, err)
		}
		pipelineName := l[p.Labels.Pipeline]
		if pipelineName == "" {
			continue
		}
		a, ok := p.discovery[pipelineName]
		if !ok {
			a = discoveryPipeline{
				name:         pipelineName,
				environments: newEnvironmentSet(),
			}
		}

		if n, ok := l[p.Labels.Environment]; ok {
			after := l[p.Labels.After]
			a.environments.Insert(environment{name: n, after: after})
		}
		p.discovery[pipelineName] = a
	}

	return nil
}

// Pipelines returns the discovered pipelines.
//
// The environments are ordered based on the configuration of pipeline after
// labels.
func (p *Parser) Pipelines() ([]Pipeline, error) {
	res := []Pipeline{}
	for _, v := range p.discovery {
		ordered, err := OrderEnvironments(v.environments.List())
		if err != nil {
			return nil, fmt.Errorf("failed parsing pipeline %q: %w", v.name, err)
		}
		p := Pipeline{
			Name:         v.name,
			Environments: ordered,
		}
		res = append(res, p)
	}

	// Sorting to ensure that the tests are stable
	sort.Slice(res, func(i, j int) bool { return res[i].Name < res[j].Name })
	return res, nil
}

// Pipeline is a Continuous-Delivery pipeline with a sequence of environments
// that an application change passes through.
type Pipeline struct {
	Name         string
	Environments []string
}

type discoveryPipeline struct {
	name         string
	environments environmentSet
}

type environment struct {
	name  string
	after string
}

func (e environment) String() string {
	return e.name
}

// environmentSet is a set of environments to simplify environment ordering.
type environmentSet map[environment]sets.Empty

func newEnvironmentSet(items ...environment) environmentSet {
	ss := environmentSet{}
	return ss.Insert(items...)
}

func (s environmentSet) Insert(items ...environment) environmentSet {
	for _, item := range items {
		s[item] = sets.Empty{}
	}
	return s
}

// List returns the contents as a sorted slice.
// WARNING: This is suboptimal as it's stringifying on each comparison, there
// aren't expected to be a huge number of environments.
func (s environmentSet) List() []environment {
	if len(s) == 0 {
		return nil
	}
	res := []environment{}
	for key := range s {
		res = append(res, key)
	}
	sort.Slice(res, func(i, j int) bool { return res[i].String() < res[j].String() })
	return res
}
