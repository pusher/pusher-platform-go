package pusherplatform

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
)

var (
	slashFoldingRegexp  regexp.Regexp
	trailingSlashRegexp regexp.Regexp
)

func init() {
	slashFoldingRegexp = *regexp.MustCompile("\\/+")
	trailingSlashRegexp = *regexp.MustCompile("\\/$")
}

const (
	hostBase = "pusherplatform.io"
)

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

	// Optional host
	// If not provided, it will be constructed on the basis
	// of the instance locator and the base host
	Host string
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
	platformVersion, cluster, instanceID, err := getinstanceLocatorComponents(options.Locator)
	if err != nil {
		return nil, err
	}

	keyID, keySecret, err := getKeyComponents(options.Key)
	if err != nil {
		return nil, err
	}

	var host string
	if options.Host != "" {
		host = options.Host
	} else {
		host = fmt.Sprintf("%s.%s", cluster, hostBase)
	}

	client := options.Client
	if options.Client == nil {
		client = NewBaseClient(BaseClientOptions{
			Host: host,
		})
	}

	return &instance{
		instanceID:      instanceID,
		serviceName:     options.ServiceName,
		serviceVersion:  options.ServiceVersion,
		cluster:         cluster,
		platformVersion: platformVersion,
		keyID:           keyID,
		keySecret:       keySecret,
		authenticator:   NewAuthenticator(instanceID, keyID, keySecret),
		client:          client,
	}, nil
}

// Request allows making HTTP requests to services
func (i *instance) Request(ctx context.Context, options RequestOptions) (*http.Response, error) {
	path := i.scopePath(options.Path)

	var jwt *string
	if options.Jwt == nil {
		tokenResponse, err := i.Authenticator().GenerateAccessToken(AuthenticateOptions{Su: true})
		if err != nil {
			return nil, err
		}

		jwt = &tokenResponse.Token
	} else {
		jwt = options.Jwt
	}

	return i.client.Request(ctx, RequestOptions{
		Method:      options.Method,
		Path:        path,
		Jwt:         jwt,
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
