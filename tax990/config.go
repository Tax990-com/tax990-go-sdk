// Package tax990 is the entry point for the Tax990 Go SDK.
package tax990

import (
	"os"
	"time"
)

// Environment selects the API deployment target.
type Environment string

const (
	// Production targets the live Tax990 API.
	Production Environment = "production"
	// Sandbox targets the local development server.
	Sandbox Environment = "sandbox"
)

// Config holds all configuration for the Tax990Client.
type Config struct {
	// ClientID is the OAuth client identifier (TAX990_CLIENT_ID).
	ClientID string
	// ClientSecret is the OAuth client secret used to sign JWS tokens (TAX990_CLIENT_SECRET).
	ClientSecret string
	// UserToken is the OAuth audience token for this client (TAX990_USER_TOKEN).
	UserToken string
	// Environment selects between Production and Sandbox endpoints.
	// Defaults to Sandbox when empty.
	Environment Environment
	// APIUrl overrides the default API base URL for the chosen environment.
	APIUrl string
	// OAuthUrl overrides the default OAuth base URL for the chosen environment.
	OAuthUrl string
	// Timeout overrides the default 30-second HTTP request timeout.
	Timeout time.Duration
}

// resolveURLs returns the effective API and OAuth base URLs for the config.
func (c *Config) resolveURLs() (apiURL, oauthURL string) {
	apiURL = os.Getenv("TAX990_API_URL")
	oauthURL = os.Getenv("TAX990_OAUTH_URL")
	if c.APIUrl != "" {
		apiURL = c.APIUrl
	}
	if c.OAuthUrl != "" {
		oauthURL = c.OAuthUrl
	}
	return apiURL, oauthURL
}
