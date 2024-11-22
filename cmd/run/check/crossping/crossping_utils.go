package crossping

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func createPingDaemonSet(
	clientset kubernetes.Interface,
	daemonSetName,
	label,
	tolerations,
	namespace string) (*appsv1.DaemonSet, error) {

	var finalTolerations []v1.Toleration

	if len(tolerations) > 0 {
		tArray := strings.Split(tolerations, ",")

		for _, t := range tArray {
			toleration := strings.Split(t, ":")
			finalTolerations = append(finalTolerations, v1.Toleration{Key: toleration[0], Operator: v1.TolerationOperator(toleration[1]), Value: toleration[2], Effect: v1.TaintEffect(toleration[3])})
		}
	}

	daemonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      daemonSetName,
			Namespace: namespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": label,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": label,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "ping-container",
							Image: "alpine", // Alpine has ping installed
							Command: []string{
								"/bin/sh", "-c", "sleep 3600",
							},
						},
					},
					RestartPolicy: v1.RestartPolicyAlways,
					Tolerations:   finalTolerations,
				},
			},
		},
	}

	daemonSet, err := clientset.AppsV1().DaemonSets(namespace).Create(context.TODO(), daemonSet, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("Failed to create daemon set:", err)
		os.Exit(1)
	}

	return daemonSet, err
}

func waitForDaemonSetPodsReady(ctx context.Context, clientset kubernetes.Interface, daemonSetName, namespace string) error {
	for {

		select {
		case <-ctx.Done():
			fmt.Println("Waiting DaemonSet aborted by user.")
			return ctx.Err() // Return the context's error to indicate it was canceled
		default:
		}

		_, err := clientset.AppsV1().DaemonSets(namespace).Get(context.TODO(), daemonSetName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("name=%s", daemonSetName),
		})
		if err != nil {
			return err
		}

		allPodsReady := true
		for _, pod := range pods.Items {
			if pod.Status.Phase != v1.PodRunning {
				allPodsReady = false
				break
			}

			for _, cond := range pod.Status.Conditions {
				if cond.Type == v1.PodReady && cond.Status != v1.ConditionTrue {
					allPodsReady = false
					break
				}
			}
		}

		if allPodsReady {
			fmt.Printf("All pods from DaemonSet %s are running and ready.\n", daemonSetName)
			time.Sleep(20 * time.Second)
			return nil
		}

		fmt.Printf("Waiting for all replicas of DaemonSet %s to be ready...\n", daemonSetName)
		time.Sleep(20 * time.Second) // Polling interval
	}
}

// pingPodFromPod runs a ping from the source pod to the target IP
func pingPodFromPod(podName, namespace, targetIP string, retryAttempts int) bool {
	for i := 0; i < retryAttempts; i++ {
		cmd := exec.Command("kubectl", "exec", podName, "-n", namespace, "--", "ping", "-c", "1", targetIP)
		out, err := cmd.CombinedOutput()
		if strings.Contains(string(out), "1 packets received") && err == nil {
			return true
		}
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("Retrying pinging from %s to %s\n", podName, targetIP)
	}
	return false
}

// deleteStatefulSet deletes the StatefulSet after the ping test
func deleteDaemonSet(clientset kubernetes.Interface, daemonSetName, namespace string) error {
	gracePeriodSeconds := int64(0)
	return clientset.AppsV1().DaemonSets(namespace).Delete(context.TODO(), daemonSetName, metav1.DeleteOptions{GracePeriodSeconds: &gracePeriodSeconds})
}
