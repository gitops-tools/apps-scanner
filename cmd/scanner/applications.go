package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gitops-tools/apps-scanner/pkg/applications"
	"github.com/gitops-tools/apps-scanner/pkg/visualise"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func newApplicationsCmd(cl client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "applications",
		Short: "List applications in the cluster",
		RunE:  listApplications(cl),
	}

	cmd.Flags().String("graphviz-file", "", "Write a graphviz of the discovered applications")
	cobra.CheckErr(viper.BindPFlag("graphviz-file", cmd.Flags().Lookup("graphviz-file")))

	return cmd
}

func listApplications(cl client.Client) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting to scan for applications")
		deploymentList := &appsv1.DeploymentList{}
		err := cl.List(context.Background(), deploymentList, client.HasLabels([]string{applications.AppLabel}))
		if err != nil {
			return fmt.Errorf("failed to list deployments: %w", err)
		}
		fmt.Printf("found %d deployments\n", len(deploymentList.Items))

		p := applications.NewParser()
		if err := p.Add(deploymentList); err != nil {
			return fmt.Errorf("failed to discover applications: %w", err)
		}

		apps := p.Applications()

		for _, parent := range parentApps(apps) {
			fmt.Printf("application %s\n", parent.Name)
			for _, app := range childApps(apps, parent.Name) {
				fmt.Printf("  child app: %s\n", app.Name)
				fmt.Println("     instances:")
				for _, e := range app.Instances {
					fmt.Printf("         %s\n", e)
				}
				fmt.Println("     components:")
				for _, s := range app.Components {
					fmt.Printf("         %s\n", s)
				}

				fmt.Println("      kustomizations:")
				for _, s := range app.Kustomizations {
					fmt.Printf("         %s\n", s)
				}
			}
		}

		if filename := viper.GetString("graphviz-file"); filename != "" {
			if err := writeGraph(apps, filename); err != nil {
				return err
			}
		}
		return nil
	}
}

func writeGraph(apps []applications.Application, filename string) error {
	graph := visualise.NewDOT(apps)
	if err := os.WriteFile(filename, []byte(graph.String()), 0644); err != nil {
		return err
	}
	return nil
}

func parentApps(apps []applications.Application) []applications.Application {
	res := []applications.Application{}
	for _, v := range apps {
		if len(v.Parents) == 0 {
			res = append(res, v)
		}
	}
	return res
}

func childApps(apps []applications.Application, parent string) []applications.Application {
	res := []applications.Application{}
	for _, child := range apps {
		if hasParentApplication(child, parent) {
			res = append(res, child)
		}
	}
	return res
}

func hasParentApplication(app applications.Application, parent string) bool {
	for _, v := range app.Parents {
		if v.Name == parent {
			return true
		}
	}
	return false
}
