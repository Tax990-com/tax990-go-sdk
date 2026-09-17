package tax990

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ResolveURLs_Sandbox(t *testing.T) {
	cfg := Config{ClientID: "c", ClientSecret: "s", UserToken: "t"}
	api, oauth := cfg.resolveURLs()
	assert.Equal(t, "http://localhost:9005", api)
	assert.Equal(t, "http://localhost:4000", oauth)
}

func TestConfig_ResolveURLs_Production(t *testing.T) {
	cfg := Config{ClientID: "c", ClientSecret: "s", UserToken: "t", Environment: Production}
	api, oauth := cfg.resolveURLs()
	assert.Equal(t, "https://api.tax990.com", api)
	assert.Equal(t, "https://oauth.tax990.com", oauth)
}

func TestConfig_ResolveURLs_CustomOverride(t *testing.T) {
	cfg := Config{
		ClientID:     "c",
		ClientSecret: "s",
		UserToken:    "t",
		APIUrl:       "https://custom-api.example.com",
		OAuthUrl:     "https://custom-oauth.example.com",
	}
	api, oauth := cfg.resolveURLs()
	assert.Equal(t, "https://custom-api.example.com", api)
	assert.Equal(t, "https://custom-oauth.example.com", oauth)
}

func TestConfig_ResolveURLs_UnknownEnvironment(t *testing.T) {
	cfg := Config{ClientID: "c", ClientSecret: "s", UserToken: "t", Environment: "unknown"}
	api, oauth := cfg.resolveURLs()
	assert.Equal(t, "http://localhost:9005", api)
	assert.Equal(t, "http://localhost:4000", oauth)
}
