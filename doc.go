package pusher

/*
Package pusher provides an interface to interact with the Pusher platform.
It is mostly intended to be used as a library on top of which products
that run on the platform are built.

To interact with products that run on the platform, it is better to use the product
specific library.

This package does not contain any top level exports and is divided into several small
sub-packages that can be selectively imported to pull in the required functionality.

Client

The client package provides an HTTP/2 client that connects to the platform. For most use cases,
construction of a client directly isn't required. The interface is very basic

type Client interface {
	Request(ctx context.Context, options RequestOptions) (*http.Response, error)
}

A client can be constructed using the New function that accepts an Options argument
that lets you configure the client.

type Options struct {
	Host               string
	TLSConfig          *tls.Config
	Timeout            time.Duration
	DontFollowRedirect bool
}

Following the construction of a client, requests may then be performed via the Request method on the
client object. This method allows passing in a context that may be used to cancel the request along
with a RequestOptions argument that allows configuring the request.

ctx := context.Background()
response, err := client.Request(ctx, RequestOptions{
	Method: "GET",
	Path: "/foo/bar",
})

RequestOptions contains basic HTTP request parameters many of which are optional. The Jwt attribute
allows passing in an optional JWT string that is then added to the request as an Authorization header.

type RequestOptions struct {
	Method      string
	Path        string
	Jwt         *string
	Headers     *http.Header
	Body        io.Reader
	QueryParams *url.Values
}

Instance

Instance represents a higher level abstraction over the client, which allows making requests only to
a specific instance of a product that is built on the platform. Instances are tied to their instance
locators that can be found in the dashboard at https://dash.pusher.com.

Constructing an Instance requires passing in an Options argument that specificies the instance locator
and several other parameters.

type Options struct {
	Locator        string
	Key            string
	ServiceName    string
	ServiceVersion string
	Client         client.Client
}

The locator and key can both be found in the dashboard and are essential to being able to successfully
make requests to a service. The list of services can also be found in the dashboard. Service versions
can be found in the documentation pages associcated with the services at https://docs.pusher.com.

Note that without a valid instance locator and key, an Instance cannot be constructed.

Once an Instance object has been constructed, making a request is a case of calling the Request method
on it, which is similar in signature to the client.Client interface. Instances make use of an
underlying Client that the request is delegated to.

The instance package also includes a helper package that exposes some functions to parse
instance locators and keys.

Instance locators are of the format <platform-version>:<cluster>:<instance-id> and keys are of the
format <key>:<secret>.

Authenticator

The authenticator package exposes an Authenticator that can be used to generate JWT tokens.
The interface contains two methods

type Authenticator interface {
	Authenticate(
		payload AuthenticatePayload,
		options AuthenticateOptions,
	) (AuthenticationResponse, error)
	GenerateAccessToken(opts AuthenticateOptions) (TokenWithExpiry, error)
}

Authenticate is used within the context of a token provider - which is usually written as an endpoint
that makes use of this function to authenticate a user. It can generate tokens with the user id
provided as part of the AuthenticateOptions

type AuthenticationResponse struct {
	Status  int
	Headers http.Header
	Body    interface{}
}

*/
