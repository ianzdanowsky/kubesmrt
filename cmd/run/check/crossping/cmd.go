package crossping

import (
	"context"
	"fmt"
	"kubesmrt/pkg/auth"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

// Flags
var kubeConfigFilePathFlag string
var namespaceFlag string
var targetNamespaceFlag string
var imageFlag string
var tolerationsFlag string

// get pods nolimits
var Crossping = &cobra.Command{
	Use:   "crossping",
	Short: "Use to create a daemonset and run a crossping test between cluster pods",
	Run: func(cmd *cobra.Command, args []string) {

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			<-c
			fmt.Println("\nReceived interrupt signal, cleaning up...")
			cancel()
		}()

		clientset, err := auth.GetKubeApiAuth(kubeConfigFilePathFlag)
		if err != nil {
			fmt.Println("Failed to get kube api auth:", err)
			os.Exit(1)
		}
		RunCheckCrossping(
			ctx,
			clientset,
			namespaceFlag,
			targetNamespaceFlag,
			imageFlag)
	},
}

func init() {
	Crossping.Flags().StringVarP(&kubeConfigFilePathFlag,
		"kubeconfig", "k", "", "Kubeconfig file (optional)")
	Crossping.Flags().StringVarP(&namespaceFlag,
		"namespace", "n", "default", "Namespace")
	Crossping.Flags().StringVarP(&targetNamespaceFlag,
		"target-namespace", "T", "default", "Target Namespace")
	Crossping.Flags().StringVarP(&imageFlag,
		"image", "i", "alpine:latest",
		"Specify the container image to be used for the ping test. The image should include ping capabilities. Defaults to 'alpine:latest', which will be pulled if not available locally.")
	Crossping.Flags().StringVarP(&tolerationsFlag,
		"tolerations", "t", "", "Tolerations for the testing pods. If more than one, separate by comma. E.g. example-key:Equal:example-value:NoSchedule")
}
