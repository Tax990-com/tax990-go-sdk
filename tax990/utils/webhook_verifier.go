package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/tax990/sdk-go/tax990/models"
)

// VerifyWebhookSignature validates an HMAC-SHA256 webhook signature using a
// timing-safe comparison to prevent timing attacks.
//
// The signature must be the lowercase hex-encoded HMAC-SHA256 digest of the
// raw request payload, keyed with the shared secret.
func VerifyWebhookSignature(opts models.WebhookVerifyOptions) bool {
	mac := hmac.New(sha256.New, []byte(opts.Secret))
	mac.Write([]byte(opts.Payload))
	expected := hex.EncodeToString(mac.Sum(nil))

	// Constant-time comparison prevents timing side-channel attacks.
	return hmac.Equal([]byte(expected), []byte(opts.Signature))
}
