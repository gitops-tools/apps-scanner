package main

import "github.com/spf13/cobra"

func makeRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "scanner <command>",
		Short:         "Scan repositories",
		Long:          "Scan and log information from clusters based on labels",
		SilenceErrors: true,
	}
}
