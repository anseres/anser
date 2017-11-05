package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(genCmd)
	genCmd.PersistentFlags().StringVar(&protoFile, "proto", "", "proto file path")
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "gen service framework",
	Run:   gen,
}

var (
	protoFile string
)

func gen(cmd *cobra.Command, args []string) {
	fmt.Printf("gen services\n")
}
