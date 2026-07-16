// Package fixtures provides shared test data for the Tax990 SDK test suite.
package fixtures

// Auth fixture constants.
const (
	ClientID     = "test-client-id"
	ClientSecret = "test-client-secret"
	UserToken    = "test-user-token"
	AccessToken  = "eyJhbGciOiJSUzI1NiJ9.test.mock"
	JWSToken     = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mock.signature"
)

// TokenResponseBody is a valid JSON response from GET /Auth/GetTax990Token.
const TokenResponseBody = `{
	"statusCode": 200,
	"status": "OK",
	"message": "Successful API call.",
	"response": {
		"AccessToken": "eyJhbGciOiJSUzI1NiJ9.test.mock",
		"TokenType": "Bearer",
		"ExpiresIn": 3600,
		"Errors": null
	}
}`

// TokenResponseBodyWithError is a token response that contains an error payload.
const TokenResponseBodyWithError = `{
	"statusCode": 401,
	"status": "Unauthorized",
	"message": "Invalid credentials.",
	"response": {
		"AccessToken": "",
		"TokenType": null,
		"ExpiresIn": 0,
		"Errors": {
			"ErrorCode": "AUTH001",
			"ErrorName": "Unauthorized",
			"ErrorMessage": "Invalid client credentials."
		}
	}
}`

// GenerateJWSResponseBody is a valid JSON response from POST /Auth/GenerateJWS.
const GenerateJWSResponseBody = `{
	"statusCode": 200,
	"status": "Success",
	"message": "JWS Generated Successfully",
	"response": {
		"JWSToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mock.signature"
	}
}`
