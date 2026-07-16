package resources

import (
	"context"

	httpclient "github.com/tax990/sdk-go/tax990/http"
	"github.com/tax990/sdk-go/tax990/models"
)

// OrganizationResource is the interface for organization (Business entity) queries.
// It wraps Form 990-N list and get endpoints, filtering by business data.
type OrganizationResource interface {
	// List returns organization records matching the given submission or business ID.
	// At least one of submissionID or businessID must be non-empty.
	List(ctx context.Context, submissionID, businessID string) (*models.ApiResponse[models.GetSuccessRecord, models.ErrorRecord], error)

	// Get retrieves a single organization record by submission and optional record ID.
	Get(ctx context.Context, submissionID, recordID string) (*models.ApiResponse[models.GetSuccessRecord, models.ErrorRecord], error)
}

type organizationResource struct {
	http *httpclient.Client
}

// NewOrganization creates an OrganizationResource backed by the given HTTP client.
func NewOrganization(http *httpclient.Client) OrganizationResource {
	return &organizationResource{http: http}
}

func (r *organizationResource) List(ctx context.Context, submissionID, businessID string) (*models.ApiResponse[models.GetSuccessRecord, models.ErrorRecord], error) {
	resp, err := httpclient.Get[models.ApiResponse[models.GetSuccessRecord, models.ErrorRecord]](
		ctx, r.http, "/v1/form990n/list",
		cleanParams(map[string]string{
			"SubmissionId": submissionID,
			"BusinessId":   businessID,
		}),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *organizationResource) Get(ctx context.Context, submissionID, recordID string) (*models.ApiResponse[models.GetSuccessRecord, models.ErrorRecord], error) {
	resp, err := httpclient.Get[models.ApiResponse[models.GetSuccessRecord, models.ErrorRecord]](
		ctx, r.http, "/v1/form990n/get",
		cleanParams(map[string]string{
			"SubmissionId": submissionID,
			"RecordId":     recordID,
		}),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
