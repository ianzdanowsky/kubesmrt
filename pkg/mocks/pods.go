package mocks

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// get pods memory
func MockGetPodsMemory() *v1.PodList {
	return &v1.PodList{
		Items: []v1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "pod-with-limits", Namespace: "default"},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-with-limits",
							Image: "busybox",
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("100m"),
									v1.ResourceMemory: resource.MustParse("256Mi"),
								},
							},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: "pod-without-limits", Namespace: "default"},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:      "container-without-limits",
							Image:     "nginx",
							Resources: v1.ResourceRequirements{
								// No Limits
							},
						},
					},
				},
			},
		},
	}
}
