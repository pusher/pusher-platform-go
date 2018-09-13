package authenticator

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

// AuthenticationResponse represents data that is returned
// when making a call to the Authenticate method
// It returns the status of the response, headers and the response body
type AuthenticationResponse struct {
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Body    interface{} `json:"body,omitempty"`
}

// AuthenticateOptions contains information to configure Authenticate method calls
type AuthenticateOptions struct {
	UserID        *string
	ServiceClaims map[string]interface{}
	Su            bool
	TokenExpiry   *time.Duration
}

// AuthenticatePayload specifies the grant type for the token
type AuthenticatePayload struct {
	GrantType string
}
