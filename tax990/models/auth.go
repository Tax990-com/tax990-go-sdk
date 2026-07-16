package models

import "time"

// StoredToken holds a cached access token and its expiry time.
type StoredToken struct {
	AccessToken string
	ExpiresAt   time.Time
}

// TokenErrorDetail carries error fields from the OAuth token response body.
type TokenErrorDetail struct {
	ErrorCode    string `json:"ErrorCode"`
	ErrorName    string `json:"ErrorName"`
	ErrorMessage string `json:"ErrorMessage"`
}

// Tax990TokenData is the inner response payload from /Auth/GetTax990Token.
type Tax990TokenData struct {
	AccessToken string            `json:"AccessToken"`
	TokenType   string            `json:"TokenType"`
	ExpiresIn   int               `json:"ExpiresIn"`
	Errors      *TokenErrorDetail `json:"Errors"`
}

// Tax990TokenResponse is the full envelope from /Auth/GetTax990Token.
type Tax990TokenResponse struct {
	StatusCode int             `json:"statusCode"`
	Status     string          `json:"status"`
	Message    string          `json:"message"`
	Response   Tax990TokenData `json:"response"`
}

// GenerateJWSData is the inner payload from /Auth/GenerateJWS.
type GenerateJWSData struct {
	JWSToken string `json:"JWSToken"`
}

// GenerateJWSResponse is the full envelope from /Auth/GenerateJWS.
type GenerateJWSResponse struct {
	StatusCode int             `json:"statusCode"`
	Status     string          `json:"status"`
	Message    string          `json:"message"`
	Response   GenerateJWSData `json:"response"`
}

// GenerateJWSRequest is the request body for POST /Auth/GenerateJWS.
type GenerateJWSRequest struct {
	ClientId       string `json:"ClientId"`
	ClientSecretId string `json:"ClientSecretId"`
	UserToken      string `json:"UserToken"`
}
