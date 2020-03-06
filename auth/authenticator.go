// Package auth exposes the Authenticator interface that allows for token generation.
//
// There is usually no need to construct an Authenticator by itself since the interface is
// exposed via an Instance and this should be the primary entry point.
package auth

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/pusher/jwt-go"
)

const (
	defaultTokenExpiry         = 24 * time.Hour
	clientCredentialsGrantType = "client_credentials"
	tokenType                  = "Bearer"
)

// Authenticator specifies the public facing interface
// for performing authentication and token generation.
type Authenticator interface {
	Do(payload Payload, options Options) (*Response, error)
	GenerateAccessToken(options Options) (TokenWithExpiry, error)
}

type authenticator struct {
	instanceID string
	keyID      string
	keySecret  string
}

// New returns a new instance of an authenticator that conforms to the Authenticator interface.
func New(instanceID, keyID, keySecret string) Authenticator {
	return &authenticator{
		instanceID,
		keyID,
		keySecret,
	}
}

// Do generates access tokens based on the options provided and returns a Response.
func (auth *authenticator) Do(
	payload Payload,
	options Options,
) (*Response, error) {
	grantType := payload.GrantType
	if grantType != clientCredentialsGrantType {
		return &Response{
			Status: http.StatusUnprocessableEntity,
			Body: &ErrorBody{
				ErrorType:        "token_provider/invalid_grant_type",
				ErrorDescription: fmt.Sprintf("The grant type provided %s is unsupported", grantType),
			},
		}, nil
	}

	tokenWithExpiry, err := auth.GenerateAccessToken(options)
	if err != nil {
		return nil, err
	}

	return &Response{
		Status: http.StatusOK,
		Body: &TokenResponse{
			AccessToken: tokenWithExpiry.Token,
			TokenType:   tokenType,
			ExpiresIn:   tokenWithExpiry.ExpiresIn,
		},
	}, nil
}

// GenerateAccessToken returns a TokenWithExpiry based on the options provided.
func (auth *authenticator) GenerateAccessToken(options Options) (TokenWithExpiry, error) {
	now := time.Now()
	var tokenExpiry time.Duration
	if options.TokenExpiry == nil {
		tokenExpiry = defaultTokenExpiry
	} else {
		tokenExpiry = *options.TokenExpiry
	}

	tokenClaims := jwt.MapClaims{
		"instance": auth.instanceID,
		"iss":      "api_keys/" + auth.keyID,
		"iat":      now.Unix(),
		"exp":      now.Add(tokenExpiry).Unix(),
	}

	if options.UserID != nil {
		tokenClaims["sub"] = *options.UserID
	}

	if options.Su {
		tokenClaims["su"] = true
	}

	if options.ServiceClaims != nil {
		for claimName, value := range options.ServiceClaims {
			tokenClaims[claimName] = value
		}
	}

	signedToken, err := signToken(auth.keySecret, tokenClaims)
	if err != nil {
		return TokenWithExpiry{}, err
	}

	return TokenWithExpiry{
		Token:     signedToken,
		ExpiresIn: tokenExpiry.Seconds(),
	}, nil
}

// Signs a token with the secret key.
func signToken(
	keySecret string,
	jwtClaims jwt.MapClaims,
) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err = token.SignedString([]byte(keySecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
