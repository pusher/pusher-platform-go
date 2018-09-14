package auth_test

import (
	"net/http"

	"github.com/pusher/pusher-platform-go/auth"
	"github.com/pusher/pusher-platform-go/instance"
)

func Example() {
	app, err := instance.New(instance.Options{
		Locator:        "version:cluster:instance-id",
		Key:            "key:secret",
		ServiceName:    "service-name",
		ServiceVersion: "service-version",
	})
	if err != nil {
		// Do something with error
	}

	http.HandleFunc("/token", func(w http.ResponseWriter, req *http.Request) {
		// Get the user ID from url query params
		// Alternatively, this can come from the body
		userID := req.URL.Query().Get("user_id")
		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Get an auth response
		authResponse, err := app.Authenticate(auth.Payload{
			GrantType: auth.GrantTypeClientCredentials,
		}, auth.Options{
			UserID: &userID,
		})
		if err != nil {
			// Error when getting response
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Auth response contains an ErrorBody
		// Access it and respond back
		err = authResponse.Error()
		if err != nil {
			w.WriteHeader(authResponse.Status)
			w.Write([]byte(err.Error()))
			return
		}

		// Check if we have a token
		tokenResponse := authResponse.TokenResponse()
		if tokenResponse == nil {
			// Token was empty
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Things went well
		// Send back status and write the token
		w.WriteHeader(authResponse.Status)
		w.Write([]byte(tokenResponse.AccessToken))
		return
	})
}
