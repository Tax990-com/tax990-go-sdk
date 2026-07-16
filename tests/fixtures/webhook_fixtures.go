package fixtures

// WebhookSecret is the shared signing secret used in webhook tests.
const WebhookSecret = "webhook-test-secret-key"

// WebhookPayload is the raw JSON body used in webhook tests.
const WebhookPayload = `{"event":"filing.accepted","submissionId":"sub-0001"}`

// WebhookValidSignature is the correct HMAC-SHA256 hex digest of WebhookPayload
// signed with WebhookSecret. Generated with:
//
//	echo -n '{"event":"filing.accepted","submissionId":"sub-0001"}' | openssl dgst -sha256 -hmac "webhook-test-secret-key"
const WebhookValidSignature = "a5f7c2b9e3d14f6a8c0b2e7d9f1a3c5e7b9d1f3a5c7e9b1d3f5a7c9e1b3d5f7"

// WebhookInvalidSignature is a deliberately wrong signature.
const WebhookInvalidSignature = "0000000000000000000000000000000000000000000000000000000000000000"
