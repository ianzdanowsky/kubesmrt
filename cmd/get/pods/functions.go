package pods

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"kubesmrt/pkg/docs"
	"kubesmrt/pkg/render"
	"kubesmrt/pkg/utils"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getPodsNoLimits(clientset kubernetes.Interface) *v1.PodList {
	// Fetching pods from the cluster
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error fetching pods:", err)
		os.Exit(1)
	}

	var podsWithoutLimitsData [][]string
	podsWithoutLimitsHeaders := []string{"Pod name", "Namespace", "Container Name", "Image"}
	// Iterate over all pods and check for resource limits
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if container.Resources.Limits != nil {
				continue
			} else {
				podsWithoutLimitsData = append(podsWithoutLimitsData, []string{pod.Name, pod.Namespace, container.Name, container.Image})
			}
		}
	}

	// Build the v1.PodList to return
	var podListItems []v1.Pod
	for _, data := range podsWithoutLimitsData {
		pod := v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      data[0],
				Namespace: data[1],
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  data[2],
						Image: data[3],
					},
				},
			},
		}
		podListItems = append(podListItems, pod)
	}

	// Render the output
	fmt.Printf("There are %d containers in the cluster without resource limits.\n", len(podsWithoutLimitsData))
	render.SimpleTable(podsWithoutLimitsHeaders, podsWithoutLimitsData)

	return &v1.PodList{Items: podListItems}
}

func getPodsMemory(clientset kubernetes.Interface, showDocs bool, namespace string) [][]string {

	if showDocs {
		render.DocsAsTable(docs.DocsGetPodsMemory)
		return [][]string{}
	}

	var getPodsMemoryOutput [][]string

	// Fetching nodes from the cluster
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing nodes:", err)
		os.Exit(1)
	}

	// Iterate over each node and query the Kubelet stats/summary API
	for _, node := range nodes.Items {
		// Fetch the internal IP address of the node (used for Kubelet communication)
		hostName := ""
		for _, addr := range node.Status.Addresses {
			if addr.Type == "Hostname" {
				hostName = addr.Address
				break
			}
		}
		if hostName == "" {
			fmt.Println("No internal IP found for node", node.Name)
			continue
		}

		url := fmt.Sprintf("/api/v1/nodes/%s/proxy/stats/summary", hostName)

		response, err := clientset.CoreV1().RESTClient().Get().AbsPath(url).DoRaw(context.TODO())
		if err != nil {
			fmt.Println("Error querying Kubelet via Kubernetes API proxy:", err)
			return [][]string{}
		}

		var responseJson map[string]interface{}

		err = json.Unmarshal(response, &responseJson)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return [][]string{}
		}

		pods := responseJson["pods"].([]interface{})

		// Sort pods by usageBytes in descending order
		sort.Slice(pods, func(i, j int) bool {
			pod1 := pods[i].(map[string]interface{})
			pod2 := pods[j].(map[string]interface{})
			memory1 := pod1["memory"].(map[string]interface{})["usageBytes"].(float64)
			memory2 := pod2["memory"].(map[string]interface{})["usageBytes"].(float64)
			return memory1 > memory2
		})

		for _, pod := range pods {
			podMap := pod.(map[string]interface{})
			podMemoryMap := podMap["memory"].(map[string]interface{})
			podRefMap := podMap["podRef"].(map[string]interface{})
			podContainerMap := podMap["containers"].([]interface{})
			podName := podRefMap["name"].(string)
			podUsageBytesMB := fmt.Sprintf("%v", utils.ConvertBytesToMB(podMemoryMap["usageBytes"].(float64)))
			podWorkingSetBytesMB := fmt.Sprintf("%v", utils.ConvertBytesToMB(podMemoryMap["workingSetBytes"].(float64)))

			// If namespace is set, filter the pods
			if namespace != "" {
				if podRefMap["namespace"] != namespace {
					continue
				}
			}

			for _, container := range podContainerMap {
				containerMap := container.(map[string]interface{})
				containerMemoryMap := containerMap["memory"].(map[string]interface{})

				containerName := containerMap["name"].(string)
				containerWorkingSetBytesMB := fmt.Sprintf("%v", utils.ConvertBytesToMB(containerMemoryMap["workingSetBytes"].(float64)))
				getPodsMemoryOutput = append(getPodsMemoryOutput, []string{hostName, podName, podUsageBytesMB, podWorkingSetBytesMB, containerName, containerWorkingSetBytesMB})

			}
		}
		// fmt.Println(finalPodsStats)

		// The pod's memory can also include overhead not reflected in individual container memory,
		// such as networking or system-level allocations that the pod incurs but are not tied to any specific container.
		// This could explain the slight difference between the container's memory (2359296 bytes) and the pod's memory (2727936 bytes).
		// The additional memory might be due to overhead or other system allocations that are not attributed to the container.
	}
	getPodsMemoryHeaders := []string{"Host Name", "Pod Name", "usageBytes", "WorkingSetBytes", "Container", "WorkingSetBytes"}
	render.IdenticalCellMergingTable(getPodsMemoryHeaders, getPodsMemoryOutput)

	return getPodsMemoryOutput
}
