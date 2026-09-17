package models

// NonprofitOrganization holds the details of a nonprofit organization returned by EIN lookup.
type NonprofitOrganization struct {
	EIN              string  `json:"EIN"`
	OrganizationName string  `json:"OrganizationName"`
	City             *string `json:"City"`
	State            *string `json:"State"`
	Country          *string `json:"Country"`
	TaxPeriod        *string `json:"TaxPeriod"`
	AssetAmount      *string `json:"AssetAmount"`
	IncomeAmount     *string `json:"IncomeAmount"`
	RevenueAmount    *string `json:"RevenueAmount"`
	NTEECode         *string `json:"NTEECode"`
	DeductibilityCode *string `json:"DeductibilityCode"`
}

// NonprofitResponse is the response envelope from GET /v1/nonprofits/getOrganizationDetailsByEIN.
type NonprofitResponse struct {
	StatusCode    int                    `json:"StatusCode"`
	StatusNm      string                 `json:"StatusNm"`
	StatusMessage string                 `json:"StatusMessage"`
	CorrelationId string                 `json:"CorrelationId"`
	Organization  *NonprofitOrganization `json:"Organization"`
}
