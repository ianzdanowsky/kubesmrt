package memory

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

// get pods memory
var MemoryCmd = &cobra.Command{
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
	MemoryCmd.Flags().StringVarP(&kubeConfigFilePath, "kubeconfig", "k", "", "Kubeconfig file (optional)")
	MemoryCmd.Flags().BoolVarP(&showDocs, "docs", "d", false, "Show (optional)")
	MemoryCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Filter by namespace. Defaults to all namespaces.")
}
