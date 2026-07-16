// check_filing_status demonstrates how to retrieve the IRS filing status for a submission.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tax990/sdk-go/tax990"
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

	submissionID := os.Getenv("TAX990_SUBMISSION_ID")
	if submissionID == "" {
		log.Fatal("TAX990_SUBMISSION_ID environment variable is required")
	}

	ctx := context.Background()
	resp, err := client.FilingStatus.Get(ctx, submissionID, nil)
	if err != nil {
		log.Fatalf("get filing status: %v", err)
	}

	out, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(out))

	if resp.Form990NRecords != nil {
		for _, rec := range resp.Form990NRecords.SuccessRecords {
			fmt.Printf("Record %s — Status: %s\n", rec.RecordId, rec.RecordStatus)
		}
	}
}
