package check

import (
	"kubesmrt/cmd/run/check/crossping"

	"github.com/spf13/cobra"
)

// get pods
var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Use this for checks",
	Long:  `A CLI tool built with Golang to interact with Kubernetes clusters and perform custom aggregations`,
}

func init() {
	CheckCmd.AddCommand(crossping.Crossping)
}
