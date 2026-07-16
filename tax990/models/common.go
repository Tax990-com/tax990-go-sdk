// Package models defines request and response data structures for the Tax990 API.
package models

// StructuredError is an individual error detail in an API error response.
type StructuredError struct {
	Classification string  `json:"Classification"`
	Code           string  `json:"Code"`
	Message        string  `json:"Message"`
	Field          *string `json:"Field"`
}

// Form990NRecordsContainer holds success and error records from an API response.
type Form990NRecordsContainer[S any, E any] struct {
	SuccessRecords []S `json:"SuccessRecords"`
	ErrorRecords   []E `json:"ErrorRecords"`
}

// ApiResponse is the standard envelope returned by all Form990N endpoints.
type ApiResponse[S any, E any] struct {
	StatusCode      int                             `json:"StatusCode"`
	StatusNm        string                          `json:"StatusNm"`
	StatusMessage   string                          `json:"StatusMessage"`
	CorrelationId   string                          `json:"CorrelationId"`
	SubmissionId    *string                         `json:"SubmissionId"`
	Form990NRecords *Form990NRecordsContainer[S, E] `json:"Form990NRecords"`
	Errors          []StructuredError               `json:"Errors"`
}

// PaginatedResponse is a generic paginated result container.
// Per ANALYSIS.md, the current API does not paginate — this type is provided
// for forward compatibility only.
type PaginatedResponse[T any] struct {
	Items      []T
	TotalCount int
}

// ListOptions are query parameters for list endpoints.
type ListOptions struct {
	SubmissionID string
	BusinessID   string
}
