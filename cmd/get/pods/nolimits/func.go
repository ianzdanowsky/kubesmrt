package nolimits

import (
	"context"
	"fmt"
	"kubesmrt/pkg/render"
	"os"

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
	render.SimpleTable(podsWithoutLimitsHeaders, []string{}, podsWithoutLimitsData)

	return &v1.PodList{Items: podListItems}
}
