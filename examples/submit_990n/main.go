// submit_990n demonstrates how to submit a Form 990-N filing using the Tax990 Go SDK.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tax990/sdk-go/tax990"
	"github.com/tax990/sdk-go/tax990/models"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	client, err := tax990.NewClient(tax990.Config{
		ClientID:     os.Getenv("TAX990_CLIENT_ID"),
		ClientSecret: os.Getenv("TAX990_CLIENT_SECRET"),
		UserToken:    os.Getenv("TAX990_USER_TOKEN"),
		Environment:  tax990.Sandbox,
	})
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	trueVal := true
	ein := "12-3456789"
	bizNm := "My Nonprofit Org"
	seqID := "SEQ001"
	taxYr := "2023"
	begin := "2023-01-01"
	end := "2023-12-31"
	officerNm := "Jane Smith"
	addr1 := "456 Oak Ave"
	city := "Austin"
	state := "TX"
	zip := "78702"

	payload := &models.CreatePayload{
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

	ctx := context.Background()
	resp, err := client.Form990N.Submit(ctx, payload, "")
	if err != nil {
		log.Fatalf("submit 990-N: %v", err)
	}

	out, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(out))

	if resp.SubmissionId != nil {
		fmt.Printf("\nSubmission ID: %s\n", *resp.SubmissionId)
	}
	if resp.Form990NRecords != nil {
		fmt.Printf("Success records: %d\n", len(resp.Form990NRecords.SuccessRecords))
		fmt.Printf("Error records:   %d\n", len(resp.Form990NRecords.ErrorRecords))
	}
}
