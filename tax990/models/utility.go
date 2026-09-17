package models

// PingResponse is the response from GET /v1/utility/ping (no auth required).
type PingResponse struct {
	StatusCode    int    `json:"StatusCode"`
	StatusNm      string `json:"StatusNm"`
	StatusMessage string `json:"StatusMessage"`
	CorrelationId string `json:"CorrelationId"`
}

// SubmissionIdResponse is the response envelope for submission ID lookups.
type SubmissionIdResponse struct {
	StatusCode    int      `json:"StatusCode"`
	StatusNm      string   `json:"StatusNm"`
	StatusMessage string   `json:"StatusMessage"`
	CorrelationId string   `json:"CorrelationId"`
	SubmissionIds []string `json:"SubmissionIds"`
}

// RecordIdResponse is the response envelope for record ID lookups.
type RecordIdResponse struct {
	StatusCode    int      `json:"StatusCode"`
	StatusNm      string   `json:"StatusNm"`
	StatusMessage string   `json:"StatusMessage"`
	CorrelationId string   `json:"CorrelationId"`
	RecordIds     []string `json:"RecordIds"`
}

// BusinessIdResponse is the response envelope for business ID lookups.
type BusinessIdResponse struct {
	StatusCode    int      `json:"StatusCode"`
	StatusNm      string   `json:"StatusNm"`
	StatusMessage string   `json:"StatusMessage"`
	CorrelationId string   `json:"CorrelationId"`
	BusinessIds   []string `json:"BusinessIds"`
}

// RecordDetail holds detailed information about a record within a submission.
type RecordDetail struct {
	RecordId     string `json:"RecordId"`
	BusinessId   string `json:"BusinessId"`
	RecordStatus string `json:"RecordStatus"`
	CreatedTs    string `json:"CreatedTs"`
	UpdatedTs    string `json:"UpdatedTs"`
}

// RecordDetailResponse is the response envelope for record detail lookups.
type RecordDetailResponse struct {
	StatusCode    int            `json:"StatusCode"`
	StatusNm      string         `json:"StatusNm"`
	StatusMessage string         `json:"StatusMessage"`
	CorrelationId string         `json:"CorrelationId"`
	SubmissionId  string         `json:"SubmissionId"`
	RecordDetails []RecordDetail `json:"RecordDetails"`
}
