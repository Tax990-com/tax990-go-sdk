package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tax990/sdk-go/tax990"
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

	ping, err := client.Utility.Ping(context.Background())
	if err != nil {
		log.Fatalf("ping failed: %v", err)
	}

	out, _ := json.MarshalIndent(ping, "", "  ")
	fmt.Println(string(out))

	submissions, err := client.Utility.GetAllSubmissionId(context.Background())
	if err != nil {
		log.Fatalf("get submissions failed: %v", err)
	}

	out2, _ := json.MarshalIndent(submissions, "", "  ")
	fmt.Println(string(out2))
}
