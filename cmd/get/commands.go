package get

import (
	"kubesmrt/cmd/get/pods"

	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Use to get pods statistics",
	Long:  `A CLI tool built with Golang to interact with Kubernetes clusters and perform custom aggregations`,
}

func init() {
	// Adding a flag for the kubeconfig path if needed
	GetCmd.AddCommand(pods.PodsCmd)
}
