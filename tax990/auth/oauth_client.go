// Package auth implements the Tax990 two-step OAuth token flow.
package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	apierrors "github.com/tax990/sdk-go/tax990/errors"
	httpclient "github.com/tax990/sdk-go/tax990/http"
	"github.com/tax990/sdk-go/tax990/models"
)

// OAuthClient implements the Tax990 two-step token flow:
//  1. Sign a JWS locally using HS256 (claims: iss/sub/aud/iat, key: ClientSecret)
//  2. GET /Auth/GetTax990Token with `authentication: <JWS>` → RS256 access token
type OAuthClient struct {
	http         *httpclient.Client
	clientID     string
	clientSecret string
	userToken    string
}

// NewOAuthClient creates an OAuthClient that uses the given HTTP client and credentials.
func NewOAuthClient(http *httpclient.Client, clientID, clientSecret, userToken string) *OAuthClient {
	return &OAuthClient{
		http:         http,
		clientID:     clientID,
		clientSecret: clientSecret,
		userToken:    userToken,
	}
}

// SignJWSLocally produces a HS256 JWT using the ClientSecret.
// Claims: iss=clientId, sub=clientId, aud=userToken, iat=now.
func (c *OAuthClient) SignJWSLocally() string {
	headerJSON := `{"alg":"HS256","typ":"JWT"}`
	headerEncoded := base64.RawURLEncoding.EncodeToString([]byte(headerJSON))

	payload := map[string]any{
		"iss": c.clientID,
		"sub": c.clientID,
		"aud": c.userToken,
		"iat": time.Now().Unix(),
	}
	payloadBytes, _ := json.Marshal(payload)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)

	signingInput := headerEncoded + "." + payloadEncoded

	mac := hmac.New(sha256.New, []byte(c.clientSecret))
	mac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + signature
}

// GenerateJWSFromServer calls POST /Auth/GenerateJWS to obtain the JWS from the server.
// Use this when you prefer the server to produce the token rather than signing locally.
func (c *OAuthClient) GenerateJWSFromServer(ctx context.Context) (string, error) {
	resp, err := httpclient.Post[models.GenerateJWSResponse](
		ctx, c.http, "/Auth/GenerateJWS",
		models.GenerateJWSRequest{
			ClientId:       c.clientID,
			ClientSecretId: c.clientSecret,
			UserToken:      c.userToken,
		},
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("generate jws from server: %w", err)
	}
	return resp.Response.JWSToken, nil
}

// GetAccessToken performs the full two-step token acquisition:
//  1. Sign JWS locally with HS256.
//  2. GET /Auth/GetTax990Token → AccessToken (RS256, expires in ~3600s).
//
// Returns the access token string and its lifetime in seconds.
func (c *OAuthClient) GetAccessToken(ctx context.Context) (accessToken string, expiresIn int, err error) {
	jws := c.SignJWSLocally()

	resp, err := httpclient.Get[models.Tax990TokenResponse](
		ctx, c.http, "/Auth/GetTax990Token",
		nil,
		map[string]string{"authentication": jws},
	)
	if err != nil {
		return "", 0, fmt.Errorf("get access token: %w", err)
	}

	if resp.Response.Errors != nil {
		return "", 0, apierrors.NewAuthError(resp.Response.Errors.ErrorMessage, "")
	}

	ttl := resp.Response.ExpiresIn
	if ttl == 0 {
		ttl = 3600
	}
	return resp.Response.AccessToken, ttl, nil
}
