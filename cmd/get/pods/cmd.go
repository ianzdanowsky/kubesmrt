package pods

import (
	"fmt"
	"kubesmrt/pkg/auth"
	"os"

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

// get pods nolimits
var NoLimits = &cobra.Command{
	Use:   "nolimits",
	Short: "Use to get pods without resource limit set",
	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := auth.GetKubeApiAuth(kubeConfigFilePath)
		if err != nil {
			fmt.Println("Failed to get kube api auth:", err)
			os.Exit(1)
		}
		getPodsNoLimits(clientset)
	},
}

// get pods memory
var Memory = &cobra.Command{
	Use:   "memory",
	Short: "Use to get pods without resource limit set",
	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := auth.GetKubeApiAuth(kubeConfigFilePath)
		if err != nil {
			fmt.Println("Failed to get kube api auth:", err)
			os.Exit(1)
		}
		getPodsMemory(clientset, showDocs, namespace)
	},
}

func init() {
	// get pods nolimits
	PodsCmd.AddCommand(NoLimits)
	NoLimits.Flags().StringVarP(&kubeConfigFilePath, "kubeconfig", "k", "", "Kubeconfig file (optional)")

	// get pods memory
	PodsCmd.AddCommand(Memory)
	Memory.Flags().StringVarP(&kubeConfigFilePath, "kubeconfig", "k", "", "Kubeconfig file (optional)")
	Memory.Flags().BoolVarP(&showDocs, "docs", "d", false, "Show (optional)")
	Memory.Flags().StringVarP(&namespace, "namespace", "n", "", "Filter by namespace. Defaults to all namespaces.")
}
