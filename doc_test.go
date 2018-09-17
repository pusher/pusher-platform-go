package pusher_test

import (
	"context"

	"github.com/pusher/pusher-platform-go/client"
	"github.com/pusher/pusher-platform-go/instance"
)

func Example() {
	// Instantiate
	serviceInstance, err := instance.New(instance.Options{
		Locator:        "version:cluster:instance-id",
		Key:            "key:secret",
		ServiceName:    "service-name",
		ServiceVersion: "service-version",
	})
	if err != nil {
		// Do something with error
	}

	// Network requests
	ctx := context.Background()
	response, err := serviceInstance.Request(ctx, client.RequestOptions{
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
