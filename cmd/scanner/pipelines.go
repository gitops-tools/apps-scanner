package main

import (
	"context"
	"fmt"
	"strings"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta2"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gitops-tools/apps-scanner/pkg/pipelines"
)

func newPipelinesCmd(cl client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "pipelines",
		Short: "List pipelines in the cluster",
		RunE:  listPipelines(cl),
	}
}

func listPipelines(cl client.Client) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting to scan for kustomizations")
		kustomizationList := &kustomizev1.KustomizationList{}
		err := cl.List(context.Background(), kustomizationList, client.HasLabels([]string{pipelines.PipelineNameLabel}))
		if err != nil {
			return fmt.Errorf("failed to list kustomizations: %w", err)
		}
		fmt.Printf("found %d kustomizations\n", len(kustomizationList.Items))

		p := pipelines.NewParser()
		pipelines, err := p.Pipelines()
		if err != nil {
			return fmt.Errorf("failed to discover pipelines: %w", err)
		}

		for _, v := range pipelines {
			fmt.Printf("pipeline %s has stages: %s\n", v.Name, strings.Join(v.Environments, ","))
		}
		return nil
	}
}
