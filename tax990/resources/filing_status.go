package resources

import (
	"context"
	"strings"

	httpclient "github.com/tax990/sdk-go/tax990/http"
	"github.com/tax990/sdk-go/tax990/models"
)

// FilingStatusResource is the interface for querying IRS filing status.
type FilingStatusResource interface {
	// Get returns the current IRS filing status for a submission or specific records.
	// Pass an empty recordIDs slice to query all records in the submission.
	Get(ctx context.Context, submissionID string, recordIDs []string) (*models.ApiResponse[models.SuccessRecord, models.ErrorRecord], error)
}

type filingStatusResource struct {
	http *httpclient.Client
}

// NewFilingStatus creates a FilingStatusResource backed by the given HTTP client.
func NewFilingStatus(http *httpclient.Client) FilingStatusResource {
	return &filingStatusResource{http: http}
}

func (r *filingStatusResource) Get(ctx context.Context, submissionID string, recordIDs []string) (*models.ApiResponse[models.SuccessRecord, models.ErrorRecord], error) {
	resp, err := httpclient.Get[models.ApiResponse[models.SuccessRecord, models.ErrorRecord]](
		ctx, r.http, "/v1/form990n/status",
		cleanParams(map[string]string{
			"SubmissionId": submissionID,
			"RecordIds":    strings.Join(recordIDs, ","),
		}),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
