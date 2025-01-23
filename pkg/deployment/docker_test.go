package deployment

import (
	"testing"
)

func TestDockerStrategy_BuildImageTag(t *testing.T) {
	tests := []struct {
		name     string
		config   DockerConfig
		version  string
		expected string
	}{
		{
			name: "standard image tag",
			config: DockerConfig{
				ServiceName: "myapp",
				Registry:    "registry.example.com",
			},
			version:  "v1.0.0",
			expected: "registry.example.com/myapp:v1.0.0",
		},
		{
			name: "custom template",
			config: DockerConfig{
				ImageTemplate: "custom-registry.com/%s-service",
			},
			version:  "v1.0.0",
			expected: "custom-registry.com/v1.0.0-service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDockerStrategy(tt.config)
			result := d.buildImageTag(tt.version)
			if result != tt.expected {
				t.Errorf("buildImageTag() = %v, want %v", result, tt.expected)
			}
		})
	}
}
