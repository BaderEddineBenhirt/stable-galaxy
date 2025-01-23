package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/BaderEddineBenhirt/stable-galaxy/pkg/deployment"
	"github.com/BaderEddineBenhirt/stable-galaxy/pkg/rollback"
)

func main() {
	dockerConfig := deployment.DockerConfig{
		ServiceName:   os.Getenv("DOCKER_SERVICE_NAME"),
		Registry:      os.Getenv("DOCKER_REGISTRY"),
		NetworkMode:   os.Getenv("DOCKER_NETWORK_MODE"),
		ImageTemplate: os.Getenv("DOCKER_IMAGE_TEMPLATE"),
		Labels:        parseMapFromEnv("DOCKER_LABELS"),
		EnvVars:       parseMapFromEnv("DOCKER_ENV_VARS"),
		CustomArgs:    parseMapFromEnv("DOCKER_CUSTOM_ARGS"),
	}

	dockerStrat := deployment.NewDockerStrategy(dockerConfig)
	dockerRollback := rollback.NewService(buildRollbackConfig(), dockerStrat)

	if err := dockerRollback.Rollback(os.Getenv("ROLLBACK_FROM_VERSION")); err != nil {
		log.Printf("Docker rollback failed: %v", err)
	}

	clientset := setupKubernetesClient()
	k8sConfig := deployment.KubernetesConfig{
		Namespace:     os.Getenv("K8S_NAMESPACE"),
		Deployment:    os.Getenv("K8S_DEPLOYMENT"),
		ImageTemplate: os.Getenv("K8S_IMAGE_TEMPLATE"),
		Labels:        parseMapFromEnv("K8S_LABELS"),
		Annotations:   parseMapFromEnv("K8S_ANNOTATIONS"),
		Strategy:      os.Getenv("K8S_STRATEGY"),
		Context:       os.Getenv("K8S_CONTEXT"),
	}

	k8sStrat := deployment.NewKubernetesStrategy(clientset, k8sConfig)
	k8sRollback := rollback.NewService(buildRollbackConfig(), k8sStrat)

	if err := k8sRollback.Rollback(os.Getenv("ROLLBACK_FROM_VERSION")); err != nil {
		log.Printf("K8s rollback failed: %v", err)
	}
}

func buildRollbackConfig() rollback.RollbackConfig {
	config := rollback.DefaultConfig()
	config.MaxAttempts = getEnvInt("ROLLBACK_MAX_ATTEMPTS", 3)
	config.DryRun = getEnvBool("ROLLBACK_DRY_RUN", false)
	config.LogLevel = os.Getenv("ROLLBACK_LOG_LEVEL")
	config.HealthCheck.URL = os.Getenv("HEALTH_CHECK_URL")
	return config
}

func setupKubernetesClient() *kubernetes.Clientset {
	var config *rest.Config
	var err error

	if os.Getenv("K8S_IN_CLUSTER") == "true" {
		config, err = rest.InClusterConfig()
	} else {
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		log.Fatalf("Failed to create k8s config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create k8s clientset: %v", err)
	}

	return clientset
}
func parseMapFromEnv(prefix string) map[string]string {
	result := make(map[string]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 && strings.HasPrefix(pair[0], prefix) {
			key := strings.TrimPrefix(pair[0], prefix+"_")
			result[key] = pair[1]
		}
	}
	return result
}

func getEnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return intVal
}

func getEnvBool(key string, defaultVal bool) bool {
	val := strings.ToLower(os.Getenv(key))
	switch val {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return defaultVal
	}
}
