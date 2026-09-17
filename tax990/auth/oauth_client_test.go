package auth

import (
	"strings"
	"testing"

	httpclient "github.com/tax990/sdk-go/tax990/http"
	"github.com/stretchr/testify/assert"
)

func TestSignJWSLocally(t *testing.T) {
	http := httpclient.NewClient(httpclient.Options{BaseURL: "http://localhost:4000"})
	client := NewOAuthClient(http, "test-client-id", "test-client-secret", "test-user-token")

	jws := client.SignJWSLocally()

	parts := strings.Split(jws, ".")
	assert.Equal(t, 3, len(parts), "JWS must have 3 parts: header.payload.signature")
	for _, part := range parts {
		assert.NotEmpty(t, part)
	}
}

func TestSignJWSLocally_Deterministic(t *testing.T) {
	http := httpclient.NewClient(httpclient.Options{BaseURL: "http://localhost:4000"})
	client := NewOAuthClient(http, "id", "secret", "token")

	jws1 := client.SignJWSLocally()
	jws2 := client.SignJWSLocally()

	parts1 := strings.Split(jws1, ".")
	parts2 := strings.Split(jws2, ".")
	assert.Equal(t, parts1[0], parts2[0], "header should be the same")
}
