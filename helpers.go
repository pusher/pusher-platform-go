package pusherplatform

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	jwt "github.com/pusher/jwt-go"
)

// Signs a token with the secret key
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

// Check if a string slice contains the provided string
func strSliceContains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}

	return false
}

// ReadJSON reads the body of an http response as a JSON document.
func readJSON(body io.Reader, dest interface{}) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&dest)
	if err != nil {
		return err
	}

	return nil
}

// Splits the instance locator string to retrieve the
// service version, cluster and instance id
func getinstanceLocatorComponents(
	instanceLocator string,
) (apiVersion string, cluster string, instanceID string, err error) {
	components, err := getColonSeperatedComponents(instanceLocator, 3)
	if err != nil {
		return "", "", "", errors.New(
			"Incorrect instance locator format given, please get your instance locator from your user dashboard",
		)
	}
	return components[0], components[1], components[2], nil
}

// Splits the key to retrieve the public key and secret
func getKeyComponents(key string) (keyID string, keySecret string, err error) {
	components, err := getColonSeperatedComponents(key, 2)
	if err != nil {
		return "", "", errors.New(
			"Incorrect key format given, please get your key from your user dashboard",
		)
	}
	return components[0], components[1], nil
}

// Generic function to split strings by :
func getColonSeperatedComponents(s string, expectedComponents int) ([]string, error) {
	if s == "" {
		return nil, errors.New("Empty string")
	}

	components := strings.Split(s, ":")
	if len(components) != expectedComponents {
		return nil, errors.New("Incorrect format")
	}

	for _, component := range components {
		if component == "" {
			return nil, errors.New("Incorrect format")
		}
	}

	return components, nil
}
