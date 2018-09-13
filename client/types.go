package client

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// BodyNotJSONError indicates a response to an elements request has a body that
// could not be parsed as JSON.
type BodyNotJSONError struct {
	JSONDecodeError error
	StatusCode      int
	BodyBytes       []byte
}

// Implements the Error interface
func (e BodyNotJSONError) Error() string {
	return fmt.Sprintf(
		"Body is not valid JSON. Status: %v Body: %v Error: %s",
		e.StatusCode,
		e.BodyBytes,
		e.JSONDecodeError,
	)
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

// Options includes configuration options for a new base client
type Options struct {
	Host               string
	TLSConfig          *tls.Config
	Timeout            time.Duration
	DontFollowRedirect bool
}
