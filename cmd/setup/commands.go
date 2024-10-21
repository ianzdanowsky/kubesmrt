package setup

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

var kubesmrtConfigFilePath string // .kubesmrt.config
var kubeConfigFilePath string     // This represents the kubeconfig path that is saved in .kubesmrt.config

// SetupCmd represents the setup command
var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the kubeconfig file to be used by the CLI",
	Run: func(cmd *cobra.Command, args []string) {

		// Get the current OS user
		usr, err := user.Current()
		if err != nil {
			fmt.Println("Error fetching current user:", err)
			os.Exit(1)
		}

		// Use default ~/.kube/config if none is provided in the setup command
		if kubeConfigFilePath == "" {
			kubeConfigFilePath = filepath.Join(usr.HomeDir, ".kube", "config")
		}

		// Overwrite the kubeconfig path in ~/.kubesmrt.config
		if err := setup(kubeConfigFilePath); err != nil {
			fmt.Println("Error saving kubeconfig path:", err)
			os.Exit(1)
		}

		fmt.Println("Kubeconfig is set to:", kubeConfigFilePath)
	},
}

func init() {
	// Register the 'config' flag to take a kubeconfig file path
	SetupCmd.Flags().StringVarP(&kubeConfigFilePath, "config", "c", "", "Path to kubeconfig file")

	// Build the kubesmrt config file path (~/.kubesmrt.config)
	// Will be used in the setup() function to save the kubeconfig location.
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error fetching current user:", err)
		os.Exit(1)
	}
	kubesmrtConfigFilePath = filepath.Join(usr.HomeDir, ".kubesmrt.config")
}
