package instance

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/pusher/pusher-platform-go/authenticator"
	"github.com/pusher/pusher-platform-go/client"
	"github.com/pusher/pusher-platform-go/instance/helpers"
)

var (
	slashFoldingRegexp  regexp.Regexp
	trailingSlashRegexp regexp.Regexp
)

func init() {
	slashFoldingRegexp = *regexp.MustCompile("\\/+")
	trailingSlashRegexp = *regexp.MustCompile("\\/$")
}

// Instance allows making HTTP requests to a service
type Instance interface {
	Request(ctx context.Context, options client.RequestOptions) (*http.Response, error)
	Authenticator() authenticator.Authenticator
}

// Options to initialize a new instance.
type Options struct {
	// Instance locator unique to an app
	Locator string
	// Key unique to an app
	Key string
	// Service name to connect to
	ServiceName string
	// Version of service to connect to
	ServiceVersion string
	// Optional Client
	// If not provided, it will be constructed
	Client client.Client
}

type instance struct {
	instanceID      string
	serviceName     string
	serviceVersion  string
	cluster         string
	platformVersion string

	keyID     string
	keySecret string

	authenticator authenticator.Authenticator
	client        client.Client
}

// New creates a new instance satisfying the Instance interface
func New(options Options) (Instance, error) {
	locatorComponents, err := helpers.ParseInstanceLocator(options.Locator)
	if err != nil {
		return nil, err
	}

	keyComponents, err := helpers.ParseKey(options.Key)
	if err != nil {
		return nil, err
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
		authenticator: authenticator.New(
			locatorComponents.InstanceID,
			keyComponents.Key,
			keyComponents.Secret,
		),
		client: underlyingClient,
	}, nil
}

// Request allows making HTTP requests to services
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

// Authenticator exposes the Authenticator interface to allow
// authentication and token generation
func (i *instance) Authenticator() authenticator.Authenticator {
	return i.authenticator
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
