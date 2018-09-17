package auth_test

import (
	"github.com/pusher/pusher-platform-go/auth"
	"github.com/pusher/pusher-platform-go/instance"
)

func ExampleNew() {
	locator := "v1:cluster:instance-id"
	key := "key:secret"

	locatorComponents, err := instance.ParseInstanceLocator(locator)
	if err != nil {
		// Do something with err
	}

	keyComponents, err := instance.ParseKey(key)
	if err != nil {
		// Do something with err
	}

	auth.New(locatorComponents.InstanceID, keyComponents.Key, keyComponents.Secret)
}

func ExampleNew_authenticate() {
	locator := "v1:cluster:instance-id"
	key := "key:secret"

	locatorComponents, err := instance.ParseInstanceLocator(locator)
	if err != nil {
		// Do something with err
	}

	keyComponents, err := instance.ParseKey(key)
	if err != nil {
		// Do something with err
	}

	userID := "test-user"
	authenticator := auth.New(locatorComponents.InstanceID, keyComponents.Key, keyComponents.Secret)
	_, err = authenticator.Do(auth.Payload{
		GrantType: auth.GrantTypeClientCredentials,
	}, auth.Options{
		UserID: &userID,
	})
	if err != nil {
		// Do something with error
	}
}

func ExampleNew_generateAccessToken() {
	locator := "v1:cluster:instance-id"
	key := "key:secret"

	locatorComponents, err := instance.ParseInstanceLocator(locator)
	if err != nil {
		// Do something with err
	}

	keyComponents, err := instance.ParseKey(key)
	if err != nil {
		// Do something with err
	}

	userID := "test-user"
	authenticator := auth.New(locatorComponents.InstanceID, keyComponents.Key, keyComponents.Secret)
	_, err = authenticator.GenerateAccessToken(auth.Options{
		UserID: &userID,
	})
	if err != nil {
		// Do something with error
	}
}

func ExampleResponse_Error() {
	locator := "v1:cluster:instance-id"
	key := "key:secret"

	locatorComponents, err := instance.ParseInstanceLocator(locator)
	if err != nil {
		// Do something with err
	}

	keyComponents, err := instance.ParseKey(key)
	if err != nil {
		// Do something with err
	}

	userID := "test-user"
	authenticator := auth.New(locatorComponents.InstanceID, keyComponents.Key, keyComponents.Secret)
	authResponse, err := authenticator.Do(auth.Payload{
		GrantType: auth.GrantTypeClientCredentials,
	}, auth.Options{
		UserID: &userID,
	})
	if err != nil {
		// Do something with error
	}

	err = authResponse.Error()
	if err != nil {
		// Do something with error
	}
}

func ExampleResponse_TokenResponse() {
	locator := "v1:cluster:instance-id"
	key := "key:secret"

	locatorComponents, err := instance.ParseInstanceLocator(locator)
	if err != nil {
		// Do something with err
	}

	keyComponents, err := instance.ParseKey(key)
	if err != nil {
		// Do something with err
	}

	userID := "test-user"
	authenticator := auth.New(locatorComponents.InstanceID, keyComponents.Key, keyComponents.Secret)
	authResponse, err := authenticator.Do(auth.Payload{
		GrantType: auth.GrantTypeClientCredentials,
	}, auth.Options{
		UserID: &userID,
	})
	if err != nil {
		// Do something with error
	}

	err = authResponse.Error()
	if err != nil {
		// Do something with error
	}

	tokenResponse := authResponse.TokenResponse()
	if tokenResponse != nil {
		// Do something with response
	}
}
