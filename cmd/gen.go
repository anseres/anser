package cmd

import (
	"fmt"
	"os"
	"os/exec"

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

// protoc -I.  demo.proto --anser_out=.
func gen(cmd *cobra.Command, args []string) {
	fmt.Printf("gen services\n")
	if protoFile == "" {
		fmt.Printf("please input proto file\n")
		os.Exit(1)
	}

	pbCmd := exec.Command("protoc", "-I.", protoFile, "--anser_out=.")
	out, err := pbCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error: %s", string(out))
		panic(err)
	}
	pbCmd = exec.Command("protoc", "-I.", protoFile, "--go_out=plugins=grpc:.")
	out, err = pbCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error: %s", string(out))
		panic(err)
	}
}
