# pusher-platform-go

Pusher Platform SDK for Go.

Detailed documentation and examples can be found [here](https://godoc.org/github.com/pusher/pusher-platform-go).

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


## Request API

Instance objects provide a request API, which can be used to make HTTP calls to the platform.

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
	Method: "GET",
	Path: "/users",
	Jwt: &jwt,
})
if err != nil {
	...
}

// do something with response

```

## Authenticator

Instance objects also provide access to methods that can be used to generate tokens and authenticate users.

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

userID := "user-id"
authResponse, err := app.Authenticate(
	auth.Payload{auth.GrantTypeClientCredentials},
	auth.Options{
		UserID: &userID,
	},
)
if err != nil {
	// Do something with error
}

// Do something with the auth response
```

## Tests

To run tests

```
go test ./...
```

## Issues, Bugs and Feature Requests

Feel free to create an issue on Github if you find anything wrong. Please use the existing template. If you wish to contribute, please open a pull request.

## License

pusher-platform-go is released under the MIT license. See LICENSE for details.
