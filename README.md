# Tax990 Go SDK

Official Go SDK for the Tax990 Public API. Mirrors the same resource names and
method signatures as the Node.js and Python SDKs, adapted to Go conventions.

**Requires Go 1.21+.**

## Installation

```bash
go get github.com/tax990/sdk-go
```

## Quick start

```go
import (
    "context"
    "os"

    "github.com/tax990/sdk-go/tax990"
    "github.com/tax990/sdk-go/tax990/models"
)

client, err := tax990.NewClient(tax990.Config{
    ClientID:     os.Getenv("TAX990_CLIENT_ID"),
    ClientSecret: os.Getenv("TAX990_CLIENT_SECRET"),
    UserToken:    os.Getenv("TAX990_USER_TOKEN"),
    Environment:  tax990.Production, // or tax990.Sandbox
})
if err != nil {
    log.Fatal(err)
}

ctx := context.Background()

// Submit a Form 990-N filing
resp, err := client.Form990N.Submit(ctx, &models.CreatePayload{
    Form990NRecords: []models.Form990NRecord{ /* ... */ },
}, "")

// Check filing status
status, err := client.FilingStatus.Get(ctx, submissionID, nil)

// List organizations
orgs, err := client.Organizations.List(ctx, "", businessID)
```

## Configuration

| Field | Env var | Required | Default |
|---|---|---|---|
| `ClientID` | `TAX990_CLIENT_ID` | Yes | — |
| `ClientSecret` | `TAX990_CLIENT_SECRET` | Yes | — |
| `UserToken` | `TAX990_USER_TOKEN` | Yes | — |
| `Environment` | — | No | `Sandbox` |
| `APIUrl` | — | No | env default |
| `OAuthUrl` | — | No | env default |
| `Timeout` | — | No | 30s |

## Environments

| Constant | API URL | OAuth URL |
|---|---|---|
| `tax990.Production` | `https://api.tax990.com` | `https://oauth.tax990.com` |
| `tax990.Sandbox` | `http://localhost:9005` | `http://localhost:4000` |

## Authentication

The SDK implements the Tax990 two-step token flow automatically:

1. Signs a JWS locally using HS256 (HMAC-SHA256) with `ClientSecret`.
2. Exchanges the JWS for an RS256 access token via `GET /Auth/GetTax990Token`.
3. Caches the token and refreshes it 30 seconds before expiry.

No manual token management is required.

## Resources

### Form990N

```go
// Create / Submit (alias)
resp, err := client.Form990N.Create(ctx, payload, idempotencyKey)
resp, err := client.Form990N.Submit(ctx, payload, "")   // auto-generates idempotency key

// Update
resp, err := client.Form990N.Update(ctx, updatePayload)

// Get a submission or specific record
resp, err := client.Form990N.Get(ctx, submissionID, recordID)

// List by submission or business
resp, err := client.Form990N.List(ctx, submissionID, businessID)

// Delete
resp, err := client.Form990N.Delete(ctx, submissionID, recordID)

// Validate
resp, err := client.Form990N.Validate(ctx, submissionID, []string{recordID})

// Transmit to IRS
resp, err := client.Form990N.Transmit(ctx, &models.TransmitPayload{
    SubmissionId: submissionID,
    RecordIds:    []string{recordID},
})

// Get PDF URLs
resp, err := client.Form990N.GetPDF(ctx, submissionID, []string{recordID})

// Check status
resp, err := client.Form990N.Status(ctx, submissionID, nil)
```

### Organizations

```go
orgs, err := client.Organizations.List(ctx, submissionID, businessID)
org, err  := client.Organizations.Get(ctx, submissionID, recordID)
```

### FilingStatus

```go
status, err := client.FilingStatus.Get(ctx, submissionID, nil)
```

## Error handling

```go
import (
    "errors"
    apierrors "github.com/tax990/sdk-go/tax990/errors"
)

resp, err := client.Form990N.Get(ctx, id, "")
if err != nil {
    var authErr *apierrors.AuthError
    var notFound *apierrors.NotFoundError
    var rateLimit *apierrors.RateLimitError
    var valErr  *apierrors.ValidationError

    switch {
    case errors.As(err, &authErr):
        // re-authenticate
    case errors.As(err, &notFound):
        // resource missing
    case errors.As(err, &rateLimit):
        // back off and retry
    case errors.As(err, &valErr):
        for _, e := range valErr.Errors {
            fmt.Println(e.Code, e.Message)
        }
    default:
        log.Fatal(err)
    }
}
```

## Webhook verification

```go
import (
    "github.com/tax990/sdk-go/tax990/models"
    "github.com/tax990/sdk-go/tax990/utils"
)

ok := utils.VerifyWebhookSignature(models.WebhookVerifyOptions{
    Secret:    os.Getenv("TAX990_WEBHOOK_SECRET"),
    Payload:   string(requestBody),
    Signature: r.Header.Get("x-tax990-signature"),
})
```

## EIN utilities

```go
import "github.com/tax990/sdk-go/tax990/utils"

utils.ValidateEIN("12-3456789")  // true
utils.ValidateEIN("123456789")   // true
utils.FormatEIN("123456789")     // "12-3456789"
```

## Running examples

```bash
cp .env.example .env
# fill in credentials in .env

go run ./examples/submit_990n/
go run ./examples/check_filing_status/
go run ./examples/handle_webhook/
go run ./examples/list_organizations/
```

## Running tests

```bash
go test ./...
```

## License

See the repository root for license information.
