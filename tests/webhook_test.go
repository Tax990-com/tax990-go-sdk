package tests

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tax990/sdk-go/tax990/models"
	"github.com/tax990/sdk-go/tax990/utils"
	"github.com/tax990/sdk-go/tests/fixtures"
)

func TestVerifyWebhookSignature_Valid(t *testing.T) {
	// Compute the expected signature from the fixture data.
	mac := hmac.New(sha256.New, []byte(fixtures.WebhookSecret))
	mac.Write([]byte(fixtures.WebhookPayload))
	validSig := hex.EncodeToString(mac.Sum(nil))

	ok := utils.VerifyWebhookSignature(models.WebhookVerifyOptions{
		Secret:    fixtures.WebhookSecret,
		Payload:   fixtures.WebhookPayload,
		Signature: validSig,
	})

	assert.True(t, ok, "valid signature should pass verification")
}

func TestVerifyWebhookSignature_Invalid(t *testing.T) {
	ok := utils.VerifyWebhookSignature(models.WebhookVerifyOptions{
		Secret:    fixtures.WebhookSecret,
		Payload:   fixtures.WebhookPayload,
		Signature: fixtures.WebhookInvalidSignature,
	})

	assert.False(t, ok, "invalid signature should fail verification")
}

func TestVerifyWebhookSignature_EmptySecret(t *testing.T) {
	ok := utils.VerifyWebhookSignature(models.WebhookVerifyOptions{
		Secret:    "",
		Payload:   fixtures.WebhookPayload,
		Signature: fixtures.WebhookValidSignature,
	})

	assert.False(t, ok, "empty secret should not match a real signature")
}

func TestVerifyWebhookSignature_TamperedPayload(t *testing.T) {
	mac := hmac.New(sha256.New, []byte(fixtures.WebhookSecret))
	mac.Write([]byte(fixtures.WebhookPayload))
	validSig := hex.EncodeToString(mac.Sum(nil))

	ok := utils.VerifyWebhookSignature(models.WebhookVerifyOptions{
		Secret:    fixtures.WebhookSecret,
		Payload:   fixtures.WebhookPayload + "tampered",
		Signature: validSig,
	})

	assert.False(t, ok, "tampered payload should not match the original signature")
}
