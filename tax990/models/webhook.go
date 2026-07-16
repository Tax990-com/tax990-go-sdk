package models

// WebhookVerifyOptions holds the inputs for HMAC-SHA256 webhook signature verification.
type WebhookVerifyOptions struct {
	// Secret is the shared webhook signing secret.
	Secret string
	// Payload is the raw request body string to verify.
	Payload string
	// Signature is the hex-encoded HMAC-SHA256 digest to compare against.
	Signature string
}
