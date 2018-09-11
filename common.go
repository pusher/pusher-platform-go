package pusherplatform

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
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

// ErrorResponse represents information that is returned in case of an error
type ErrorResponse struct {
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Info    interface{} `json:"info"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("Error response: %d, %v", e.Status, e.Info)
}

// RequestOptions is used to configure HTTP requests
type RequestOptions struct {
	Method      string
	Path        string
	Jwt         *string
	Headers     *http.Header
	Body        io.Reader
	QueryParams *url.Values
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
