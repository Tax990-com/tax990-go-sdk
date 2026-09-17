package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tax990/sdk-go/tax990"
	"github.com/tax990/sdk-go/tax990/models"
)

func main() {
	client, err := tax990.NewClient(tax990.Config{
		ClientID:     os.Getenv("TAX990_CLIENT_ID"),
		ClientSecret: os.Getenv("TAX990_CLIENT_SECRET"),
		UserToken:    os.Getenv("TAX990_USER_TOKEN"),
		Environment:  tax990.Sandbox,
	})
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	boolTrue := true
	boolFalse := false
	payload := &models.CreatePayload{
		Form990NRecords: []models.Form990NRecord{
			{
				Business: models.Business{
					BusinessNm:   strPtr("Example Nonprofit"),
					EIN:          strPtr("123456789"),
					EmailAddress: strPtr("contact@example.org"),
					Phone:        strPtr("5551234567"),
					IsForeign:    &boolFalse,
					USAddress: &models.USAddress{
						Address1: strPtr("123 Main St"),
						City:     strPtr("Springfield"),
						State:    strPtr("IL"),
						ZipCd:    strPtr("62701"),
					},
				},
				Form990N: models.Form990NData{
					SequenceId:              strPtr("1"),
					TaxYr:                   strPtr("2024"),
					TaxPeriodBeginDt:        strPtr("2024-01-01"),
					TaxPeriodEndDt:          strPtr("2024-12-31"),
					IsGrossReceiptsUnder50K: &boolTrue,
					IsOrganizationTerminated: &boolFalse,
					WebsiteAddress:          strPtr("https://example.org"),
					PrincipalOfficer: &models.PrincipalOfficer{
						OfficerNm: strPtr("Jane Doe"),
						IsForeign: &boolFalse,
						USAddress: &models.USAddress{
							Address1: strPtr("456 Oak Ave"),
							City:     strPtr("Springfield"),
							State:    strPtr("IL"),
							ZipCd:    strPtr("62701"),
						},
					},
				},
			},
		},
	}

	result, err := client.Form990N.Create(context.Background(), payload, "")
	if err != nil {
		log.Fatalf("create failed: %v", err)
	}

	out, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(out))
}

func strPtr(s string) *string { return &s }
