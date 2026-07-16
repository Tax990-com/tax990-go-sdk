// handle_webhook demonstrates verifying an incoming Tax990 webhook request.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tax990/sdk-go/tax990/models"
	"github.com/tax990/sdk-go/tax990/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	secret := os.Getenv("TAX990_WEBHOOK_SECRET")
	if secret == "" {
		log.Fatal("TAX990_WEBHOOK_SECRET environment variable is required")
	}

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}

		signature := r.Header.Get("x-tax990-signature")
		if signature == "" {
			http.Error(w, "missing signature header", http.StatusUnauthorized)
			return
		}

		if !utils.VerifyWebhookSignature(models.WebhookVerifyOptions{
			Secret:    secret,
			Payload:   string(body),
			Signature: signature,
		}) {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}

		var event map[string]any
		if err := json.Unmarshal(body, &event); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		fmt.Printf("Received verified webhook event: %v\n", event)
		w.WriteHeader(http.StatusOK)
	})

	addr := ":8080"
	fmt.Printf("Listening for webhooks on %s/webhook\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
