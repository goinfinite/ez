package valueObject

import (
	"testing"
)

func TestNewObjectStorageProviderRegion(t *testing.T) {
	t.Run("ValidObjectStorageProviderRegion", func(t *testing.T) {
		validObjectStorageProviderRegions := []string{
			"us-east-1",
			"ap-southeast-1",
			"eu-west-1",
			"iad",
			"lhr",
			"US1",
			"eastus",
			"westus",
		}

		for _, providerRegion := range validObjectStorageProviderRegions {
			_, err := NewObjectStorageProviderRegion(providerRegion)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), providerRegion)
			}
		}
	})

	t.Run("InvalidObjectStorageProviderRegion", func(t *testing.T) {
		invalidObjectStorageProviderRegions := []string{
			"",
			"127.0.0.1",
			"goinfinite..net",
			"us_east_1",
			"UNION SELECT * FROM USERS",
			"/providerRegion\n/providerRegion",
			"?param=value",
			"/providerRegion/'; DROP TABLE users; --",
		}

		for _, providerRegion := range invalidObjectStorageProviderRegions {
			_, err := NewObjectStorageProviderRegion(providerRegion)
			if err == nil {
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", providerRegion)
			}
		}
	})
}
