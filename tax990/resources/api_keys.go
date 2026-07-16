package resources

import (
	"context"
	"errors"
)

// ApiKeysResource is the interface for API key management.
// API key endpoints are not yet available in the Tax990 Public API.
type ApiKeysResource interface {
	// Create is not yet available in the Tax990 Public API.
	Create(ctx context.Context) error
	// Revoke is not yet available in the Tax990 Public API.
	Revoke(ctx context.Context, id string) error
	// List is not yet available in the Tax990 Public API.
	List(ctx context.Context) error
}

type apiKeysResource struct{}

// NewApiKeys creates an ApiKeysResource stub.
func NewApiKeys() ApiKeysResource {
	return &apiKeysResource{}
}

var errApiKeysNotAvailable = errors.New("API key endpoints are not yet available in the Tax990 Public API")

func (r *apiKeysResource) Create(_ context.Context) error        { return errApiKeysNotAvailable }
func (r *apiKeysResource) Revoke(_ context.Context, _ string) error { return errApiKeysNotAvailable }
func (r *apiKeysResource) List(_ context.Context) error          { return errApiKeysNotAvailable }
