package client_test

import (
	"context"

	"github.com/pusher/pusher-platform-go/client"
	"github.com/pusher/pusher-platform-go/instance"
)

func ExampleNew() {
	client.New(client.Options{
		Host: "mycoolhost.io",
	})
}

func ExampleNew_request() {
	httpClient := client.New(client.Options{
		Host: "mycoolhost.io",
	})

	ctx := context.Background()
	response, err := httpClient.Request(ctx, client.RequestOptions{
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

func ExampleNew_instance() {
	// Instance with client passed in
	// This will override the client that is constructed when creating an instance
	httpClient := client.New(client.Options{
		Host: "mycoolhost.io",
	})
	_, err := instance.New(instance.Options{
		Locator:        "version:cluster:instance-id",
		Key:            "key:secret",
		ServiceName:    "service-name",
		ServiceVersion: "service-version",
		Client:         httpClient,
	})
	if err != nil {
		// Do something with error
	}
}
