package resources

import (
	"context"
	"errors"
)

// WebhookResource is the interface for webhook management.
// Webhook endpoints are not yet available in the Tax990 Public API.
type WebhookResource interface {
	// Register is not yet available in the Tax990 Public API.
	Register(ctx context.Context) error
	// Delete is not yet available in the Tax990 Public API.
	Delete(ctx context.Context, id string) error
	// List is not yet available in the Tax990 Public API.
	List(ctx context.Context) error
}

type webhookResource struct{}

// NewWebhook creates a WebhookResource stub.
func NewWebhook() WebhookResource {
	return &webhookResource{}
}

var errWebhookNotAvailable = errors.New("webhook endpoints are not yet available in the Tax990 Public API")

func (r *webhookResource) Register(_ context.Context) error        { return errWebhookNotAvailable }
func (r *webhookResource) Delete(_ context.Context, _ string) error { return errWebhookNotAvailable }
func (r *webhookResource) List(_ context.Context) error            { return errWebhookNotAvailable }
