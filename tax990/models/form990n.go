package models

// USAddress is a United States mailing address.
type USAddress struct {
	Address1 *string `json:"Address1"`
	Address2 *string `json:"Address2"`
	City     *string `json:"City"`
	State    *string `json:"State"`
	ZipCd    *string `json:"ZipCd"`
}

// ForeignAddress is a non-US mailing address.
type ForeignAddress struct {
	Address1          *string `json:"Address1"`
	Address2          *string `json:"Address2"`
	City              *string `json:"City"`
	ProvinceOrStateNm *string `json:"ProvinceOrStateNm"`
	Country           *string `json:"Country"`
	PostalCd          *string `json:"PostalCd"`
}

// Business holds the organization entity data for a Form 990-N filing.
type Business struct {
	BusinessId     *string         `json:"BusinessId"`
	BusinessNm     *string         `json:"BusinessNm"`
	EIN            *string         `json:"EIN"`
	DBANm          *string         `json:"DBANm"`
	InCareOfNm     *string         `json:"InCareOfNm"`
	EmailAddress   *string         `json:"EmailAddress"`
	Phone          *string         `json:"Phone"`
	IsForeign      *bool           `json:"IsForeign"`
	USAddress      *USAddress      `json:"USAddress"`
	ForeignAddress *ForeignAddress `json:"ForeignAddress"`
}

// PrincipalOfficer holds officer contact information for a Form 990-N filing.
type PrincipalOfficer struct {
	OfficerNm      *string         `json:"OfficerNm"`
	IsForeign      *bool           `json:"IsForeign"`
	USAddress      *USAddress      `json:"USAddress"`
	ForeignAddress *ForeignAddress `json:"ForeignAddress"`
}

// Form990NData contains the tax form fields for a single 990-N record.
type Form990NData struct {
	SequenceId               *string           `json:"SequenceId"`
	RecordId                 *string           `json:"RecordId"`
	TaxYr                    *string           `json:"TaxYr"`
	TaxPeriodBeginDt         *string           `json:"TaxPeriodBeginDt"`
	TaxPeriodEndDt           *string           `json:"TaxPeriodEndDt"`
	IsGrossReceiptsUnder50K  *bool             `json:"IsGrossReceiptsUnder50K"`
	IsOrganizationTerminated *bool             `json:"IsOrganizationTerminated"`
	WebsiteAddress           *string           `json:"WebsiteAddress"`
	PrincipalOfficer         *PrincipalOfficer `json:"PrincipalOfficer"`
}

// Form990NRecord is a single business + form record within a submission.
type Form990NRecord struct {
	Business Business     `json:"Business"`
	Form990N Form990NData `json:"Form990N"`
}

// CreatePayload is the request body for POST /v1/form990n/create.
type CreatePayload struct {
	Form990NRecords []Form990NRecord `json:"Form990NRecords"`
}

// UpdatePayload is the request body for POST /v1/form990n/update.
type UpdatePayload struct {
	SubmissionId    string           `json:"SubmissionId"`
	Form990NRecords []Form990NRecord `json:"Form990NRecords"`
}

// TransmitPayload is the request body for POST /v1/form990n/transmit.
type TransmitPayload struct {
	SubmissionId string   `json:"SubmissionId"`
	RecordIds    []string `json:"RecordIds,omitempty"`
}

// RecordStatus represents the filing status string returned by the API.
type RecordStatus = string

const (
	RecordStatusCreated     RecordStatus = "Created"
	RecordStatusUpdated     RecordStatus = "Updated"
	RecordStatusDeleted     RecordStatus = "Deleted"
	RecordStatusTransmitted RecordStatus = "Transmitted"
	RecordStatusAccepted    RecordStatus = "Accepted"
	RecordStatusRejected    RecordStatus = "Rejected"
	RecordStatusInProgress  RecordStatus = "In-Progress"
	RecordStatusFailed      RecordStatus = "Failed"
)

// SuccessRecord is returned in the SuccessRecords array for create/update/delete.
type SuccessRecord struct {
	SequenceId   string `json:"SequenceId"`
	RecordId     string `json:"RecordId"`
	BusinessId   string `json:"BusinessId"`
	RecordStatus string `json:"RecordStatus"`
	CreatedTs    string `json:"CreatedTs"`
	UpdatedTs    string `json:"UpdatedTs"`
}

// ErrorRecord is returned in the ErrorRecords array when a record fails.
type ErrorRecord struct {
	SequenceId   *string           `json:"SequenceId"`
	RecordId     *string           `json:"RecordId"`
	BusinessId   *string           `json:"BusinessId"`
	RecordStatus string            `json:"RecordStatus"`
	Errors       []StructuredError `json:"Errors"`
}

// GetForm990NData is the Form990N sub-object inside a GET response record.
type GetForm990NData struct {
	TaxYear                  *string           `json:"TaxYear"`
	TaxPeriodBeginDate       *string           `json:"TaxPeriodBeginDate"`
	TaxPeriodEndDate         *string           `json:"TaxPeriodEndDate"`
	IsGrossReceiptsUnder50K  *bool             `json:"IsGrossReceiptsUnder50K"`
	IsOrganizationTerminated *bool             `json:"IsOrganizationTerminated"`
	WebsiteAddress           *string           `json:"WebsiteAddress"`
	PrincipalOfficer         *PrincipalOfficer `json:"PrincipalOfficer"`
}

// GetSuccessRecord is a full record returned by GET /v1/form990n/get and list.
type GetSuccessRecord struct {
	SuccessRecord
	Business Business        `json:"Business"`
	Form990N GetForm990NData `json:"Form990N"`
}

// ValidationWarning is a single warning or error from the validate endpoint.
type ValidationWarning struct {
	ErrorCode string `json:"ErrorCode"`
	Name      string `json:"Name"`
	Message   string `json:"Message"`
}

// ValidateSuccessRecord is returned when a record passes validation.
type ValidateSuccessRecord struct {
	SequenceId string              `json:"SequenceId"`
	RecordId   string              `json:"RecordId"`
	Warnings   []ValidationWarning `json:"Warnings"`
}

// ValidateErrorRecord is returned when a record fails validation.
type ValidateErrorRecord struct {
	SequenceId string              `json:"SequenceId"`
	RecordId   string              `json:"RecordId"`
	Errors     []ValidationWarning `json:"Errors"`
}

// TransmitSuccessRecord is returned for a successfully transmitted record.
type TransmitSuccessRecord struct {
	SequenceId string  `json:"SequenceId"`
	RecordId   string  `json:"RecordId"`
	Status     string  `json:"Status"`
	StatusTs   *string `json:"StatusTs"`
}

// TransmitErrorRecord is returned for a record that failed transmission.
type TransmitErrorRecord struct {
	SequenceId   string `json:"SequenceId"`
	RecordId     string `json:"RecordId"`
	Status       string `json:"Status"`
	ErrorMessage string `json:"ErrorMessage"`
}

// PDFRecord holds the PDF URL for a single filing record.
type PDFRecord struct {
	RecordId string `json:"RecordId"`
	PDFUrl   string `json:"PDFUrl"`
}

// PDFErrorRecord is an error entry in a PDF response.
type PDFErrorRecord struct {
	RecordId string `json:"RecordId"`
	Message  string `json:"Message"`
}

// PDFResponse is the response envelope from GET /v1/form990n/getPDF.
type PDFResponse struct {
	StatusCode      int              `json:"StatusCode"`
	StatusName      string           `json:"StatusName"`
	StatusMessage   string           `json:"StatusMessage"`
	CorrelationId   string           `json:"CorrelationId"`
	SubmissionId    string           `json:"SubmissionId"`
	Form990NRecords []PDFRecord      `json:"Form990NRecords"`
	Errors          []PDFErrorRecord `json:"Errors"`
}

// FilingStatusId maps numeric DB values to display statuses.
type FilingStatusId int

const (
	FilingStatusInProgress  FilingStatusId = 1
	FilingStatusTransmitted FilingStatusId = 2
	FilingStatusAccepted    FilingStatusId = 3
	FilingStatusRejected    FilingStatusId = 4
)

// FormType identifies the IRS form variant.
type FormType int

const (
	FormType990   FormType = 1
	FormType990EZ FormType = 2
	FormType990N  FormType = 3
	FormType990PF FormType = 4
)
