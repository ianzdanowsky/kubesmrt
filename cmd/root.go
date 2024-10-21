package cmd

import (
	"fmt"
	"os"

	"kubesmrt/cmd/get"
	"kubesmrt/cmd/setup"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubesmrt",
	Short: "Custom Kubernetes CLI to perform pod aggregations",
	Long:  `A CLI tool built with Golang to interact with Kubernetes clusters and perform custom aggregations`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(setup.SetupCmd)
	rootCmd.AddCommand(get.GetCmd)
}
