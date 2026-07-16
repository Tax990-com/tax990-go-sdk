# Changelog

All notable changes to the Tax990 Go SDK are documented here.

## [0.1.0] — Initial release

### Added
- `Tax990Client` entry point with `NewClient(Config)` factory.
- `Config` with `Production` / `Sandbox` environment constants.
- Two-step OAuth token flow: local HS256 JWS signing → RS256 access token exchange.
- `TokenManager` with automatic pre-expiry refresh and concurrent-safe deduplication.
- `net/http` wrapper with exponential-backoff retry on 429 and 5xx responses.
- Auto-injected `x-correlation-id` header on every request.
- `Form990NResource` — Create, Submit, Update, Get, List, Delete, Validate, Transmit, GetPDF, Status.
- `OrganizationResource` — List, Get (wraps Form990N endpoints).
- `FilingStatusResource` — Get (wraps /status endpoint).
- `WebhookResource` and `ApiKeysResource` stubs (endpoints not yet in Public API).
- `VerifyWebhookSignature` utility — HMAC-SHA256 with timing-safe comparison.
- `ValidateEIN` and `FormatEIN` utilities.
- Typed errors: `Tax990Error`, `AuthError`, `ValidationError`, `RateLimitError`, `NotFoundError`.
- Generic `ApiResponse[S, E]` and `PaginatedResponse[T]` model types (Go 1.21+).
- Full test suite using `net/http/httptest` — no external test dependencies.
- Runnable examples: submit_990n, check_filing_status, handle_webhook, list_organizations.
