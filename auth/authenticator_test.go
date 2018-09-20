package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	jwt "github.com/pusher/jwt-go"
)

func TestAuthenticateSuccess(t *testing.T) {
	userID := "test-user"
	authenticator := New("instance-id", "key", "secret")

	t.Run("Basic authenticate", func(t *testing.T) {
		authResponse, err := authenticator.Do(
			Payload{GrantType: "client_credentials"},
			Options{
				UserID: &userID,
			},
		)
		if err != nil {
			t.Fatalf("Expected no error, but got %+v", err)
		}

		if authResponse.Status != http.StatusOK {
			t.Fatalf("Expected a 200 status, but got %v", authResponse.Status)
		}

		if err := authResponse.Error(); err != nil {
			t.Fatalf("Expected no error in auth response, but got %s", err.Error())
		}

		if tokenResponse := authResponse.TokenResponse(); tokenResponse != nil {
			if tokenResponse.ExpiresIn != 24*60*60 {
				t.Fatalf("Expected token to expire in a day, but got %v", tokenResponse.ExpiresIn)
			}

			if tokenResponse.TokenType != "Bearer" {
				t.Fatalf("Expected token type to be Bearer, but got %s", tokenResponse.TokenType)
			}

			token := tokenResponse.AccessToken
			parsedToken, err := parseToken(token)
			if err != nil {
				t.Fatalf("Expected no error when parsing token, but got %+v", err)
			}

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				t.Fatal("Expected token claims to be of type MapClaims")
			}

			userIDClaim, ok := claims["sub"]
			if !ok {
				t.Fatal("Expected `sub` claim to exist, but it didn't")
			}

			if userID != userIDClaim {
				t.Fatalf("Expected `sub` claim to be %s, but got %s", userID, userIDClaim)
			}

			issClaim, ok := claims["iss"]
			if !ok {
				t.Fatal("Expected `iss` claim to exist, but it didn't")
			}

			if issClaim != "api_keys/key" {
				t.Fatalf("Expected `iss` claim to be api_keys/key, but got %s", issClaim)
			}

			instanceClaim, ok := claims["instance"]
			if !ok {
				t.Fatal("Expected `instance` claim to exist, but it didn't")
			}

			if instanceClaim != "instance-id" {
				t.Fatalf("Expected `instance` claim to be instance-id, but got %s", instanceClaim)
			}
		}
	})

	t.Run("Authenticate with Su claim", func(t *testing.T) {
		authResponse, err := authenticator.Do(
			Payload{GrantType: "client_credentials"},
			Options{
				Su: true,
			},
		)
		if err != nil {
			t.Fatalf("Expected no error, but got %+v", authResponse.Status)
		}

		if authResponse.Status != http.StatusOK {
			t.Fatalf("Expected a 200 status, but got %v", authResponse.Status)
		}

		if tokenResponse := authResponse.TokenResponse(); tokenResponse != nil {
			token := tokenResponse.AccessToken
			parsedToken, err := parseToken(token)
			if err != nil {
				t.Fatalf("Expected no error when parsing token, but got %+v", err)
			}

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				t.Fatal("Expected token claims to be of type MapClaims")
			}

			_, ok = claims["sub"]
			if ok {
				t.Fatal("Expected `sub` claim to not be present, but it was")
			}

			_, ok = claims["su"]
			if !ok {
				t.Fatal("Expected `su` claim to be present, but it wasn't")
			}
		}
	})

	t.Run("Authenticate with token claims", func(t *testing.T) {
		authResponse, err := authenticator.Do(
			Payload{GrantType: "client_credentials"},
			Options{
				UserID: &userID,
				ServiceClaims: map[string]interface{}{
					"foo": "bar",
				},
			},
		)
		if err != nil {
			t.Fatalf("Expected no error, but got %+v", authResponse.Status)
		}

		if authResponse.Status != http.StatusOK {
			t.Fatalf("Expected a 200 status, but got %v", authResponse.Status)
		}

		if tokenResponse := authResponse.TokenResponse(); tokenResponse != nil {
			token := tokenResponse.AccessToken
			parsedToken, err := parseToken(token)
			if err != nil {
				t.Fatalf("Expected no error when parsing token, but got %+v", err)
			}

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				t.Fatal("Expected token claims to be of type MapClaims")
			}

			fooClaim, ok := claims["foo"]
			if !ok {
				t.Fatal("Expected `foo` claim to be present, but it wasn't")
			}

			if fooClaim != "bar" {
				t.Fatalf("Expected `foo` claim value to be bar, but got %s", fooClaim)
			}
		}
	})

	t.Run("Authenticate with custom token expiry time", func(t *testing.T) {
		expiry := 1 * time.Hour
		authResponse, err := authenticator.Do(
			Payload{GrantType: "client_credentials"},
			Options{
				UserID:      &userID,
				TokenExpiry: &expiry,
			},
		)
		if err != nil {
			t.Fatalf("Expected no error, but got %+v", authResponse.Status)
		}

		if authResponse.Status != http.StatusOK {
			t.Fatalf("Expected a 200 status, but got %v", authResponse.Status)
		}

		if tokenResponse := authResponse.TokenResponse(); tokenResponse != nil {
			if tokenResponse.ExpiresIn != 60*60 {
				t.Fatalf("Expected token expiry to be 1 hour, but got %v", tokenResponse.ExpiresIn)
			}
		}
	})
}

func TestAuthenticateFailure(t *testing.T) {
	userID := "user-id"
	authenticator := New("instance-id", "key", "secret")
	authResponse, err := authenticator.Do(
		Payload{GrantType: "custom_grant_type"},
		Options{
			UserID: &userID,
		},
	)
	if err != nil {
		t.Fatalf("Expected no error, but got %+v", authResponse.Status)
	}

	if authResponse.Status != http.StatusUnprocessableEntity {
		t.Fatalf("Expected a 422 status, but got %v", authResponse.Status)
	}

	if errorBody := authResponse.Error(); errorBody != nil {
		if errorBody.ErrorType != "token_provider/invalid_grant_type" {
			t.Fatalf(
				"Expected error type to be an invalid_grant_type, but got %s",
				errorBody.ErrorType,
			)
		}
	}
}

// Helpers

func parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})
}
