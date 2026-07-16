// Package apierrors defines typed errors returned by the Tax990 SDK.
package apierrors

import (
	"fmt"

	"github.com/tax990/sdk-go/tax990/models"
)

// Tax990Error is the base error type for all Tax990 API errors.
type Tax990Error struct {
	Code       string
	Message    string
	StatusCode int
	RequestID  string
}

func (e *Tax990Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// AuthError is returned when authentication fails (HTTP 401).
type AuthError struct {
	Tax990Error
}

// ValidationError is returned when request validation fails (HTTP 400).
// Errors contains the structured field-level error details from the API.
type ValidationError struct {
	Tax990Error
	Errors []models.StructuredError
}

// RateLimitError is returned when the rate limit is exceeded (HTTP 429).
type RateLimitError struct {
	Tax990Error
}

// NotFoundError is returned when the requested resource is not found (HTTP 404).
type NotFoundError struct {
	Tax990Error
}

// NewTax990Error creates a Tax990Error with the given code, message, status, and request ID.
func NewTax990Error(code, message string, statusCode int, requestID string) *Tax990Error {
	return &Tax990Error{Code: code, Message: message, StatusCode: statusCode, RequestID: requestID}
}

// NewAuthError creates an AuthError for authentication failures.
func NewAuthError(message, requestID string) *AuthError {
	return &AuthError{Tax990Error{Code: "AUTH_ERROR", Message: message, StatusCode: 401, RequestID: requestID}}
}

// NewValidationError creates a ValidationError with structured error details.
func NewValidationError(message, requestID string, errs []models.StructuredError) *ValidationError {
	return &ValidationError{
		Tax990Error: Tax990Error{Code: "VALIDATION_ERROR", Message: message, StatusCode: 400, RequestID: requestID},
		Errors:      errs,
	}
}

// NewRateLimitError creates a RateLimitError for rate-limit responses.
func NewRateLimitError(message, requestID string) *RateLimitError {
	return &RateLimitError{Tax990Error{Code: "RATE_LIMIT_ERROR", Message: message, StatusCode: 429, RequestID: requestID}}
}

// NewNotFoundError creates a NotFoundError for missing resources.
func NewNotFoundError(message, requestID string) *NotFoundError {
	return &NotFoundError{Tax990Error{Code: "NOT_FOUND", Message: message, StatusCode: 404, RequestID: requestID}}
}
