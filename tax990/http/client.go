package httpclient

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"
	"net/url"
	"strings"
	"time"

	apierrors "github.com/tax990/sdk-go/tax990/errors"
	"github.com/tax990/sdk-go/tax990/models"
)

// Options configures the HTTP client.
type Options struct {
	// BaseURL is the scheme+host prefix for all requests.
	BaseURL string
	// Timeout overrides the default 30-second request timeout.
	Timeout time.Duration
	// GetToken, when set, returns a Bearer token that is attached to every request.
	GetToken func(ctx context.Context) (string, error)
}

// Client is a thin net/http wrapper with retry, correlation-ID injection,
// and automatic Bearer token attachment.
type Client struct {
	baseURL    string
	httpClient *nethttp.Client
	getToken   func(ctx context.Context) (string, error)
}

// NewClient creates a new Client with the given options.
func NewClient(opts Options) *Client {
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &Client{
		baseURL:    strings.TrimRight(opts.BaseURL, "/"),
		httpClient: &nethttp.Client{Timeout: timeout},
		getToken:   opts.GetToken,
	}
}

// Get executes a GET request and decodes the JSON response into T.
func Get[T any](ctx context.Context, c *Client, path string, params map[string]string, extraHeaders map[string]string) (T, error) {
	return do[T](ctx, c, nethttp.MethodGet, path, params, nil, extraHeaders)
}

// Post executes a POST request with a JSON body and decodes the JSON response into T.
func Post[T any](ctx context.Context, c *Client, path string, body any, extraHeaders map[string]string) (T, error) {
	return do[T](ctx, c, nethttp.MethodPost, path, nil, body, extraHeaders)
}

// Delete executes a DELETE request and decodes the JSON response into T.
func Delete[T any](ctx context.Context, c *Client, path string, params map[string]string, extraHeaders map[string]string) (T, error) {
	return do[T](ctx, c, nethttp.MethodDelete, path, params, nil, extraHeaders)
}

// do is the core request executor with retry logic.
func do[T any](
	ctx context.Context,
	c *Client,
	method, path string,
	params map[string]string,
	body any,
	extraHeaders map[string]string,
) (T, error) {
	var zero T
	var lastStatus int

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return zero, ctx.Err()
			case <-time.After(retryDelay(attempt - 1)):
			}
		}

		raw, status, err := c.execute(ctx, method, path, params, body, extraHeaders)
		lastStatus = status
		if err != nil {
			if retryPolicy(status, attempt) {
				continue
			}
			return zero, err
		}

		var result T
		if jsonErr := json.Unmarshal(raw, &result); jsonErr != nil {
			return zero, fmt.Errorf("tax990: decode response: %w", jsonErr)
		}
		return result, nil
	}

	return zero, apierrors.NewTax990Error(
		fmt.Sprintf("%d", lastStatus),
		"request failed after retries",
		lastStatus,
		"",
	)
}

// execute sends a single HTTP request and returns the raw body bytes plus the status code.
func (c *Client) execute(
	ctx context.Context,
	method, path string,
	params map[string]string,
	body any,
	extraHeaders map[string]string,
) ([]byte, int, error) {
	fullURL := c.baseURL + path
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			if v != "" {
				q.Set(k, v)
			}
		}
		if encoded := q.Encode(); encoded != "" {
			fullURL += "?" + encoded
		}
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("tax990: marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := nethttp.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("tax990: build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-correlation-id", newUUID())

	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

	if c.getToken != nil {
		token, err := c.getToken(ctx)
		if err != nil {
			return nil, 0, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("tax990: http request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("tax990: read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, mapHTTPError(resp.StatusCode, raw)
	}
	return raw, resp.StatusCode, nil
}

// mapHTTPError parses an error response body into the appropriate typed error.
func mapHTTPError(statusCode int, body []byte) error {
	var envelope struct {
		StatusMessage string                  `json:"StatusMessage"`
		CorrelationId string                  `json:"CorrelationId"`
		Message       string                  `json:"message"`
		Errors        []models.StructuredError `json:"Errors"`
	}
	_ = json.Unmarshal(body, &envelope)

	msg := envelope.StatusMessage
	if msg == "" {
		msg = envelope.Message
	}
	if msg == "" {
		msg = fmt.Sprintf("HTTP %d", statusCode)
	}
	rid := envelope.CorrelationId

	switch statusCode {
	case 401:
		return apierrors.NewAuthError(msg, rid)
	case 404:
		return apierrors.NewNotFoundError(msg, rid)
	case 429:
		return apierrors.NewRateLimitError(msg, rid)
	case 400:
		if len(envelope.Errors) > 0 {
			return apierrors.NewValidationError(msg, rid, envelope.Errors)
		}
		return apierrors.NewTax990Error("VALIDATION_ERROR", msg, 400, rid)
	default:
		return apierrors.NewTax990Error(fmt.Sprintf("%d", statusCode), msg, statusCode, rid)
	}
}

// newUUID generates a random UUID v4.
func newUUID() string {
	var b [16]byte
	rand.Read(b[:]) //nolint:errcheck // crypto/rand.Read never returns an error on supported platforms
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
