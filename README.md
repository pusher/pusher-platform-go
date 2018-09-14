# pusher-platform-go

Pusher Platform SDK for Go.

## Installation

```
go get github.com/pusher/pusher-platform-go
```

## Usage

In order to access Pusher Platform, instantiate an object first. This can be done like so

```go
import "github.com/pusher/pusher-platform-go/instance"

app, err := instance.New(instance.Options{
	Locator: "<YOUR-INSTANCE-LOCATOR>",
	Key: "<YOUR-KEY>",
	ServiceName: "<SERVICE-NAME-TO-CONNECT-TO>",
	ServiceVersion: "<SERVICE-VERSION>",
})
if err != nil {
	...
}
```

The `Locator` and `Key` can be found in the [dashboard](https://dash.pusher.com). The `ServiceVersion` and `ServiceVersion` represent the name of the service to connect to and the version of the service to connect to respectively.

If you'd like to specify a custom host, `instance.Options` contains a `Client` property which accepts a `client.Client`. This allows
for an externally constructed `Client` to be passed when creating the `Instance`.

```go
import (
	"github.com/pusher/pusher-platform-go/client"
	"github.com/pusher/pusher-platform-go/instance"
)

underlyingClient := client.New(client.Options{
	Host: "mycoolhost.io"
})

app, err := instance.New(instance.Options{
	Locator: "<YOUR-INSTANCE-LOCATOR>",
	Key: "<YOUR-KEY>",
	ServiceName: "<SERVICE-NAME-TO-CONNECT-TO>",
	ServiceVersion: "<SERVICE-VERSION>",
	Client: underlyingClient,
})
if err != nil {
	...
}
```

## Authentication

Instance objects provide access to the `Authenticator` which can be used to build authentication endpoints. Authentication endpoints issue access tokens used by Pusher Platform clients to access the API.

This can be done like so

```go
import (
	"github.com/pusher/pusher-platform-go/instance"
	"github.com/pusher/pusher-platform-go/auth"
)

app, err := instance.New(instance.Options{
	Locator: "<YOUR-INSTANCE-LOCATOR>",
	Key: "<YOUR-KEY>",
	ServiceName: "<SERVICE-NAME-TO-CONNECT-TO>",
	ServiceVersion: "<SERVICE-VERSION>",
})
if err != nil {
	...
}

userID := "abc"
authResponse, err := app.Authenticate(auth.Payload{
	GrantType: "client_credentials"
}, auth.Options{
	UserID: &userID,
})
if err != nil {
	...
}

// In your HTTP handler that acts as the token providing endpoint
if err := authResponse.Error(); err != nil {
	w.WriteHeader(authResponse.Status)
} else {
	if tokenResponse := authResponse.TokenResponse(); tokenResponse != nil {
		token := tokenResponse.Token
		w.WriteHeader(authResponse.Status)
		w.Write(token)
	}
}

```

The struct returned by `Authenticate` is an `AuthenticationResponse` which contains a `Body` field that can be used to retrieve the a `TokenResponse` or an `ErrorBody`. The `Status` is `200` if the call was successful, in which case a `TokenResponse` is expected, otherwise a suitable status code is returned along with an `ErrorBody`. In addition to this, there is also a `Headers` field that returns headers which can be returned back to the client.

## Request API

Instance objects provide a low-level request API, which can be used to make HTTP calls to the platform.

```go
import (
	"context"

	"github.com/pusher/pusher-platform-go/instance"
	"github.com/pusher/pusher-platform-go/client"
)

app, err := instance.New(instance.Options{
	Locator: "<YOUR-INSTANCE-LOCATOR>",
	Key: "<YOUR-KEY>",
	ServiceName: "<SERVICE-NAME-TO-CONNECT-TO>",
	ServiceVersion: "<SERVICE-VERSION>",
})
if err != nil {
	...
}

ctx := context.background()
jwt := "your-jwt-token"
resp, err := app.Request(ctx, client.RequestOptions{
	Method: "POST",
	Path: "/users",
	Jwt: &jwt,
})
if err != nil {
	...
}

// do something with response

```

The `Request` method always takes a `RequestOptions` that can be used to configure the request.

## Tests

To run tests

```
go test ./...
```

## Issues, Bugs and Feature Requests

Feel free to create an issue on Github if you find anything wrong. Please use the existing template. If you wish to contribute, please open a pull request.

## License

pusher-platform-go is released under the MIT license. See LICENSE for details.
