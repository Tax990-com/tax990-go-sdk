// list_organizations demonstrates how to list organizations (business entities)
// associated with a submission or business ID.
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
	businessID := os.Getenv("TAX990_BUSINESS_ID")

	if submissionID == "" && businessID == "" {
		log.Fatal("at least one of TAX990_SUBMISSION_ID or TAX990_BUSINESS_ID is required")
	}

	ctx := context.Background()
	resp, err := client.Organizations.List(ctx, submissionID, businessID)
	if err != nil {
		log.Fatalf("list organizations: %v", err)
	}

	out, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(out))

	if resp.Form990NRecords != nil {
		fmt.Printf("\nFound %d organization(s):\n", len(resp.Form990NRecords.SuccessRecords))
		for _, rec := range resp.Form990NRecords.SuccessRecords {
			bizNm := ""
			if rec.Business.BusinessNm != nil {
				bizNm = *rec.Business.BusinessNm
			}
			ein := ""
			if rec.Business.EIN != nil {
				ein = *rec.Business.EIN
			}
			fmt.Printf("  - %s (EIN: %s) — Status: %s\n", bizNm, ein, rec.RecordStatus)
		}
	}
}
