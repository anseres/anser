package cmd

import "github.com/spf13/cobra"

// RootCmd is tool command for anser
var RootCmd = &cobra.Command{
	Use:   "anser",
	Short: "Anser is tool for help write grpc server and manage k8s cluster",
}
