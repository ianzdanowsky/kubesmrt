package pods

import (
	"kubesmrt/cmd/get/pods/memory"
	"kubesmrt/cmd/get/pods/nolimits"

	"github.com/spf13/cobra"
)

// Variables to store the flag values
var kubeConfigFilePath string
var showDocs bool
var namespace string

// get pods
var PodsCmd = &cobra.Command{
	Use:     "pod",
	Aliases: []string{"pods"},
	Short:   "Use to get pods statistics",
	Long:    `A CLI tool built with Golang to interact with Kubernetes clusters and perform custom aggregations`,
}

func init() {
	// get pods nolimits
	PodsCmd.AddCommand(nolimits.NoLimitsCmd)
	// get pods memory
	PodsCmd.AddCommand(memory.MemoryCmd)
}
