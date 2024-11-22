package run

import (
	"kubesmrt/cmd/run/check"

	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Use to run",
	Long:  `A CLI tool built with Golang to interact with Kubernetes clusters and perform custom aggregations`,
}

func init() {
	RunCmd.AddCommand(check.CheckCmd)
}
