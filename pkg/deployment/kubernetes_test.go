package deployment

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestKubernetesStrategy_Rollback(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-app",
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "test-container",
							Image: "test-app:v1.0.0",
						},
					},
				},
			},
		},
	}

	_, err := clientset.AppsV1().Deployments("default").Create(
		context.TODO(),
		deployment,
		metav1.CreateOptions{},
	)
	if err != nil {
		t.Fatalf("error creating deployment: %v", err)
	}

	config := KubernetesConfig{
		Namespace:  "default",
		Deployment: "test-app",
	}

	k8s := NewKubernetesStrategy(clientset, config)
	err = k8s.Rollback("v1.0.0", "v0.9.0")
	if err != nil {
		t.Errorf("Rollback failed: %v", err)
	}

	updated, err := clientset.AppsV1().Deployments("default").Get(
		context.TODO(),
		"test-app",
		metav1.GetOptions{},
	)
	if err != nil {
		t.Fatalf("error getting deployment: %v", err)
	}

	if len(updated.Spec.Template.Spec.Containers) == 0 {
		t.Fatal("no containers found in updated deployment")
	}

	expectedImage := fmt.Sprintf("%s:v0.9.0", config.Deployment)
	if updated.Spec.Template.Spec.Containers[0].Image != expectedImage {
		t.Errorf("container image = %s, want %s",
			updated.Spec.Template.Spec.Containers[0].Image,
			expectedImage)
	}
}

type mockExecutor struct {
	commands []string
	err      error
}

func (m *mockExecutor) Command(name string, args ...string) *exec.Cmd {
	m.commands = append(m.commands, name+" "+strings.Join(args, " "))
	return exec.Command("echo", "mock")
}
