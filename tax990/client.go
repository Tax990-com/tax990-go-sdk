package tax990

import (
	"fmt"
	"time"

	"github.com/tax990/sdk-go/tax990/auth"
	httpclient "github.com/tax990/sdk-go/tax990/http"
	"github.com/tax990/sdk-go/tax990/resources"
)

// Tax990Client is the top-level SDK entry point.
// Create one with NewClient and use its fields to access each API resource.
//
//	client, err := tax990.NewClient(tax990.Config{
//	    ClientID:     os.Getenv("TAX990_CLIENT_ID"),
//	    ClientSecret: os.Getenv("TAX990_CLIENT_SECRET"),
//	    UserToken:    os.Getenv("TAX990_USER_TOKEN"),
//	    Environment:  tax990.Production,
//	})
type Tax990Client struct {
	// Form990N provides Form 990-N filing operations (create, update, get, list,
	// delete, validate, transmit, status, getPDF).
	Form990N resources.Form990NResource

	// Organizations wraps Form990N endpoints for business entity queries.
	Organizations resources.OrganizationResource

	// FilingStatus wraps the /status endpoint for IRS acknowledgement queries.
	FilingStatus resources.FilingStatusResource

	// Webhooks is a stub — webhook endpoints are not yet in the Public API.
	Webhooks resources.WebhookResource

	// ApiKeys is a stub — API key endpoints are not yet in the Public API.
	ApiKeys resources.ApiKeysResource
}

// NewClient validates the configuration and returns an initialised Tax990Client.
// Returns an error if any required credential is missing.
func NewClient(cfg Config) (*Tax990Client, error) {
	if cfg.ClientID == "" {
		return nil, fmt.Errorf("tax990: ClientID is required")
	}
	if cfg.ClientSecret == "" {
		return nil, fmt.Errorf("tax990: ClientSecret is required")
	}
	if cfg.UserToken == "" {
		return nil, fmt.Errorf("tax990: UserToken is required")
	}

	apiURL, oauthURL := cfg.resolveURLs()

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	// OAuth HTTP client — no bearer token, used only to acquire the token.
	oauthHTTP := httpclient.NewClient(httpclient.Options{
		BaseURL: oauthURL,
		Timeout: timeout,
	})

	oauthClient := auth.NewOAuthClient(oauthHTTP, cfg.ClientID, cfg.ClientSecret, cfg.UserToken)
	tokenManager := auth.NewTokenManager(oauthClient)

	// API HTTP client — bearer token injected automatically via tokenManager.
	apiHTTP := httpclient.NewClient(httpclient.Options{
		BaseURL:  apiURL,
		Timeout:  timeout,
		GetToken: tokenManager.GetToken,
	})

	return &Tax990Client{
		Form990N:      resources.NewForm990N(apiHTTP),
		Organizations: resources.NewOrganization(apiHTTP),
		FilingStatus:  resources.NewFilingStatus(apiHTTP),
		Webhooks:      resources.NewWebhook(),
		ApiKeys:       resources.NewApiKeys(),
	}, nil
}
