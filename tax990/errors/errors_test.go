package apierrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthError(t *testing.T) {
	err := NewAuthError("auth failed", "req-123")
	assert.Equal(t, "AUTH_ERROR", err.Code)
	assert.Equal(t, 401, err.StatusCode)
	assert.Equal(t, "req-123", err.RequestID)
	assert.Contains(t, err.Error(), "auth failed")
}

func TestRateLimitError(t *testing.T) {
	err := NewRateLimitError("too many requests", "req-456")
	assert.Equal(t, "RATE_LIMIT_ERROR", err.Code)
	assert.Equal(t, 429, err.StatusCode)
}

func TestNotFoundError(t *testing.T) {
	err := NewNotFoundError("not found", "req-789")
	assert.Equal(t, "NOT_FOUND", err.Code)
	assert.Equal(t, 404, err.StatusCode)
}

func TestValidationError(t *testing.T) {
	err := NewValidationError("invalid", "req-101", nil)
	assert.Equal(t, "VALIDATION_ERROR", err.Code)
	assert.Equal(t, 400, err.StatusCode)
}
