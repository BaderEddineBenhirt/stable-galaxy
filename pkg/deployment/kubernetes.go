package deployment

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubernetesConfig struct {
	Namespace     string
	Deployment    string
	ImageTemplate string
	Labels        map[string]string
	Annotations   map[string]string
	Resources     struct {
		Limits   map[string]string
		Requests map[string]string
	}
	Replicas      *int32
	Strategy      string
	ConfigPath    string
	Context       string
	CustomOptions map[string]interface{}
}
type KubernetesClientset interface {
	kubernetes.Interface
}

type KubernetesStrategy struct {
	clientset KubernetesClientset
	config    KubernetesConfig
}

func NewKubernetesStrategy(clientset KubernetesClientset, config KubernetesConfig) *KubernetesStrategy {
	return &KubernetesStrategy{
		clientset: clientset,
		config:    config,
	}
}

func (k *KubernetesStrategy) buildImage(version string) string {
	if k.config.ImageTemplate != "" {
		return fmt.Sprintf(k.config.ImageTemplate, version)
	}
	return fmt.Sprintf("%s:%s", k.config.Deployment, version)
}

func (k *KubernetesStrategy) updateDeployment(deployment *appsv1.Deployment, version string) {
	for i := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[i].Image = k.buildImage(version)
	}

	if deployment.Labels == nil {
		deployment.Labels = make(map[string]string)
	}
	for key, value := range k.config.Labels {
		deployment.Labels[key] = value
	}

	if deployment.Annotations == nil {
		deployment.Annotations = make(map[string]string)
	}
	for key, value := range k.config.Annotations {
		deployment.Annotations[key] = value
	}

	if len(k.config.Resources.Limits) > 0 || len(k.config.Resources.Requests) > 0 {
		for i := range deployment.Spec.Template.Spec.Containers {
			container := &deployment.Spec.Template.Spec.Containers[i]
			if len(k.config.Resources.Limits) > 0 {
				container.Resources.Limits = convertToResourceList(k.config.Resources.Limits)
			}
			if len(k.config.Resources.Requests) > 0 {
				container.Resources.Requests = convertToResourceList(k.config.Resources.Requests)
			}
		}
	}

	if k.config.Replicas != nil {
		deployment.Spec.Replicas = k.config.Replicas
	}

	if k.config.Strategy != "" {
		deployment.Spec.Strategy.Type = appsv1.DeploymentStrategyType(k.config.Strategy)
	}
}

func (k *KubernetesStrategy) Rollback(from, to string) error {
	deployment, err := k.clientset.AppsV1().Deployments(k.config.Namespace).Get(context.TODO(), k.config.Deployment, metav1.GetOptions{})
	if err != nil {
		return err
	}

	k.updateDeployment(deployment, to)
	_, err = k.clientset.AppsV1().Deployments(k.config.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	return err
}

func (k *KubernetesStrategy) Deploy(version string) error {
	deployment, err := k.clientset.AppsV1().Deployments(k.config.Namespace).Get(context.TODO(), k.config.Deployment, metav1.GetOptions{})
	if err != nil {
		return err
	}

	k.updateDeployment(deployment, version)
	_, err = k.clientset.AppsV1().Deployments(k.config.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	return err
}

func (k *KubernetesStrategy) GetCurrentVersion() (string, error) {
	deployment, err := k.clientset.AppsV1().Deployments(k.config.Namespace).Get(context.TODO(), k.config.Deployment, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if len(deployment.Spec.Template.Spec.Containers) > 0 {
		return parseVersionFromImage(deployment.Spec.Template.Spec.Containers[0].Image), nil
	}
	return "", fmt.Errorf("no containers found in deployment")
}

func convertToResourceList(resources map[string]string) corev1.ResourceList {
	result := make(corev1.ResourceList)
	for key, value := range resources {
		quantity := resource.MustParse(value)
		result[corev1.ResourceName(key)] = quantity
	}
	return result
}
