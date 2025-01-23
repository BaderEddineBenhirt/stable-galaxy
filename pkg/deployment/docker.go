package deployment

import (
	"fmt"
	"os/exec"
)

type DockerConfig struct {
	ServiceName   string
	Registry      string
	NetworkMode   string
	Constraints   []string
	Labels        map[string]string
	EnvVars       map[string]string
	CustomArgs    map[string]string
	ImageTemplate string
	ConfigPath    string
}

type DockerStrategy struct {
	config DockerConfig
}

func NewDockerStrategy(config DockerConfig) *DockerStrategy {
	return &DockerStrategy{
		config: config,
	}
}

func (d *DockerStrategy) buildImageTag(version string) string {
	if d.config.ImageTemplate != "" {
		return fmt.Sprintf(d.config.ImageTemplate, version)
	}
	return fmt.Sprintf("%s/%s:%s", d.config.Registry, d.config.ServiceName, version)
}

func (d *DockerStrategy) buildUpdateCommand(imageTag string) *exec.Cmd {
	args := []string{"service", "update", "--image", imageTag}

	if d.config.NetworkMode != "" {
		args = append(args, "--network", d.config.NetworkMode)
	}

	for _, constraint := range d.config.Constraints {
		args = append(args, "--constraint", constraint)
	}

	for k, v := range d.config.Labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", k, v))
	}

	for k, v := range d.config.EnvVars {
		args = append(args, "--env", fmt.Sprintf("%s=%s", k, v))
	}

	for k, v := range d.config.CustomArgs {
		args = append(args, k, v)
	}

	args = append(args, d.config.ServiceName)

	return exec.Command("docker", args...)
}

func (d *DockerStrategy) Rollback(from, to string) error {
	imageTag := d.buildImageTag(to)
	cmd := d.buildUpdateCommand(imageTag)
	return cmd.Run()
}

func (d *DockerStrategy) Deploy(version string) error {
	imageTag := d.buildImageTag(version)
	cmd := d.buildUpdateCommand(imageTag)
	return cmd.Run()
}

func (d *DockerStrategy) GetCurrentVersion() (string, error) {
	cmd := exec.Command("docker", "service", "inspect", "--format", "{{.Spec.TaskTemplate.ContainerSpec.Image}}", d.config.ServiceName)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return parseVersionFromImage(string(output)), nil
}
