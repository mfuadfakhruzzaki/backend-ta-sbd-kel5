package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// StandardError represents a standardized error with type, code and metadata
type StandardError struct {
	// Original error
	Err error
	// Type is the error type (e.g., "validation_error", "not_found")
	Type string
	// Message is a human-readable message
	Message string
	// Code is the HTTP status code
	Code int
	// Metadata contains additional error context
	Metadata map[string]interface{}
}

// Error implements the error interface
func (e *StandardError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *StandardError) Unwrap() error {
	return e.Err
}

// New creates a new error with a message
func New(message string) error {
	return errors.New(message)
}

// Wrap wraps an error with additional context
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// ValidationError creates a new validation error
func ValidationError(message string, err error) *StandardError {
	return &StandardError{
		Err:     err,
		Type:    "validation_error",
		Message: message,
		Code:    http.StatusBadRequest,
	}
}

// NotFoundError creates a new not found error
func NotFoundError(message string, err error) *StandardError {
	return &StandardError{
		Err:     err,
		Type:    "not_found_error",
		Message: message,
		Code:    http.StatusNotFound,
	}
}

// UnauthorizedError creates a new unauthorized error
func UnauthorizedError(message string, err error) *StandardError {
	return &StandardError{
		Err:     err,
		Type:    "unauthorized_error",
		Message: message,
		Code:    http.StatusUnauthorized,
	}
}

// ForbiddenError creates a new forbidden error
func ForbiddenError(message string, err error) *StandardError {
	return &StandardError{
		Err:     err,
		Type:    "forbidden_error",
		Message: message,
		Code:    http.StatusForbidden,
	}
}

// InternalError creates a new internal server error
func InternalError(message string, err error) *StandardError {
	return &StandardError{
		Err:     err,
		Type:    "internal_error",
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}

// WithMetadata adds metadata to a StandardError
func (e *StandardError) WithMetadata(key string, value interface{}) *StandardError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}

// StatusCode returns the HTTP status code for the error
func (e *StandardError) StatusCode() int {
	return e.Code
}

// Is compares the error by type
func (e *StandardError) Is(target error) bool {
	t, ok := target.(*StandardError)
	if !ok {
		return false
	}
	return e.Type == t.Type
}

// AsStandardError tries to convert an error to a StandardError
func AsStandardError(err error) (*StandardError, bool) {
	var stdErr *StandardError
	if errors.As(err, &stdErr) {
		return stdErr, true
	}
	return nil, false
} 