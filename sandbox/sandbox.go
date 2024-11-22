package main

import (
	"context"
	"fmt"
	"kubesmrt/pkg/auth"
	"kubesmrt/pkg/render"
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
	// nodeData, err := json.MarshalIndent(nodes, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error marshalling nodes data to JSON:", err)
	// 	return
	// }

	// Output all the nodes data to a JSON file named "nodes_data.json"
	// render.SaveOutputToJson(nodeData, "nodes_data")

	fmt.Println("Node data saved to nodes_data.json")
}
