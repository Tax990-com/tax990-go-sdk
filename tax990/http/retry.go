// Package httpclient provides a net/http wrapper with retry and timeout support.
package httpclient

import (
	"math"
	"time"
)

// retryPolicy decides whether a failed request should be retried.
// Retries on network errors (status 0), 5xx, and 429.
func retryPolicy(statusCode int, attempt int) bool {
	if attempt >= maxRetries {
		return false
	}
	return statusCode == 0 || statusCode == 429 || statusCode >= 500
}

// retryDelay returns the exponential backoff duration for the given attempt (0-indexed).
// Delay = 100ms * 2^attempt, capped at 2 seconds.
func retryDelay(attempt int) time.Duration {
	ms := 100 * math.Pow(2, float64(attempt))
	if ms > 2000 {
		ms = 2000
	}
	return time.Duration(ms) * time.Millisecond
}

const maxRetries = 3
