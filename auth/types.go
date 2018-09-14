package auth

import (
	"net/http"
	"time"
)

// TokenWithExpiry represents a token that has an expiry time
type TokenWithExpiry struct {
	Token     string
	ExpiresIn float64
}

// ErrorBody is the corresponding structure of a platform error
type ErrorBody struct {
	ErrorType        string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorURI         string `json:"error_uri,omitempty"`
}

// Error conforms to the Error interface
func (e *ErrorBody) Error() string {
	return e.ErrorDescription
}

// TokenResponse represents information that is returned on
// generation of a token
type TokenResponse struct {
	AccessToken string  `json:"access_token"`
	TokenType   string  `json:"token_type"`
	ExpiresIn   float64 `json:"expires_in"`
}

// Response represents data that is returned
// when making a call to the Authenticate method
// It returns the status of the response, headers and the response body
type Response struct {
	Status  int
	Headers http.Header
	body    interface{}
}

// Error returns the ErrorBody of the authentication response
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

// TokenResponse returns the token returned by the response
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

// Options contains information to configure Authenticate method calls
type Options struct {
	UserID        *string
	ServiceClaims map[string]interface{}
	Su            bool
	TokenExpiry   *time.Duration
}

// Payload specifies the grant type for the token
type Payload struct {
	GrantType string
}
