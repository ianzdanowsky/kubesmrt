package nolimits

import (
	"context"
	"kubesmrt/pkg/mocks"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// Unit test for getPodsNoLimits
func TestGetPodsNoLimits(t *testing.T) {
	// Create a fake Kubernetes clientset
	fakeClientset := fake.NewSimpleClientset()

	// Add mock pods to the clientset
	fakeClientset.CoreV1().Pods("default").Create(context.TODO(), &mocks.MockGetPodsMemory().Items[0], metav1.CreateOptions{})
	fakeClientset.CoreV1().Pods("default").Create(context.TODO(), &mocks.MockGetPodsMemory().Items[1], metav1.CreateOptions{})

	// Test the function
	noLimitsPodList := getPodsNoLimits(fakeClientset)

	// Assertion for the expected output
	for _, pod := range noLimitsPodList.Items {
		for _, container := range pod.Spec.Containers {
			if container.Name != "container-without-limits" {
				t.Errorf("Expected 'containter-without-limits', but got %s", container.Name)
			}
		}
	}

}
