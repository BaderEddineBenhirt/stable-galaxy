package errors

import (
	"fmt"
)

type ErrorType string

const (
	ErrorTypeValidation    ErrorType = "ValidationError"
	ErrorTypeDeployment    ErrorType = "DeploymentError"
	ErrorTypeHealthCheck   ErrorType = "HealthCheckError"
	ErrorTypeConfiguration ErrorType = "ConfigurationError"
	ErrorTypeNetwork       ErrorType = "NetworkError"
)

type RollbackError struct {
	Type    ErrorType
	Message string
	Cause   error
	Meta    map[string]interface{}
}

func (e *RollbackError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func NewValidationError(msg string, cause error) *RollbackError {
	return &RollbackError{
		Type:    ErrorTypeValidation,
		Message: msg,
		Cause:   cause,
	}
}

func NewDeploymentError(msg string, cause error, meta map[string]interface{}) *RollbackError {
	return &RollbackError{
		Type:    ErrorTypeDeployment,
		Message: msg,
		Cause:   cause,
		Meta:    meta,
	}
}

func NewHealthCheckError(msg string, cause error) *RollbackError {
	return &RollbackError{
		Type:    ErrorTypeHealthCheck,
		Message: msg,
		Cause:   cause,
	}
}
