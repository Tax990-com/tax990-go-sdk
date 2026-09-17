package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/tax990/sdk-go/tax990"
	"github.com/tax990/sdk-go/tax990/models"
)

var client *tax990.Tax990Client

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusInternalServerError, map[string]any{
		"StatusCode":    500,
		"StatusMessage": err.Error(),
	})
}

func csv(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, idempotency-key")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			return
		}
		h(w, r)
	}
}

func main() {
	godotenv.Load()
	c, err := tax990.NewClient(tax990.Config{
		ClientID:     os.Getenv("TAX990_CLIENT_ID"),
		ClientSecret: os.Getenv("TAX990_CLIENT_SECRET"),
		UserToken:    os.Getenv("TAX990_USER_TOKEN"),
		APIUrl:       os.Getenv("TAX990_API_URL"),
		OAuthUrl:     os.Getenv("TAX990_OAUTH_URL"),
	})
	if err != nil {
		panic(err)
	}
	client = c

	mux := http.NewServeMux()

	// ─── Auth ────────────────────────────────────────────────────
	mux.HandleFunc("/api/auth/token", withCORS(func(w http.ResponseWriter, r *http.Request) {
		result, err := client.Utility.Ping(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/auth/server-time", withCORS(func(w http.ResponseWriter, r *http.Request) {
		result, err := client.Utility.Ping(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))

	// ─── Form 990-N ──────────────────────────────────────────────
	mux.HandleFunc("/api/form990n/create", withCORS(func(w http.ResponseWriter, r *http.Request) {
		var payload models.CreatePayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, err)
			return
		}
		result, err := client.Form990N.Create(r.Context(), &payload, r.Header.Get("idempotency-key"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/form990n/update", withCORS(func(w http.ResponseWriter, r *http.Request) {
		var payload models.UpdatePayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, err)
			return
		}
		result, err := client.Form990N.Update(r.Context(), &payload)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/form990n/get", withCORS(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		result, err := client.Form990N.Get(r.Context(), q.Get("SubmissionId"), q.Get("RecordId"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/form990n/list", withCORS(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		result, err := client.Form990N.List(r.Context(), q.Get("SubmissionId"), q.Get("BusinessId"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/form990n/delete", withCORS(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		result, err := client.Form990N.Delete(r.Context(), q.Get("SubmissionId"), q.Get("RecordId"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/form990n/validate", withCORS(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		result, err := client.Form990N.Validate(r.Context(), q.Get("SubmissionId"), csv(q.Get("RecordIds")))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/form990n/transmit", withCORS(func(w http.ResponseWriter, r *http.Request) {
		var payload models.TransmitPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, err)
			return
		}
		result, err := client.Form990N.Transmit(r.Context(), &payload)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/form990n/getPDF", withCORS(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		result, err := client.Form990N.GetPDF(r.Context(), q.Get("SubmissionId"), csv(q.Get("RecordIds")))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/form990n/status", withCORS(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		result, err := client.Form990N.Status(r.Context(), q.Get("SubmissionId"), csv(q.Get("RecordIds")))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))

	// ─── Utility ─────────────────────────────────────────────────
	mux.HandleFunc("/api/utility/ping", withCORS(func(w http.ResponseWriter, r *http.Request) {
		result, err := client.Utility.Ping(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/utility/getAllSubmissionId", withCORS(func(w http.ResponseWriter, r *http.Request) {
		result, err := client.Utility.GetAllSubmissionId(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/utility/getSubmissionIdByBusinessId", withCORS(func(w http.ResponseWriter, r *http.Request) {
		result, err := client.Utility.GetSubmissionIdByBusinessId(r.Context(), r.URL.Query().Get("businessId"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/utility/getSubmissionIdByRecordId", withCORS(func(w http.ResponseWriter, r *http.Request) {
		result, err := client.Utility.GetSubmissionIdByRecordId(r.Context(), r.URL.Query().Get("recordId"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/utility/getRecordIds", withCORS(func(w http.ResponseWriter, r *http.Request) {
		result, err := client.Utility.GetRecordIds(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/utility/getRecordIdBySubmissionId", withCORS(func(w http.ResponseWriter, r *http.Request) {
		result, err := client.Utility.GetRecordIdBySubmissionId(r.Context(), r.URL.Query().Get("submissionId"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/utility/getRecordDetailBySubmissionId", withCORS(func(w http.ResponseWriter, r *http.Request) {
		result, err := client.Utility.GetRecordDetailBySubmissionId(r.Context(), r.URL.Query().Get("submissionId"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/utility/getAllBusinessId", withCORS(func(w http.ResponseWriter, r *http.Request) {
		result, err := client.Utility.GetAllBusinessId(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))
	mux.HandleFunc("/api/utility/getBusinessIdBySubmissionId", withCORS(func(w http.ResponseWriter, r *http.Request) {
		result, err := client.Utility.GetBusinessIdBySubmissionId(r.Context(), r.URL.Query().Get("submissionId"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))

	// ─── Nonprofits ──────────────────────────────────────────────
	mux.HandleFunc("/api/nonprofits/getOrganizationDetailsByEIN", withCORS(func(w http.ResponseWriter, r *http.Request) {
		result, err := client.Nonprofits.GetOrganizationDetailsByEIN(r.Context(), r.URL.Query().Get("ein"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	}))

	port := envOr("PORT", "4100")
	http.ListenAndServe(":"+port, mux)
}
