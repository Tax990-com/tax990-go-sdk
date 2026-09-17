package tax990

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient_MissingClientID(t *testing.T) {
	_, err := NewClient(Config{ClientSecret: "s", UserToken: "t"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ClientID")
}

func TestNewClient_MissingClientSecret(t *testing.T) {
	_, err := NewClient(Config{ClientID: "c", UserToken: "t"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ClientSecret")
}

func TestNewClient_MissingUserToken(t *testing.T) {
	_, err := NewClient(Config{ClientID: "c", ClientSecret: "s"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UserToken")
}

func TestNewClient_Success(t *testing.T) {
	client, err := NewClient(Config{ClientID: "c", ClientSecret: "s", UserToken: "t"})
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.Form990N)
	assert.NotNil(t, client.Organizations)
	assert.NotNil(t, client.FilingStatus)
	assert.NotNil(t, client.Utility)
	assert.NotNil(t, client.Nonprofits)
}
