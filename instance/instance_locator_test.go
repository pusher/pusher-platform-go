package instance

import (
	"testing"
)

func TestParseInstanceLocatorSuccess(t *testing.T) {
	platformVersion := "v1"
	cluster := "local"
	instanceID := "aaa-bbb-ccc"
	instanceLocator := platformVersion + ":" + cluster + ":" + instanceID

	components, err := ParseInstanceLocator(instanceLocator)
	if err != nil {
		t.Fatalf("Expected no error, but got %+v", err)
	}

	if components.Cluster != cluster {
		t.Fatalf("Expected cluster to be %s, but got %s", cluster, components.Cluster)
	}

	if components.InstanceID != instanceID {
		t.Fatalf("Expected instance id to be %s, but got %s", instanceID, components.InstanceID)
	}

	if components.PlatformVersion != platformVersion {
		t.Fatalf("Expected platform version to be %s, but got %s", platformVersion, components.PlatformVersion)
	}
}

func TestParseInstanceLocatorIncorrectFormat(t *testing.T) {
	_, err := ParseInstanceLocator("invalid:format")
	if err != nil {
		if err.Error() != "Instance locator must be of the format <version>:<cluster>:<instance-id>" {
			t.Fatal("Expected incorrect format error")
		}
	}
}

func TestParseKeySuccess(t *testing.T) {
	keyID := "key"
	keySecret := "secret"
	key := keyID + ":" + keySecret

	components, err := ParseKey(key)
	if err != nil {
		t.Fatalf("Expected no error, but got %+v", err)
	}

	if components.Key != keyID {
		t.Fatalf("Expected keyID to be %s, but got %s", keyID, components.Key)
	}

	if components.Secret != keySecret {
		t.Fatalf("Expected keySecret to be %s, but got %s", keySecret, components.Secret)
	}
}

func TestParseKeyIncorrectFormat(t *testing.T) {
	_, err := ParseKey("invalid")
	if err != nil {
		if err.Error() != "Key must be of the format <key>:<secret>" {
			t.Fatal("Expected incorrect format error")
		}
	}
}
