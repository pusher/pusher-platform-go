package helpers

import (
	"errors"
	"fmt"
	"strings"
)

const hostBase = "pusherplatform.io"

// InstanceLocatorComponents contains information contained
// within an instance locator string which is of the format
// <version>:<cluster>:<instance-id>
type InstanceLocatorComponents struct {
	PlatformVersion string
	Cluster         string
	InstanceID      string
}

// Host returns the host based on the cluster
// specified in the instance locator
func (i InstanceLocatorComponents) Host() string {
	return fmt.Sprintf("%s.%s", i.Cluster, hostBase)
}

// KeyComponents contains the key and secret part
// of the key which is of the format <key>:<secret>
type KeyComponents struct {
	Key    string
	Secret string
}

// ParseInstanceLocator splits the instance locator string to retrieve the
// service version, cluster and instance id which is returned as
// an `InstanceLocatorComponents` type
func ParseInstanceLocator(instanceLocator string) (InstanceLocatorComponents, error) {
	components, err := getColonSeperatedComponents(instanceLocator, 3)
	if err != nil {
		return InstanceLocatorComponents{}, errors.New(
			"Instance locator must be of the format <version>:<cluster>:<instance-id>",
		)
	}

	return InstanceLocatorComponents{
		PlatformVersion: components[0],
		Cluster:         components[1],
		InstanceID:      components[2],
	}, nil
}

// ParseKey splits the key to retrieve the public key and secret
// which is returned as part of `KeyComponents`
func ParseKey(key string) (KeyComponents, error) {
	components, err := getColonSeperatedComponents(key, 2)
	if err != nil {
		return KeyComponents{}, errors.New(
			"Key must be of the format <key>:<secret>",
		)
	}

	return KeyComponents{
		Key:    components[0],
		Secret: components[1],
	}, nil
}

// Generic function to split strings by `:`
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
