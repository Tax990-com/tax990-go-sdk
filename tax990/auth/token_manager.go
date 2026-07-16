package auth

import (
	"context"
	"sync"
	"time"

	"github.com/tax990/sdk-go/tax990/models"
)

// expiryBuffer is deducted from the token's lifetime before considering it expired,
// ensuring a refresh happens before the token actually expires.
const expiryBuffer = 30 * time.Second

// TokenManager caches an access token and auto-refreshes it before expiry.
// Concurrent callers that arrive while a refresh is in-flight all wait for
// the single in-flight refresh rather than spawning multiple token requests.
type TokenManager struct {
	mu          sync.Mutex
	token       *models.StoredToken
	oauthClient *OAuthClient
}

// NewTokenManager creates a TokenManager backed by the given OAuthClient.
func NewTokenManager(oauthClient *OAuthClient) *TokenManager {
	return &TokenManager{oauthClient: oauthClient}
}

// GetToken returns the cached access token, refreshing it if it is expired or absent.
func (tm *TokenManager) GetToken(ctx context.Context) (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.token != nil && !tm.isExpired(tm.token) {
		return tm.token.AccessToken, nil
	}

	accessToken, expiresIn, err := tm.oauthClient.GetAccessToken(ctx)
	if err != nil {
		return "", err
	}

	tm.token = &models.StoredToken{
		AccessToken: accessToken,
		ExpiresAt:   time.Now().Add(time.Duration(expiresIn)*time.Second - expiryBuffer),
	}
	return tm.token.AccessToken, nil
}

// ClearToken evicts the cached token, forcing the next GetToken call to fetch a new one.
// This is primarily useful in tests.
func (tm *TokenManager) ClearToken() {
	tm.mu.Lock()
	tm.token = nil
	tm.mu.Unlock()
}

func (tm *TokenManager) isExpired(t *models.StoredToken) bool {
	return time.Now().After(t.ExpiresAt)
}
