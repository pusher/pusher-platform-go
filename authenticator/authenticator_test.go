package authenticator

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
		authResponse, err := authenticator.Authenticate(
			AuthenticatePayload{GrantType: "client_credentials"},
			AuthenticateOptions{
				UserID: &userID,
			},
		)
		if err != nil {
			t.Fatalf("Expected no error, but got %+v", err)
		}

		if authResponse.Status != http.StatusOK {
			t.Fatalf("Expected a 200 status, but got %v", authResponse.Status)
		}

		if tokenResponse, ok := authResponse.Body.(TokenResponse); !ok {
			t.Fatal("Expected auth response to contain a token response")

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
		}
	})

	t.Run("Authenticate with Su claim", func(t *testing.T) {
		authResponse, err := authenticator.Authenticate(
			AuthenticatePayload{GrantType: "client_credentials"},
			AuthenticateOptions{
				Su: true,
			},
		)
		if err != nil {
			t.Fatalf("Expected no error, but got %+v", authResponse.Status)
		}

		if authResponse.Status != http.StatusOK {
			t.Fatalf("Expected a 200 status, but got %v", authResponse.Status)
		}

		if tokenResponse, ok := authResponse.Body.(TokenResponse); !ok {
			t.Fatal("Expected auth response to contain a token response")

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
		authResponse, err := authenticator.Authenticate(
			AuthenticatePayload{GrantType: "client_credentials"},
			AuthenticateOptions{
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

		if tokenResponse, ok := authResponse.Body.(TokenResponse); !ok {
			t.Fatal("Expected auth response to contain a token response")

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
		authResponse, err := authenticator.Authenticate(
			AuthenticatePayload{GrantType: "client_credentials"},
			AuthenticateOptions{
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

		if tokenResponse, ok := authResponse.Body.(TokenResponse); !ok {
			t.Fatal("Expected auth response to contain a token response")

			if tokenResponse.ExpiresIn != 60*60 {
				t.Fatalf("Expected token expiry to be 1 hour, but got %v", tokenResponse.ExpiresIn)
			}
		}
	})
}

func TestAuthenticateFailure(t *testing.T) {
	userID := "user-id"
	authenticator := New("instance-id", "key", "secret")
	authResponse, err := authenticator.Authenticate(
		AuthenticatePayload{GrantType: "custom_grant_type"},
		AuthenticateOptions{
			UserID: &userID,
		},
	)
	if err != nil {
		t.Fatalf("Expected no error, but got %+v", authResponse.Status)
	}

	if authResponse.Status != http.StatusUnprocessableEntity {
		t.Fatalf("Expected a 422 status, but got %v", authResponse.Status)
	}

	if errorBody, ok := authResponse.Body.(ErrorBody); !ok {
		t.Fatal("Expected auth response to have an error body")

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
