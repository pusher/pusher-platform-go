package instance_test

import (
	"context"
	"fmt"

	"github.com/pusher/pusher-platform-go/auth"
	"github.com/pusher/pusher-platform-go/client"
	"github.com/pusher/pusher-platform-go/instance"
)

func ExampleNew() {
	instance.New(instance.Options{
		Locator:        "version:cluster:instance-id",
		Key:            "key:secret",
		ServiceName:    "service-name",
		ServiceVersion: "service-version",
	})
}

func ExampleNew_request() {
	app, err := instance.New(instance.Options{
		Locator:        "version:cluster1:instance-id1",
		Key:            "key1:secret1",
		ServiceName:    "service-name1",
		ServiceVersion: "service-version1",
	})
	if err != nil {
		// Do something with error
	}

	ctx := context.Background()
	response, err := app.Request(ctx, client.RequestOptions{
		Method: "GET",
		Path:   "/foo/bar",
	})
	if err != nil {
		// Do something with error
	}

	if response.StatusCode == 200 {
		// Do something with response
	}
}

func ExampleNew_authenticate() {
	app, err := instance.New(instance.Options{
		Locator:        "version:cluster:instance-id",
		Key:            "key:secret",
		ServiceName:    "service-name",
		ServiceVersion: "service-version",
	})
	if err != nil {
		// Do something with error
	}

	// For a more detailed example: check out the auth package
	userID := "test-user"
	authResponse, err := app.Authenticate(auth.Payload{
		GrantType: "client_credentials",
	}, auth.Options{
		UserID: &userID,
	})
	if err != nil {
		// Do something with err
	}

	err = authResponse.Error()
	if err != nil {
		// Do someting with response error
		// This should usually be a write to an HTTP stream with the headers and status
		fmt.Printf("Response status: %v", authResponse.Status)
		fmt.Printf("Response headers: %v", authResponse.Headers)
		fmt.Printf("Response error: %v", err.Error())
	}

	tokenResponse := authResponse.TokenResponse()
	if tokenResponse != nil {
		// Send the token back to a client
		fmt.Printf("Token: %v", tokenResponse.AccessToken)
		fmt.Printf("Token expiry: %v", tokenResponse.ExpiresIn)
		fmt.Printf("Token type: %v", tokenResponse.TokenType)
	}
}

func ExampleNew_generateAccessToken() {
	app, err := instance.New(instance.Options{
		Locator:        "version:cluster:instance-id",
		Key:            "key:secret",
		ServiceName:    "service-name",
		ServiceVersion: "service-version",
	})
	if err != nil {
		// Do something with error
	}

	userID := "test-user"
	token, err := app.GenerateAccessToken(auth.Options{
		UserID: &userID,
	})
	if err != nil {
		// Do something with error
	}

	fmt.Println(token.Token)
}

func ExampleParseInstanceLocator() {
	components, err := instance.ParseInstanceLocator("version:cluster:instance-id")
	if err != nil {
		// Do something with error
	}

	fmt.Printf("Platform version: %s", components.PlatformVersion)
	fmt.Printf("Cluster: %s", components.Cluster)
	fmt.Printf("Instance ID: %s", components.InstanceID)
}

func ExampleParseKey() {
	components, err := instance.ParseKey("key:secret")
	if err != nil {
		// Do something with error
	}

	fmt.Printf("Key: %s", components.Key)
	fmt.Printf("Secret: %s", components.Secret)
}
