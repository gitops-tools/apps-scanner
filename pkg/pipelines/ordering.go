package pipelines

import (
	"fmt"

	"github.com/heimdalr/dag"
)

// OrderEnvironments takes a set pairs of named environments and their preceeding environment and
// calculates the ordering.
//
// If more than one environment has no dependent environments, the ordering will
// be non-deterministic, all dependent environments will be in order.
func OrderEnvironments(o []environment) ([]string, error) {
	d := dag.NewDAG()
	environmentsToIDs := map[string]string{}
	idsToEnvironments := map[string]string{}

	// walk the list adding in the environments.
	for _, v := range o {
		vtx, err := d.AddVertex(v.name)
		if err != nil {
			return nil, fmt.Errorf("failed to order pipeline environments: %w", err)
		}
		environmentsToIDs[v.name] = vtx
		idsToEnvironments[vtx] = v.name
	}

	// walk the list adding in the links between environments.
	for _, v := range o {
		if v.after == "" {
			continue
		}
		after, ok := environmentsToIDs[v.after]
		if !ok {
			return nil, fmt.Errorf("reference to unknown environment %q", v.after)
		}
		d.AddEdge(after, environmentsToIDs[v.name])
	}

	roots := d.GetRoots()
	result := []string{}
	for k := range roots {
		// Add the root as the first element in the environments
		// because it's not a descendant...
		result = append(result, idsToEnvironments[k])
		deps, err := d.GetOrderedDescendants(k)
		if err != nil {
			return nil, fmt.Errorf("failed to order pipeline environments: %w", err)
		}
		for _, v := range deps {
			result = append(result, idsToEnvironments[v])
		}
	}
	return result, nil
}
