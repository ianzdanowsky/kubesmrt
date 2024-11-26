## kubesmrt

kubesmrt is a CLI tool built in Go, designed to provide aggregated and smart data about your Kubernetes cluster by querying Kubernetes resources and underlying layers. The tool simplifies the process of gathering critical information and offers additional insights beyond standard Kubernetes commands. It follows a similar command pattern to kubectl, making it intuitive and familiar for Kubernetes users.


### Features
- Fetch detailed pod statistics, including resource limits and usage.
- Query aggregated data about nodes, pods, and services within the cluster.
- Provide deeper insights by querying Kubernetes and lower layers such as the Kubelet.
- Easy-to-use CLI interface with command flags and subcommands.

### Getting Started
#### Prerequisites
To use kubesmrt, you need a Kubernetes cluster and a valid kubeconfig file (usually located at `~/.kube/config`).

For development purposes, you should have Minikube or a similar Kubernetes environment running locally.

- [Minikube: Installation Guide](https://minikube.sigs.k8s.io/docs/start/)
- [kubectl: Install kubectl to interact with your Kubernetes cluster. Installation Guide](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

#### Installation
Clone the repository and install dependencies using go mod:

```bash
git clone https://github.com/yourusername/kubesmrt.git
cd kubesmrt
go mod tidy
```

### Running the CLI
Make sure you have a valid Kubernetes environment and kubeconfig file. Replace `~/.kube/config` below with the path to your kubeconfig file (e.g., root kubeconfig file).

```bash
./kubesmrt setup --config ~/.kube/config
```

This command will save the kubeconfig location to `~/.kubesmrt.config`, which will be used in subsequent commands.

```bash
./kubesmrt get pods nolimits
```

By providing the `--kubeconfig` flag with the desired kubeconfig file path, you can dynamically change the configuration used by the CLI tool.


```bash
./kubesmrt get pods nolimits --kubeconfig ~/.kube/config
```

### Development Environment
To develop and test kubesmrt, you should have a local Kubernetes environment, such as Minikube, running on your machine. The CLI relies on the Kubernetes API and kubeconfig for authentication and querying cluster resources.

Steps for setting up the dev environment:

1. Install Minikube or another local Kubernetes solution.
2. Ensure your kubeconfig is configured correctly and points to the running Kubernetes environment.
3. Use `go test` to run the unit tests. We recommend using `go clean -testcache && go test ./...` to ensure a clean test environment.

Example:

```bash
minikube start
./kubesmrt setup --config ~/.kube/config
./kubesmrt get pods nolimits
```

#### Sandbox for testing functions

The `sandbox` package provides a place to quickly test functions and Kubernetes API calls without affecting the main application. This is useful for rapid development, experimentation, and debugging.

You can run the sandbox code by executing:

```bash
go run sandbox/sandbox.go
```

The following code demonstrates how to fetch and print the names of all nodes in your Kubernetes cluster using the Kubernetes API.


```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"kubesmrt/utils"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	var kubeconfigFilePath string
	clientset, err := auth.GetKubeApiAuth(kubeconfigFilePath)
	if err != nil {
		fmt.Println("Failed to get kubeapi auth:", err)
		os.Exit(1)
	}

	// Fetching nodes from the cluster
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing nodes:", err)
		os.Exit(1)
	}

	// Prepare data for table output
	headers := []string{"Node Name", "Ready"}
	var data [][]string
	var footer []string
	for _, node := range nodes.Items {
		nodeName := node.Name
		for _, condition := range node.Status.Conditions {
			if condition.Type == "Ready" {
				data = append(data, []string{nodeName, string(condition.Status)})
			}
		}
	}

	// Render the node data in table format
	fmt.Println("Node Information:")
	render.SimpleTable(headers, footer, data)

	// Save the full node data to a JSON file
	nodeData, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling nodes data to JSON:", err)
		return
	}

	// Output all the nodes data to a JSON file named "nodes_data.json"
	utils.SaveOutputToJson(nodeData, "nodes_data")

	fmt.Println("Node data saved to nodes_data.json")
}
```

### Folder Structure
The project follows a logical structure to separate CLI commands, utilities, and core functionalities:

```bash
├── cmd
│   ├── get                     # Logic for the "get" commands
│   ├── root.go                 # Entry point for root command setup
│   ├── run                     # Logic for the "run" commands
│   └── setup                   # Logic for the "setup" command
├── k8s-manifests               # Kubernetes manifest files used for testing
├── pkg                         # Package directory for reusable code
│   ├── auth                    # Authentication helpers for kubeconfig
│   ├── docs                    # Documentation-related utilities
│   ├── mocks                   # Mocks for unit testing
│   ├── render                  # Utilities for rendering table output in the CLI
│   └── utils                   # General utility functions
└── sandbox                     # Experimental code and testing ground
```

#### Key Packages Used
1. Cobra: A powerful library used to create CLI applications in Go.
2. Kubernetes Client-Go: The official Go client for the Kubernetes API.
3. Testify for Unit Testing: A package used for writing unit tests in Go.

#### Commands Overview
- `setup`: Configures the kubeconfig path used by the CLI tool.
- `get pods nolimits`: Fetches and displays a list of pods that have no resource limits configured.

#### Contributing
Contributions are welcome! Please submit pull requests with a description of your changes and any relevant tests. If you're adding new commands or features, ensure they are accompanied by adequate documentation and unit tests.

#### Creating a command

This guide outlines the structure and steps for creating a new command in the kubesmrt CLI tool. Follow these conventions to ensure consistency and maintainability.

Every command should follow the below structure:

```bash
cmd/<parent_command>/<subcommand>/
├── cmd.go          # Implements the Cobra command - defines the command name.
├── func.go         # Contains the main logic for the command, ideally one function.
├── func_test.go    # Unit tests for the main function in `func.go`
├── utils.go        # Helper/side functions specific to this command
└── utils_test.go   # Unit tests for the helper functions
```

The `cmd.go` file will use the `func init()` method to set up flags, retrieve their values, and pass them as arguments to the main function in `func.go`. See example below:

The `func.go` itself, should ideally have one meaningful function (e.g., getPodsNoLimits) and should not deal with user input parsing or Cobra-specific details.

cmd.go
```bash
var namespaceFlag string

func init() {
	NoLimitsCmd.Flags().StringVarP(&namespaceFlag, "namespace", "n", "default", "Namespace to filter pods")
}

var NoLimitsCmd = &cobra.Command{
	Use:   "nolimits",
	Short: "Get pods with no resource limits set",
	Run: func(cmd *cobra.Command, args []string) {
		err := getPodsNoLimits(namespaceFlag)
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}
```

func.go
```bash
func getPodsNoLimits(namespace string) error {
	// Main logic to fetch and display pods with no resource limits
	fmt.Printf("Fetching pods with no limits in namespace: %s\n", namespace)
	// Implement functionality here
	return nil
}
```

#### To Do:
1 - Create a k8s-templates directory to add k8s resource yaml that will support the command testing.
3 - Create a constant file (constants.go) to keep all the tables titles and others.

#### License
This project is licensed under the MIT License. See the LICENSE file for more details.
