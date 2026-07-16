// Package resources implements the API resource clients for the Tax990 SDK.
package resources

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"

	httpclient "github.com/tax990/sdk-go/tax990/http"
	"github.com/tax990/sdk-go/tax990/models"
)

// Form990NResource is the interface for all Form 990-N filing operations.
type Form990NResource interface {
	// Create submits new Form 990-N records. Supports idempotency via idempotencyKey
	// (pass empty string to auto-generate one).
	Create(ctx context.Context, payload *models.CreatePayload, idempotencyKey string) (*models.ApiResponse[models.SuccessRecord, models.ErrorRecord], error)

	// Submit is an alias for Create.
	Submit(ctx context.Context, payload *models.CreatePayload, idempotencyKey string) (*models.ApiResponse[models.SuccessRecord, models.ErrorRecord], error)

	// Update modifies existing Form 990-N records in a submission.
	Update(ctx context.Context, payload *models.UpdatePayload) (*models.ApiResponse[models.SuccessRecord, models.ErrorRecord], error)

	// Get retrieves a submission or a specific record within it.
	// Pass an empty recordID to retrieve all records in the submission.
	Get(ctx context.Context, submissionID, recordID string) (*models.ApiResponse[models.GetSuccessRecord, models.ErrorRecord], error)

	// List returns records matching the submission or business filter.
	// At least one of submissionID or businessID must be non-empty.
	List(ctx context.Context, submissionID, businessID string) (*models.ApiResponse[models.GetSuccessRecord, models.ErrorRecord], error)

	// Delete removes a submission or a specific record within it.
	// Pass an empty recordID to delete all records in the submission.
	Delete(ctx context.Context, submissionID, recordID string) (*models.ApiResponse[models.SuccessRecord, models.ErrorRecord], error)

	// Validate checks records for errors and warnings without transmitting.
	Validate(ctx context.Context, submissionID string, recordIDs []string) (*models.ApiResponse[models.ValidateSuccessRecord, models.ValidateErrorRecord], error)

	// Transmit sends validated records to the IRS.
	Transmit(ctx context.Context, payload *models.TransmitPayload) (*models.ApiResponse[models.TransmitSuccessRecord, models.TransmitErrorRecord], error)

	// GetPDF retrieves PDF download URLs for transmitted records.
	GetPDF(ctx context.Context, submissionID string, recordIDs []string) (*models.PDFResponse, error)

	// Status returns the current filing status for a submission or specific records.
	Status(ctx context.Context, submissionID string, recordIDs []string) (*models.ApiResponse[models.SuccessRecord, models.ErrorRecord], error)
}

type form990NResource struct {
	http *httpclient.Client
}

// NewForm990N creates a Form990NResource backed by the given HTTP client.
func NewForm990N(http *httpclient.Client) Form990NResource {
	return &form990NResource{http: http}
}

func (r *form990NResource) Create(ctx context.Context, payload *models.CreatePayload, idempotencyKey string) (*models.ApiResponse[models.SuccessRecord, models.ErrorRecord], error) {
	key := idempotencyKey
	if key == "" {
		key = genUUID()
	}
	resp, err := httpclient.Post[models.ApiResponse[models.SuccessRecord, models.ErrorRecord]](
		ctx, r.http, "/v1/form990n/create", payload,
		map[string]string{"idempotency-key": key},
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *form990NResource) Submit(ctx context.Context, payload *models.CreatePayload, idempotencyKey string) (*models.ApiResponse[models.SuccessRecord, models.ErrorRecord], error) {
	return r.Create(ctx, payload, idempotencyKey)
}

func (r *form990NResource) Update(ctx context.Context, payload *models.UpdatePayload) (*models.ApiResponse[models.SuccessRecord, models.ErrorRecord], error) {
	resp, err := httpclient.Post[models.ApiResponse[models.SuccessRecord, models.ErrorRecord]](
		ctx, r.http, "/v1/form990n/update", payload, nil,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *form990NResource) Get(ctx context.Context, submissionID, recordID string) (*models.ApiResponse[models.GetSuccessRecord, models.ErrorRecord], error) {
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

func (r *form990NResource) List(ctx context.Context, submissionID, businessID string) (*models.ApiResponse[models.GetSuccessRecord, models.ErrorRecord], error) {
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

func (r *form990NResource) Delete(ctx context.Context, submissionID, recordID string) (*models.ApiResponse[models.SuccessRecord, models.ErrorRecord], error) {
	resp, err := httpclient.Delete[models.ApiResponse[models.SuccessRecord, models.ErrorRecord]](
		ctx, r.http, "/v1/form990n/delete",
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

func (r *form990NResource) Validate(ctx context.Context, submissionID string, recordIDs []string) (*models.ApiResponse[models.ValidateSuccessRecord, models.ValidateErrorRecord], error) {
	resp, err := httpclient.Get[models.ApiResponse[models.ValidateSuccessRecord, models.ValidateErrorRecord]](
		ctx, r.http, "/v1/form990n/validate",
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

func (r *form990NResource) Transmit(ctx context.Context, payload *models.TransmitPayload) (*models.ApiResponse[models.TransmitSuccessRecord, models.TransmitErrorRecord], error) {
	resp, err := httpclient.Post[models.ApiResponse[models.TransmitSuccessRecord, models.TransmitErrorRecord]](
		ctx, r.http, "/v1/form990n/transmit", payload, nil,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *form990NResource) GetPDF(ctx context.Context, submissionID string, recordIDs []string) (*models.PDFResponse, error) {
	resp, err := httpclient.Get[models.PDFResponse](
		ctx, r.http, "/v1/form990n/getPDF",
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

func (r *form990NResource) Status(ctx context.Context, submissionID string, recordIDs []string) (*models.ApiResponse[models.SuccessRecord, models.ErrorRecord], error) {
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

// cleanParams removes empty-string values so they are not sent as query parameters.
func cleanParams(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if v != "" {
			out[k] = v
		}
	}
	return out
}

// genUUID generates a random UUID v4 using crypto/rand.
func genUUID() string {
	var b [16]byte
	rand.Read(b[:]) //nolint:errcheck
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
