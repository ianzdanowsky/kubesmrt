package nolimits

import (
	"fmt"
	"kubesmrt/pkg/auth"
	"os"

	"github.com/spf13/cobra"
)

// Variables to store the flag values
var kubeConfigFilePath string

// get pods nolimits
var NoLimitsCmd = &cobra.Command{
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

func init() {
	NoLimitsCmd.Flags().StringVarP(&kubeConfigFilePath, "kubeconfig", "k", "", "Kubeconfig file (optional)")

}
