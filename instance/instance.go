// Package instance allows making instance scoped requests.
//
// Instance provides a higher level abstraction over Client and exposes the Authenticator interface.
// It is the primary entrypoint that should be used to interact with the platform.
// Instances are tied to their instance locators that can be found in the dashboard at https://dash.pusher.com.
package instance

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/pusher/pusher-platform-go/auth"
	"github.com/pusher/pusher-platform-go/client"
)

var (
	slashFoldingRegexp  regexp.Regexp
	trailingSlashRegexp regexp.Regexp
)

func init() {
	slashFoldingRegexp = *regexp.MustCompile("\\/+")
	trailingSlashRegexp = *regexp.MustCompile("\\/$")
}

// Instance allows making HTTP requests to a service.
// It also allows access to the authenticator interface.
type Instance interface {
	Request(ctx context.Context, options client.RequestOptions) (*http.Response, error)
	Authenticate(payload auth.Payload, options auth.Options) (*auth.Response, error)
	GenerateAccessToken(options auth.Options) (auth.TokenWithExpiry, error)
}

// Options to initialize a new instance.
type Options struct {
	Locator        string        // Instance locator unique to an app
	Key            string        // Key unique to an app
	ServiceName    string        // Service name to connect to
	ServiceVersion string        // Version of service to connect to
	Client         client.Client // Optional Client, if not provided will be constructed
}

type instance struct {
	instanceID      string
	serviceName     string
	serviceVersion  string
	cluster         string
	platformVersion string

	keyID     string
	keySecret string

	authenticator auth.Authenticator
	client        client.Client
}

// New creates a new instance satisfying the Instance interface.
//
// Instance locator, key, service name and service version are all required.
// It will return an error if any of these are not provided.
func New(options Options) (Instance, error) {
	locatorComponents, err := ParseInstanceLocator(options.Locator)
	if err != nil {
		return nil, err
	}

	keyComponents, err := ParseKey(options.Key)
	if err != nil {
		return nil, err
	}

	if options.ServiceName == "" {
		return nil, errors.New("No service name provided")
	}

	if options.ServiceVersion == "" {
		return nil, errors.New("No service version provided")
	}

	underlyingClient := options.Client
	if options.Client == nil {
		underlyingClient = client.New(client.Options{
			Host: locatorComponents.Host(),
		})
	}

	return &instance{
		instanceID:      locatorComponents.InstanceID,
		serviceName:     options.ServiceName,
		serviceVersion:  options.ServiceVersion,
		cluster:         locatorComponents.Cluster,
		platformVersion: locatorComponents.PlatformVersion,
		keyID:           keyComponents.Key,
		keySecret:       keyComponents.Secret,
		authenticator: auth.New(
			locatorComponents.InstanceID,
			keyComponents.Key,
			keyComponents.Secret,
		),
		client: underlyingClient,
	}, nil
}

// Request allows making HTTP requests to services.
func (i *instance) Request(
	ctx context.Context,
	options client.RequestOptions,
) (*http.Response, error) {
	return i.client.Request(ctx, client.RequestOptions{
		Method:      options.Method,
		Path:        i.scopePath(options.Path),
		Jwt:         options.Jwt,
		Headers:     options.Headers,
		Body:        options.Body,
		QueryParams: options.QueryParams,
	})
}

// Authenticate exposes the Authenticator interface to allow
// authentication and token generation.
func (i *instance) Authenticate(payload auth.Payload, options auth.Options) (*auth.Response, error) {
	return i.authenticator.Do(payload, options)
}

// GenerateAccessToken exposes the Authenticator interface to allow token generation.
func (i *instance) GenerateAccessToken(options auth.Options) (auth.TokenWithExpiry, error) {
	return i.authenticator.GenerateAccessToken(options)
}

func (i *instance) scopePath(path string) string {
	return trailingSlashRegexp.ReplaceAllString(
		slashFoldingRegexp.ReplaceAllString(
			fmt.Sprintf("/services/%s/%s/%s/%s",
				i.serviceName,
				i.serviceVersion,
				i.instanceID,
				path),
			"/",
		),
		"/",
	)
}
