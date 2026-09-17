package resources

import (
	"context"

	httpclient "github.com/tax990/sdk-go/tax990/http"
	"github.com/tax990/sdk-go/tax990/models"
)

// NonprofitsResource is the interface for nonprofit organization lookups.
type NonprofitsResource interface {
	// GetOrganizationDetailsByEIN retrieves nonprofit organization details by EIN.
	GetOrganizationDetailsByEIN(ctx context.Context, ein string) (*models.NonprofitResponse, error)
}

type nonprofitsResource struct {
	http *httpclient.Client
}

// NewNonprofits creates a NonprofitsResource backed by the given HTTP client.
func NewNonprofits(http *httpclient.Client) NonprofitsResource {
	return &nonprofitsResource{http: http}
}

func (r *nonprofitsResource) GetOrganizationDetailsByEIN(ctx context.Context, ein string) (*models.NonprofitResponse, error) {
	resp, err := httpclient.Get[models.NonprofitResponse](
		ctx, r.http, "/v1/nonprofits/getOrganizationDetailsByEIN",
		cleanParams(map[string]string{"ein": ein}),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
