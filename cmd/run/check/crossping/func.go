package crossping

import (
	"context"
	"fmt"
	"kubesmrt/pkg/render"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func RunCheckCrossping(ctx context.Context,
	clientset kubernetes.Interface,
	namespaceFlag string,
	targetNamespaceFlag string,
	imageFlag string) {

	var getPodsCrosspingRows [][]string
	getPodsCrosspingHeader := []string{"ORIG HOST", "ORIG POD", "ORIG NS", "TARGET HOST", "TARGET POD", "TARGET NS", "TARGET POD IP", "RESULT"}

	originDaemonSetName := "orig-crossping"
	originDaemonSetLabel := "ping-test-origin"
	targetDaemonSetName := "target-crossping"
	targetDaemonSetLabel := "ping-test-target"

	_, err := createPingDaemonSet(clientset, originDaemonSetName, originDaemonSetLabel, tolerationsFlag, namespaceFlag)
	if err != nil {
		fmt.Println("Error creating origin DaemonSet:", err)
		return
	}

	_, err = createPingDaemonSet(clientset, targetDaemonSetName, targetDaemonSetLabel, tolerationsFlag, targetNamespaceFlag)
	if err != nil {
		fmt.Println("Error creating target DaemonSet:", err)
		return
	}

	// Cleanup DaemonSet after function ends
	defer func() {
		deleteDaemonSet(clientset, originDaemonSetName, namespaceFlag)
		deleteDaemonSet(clientset, targetDaemonSetName, targetNamespaceFlag)
	}()

	err = waitForDaemonSetPodsReady(ctx, clientset, originDaemonSetName, namespaceFlag)
	if err != nil {
		fmt.Println("Error waiting for origin DaemonSet pods to be ready:", err)
		return
	}

	err = waitForDaemonSetPodsReady(ctx, clientset, targetDaemonSetName, targetNamespaceFlag)
	if err != nil {
		fmt.Println("Error waiting for target DaemonSet pods to be ready:", err)
		return
	}

	originDaemonSetPods, err := clientset.CoreV1().Pods(namespaceFlag).List(ctx, metav1.ListOptions{
		LabelSelector: "app=ping-test-origin",
	})
	if err != nil {
		fmt.Println("Failed to list origin DaemonSet pods:", err)
		return
	}

	targetDaemonSetPods, err := clientset.CoreV1().Pods(targetNamespaceFlag).List(ctx, metav1.ListOptions{
		LabelSelector: "app=ping-test-target",
	})
	if err != nil {
		fmt.Println("Failed to list target DaemonSet pods:", err)
		return
	}

	fmt.Println("Starting the crossping test, this might take sometime.")

	// Set the retry attempts for pinging
	maxRetry := 2

	for _, originPod := range originDaemonSetPods.Items {
		// Check if the context is done (e.g cancelled) before each loop. Clean if it is.
		select {
		case <-ctx.Done():
			fmt.Println("Command canceled, cleaning up resources...")
			deleteDaemonSet(clientset, originDaemonSetName, namespaceFlag)
			deleteDaemonSet(clientset, targetDaemonSetName, targetNamespaceFlag)
			return
		default:
		}

		for _, targetPod := range targetDaemonSetPods.Items {
			if targetPod.Status.PodIP == "" || originPod.Name == targetPod.Name {
				continue
			}

			result := pingPodFromPod(originPod.Name, namespaceFlag, targetPod.Status.PodIP, maxRetry)
			if result {
				getPodsCrosspingRows = append(getPodsCrosspingRows, []string{
					originPod.Spec.NodeName, originPod.Name, originPod.Namespace,
					targetPod.Spec.NodeName, targetPod.Name, targetPod.Namespace, targetPod.Status.PodIP, "SUCCESS"})
			} else {
				getPodsCrosspingRows = append(getPodsCrosspingRows, []string{
					originPod.Spec.NodeName, originPod.Name, originPod.Namespace,
					targetPod.Spec.NodeName, targetPod.Name, targetPod.Namespace, targetPod.Status.PodIP, "FAILED"})
			}

			// Check context again before each ping to handle interruption mid-loop
			if ctx.Err() != nil {
				fmt.Println("Context canceled during operation, cleaning up resources...")
				return
			}
		}
	}

	// Render the output as a table
	render.SimpleTable(getPodsCrosspingHeader, []string{}, getPodsCrosspingRows)
}
