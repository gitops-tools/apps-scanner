package main

import (
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta2"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kustomizev1.AddToScheme(scheme))
}

func main() {
	cfg, err := config.GetConfig()
	cobra.CheckErr(err)

	cl, err := client.New(cfg, client.Options{Scheme: scheme})
	cobra.CheckErr(err)

	rootCmd := makeRootCmd()
	rootCmd.AddCommand(newApplicationsCmd(cl))
	rootCmd.AddCommand(newPipelinesCmd(cl))

	cobra.CheckErr(rootCmd.Execute())
}
