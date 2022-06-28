package visualise

import (
	"github.com/emicklei/dot"

	"github.com/gitops-tools/apps-scanner/pkg/applications"
)

// NewDOT converts a set of Applications to a graph of dependencies.
func NewDOT(apps []applications.Application) *dot.Graph {
	g := dot.NewGraph(dot.Directed)
	for _, app := range apps {
		g.Node(app.Name)
	}

	for _, app := range apps {
		for _, p := range app.Parents {
			parentNode := g.Node(p.Name)
			appNode := g.Node(app.Name)
			g.Edge(appNode, parentNode)
		}
	}

	return g
}
