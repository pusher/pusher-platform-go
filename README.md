# pusher-platform-go

Pusher Platform SDK for Go.

## Installation

```
go get github.com/pusher/pusher-platform-go
```

## Usage

In order to access Pusher Platform, instantiate an object first. This can be done like so

```go
import platform "github.com/pusher/pusher-platform-go"

instance := platform.NewInstance(platform.InstanceOptions{
	Locator: "<YOUR-INSTANCE-LOCATOR>",
	Key: "<YOUR-KEY>",
	ServiceName: "<SERVICE-NAME-TO-CONNECT-TO>",
	ServiceVersion: "<SERVICE-VERSION>",
})
```

The `Locator` and `Key` can be found in the [dashboard](https://dash.pusher.com). The `ServiceVersion` and `ServiceVersion` represent the name of the service to connect to and the version of the service to connect to respectively.

It is also possible to specify the `Host` to connect to, which will override the cluster value obtained from the instance locator.

## Authentication

Instance objects provide access to the `Authenticator` which can be used to build authentication endpoints. Authentication endpoints issue access tokens used by Pusher Platform clients to access the API.

This can be done like so

```go
userID := "abc"
authResponse, err := instance.Authenticator().Authenticate(platform.AuthenticatePayload{
	GrantType: "client_credentials"
}, platform.AuthenticateOptions{
	UserID: &userID
})
if err != nil {
	...
}

if status == 200 {
	tokenResponse := authResponse.Body.(platform.TokenResponse)
	token := tokenResponse.AccessToken
	...
} else {
	errorBody := authResponse.Body.(platform.ErrorBody)
	err := errorBody.ErrorType
	...
}

```

The struct returned by `Authenticate` is an `AuthenticationResponse` which contains a `Body` field that can be used to retrieve the a `TokenResponse` or an `ErrorBody`. The `Status` is `200` if the call was successful, in which case a `TokenResponse` is expected, otherwise a suitable status code is returned along with an `ErrorBody`. In addition to this, there is also a `Headers` field that returns headers which can be returned back to the client.

## Request API

Instance objects provide a low-level request API, which can be used to make HTTP calls to the platform.

```go
import (
	"context"

	platform "github.com/pusher/pusher-platform-go"
)

instance := platform.NewInstance(platform.InstanceOptions{
	Locator: "<YOUR-INSTANCE-LOCATOR>",
	Key: "<YOUR-KEY>",
	ServiceName: "<SERVICE-NAME-TO-CONNECT-TO>",
	ServiceVersion: "<SERVICE-VERSION>",
})

ctx := context.background()
resp, err := instsance.Request(ctx, RequestOptions{
	Method: "POST",
	Path: "/users",
})
if err != nil {
	...
}

// do something with response

```

The `Request` method always takes a `RequestOptions` that can be used to configure the request.

## Issues, Bugs and Feature Requests

Feel free to create an issue on Github if you find anything wrong. Please use the existing template. If you wish to contribute, please open a pull request.

## License

pusher-platform-go is released under the MIT license. See LICENSE for details.
