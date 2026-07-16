package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tax990/sdk-go/tax990/auth"
	apierrors "github.com/tax990/sdk-go/tax990/errors"
	httpclient "github.com/tax990/sdk-go/tax990/http"
	"github.com/tax990/sdk-go/tests/fixtures"
)

func TestOAuthClient_SignJWSLocally(t *testing.T) {
	oauthHTTP := httpclient.NewClient(httpclient.Options{BaseURL: "http://unused"})
	oc := auth.NewOAuthClient(oauthHTTP, fixtures.ClientID, fixtures.ClientSecret, fixtures.UserToken)

	jws := oc.SignJWSLocally()

	parts := strings.Split(jws, ".")
	assert.Len(t, parts, 3, "JWS must have header.payload.signature format")
	assert.NotEmpty(t, parts[0])
	assert.NotEmpty(t, parts[1])
	assert.NotEmpty(t, parts[2])
}

func TestOAuthClient_GetAccessToken_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/Auth/GetTax990Token", r.URL.Path)
		assert.NotEmpty(t, r.Header.Get("authentication"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fixtures.TokenResponseBody))
	}))
	defer server.Close()

	oauthHTTP := httpclient.NewClient(httpclient.Options{BaseURL: server.URL})
	oc := auth.NewOAuthClient(oauthHTTP, fixtures.ClientID, fixtures.ClientSecret, fixtures.UserToken)

	token, expiresIn, err := oc.GetAccessToken(context.Background())

	require.NoError(t, err)
	assert.Equal(t, fixtures.AccessToken, token)
	assert.Equal(t, 3600, expiresIn)
}

func TestOAuthClient_GetAccessToken_ErrorInBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Server returns 200 but with an error payload in the body.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fixtures.TokenResponseBodyWithError))
	}))
	defer server.Close()

	oauthHTTP := httpclient.NewClient(httpclient.Options{BaseURL: server.URL})
	oc := auth.NewOAuthClient(oauthHTTP, fixtures.ClientID, fixtures.ClientSecret, fixtures.UserToken)

	_, _, err := oc.GetAccessToken(context.Background())

	require.Error(t, err)
	var authErr *apierrors.AuthError
	assert.ErrorAs(t, err, &authErr)
}

func TestOAuthClient_GenerateJWSFromServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/Auth/GenerateJWS", r.URL.Path)

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, fixtures.ClientID, body["ClientId"])
		assert.Equal(t, fixtures.ClientSecret, body["ClientSecretId"])
		assert.Equal(t, fixtures.UserToken, body["UserToken"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fixtures.GenerateJWSResponseBody))
	}))
	defer server.Close()

	oauthHTTP := httpclient.NewClient(httpclient.Options{BaseURL: server.URL})
	oc := auth.NewOAuthClient(oauthHTTP, fixtures.ClientID, fixtures.ClientSecret, fixtures.UserToken)

	jws, err := oc.GenerateJWSFromServer(context.Background())

	require.NoError(t, err)
	assert.Equal(t, fixtures.JWSToken, jws)
}

func TestTokenManager_CachesToken(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fixtures.TokenResponseBody))
	}))
	defer server.Close()

	oauthHTTP := httpclient.NewClient(httpclient.Options{BaseURL: server.URL})
	oc := auth.NewOAuthClient(oauthHTTP, fixtures.ClientID, fixtures.ClientSecret, fixtures.UserToken)
	tm := auth.NewTokenManager(oc)

	ctx := context.Background()
	t1, err1 := tm.GetToken(ctx)
	t2, err2 := tm.GetToken(ctx)

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.Equal(t, fixtures.AccessToken, t1)
	assert.Equal(t, t1, t2)
	assert.Equal(t, 1, callCount, "second GetToken should use cached token")
}

func TestTokenManager_ClearToken_ForcesRefresh(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fixtures.TokenResponseBody))
	}))
	defer server.Close()

	oauthHTTP := httpclient.NewClient(httpclient.Options{BaseURL: server.URL})
	oc := auth.NewOAuthClient(oauthHTTP, fixtures.ClientID, fixtures.ClientSecret, fixtures.UserToken)
	tm := auth.NewTokenManager(oc)

	ctx := context.Background()
	_, _ = tm.GetToken(ctx)
	tm.ClearToken()
	_, _ = tm.GetToken(ctx)

	assert.Equal(t, 2, callCount, "ClearToken should force a new token fetch")
}

func TestHTTPClient_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"StatusCode":401,"StatusNm":"Unauthorized","StatusMessage":"Invalid token.","CorrelationId":"corr-xyz"}`))
	}))
	defer server.Close()

	c := httpclient.NewClient(httpclient.Options{BaseURL: server.URL})
	_, err := httpclient.Get[map[string]any](context.Background(), c, "/test", nil, nil)

	require.Error(t, err)
	var authErr *apierrors.AuthError
	assert.ErrorAs(t, err, &authErr)
	assert.Equal(t, 401, authErr.StatusCode)
}
