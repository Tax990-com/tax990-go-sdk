package resources

import (
	"context"

	httpclient "github.com/tax990/sdk-go/tax990/http"
	"github.com/tax990/sdk-go/tax990/models"
)

// UtilityResource is the interface for utility endpoints (ID lookups, health checks).
type UtilityResource interface {
	// Ping checks API connectivity. No authentication required.
	Ping(ctx context.Context) (*models.PingResponse, error)

	// GetAllSubmissionId returns all submission IDs for the authenticated user.
	GetAllSubmissionId(ctx context.Context) (*models.SubmissionIdResponse, error)

	// GetSubmissionIdByBusinessId returns submission IDs associated with a business ID.
	GetSubmissionIdByBusinessId(ctx context.Context, businessId string) (*models.SubmissionIdResponse, error)

	// GetSubmissionIdByRecordId returns the submission ID associated with a record ID.
	GetSubmissionIdByRecordId(ctx context.Context, recordId string) (*models.SubmissionIdResponse, error)

	// GetRecordIds returns all record IDs for the authenticated user.
	GetRecordIds(ctx context.Context) (*models.RecordIdResponse, error)

	// GetRecordIdBySubmissionId returns record IDs within a submission.
	GetRecordIdBySubmissionId(ctx context.Context, submissionId string) (*models.RecordIdResponse, error)

	// GetRecordDetailBySubmissionId returns detailed record information for a submission.
	GetRecordDetailBySubmissionId(ctx context.Context, submissionId string) (*models.RecordDetailResponse, error)

	// GetAllBusinessId returns all business IDs for the authenticated user.
	GetAllBusinessId(ctx context.Context) (*models.BusinessIdResponse, error)

	// GetBusinessIdBySubmissionId returns the business ID associated with a submission.
	GetBusinessIdBySubmissionId(ctx context.Context, submissionId string) (*models.BusinessIdResponse, error)
}

type utilityResource struct {
	http *httpclient.Client
}

// NewUtility creates a UtilityResource backed by the given HTTP client.
func NewUtility(http *httpclient.Client) UtilityResource {
	return &utilityResource{http: http}
}

func (r *utilityResource) Ping(ctx context.Context) (*models.PingResponse, error) {
	resp, err := httpclient.Get[models.PingResponse](
		ctx, r.http, "/v1/utility/ping", nil, nil,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *utilityResource) GetAllSubmissionId(ctx context.Context) (*models.SubmissionIdResponse, error) {
	resp, err := httpclient.Get[models.SubmissionIdResponse](
		ctx, r.http, "/v1/utility/getAllSubmissionId", nil, nil,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *utilityResource) GetSubmissionIdByBusinessId(ctx context.Context, businessId string) (*models.SubmissionIdResponse, error) {
	resp, err := httpclient.Get[models.SubmissionIdResponse](
		ctx, r.http, "/v1/utility/getSubmissionIdByBusinessId",
		cleanParams(map[string]string{"businessId": businessId}),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *utilityResource) GetSubmissionIdByRecordId(ctx context.Context, recordId string) (*models.SubmissionIdResponse, error) {
	resp, err := httpclient.Get[models.SubmissionIdResponse](
		ctx, r.http, "/v1/utility/getSubmissionIdByRecordId",
		cleanParams(map[string]string{"recordId": recordId}),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *utilityResource) GetRecordIds(ctx context.Context) (*models.RecordIdResponse, error) {
	resp, err := httpclient.Get[models.RecordIdResponse](
		ctx, r.http, "/v1/utility/getRecordIds", nil, nil,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *utilityResource) GetRecordIdBySubmissionId(ctx context.Context, submissionId string) (*models.RecordIdResponse, error) {
	resp, err := httpclient.Get[models.RecordIdResponse](
		ctx, r.http, "/v1/utility/getRecordIdBySubmissionId",
		cleanParams(map[string]string{"submissionId": submissionId}),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *utilityResource) GetRecordDetailBySubmissionId(ctx context.Context, submissionId string) (*models.RecordDetailResponse, error) {
	resp, err := httpclient.Get[models.RecordDetailResponse](
		ctx, r.http, "/v1/utility/getRecordDetailBySubmissionId",
		cleanParams(map[string]string{"submissionId": submissionId}),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *utilityResource) GetAllBusinessId(ctx context.Context) (*models.BusinessIdResponse, error) {
	resp, err := httpclient.Get[models.BusinessIdResponse](
		ctx, r.http, "/v1/utility/getAllBusinessId", nil, nil,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *utilityResource) GetBusinessIdBySubmissionId(ctx context.Context, submissionId string) (*models.BusinessIdResponse, error) {
	resp, err := httpclient.Get[models.BusinessIdResponse](
		ctx, r.http, "/v1/utility/getBusinessIdBySubmissionId",
		cleanParams(map[string]string{"submissionId": submissionId}),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
