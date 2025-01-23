/*
Package deployment implements deployment strategies for different platforms.

It provides a common interface for managing deployments across Docker and Kubernetes
environments, allowing for consistent rollback behavior regardless of the underlying
platform.

Example Docker usage:

	config := DockerConfig{
	    ServiceName: "myapp",
	    Registry:   "registry.example.com",
	}
	strategy := NewDockerStrategy(config)

Example Kubernetes usage:

	config := KubernetesConfig{
	    Namespace:  "default",
	    Deployment: "myapp",
	}
	strategy := NewKubernetesStrategy(clientset, config)
*/
package deployment
