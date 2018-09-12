package pusherplatform

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const authorizationHeader = "Authorization"

// BaseClient is a low level interface for lcinets of the elements protocol
type BaseClient interface {
	Request(ctx context.Context, options RequestOptions) (*http.Response, error)
}

// BaseClientOptions includes configuration options for a new base client
type BaseClientOptions struct {
	Host               string
	Jwt                *string
	TLSConfig          *tls.Config
	Timeout            time.Duration
	DontFollowRedirect bool
}

// NewBaseClient builds a new BaseClient
func NewBaseClient(options BaseClientOptions) BaseClient {
	return newBaseClient(options)
}

// Request allows making HTTP calls
func (c *baseClient) Request(ctx context.Context, options RequestOptions) (*http.Response, error) {
	request, err := buildRequest(ctx, c.jwt, c.schema, c.host, options)
	if err != nil {
		return nil, err
	}

	return sendRequest(c.http, request, c.options.DontFollowRedirect)
}

// Implements the BaseClient interface
type baseClient struct {
	host    string
	schema  string
	jwt     *string
	http    http.Client
	options BaseClientOptions
}

func newBaseClient(options BaseClientOptions) *baseClient {
	c := new(baseClient)
	c.host = options.Host
	c.jwt = options.Jwt
	c.schema = "https"

	if options.TLSConfig != nil {
		transport := &http.Transport{
			Proxy:              http.ProxyFromEnvironment,
			TLSClientConfig:    options.TLSConfig,
			DisableCompression: true,
		}

		c.http = http.Client{
			Transport: transport,
			Timeout:   options.Timeout,
		}
	}

	// Control how the http client handles redirect responses.
	// See: https://godoc.org/net/http#Client
	// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Redirections
	if options.DontFollowRedirect {
		c.http.CheckRedirect = func(r *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	c.options = options

	return c
}

func buildRequest(
	ctx context.Context,
	jwtPtr *string,
	schema string,
	host string,
	options RequestOptions,
) (*http.Request, error) {
	if options.Headers == nil {
		options.Headers = &http.Header{}
	}

	options.Headers.Set("Accept-Encoding", "")

	var jwt string
	if options.Jwt != nil {
		jwt = *options.Jwt
	} else if jwtPtr != nil {
		jwt = *jwtPtr
	}
	if jwt != "" {
		options.Headers.Set(authorizationHeader, fmt.Sprintf("Bearer %s", jwt))
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
	http http.Client,
	request *http.Request, dontFollowRedirect bool,
) (*http.Response, error) {
	response, err := http.Do(request)
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
		err := readJSON(bodyReader, info)
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
