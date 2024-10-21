package auth

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubeApiAuth retrieves a Kubernetes clientset using the specified kubeconfig path.
// If no kubeconfig path is provided, it attempts to fetch it from ~/.kubesmrt.config.
// If the kubeconfig path cannot be found, it prompts the user to run the setup command.
//
// Parameters:
//
//	kubeconfigPath (string): Path to the kubeconfig file. If empty, the path is fetched from ~/.kubesmrt.config.
//
// Returns:
//
//	*kubernetes.Clientset: A clientset for interacting with the Kubernetes API.
//	error: Any error encountered while loading the kubeconfig or creating the clientset.
func GetKubeApiAuth(kubeconfigPath string) (*kubernetes.Clientset, error) {
	if kubeconfigPath == "" {
		var err error
		kubeconfigPath, err = getKubeConfigPath() // Fetch the path
		if err != nil {
			fmt.Println("Failed to get kubeconfig path. Please run 'kubesmrt setup' to configure it.")
			os.Exit(1)
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Println("Error loading kubeconfig:", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating Kubernetes client:", err)
		return nil, err
	}
	return clientset, nil
}

// GetKubeConfigPath checks if ~/.kubesmrt.config exists and contains a valid kubeconfig path.
// If not, it indicates that the user needs to run the setup command.
func getKubeConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error fetching current user:", err)
		os.Exit(1)
	}
	kubesmrtConfigFilePath := filepath.Join(usr.HomeDir, ".kubesmrt.config")

	if fileExists(kubesmrtConfigFilePath) {
		file, err := os.Open(kubesmrtConfigFilePath)
		if err != nil {
			return "", fmt.Errorf("error reading config file: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			return scanner.Text(), nil
		}
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("error reading config file: %v", err)
		}
	}

	return "", fmt.Errorf("kubeconfig not found. Please run 'setup' to configure the kubeconfig path")
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}
