package models

// OrganizationListParams are query parameters for the organization list endpoint.
type OrganizationListParams struct {
	SubmissionID string
	BusinessID   string
}

// OrganizationGetParams are query parameters for the organization get endpoint.
type OrganizationGetParams struct {
	SubmissionID string
	RecordID     string
}
