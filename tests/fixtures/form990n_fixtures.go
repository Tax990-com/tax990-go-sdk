package fixtures

import "github.com/tax990/sdk-go/tax990/models"

// SubmissionID is the fixture submission identifier used in tests.
const SubmissionID = "sub-00000000-0000-0000-0000-000000000001"

// RecordID is the fixture record identifier used in tests.
const RecordID = "rec-00000000-0000-0000-0000-000000000001"

// BusinessID is the fixture business identifier used in tests.
const BusinessID = "biz-00000000-0000-0000-0000-000000000001"

// SampleCreatePayload returns a minimal valid CreatePayload for tests.
func SampleCreatePayload() *models.CreatePayload {
	trueVal := true
	ein := "12-3456789"
	bizNm := "Test Nonprofit Org"
	seqID := "SEQ001"
	taxYr := "2023"
	begin := "2023-01-01"
	end := "2023-12-31"
	officerNm := "Jane Doe"
	addr1 := "123 Main St"
	city := "Austin"
	state := "TX"
	zip := "78701"

	return &models.CreatePayload{
		Form990NRecords: []models.Form990NRecord{
			{
				Business: models.Business{
					BusinessNm: &bizNm,
					EIN:        &ein,
					IsForeign:  new(bool),
					USAddress: &models.USAddress{
						Address1: &addr1,
						City:     &city,
						State:    &state,
						ZipCd:    &zip,
					},
				},
				Form990N: models.Form990NData{
					SequenceId:              &seqID,
					TaxYr:                   &taxYr,
					TaxPeriodBeginDt:        &begin,
					TaxPeriodEndDt:          &end,
					IsGrossReceiptsUnder50K: &trueVal,
					PrincipalOfficer: &models.PrincipalOfficer{
						OfficerNm: &officerNm,
						IsForeign: new(bool),
						USAddress: &models.USAddress{
							Address1: &addr1,
							City:     &city,
							State:    &state,
							ZipCd:    &zip,
						},
					},
				},
			},
		},
	}
}

// CreateSuccessResponseBody is a valid JSON response for POST /v1/form990n/create.
const CreateSuccessResponseBody = `{
	"StatusCode": 200,
	"StatusNm": "OK",
	"StatusMessage": "Records created successfully.",
	"CorrelationId": "corr-0001",
	"SubmissionId": "sub-00000000-0000-0000-0000-000000000001",
	"Form990NRecords": {
		"SuccessRecords": [
			{
				"SequenceId": "SEQ001",
				"RecordId": "rec-00000000-0000-0000-0000-000000000001",
				"BusinessId": "biz-00000000-0000-0000-0000-000000000001",
				"RecordStatus": "Created",
				"CreatedTs": "2024-01-01T00:00:00Z",
				"UpdatedTs": "2024-01-01T00:00:00Z"
			}
		],
		"ErrorRecords": null
	},
	"Errors": null
}`

// ValidationErrorResponseBody is a 400 error response body.
const ValidationErrorResponseBody = `{
	"StatusCode": 400,
	"StatusNm": "BadRequest",
	"StatusMessage": "A validation error has occurred.",
	"CorrelationId": "corr-0002",
	"SubmissionId": null,
	"Form990NRecords": null,
	"Errors": [
		{
			"Classification": "validation",
			"Code": "F990N001",
			"Message": "EIN is invalid.",
			"Field": "EIN"
		}
	]
}`

// SampleTransmitPayload is a minimal valid TransmitPayload for tests.
var SampleTransmitPayload = models.TransmitPayload{
	SubmissionId: SubmissionID,
	RecordIds:    []string{RecordID},
}

// NotFoundResponseBody is a 404 error response body.
const NotFoundResponseBody = `{
	"StatusCode": 404,
	"StatusNm": "NotFound",
	"StatusMessage": "Submission not found.",
	"CorrelationId": "corr-0003",
	"SubmissionId": null,
	"Form990NRecords": null,
	"Errors": null
}`
