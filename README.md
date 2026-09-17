# Tax990 Go SDK 2.0

## Overview

The Tax990 Go SDK 2.0.x is a Go integration package for the Tax990 Public API. It enables
businesses and software providers to integrate IRS Form 990-N e-filing directly into their
applications, without hand-rolling OAuth, request signing, or response parsing.

This SDK provides:

- **`Tax990Client`** — a single entry point exposing typed resources for Form 990-N filing,
  utility ID lookups, and nonprofit organization lookups
- **Automatic OAuth 2.0 token management** — signs and refreshes access tokens transparently
  between calls
- **Idiomatic Go** — `context.Context` on every call, typed errors, stdlib `net/http` underneath
- **No third-party runtime dependencies** beyond `godotenv` for local example configuration

A separate React UI ([`../frontend`](../frontend)) is included in this repository as a shared
playground for exercising any of the four language SDKs side by side. See
[`../UI_INTEGRATION.md`](../UI_INTEGRATION.md) for a full `net/http` bridge server built on this
exact SDK, and [`../TESTING.md`](../TESTING.md) for the test walkthrough.

🔗 Full API Reference: [developer.tax990.com](https://developer.tax990.com)

## Project Structure

```
go/
├── tax990/
│   ├── auth/                # OAuthClient, TokenManager
│   ├── errors/                # apierrors package — Tax990Error and subtypes
│   ├── http/                   # HTTP client (bearer injection, retries)
│   ├── models/                   # Request/response structs
│   ├── resources/                 # Form990NResource, UtilityResource, NonprofitsResource, ...
│   ├── utils/                       # EIN validation, webhook signature verification
│   ├── client.go                     # Tax990Client, NewClient — package entry point
│   └── config.go                      # Config, Environment
├── examples/
│   ├── submit_990n/
│   ├── check_filing_status/
│   ├── check_status/
│   ├── list_organizations/
│   ├── lookup_nonprofit/
│   ├── handle_webhook/
│   └── utility_ping/
├── tests/
│   └── fixtures/
└── go.mod
```

## Installation

```bash
go get github.com/tax990/sdk-go
```

Or, to build from source inside this repository:

```bash
cd go
go build ./...
```

## Quick Start

```go
package main

import (
	"context"
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
		// API and OAuth URLs are read from TAX990_API_URL / TAX990_OAUTH_URL env vars
	})
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ping, err := client.Utility.Ping(context.Background())
	if err != nil {
		log.Fatalf("ping failed: %v", err)
	}
	fmt.Printf("%+v\n", ping)
}
```

## Available API Modules

### Authentication

Handled automatically. `Tax990Client` signs a JWS locally with `ClientSecret`, exchanges it for an
access token against the OAuth endpoint, and caches/refreshes it before expiry — no manual token
handling required.

### Form 990-N (`client.Form990N`)

Create, validate, and transmit IRS Form 990-N e-Postcard filings.

| Method | Endpoint | Description |
|---|---|---|
| `Create(ctx, payload, idempotencyKey)` | `POST /v1/form990n/create` | Create and save a new Form 990-N filing |
| `Submit(ctx, payload, idempotencyKey)` | `POST /v1/form990n/create` | Alias for `Create` |
| `Update(ctx, payload)` | `POST /v1/form990n/update` | Update an existing filing |
| `Get(ctx, submissionID, recordID)` | `GET /v1/form990n/get` | Retrieve a saved filing by SubmissionId |
| `List(ctx, submissionID, businessID)` | `GET /v1/form990n/list` | Paginated list of filings |
| `Delete(ctx, submissionID, recordID)` | `DELETE /v1/form990n/delete` | Delete an untransmitted filing |
| `Validate(ctx, submissionID, recordIDs)` | `GET /v1/form990n/validate` | Validate records before transmit |
| `Transmit(ctx, payload)` | `POST /v1/form990n/transmit` | E-file to the IRS |
| `GetPDF(ctx, submissionID, recordIDs)` | `GET /v1/form990n/getPDF` | Download filing PDF copies |
| `Status(ctx, submissionID, recordIDs)` | `GET /v1/form990n/status` | Check IRS acknowledgement status |

**Key fields:** `TaxYr`, `TaxPeriodBeginDt`/`EndDt`, `IsGrossReceiptsUnder50K`,
`IsOrganizationTerminated`, `PrincipalOfficer`, `Business.USAddress`/`ForeignAddress`.

### Utility (`client.Utility`)

Health checks and cross-reference ID lookups.

| Method | Endpoint | Description |
|---|---|---|
| `Ping(ctx)` | `GET /v1/utility/ping` | Health check |
| `GetAllSubmissionId(ctx)` | `GET /v1/utility/getAllSubmissionId` | Get all submission IDs |
| `GetSubmissionIdByBusinessId(ctx, businessId)` | `GET /v1/utility/getSubmissionIdByBusinessId` | Look up submission by business ID |
| `GetSubmissionIdByRecordId(ctx, recordId)` | `GET /v1/utility/getSubmissionIdByRecordId` | Look up submission by record ID |
| `GetRecordIds(ctx)` | `GET /v1/utility/getRecordIds` | Get all record IDs |
| `GetRecordIdBySubmissionId(ctx, submissionId)` | `GET /v1/utility/getRecordIdBySubmissionId` | Get records for a submission |
| `GetRecordDetailBySubmissionId(ctx, submissionId)` | `GET /v1/utility/getRecordDetailBySubmissionId` | Get record details for a submission |
| `GetAllBusinessId(ctx)` | `GET /v1/utility/getAllBusinessId` | Get all business IDs |
| `GetBusinessIdBySubmissionId(ctx, submissionId)` | `GET /v1/utility/getBusinessIdBySubmissionId` | Get business ID for a submission |

### Nonprofits (`client.Nonprofits`)

| Method | Endpoint | Description |
|---|---|---|
| `GetOrganizationDetailsByEIN(ctx, ein)` | `GET /v1/nonprofits/getOrganizationDetailsByEIN` | Look up nonprofit organization details by EIN |

Also available: `client.Organizations` (business-entity queries over the same Form 990-N data) and
`client.FilingStatus` (a status-only convenience wrapper). `client.Webhooks` and `client.ApiKeys`
are stubs — those endpoints are not yet live on the Public API.

## Environment Variables

Set these in `go/.env` (loaded automatically by the bridge via `godotenv`) or export them in your shell.

| Variable | Required | Description |
|---|---|---|
| `TAX990_CLIENT_ID` | ✅ | OAuth client identifier |
| `TAX990_CLIENT_SECRET` | ✅ | OAuth client secret, used to sign the JWS |
| `TAX990_USER_TOKEN` | ✅ | OAuth audience token for this client |
| `TAX990_API_URL` | ✅ | Public API base URL (e.g. `https://api.tax990.com`) |
| `TAX990_OAUTH_URL` | ✅ | OAuth API base URL (e.g. `https://oauth.tax990.com`) |

`Config.APIUrl` / `OAuthUrl` can be set directly to override the env vars at the call site.

## Typical Workflow

1. **Instantiate** → `tax990.NewClient(config)`
2. **Create a filing** → `client.Form990N.Create(ctx, payload, "")` → store the returned `SubmissionId`
3. **Validate (optional)** → `client.Form990N.Validate(ctx, ...)` to catch errors before transmit
4. **Review a draft** → `client.Form990N.GetPDF(ctx, ...)` for a pre-transmission preview
5. **Transmit** → `client.Form990N.Transmit(ctx, ...)` to e-file with the IRS
6. **Track status** → `client.Form990N.Status(ctx, ...)` for acknowledgement status
7. **Look up organizations** → `client.Nonprofits.GetOrganizationDetailsByEIN(ctx, ein)` as needed

## Error Handling

All API errors are typed in the `apierrors` package (`github.com/tax990/sdk-go/tax990/errors`):

```go
import apierrors "github.com/tax990/sdk-go/tax990/errors"

result, err := client.Form990N.Create(ctx, payload, "")
if err != nil {
	var validationErr *apierrors.ValidationError
	var authErr *apierrors.AuthError
	switch {
	case errors.As(err, &validationErr):
		for _, e := range validationErr.Errors {
			fmt.Printf("[%s] %s: %s\n", e.Code, e.Field, e.Message)
		}
	case errors.As(err, &authErr):
		log.Printf("authentication failed: %v", authErr)
	default:
		log.Printf("API error: %v", err)
	}
}
```

## Testing

```bash
cd go
go test ./...
```

## Documentation

🔗 [Tax990 Public API Docs](https://developer.tax990.com)

## Tech Stack

| Layer | Technology |
|---|---|
| Runtime | Go 1.21+ |
| HTTP | stdlib `net/http` |
| Auth | OAuth 2.0 Bearer tokens, JWS (HS256) |
| Tests | stdlib `testing`, `testify` |

## License

MIT — internal SDK for Tax990 Public API integration.
