package main

import (
	"context"
	"fmt"

	"kubesmrt/pkg/auth"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	var kubeconfigFilePath string

	clientset, err := auth.GetKubeApiAuth(kubeconfigFilePath)
	if err != nil {
		fmt.Println("Failed to get kubeapi auth")
		return
	}

	replicaCount := int32(1)

	d := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "test-deployment",
			Labels: map[string]string{"app": "test-deployment"},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicaCount,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "test-deployment"},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "test-deployment"},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "test-container",
							Image: "nginx",
						},
					},
				},
			},
		},
	}

	ctx := context.TODO()

	opts := metav1.CreateOptions{}
	_, err = clientset.AppsV1().Deployments("default").Create(ctx, d, opts)
	if err != nil {
		fmt.Println(err)
	}

}
