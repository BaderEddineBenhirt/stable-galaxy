package rollback

import (
	"errors"
	"testing"
	"time"

	"github.com/BaderEddineBenhirt/stable-galaxy/pkg/logging"
)

type mockStrategy struct {
	rollbackCalls []string
	shouldFail    bool
}

func (m *mockStrategy) Rollback(from, to string) error {
	m.rollbackCalls = append(m.rollbackCalls, from+"->"+to)
	if m.shouldFail {
		return errors.New("mock rollback failed")
	}
	return nil
}

func (m *mockStrategy) Deploy(version string) error {
	return nil
}

func (m *mockStrategy) GetCurrentVersion() (string, error) {
	return "v1.0.0", nil
}

func TestRollbackService(t *testing.T) {
	logger := logging.NewLogger("error", true)

	tests := []struct {
		name        string
		config      RollbackConfig
		versions    []string
		fromVersion string
		mockFail    bool
		wantErr     bool
		wantRetries int
	}{
		{
			name: "successful rollback",
			config: RollbackConfig{
				MaxAttempts:     3,
				BackoffDuration: time.Millisecond,
			},
			versions:    []string{"v0.9.0", "v1.0.0"},
			fromVersion: "v1.0.0",
			mockFail:    false,
			wantErr:     false,
			wantRetries: 1,
		},
		{
			name: "failed rollback with retries",
			config: RollbackConfig{
				MaxAttempts:     3,
				BackoffDuration: time.Millisecond,
			},
			versions:    []string{"v0.9.0", "v1.0.0"},
			fromVersion: "v1.0.0",
			mockFail:    true,
			wantErr:     true,
			wantRetries: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockStrategy{shouldFail: tt.mockFail}
			svc := NewService(tt.config, mock, logger)

			for _, v := range tt.versions {
				svc.RegisterVersion(v)
			}

			err := svc.Rollback(tt.fromVersion)

			if (err != nil) != tt.wantErr {
				t.Errorf("Rollback() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(mock.rollbackCalls) != tt.wantRetries {
				t.Errorf("Expected %d retries, got %d", tt.wantRetries, len(mock.rollbackCalls))
			}
		})
	}
}

func TestRollbackHooks(t *testing.T) {
	logger := logging.NewLogger("error", true)
	preHookCalled := false
	postHookCalled := false
	failureHookCalled := false

	config := RollbackConfig{
		MaxAttempts: 1,
		PreRollbackHook: func() error {
			preHookCalled = true
			return nil
		},
		PostRollbackHook: func() error {
			postHookCalled = true
			return nil
		},
		OnFailureHook: func(err error) {
			failureHookCalled = true
		},
	}

	mock := &mockStrategy{shouldFail: false}
	svc := NewService(config, mock, logger)

	svc.RegisterVersion("v0.9.0")
	svc.RegisterVersion("v1.0.0")

	err := svc.Rollback("v1.0.0")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !preHookCalled {
		t.Error("Pre-rollback hook was not called")
	}
	if !postHookCalled {
		t.Error("Post-rollback hook was not called")
	}
	if failureHookCalled {
		t.Error("Failure hook was called unexpectedly")
	}
}
