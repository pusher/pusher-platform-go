package pusherplatform

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/pusher/jwt-go"
)

const (
	defaultTokenExpiry         = 24 * time.Hour
	clientCredentialsGrantType = "client_credentials"
	tokenType                  = "bearer"
)

// Authenticator specifies the public facing interface
// for performing authentication and token generation
type Authenticator interface {
	Authenticate(
		payload AuthenticatePayload,
		options AuthenticateOptions,
	) (AuthenticationResponse, error)
	GenerateAccessToken(opts AuthenticateOptions) (TokenWithExpiry, error)
}

type authenticator struct {
	instanceID string
	keyID      string
	keySecret  string
}

// NewAuthenticator returns a new instance of an authenticator
// that conforms to the Authenticator interface
func NewAuthenticator(instanceID, keyID, keySecret string) Authenticator {
	return &authenticator{
		instanceID,
		keySecret,
		keySecret,
	}
}

// Authenticate generates access tokens based on the options provided
// and returns an AuthenticationResponse
func (auth *authenticator) Authenticate(
	payload AuthenticatePayload,
	options AuthenticateOptions,
) (AuthenticationResponse, error) {
	grantType := payload.GrantType
	if grantType != clientCredentialsGrantType {
		return AuthenticationResponse{
			Status: http.StatusUnprocessableEntity,
			Body: ErrorBody{
				ErrorType:        "token_provider/invalid_grant_type",
				ErrorDescription: fmt.Sprintf("The grant type provided %s is unsupported", grantType),
			},
		}, nil
	}

	tokenWithExpiry, err := auth.GenerateAccessToken(options)
	if err != nil {
		return AuthenticationResponse{}, err
	}

	return AuthenticationResponse{
		Status: http.StatusOK,
		Body: TokenResponse{
			AccessToken: tokenWithExpiry.Token,
			TokenType:   tokenType,
			ExpiresIn:   tokenWithExpiry.ExpiresIn,
		},
	}, nil
}

// GenerateAccessToken returns a TokenWithExpiry based on the options provided
func (auth *authenticator) GenerateAccessToken(options AuthenticateOptions) (TokenWithExpiry, error) {
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
		"exp":      now.Add(tokenExpiry),
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
