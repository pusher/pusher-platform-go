package auth

import (
	"net/http"
	"time"
)

const GrantTypeClientCredentials = "client_credentials"

// TokenWithExpiry represents a token that has an expiry time.
type TokenWithExpiry struct {
	Token     string  // Token string
	ExpiresIn float64 // Expiry in seconds
}

// ErrorBody is the corresponding structure of a platform error.
type ErrorBody struct {
	ErrorType        string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorURI         string `json:"error_uri,omitempty"`
}

// Error conforms to the Error interface
func (e *ErrorBody) Error() string {
	return e.ErrorDescription
}

// TokenResponse represents information that is returned on generation of a token.
type TokenResponse struct {
	AccessToken string  `json:"access_token"`
	TokenType   string  `json:"token_type"`
	ExpiresIn   float64 `json:"expires_in"`
}

// Response represents data that is returned when making a call to the Authenticate method.
//
// It returns the status of the response, headers and the response body
type Response struct {
	Status  int
	Headers http.Header
	body    interface{}
}

// Error returns the ErrorBody of the authentication response.
func (a *Response) Error() *ErrorBody {
	if a.Status != http.StatusOK {
		errorBody, ok := a.body.(*ErrorBody)
		if !ok {
			return nil
		}

		return errorBody
	}

	return nil
}

// TokenResponse returns the token returned by the response.
//
// It is important to check if the Response has an associated
// ErrorBody by calling Error() before accessing the TokenResponse.
func (a *Response) TokenResponse() *TokenResponse {
	if a.Status != http.StatusOK {
		return nil
	}

	tokenResponse, ok := a.body.(*TokenResponse)
	if !ok {
		return nil
	}

	return tokenResponse
}

// Options contains information to configure Authenticate method calls.
type Options struct {
	UserID        *string                // Optional user id
	ServiceClaims map[string]interface{} // Optional JWT service claims
	Su            bool                   // Indicates if token should contain the `su` claim
	TokenExpiry   *time.Duration         // Optional token expiry (defaults to 24 hours)
}

// Payload specifies the grant type for the token.
// Currently the only supported grant type is "client_credentials",
// passing anything else other than this will return an error
type Payload struct {
	GrantType string
}
