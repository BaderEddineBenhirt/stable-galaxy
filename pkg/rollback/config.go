package rollback

import (
	"context"
	"time"
)

type RollbackConfig struct {
	MaxAttempts     int
	BackoffDuration time.Duration
	Timeout         time.Duration
	Context         context.Context
	DryRun          bool

	DeploymentConfig struct {
		Type          string
		CustomOptions map[string]string
	}

	PreRollbackHook  func() error
	PostRollbackHook func() error
	OnFailureHook    func(error)

	ValidateVersion    func(string) bool
	VersionConstraints struct {
		MinVersion string
		MaxVersion string
		Blacklist  []string
	}

	MetricsEnabled bool
	HealthCheck    struct {
		URL           string
		Timeout       time.Duration
		RetryAttempts int
		SuccessStatus int
		CustomHeaders map[string]string
	}

	Notifications struct {
		Enabled  bool
		Channels []string
		Webhook  string
	}

	LogLevel    string
	LogFilePath string
}

func DefaultConfig() RollbackConfig {
	return RollbackConfig{
		MaxAttempts:     3,
		BackoffDuration: time.Second * 5,
		Timeout:         time.Minute * 5,
		Context:         context.Background(),
		DryRun:          false,
		DeploymentConfig: struct {
			Type          string
			CustomOptions map[string]string
		}{
			CustomOptions: make(map[string]string),
		},
		HealthCheck: struct {
			URL           string
			Timeout       time.Duration
			RetryAttempts int
			SuccessStatus int
			CustomHeaders map[string]string
		}{
			Timeout:       time.Second * 30,
			RetryAttempts: 3,
			SuccessStatus: 200,
			CustomHeaders: make(map[string]string),
		},
		Notifications: struct {
			Enabled  bool
			Channels []string
			Webhook  string
		}{
			Enabled:  true,
			Channels: []string{"slack"},
		},
		LogLevel: "info",
	}
}
