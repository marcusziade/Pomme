package utils

import (
	"fmt"
	"strings"
)

// APIError represents an error from the App Store Connect API
type APIError struct {
	StatusCode int
	ErrorCode  string
	Message    string
	Detail     string
}

// Error implements the error interface
func (e *APIError) Error() string {
	result := fmt.Sprintf("API Error [%s]: %s", e.ErrorCode, e.Message)
	if e.Detail != "" {
		result += fmt.Sprintf(" - %s", e.Detail)
	}
	result += fmt.Sprintf(" (Status: %d)", e.StatusCode)
	return result
}

// NewAPIError creates a new API error
func NewAPIError(statusCode int, errorCode, message, detail string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
		Detail:     detail,
	}
}

// ConfigError represents a configuration error
type ConfigError struct {
	Message string
	Field   string
}

// Error implements the error interface
func (e *ConfigError) Error() string {
	var sb strings.Builder
	sb.WriteString("Configuration error")
	if e.Field != "" {
		sb.WriteString(fmt.Sprintf(" for '%s'", e.Field))
	}
	sb.WriteString(": ")
	sb.WriteString(e.Message)
	return sb.String()
}

// NewConfigError creates a new configuration error
func NewConfigError(message, field string) *ConfigError {
	return &ConfigError{
		Message: message,
		Field:   field,
	}
}

// AuthError represents an authentication error
type AuthError struct {
	Message string
	Cause   error
}

// Error implements the error interface
func (e *AuthError) Error() string {
	result := fmt.Sprintf("Authentication error: %s", e.Message)
	if e.Cause != nil {
		result += fmt.Sprintf(" (%v)", e.Cause)
	}
	return result
}

// NewAuthError creates a new authentication error
func NewAuthError(message string, cause error) *AuthError {
	return &AuthError{
		Message: message,
		Cause:   cause,
	}
}
