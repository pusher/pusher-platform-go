package pusherplatform

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/pusher/pusher-platform-go/helpers"
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
	Request(ctx context.Context, options RequestOptions) (*http.Response, error)
	Authenticator() Authenticator
}

// InstanceOptions to initialize a new instance.
type InstanceOptions struct {
	// Instance locator unique to an app
	Locator string
	// Key unique to an app
	Key string
	// Service name to connect to
	ServiceName string
	// Version of service to connect to
	ServiceVersion string
	// Optional BaseClient
	// If not provided, it will be constructed
	Client BaseClient
}

type instance struct {
	instanceID      string
	serviceName     string
	serviceVersion  string
	cluster         string
	platformVersion string

	keyID     string
	keySecret string

	authenticator Authenticator
	client        BaseClient
}

// NewInstance creates a new instance satisfying the Instance interface
func NewInstance(options InstanceOptions) (Instance, error) {
	locatorComponents, err := helpers.ParseInstanceLocator(options.Locator)
	if err != nil {
		return nil, err
	}

	keyComponents, err := helpers.ParseKey(options.Key)
	if err != nil {
		return nil, err
	}

	client := options.Client
	if options.Client == nil {
		client = NewBaseClient(BaseClientOptions{
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
		authenticator: NewAuthenticator(
			locatorComponents.InstanceID,
			keyComponents.Key,
			keyComponents.Secret,
		),
		client: client,
	}, nil
}

// Request allows making HTTP requests to services
func (i *instance) Request(ctx context.Context, options RequestOptions) (*http.Response, error) {
	return i.client.Request(ctx, RequestOptions{
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
func (i *instance) Authenticator() Authenticator {
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

// readJSON reads the body of an http response as a JSON document.
func readJSON(body io.Reader, dest interface{}) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(dest)
	if err != nil {
		return err
	}

	return nil
}
