// Package client provides an HTTP/2 client that connects to the Pusher platform.
//
// For most use cases, construction of a client directly isn't required.
// The interface is very basic and only allows performing HTTP requests to the Pusher platform.
// For regular use cases, the use of Instance is encouraged.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const authorizationHeader = "Authorization"

// Client is a low level interface for clients of the elements protocol
type Client interface {
	Request(ctx context.Context, options RequestOptions) (*http.Response, error)
}

// New builds a new Client
func New(options Options) Client {
	return newClient(options)
}

// Request allows making HTTP calls
func (c *client) Request(ctx context.Context, options RequestOptions) (*http.Response, error) {
	request, err := buildRequest(ctx, c.schema, c.host, options)
	if err != nil {
		return nil, err
	}

	return sendRequest(c.underlyingClient, request, c.options.DontFollowRedirect)
}

// Implements the Client interface
type client struct {
	host             string
	schema           string
	underlyingClient http.Client
	options          Options
}

func newClient(options Options) *client {
	c := new(client)
	c.host = options.Host
	c.schema = "https"

	if options.TLSConfig != nil {
		transport := &http.Transport{
			Proxy:              http.ProxyFromEnvironment,
			TLSClientConfig:    options.TLSConfig,
			DisableCompression: true,
		}

		c.underlyingClient = http.Client{
			Transport: transport,
			Timeout:   options.Timeout,
		}
	}

	// Control how the http client handles redirect responses.
	// See: https://godoc.org/net/http#Client
	// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Redirections
	if options.DontFollowRedirect {
		c.underlyingClient.CheckRedirect = func(r *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	c.options = options

	return c
}

func buildRequest(
	ctx context.Context,
	schema string,
	host string,
	options RequestOptions,
) (*http.Request, error) {
	if options.Headers == nil {
		options.Headers = &http.Header{}
	}

	options.Headers.Set("Accept-Encoding", "")

	if options.Jwt != nil {
		options.Headers.Set(authorizationHeader, fmt.Sprintf("Bearer %s", *options.Jwt))
	}

	request, err := http.NewRequest(
		options.Method,
		fmt.Sprintf("%s://%s%s", schema, host, options.Path),
		options.Body,
	)
	if err != nil {
		return nil, err
	}

	request.Header = *options.Headers
	request = request.WithContext(ctx)
	if options.QueryParams != nil {
		request.URL.RawQuery = options.QueryParams.Encode()
	}

	return request, nil
}

func sendRequest(
	httpClient http.Client,
	request *http.Request,
	dontFollowRedirect bool,
) (*http.Response, error) {
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	statusCode := response.StatusCode
	switch {
	case statusCode >= 200 && statusCode <= 299:
		return response, nil
	case statusCode >= 300 && statusCode <= 399:
		if dontFollowRedirect {
			return response, nil
		}
		_ = response.Body.Close()

		return nil, fmt.Errorf("Unsupported Redirect Response: %v", statusCode)
	case statusCode >= 400 && statusCode <= 599:
		responseBody := response.Body
		bodyBytes, _ := ioutil.ReadAll(responseBody)
		bodyReader := bytes.NewReader(bodyBytes)

		var info interface{}
		err := readJSON(bodyReader, &info)
		_ = responseBody.Close()
		if err != nil {
			return nil, BodyNotJSONError{err, statusCode, bodyBytes}
		}

		return nil, &ErrorResponse{
			Status:  statusCode,
			Headers: response.Header,
			Info:    info,
		}
	}
	_ = response.Body.Close()

	return nil, fmt.Errorf("Unsupported Response Code: %v", statusCode)
}

// readJSON reads the body of an http response as a JSON document.
func readJSON(body io.Reader, dest interface{}) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(dest)
	if err != nil {
		return err
	}

	return nil
}
