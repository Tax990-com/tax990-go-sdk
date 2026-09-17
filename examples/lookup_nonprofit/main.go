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

	ein := "123456789"
	if len(os.Args) > 1 {
		ein = os.Args[1]
	}

	result, err := client.Nonprofits.GetOrganizationDetailsByEIN(context.Background(), ein)
	if err != nil {
		log.Fatalf("nonprofit lookup failed: %v", err)
	}

	out, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(out))
}
