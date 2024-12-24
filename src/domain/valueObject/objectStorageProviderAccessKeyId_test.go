package valueObject

import (
	"testing"
)

func TestNewObjectStorageProviderAccessKeyId(t *testing.T) {
	t.Run("ValidObjectStorageProviderAccessKeyId", func(t *testing.T) {
		validObjectStorageProviderAccessKeyIds := []string{
			"o0w157XqXTvGH1ZoRiTI",
			"0GIVRXLGvTs60L8tGzbD",
			"0OAkLj6VOU1Mpn4QP3J1",
			"12YYWGSY4C4E4ACB25CS",
			"5OWHGS3AR6UJETNVINO5",
			"AKIA8L37ZFRNZXQYHT55",
			"AKIAWA3DFKDJZNJAYM6Q",
			"f633r17391cf5227115809c8580d4eb041679f77",
			"e3b0c44298fc1c149afbf4c8996fb9255",
		}

		for _, keyId := range validObjectStorageProviderAccessKeyIds {
			_, err := NewObjectStorageProviderAccessKeyId(keyId)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), keyId)
			}
		}
	})

	t.Run("InvalidObjectStorageProviderAccessKeyId", func(t *testing.T) {
		invalidObjectStorageProviderAccessKeyIds := []string{
			"",
			"0GIV@XLGvTs#0L8t=zbD",
			"0GIV<alert>XLGvTs60L8tGzbD",
			"UNION SELECT * FROM USERS",
			"/keyId\n/keyId",
			"?param=value",
			"/keyId/'; DROP TABLE users; --",
		}

		for _, keyId := range invalidObjectStorageProviderAccessKeyIds {
			_, err := NewObjectStorageProviderAccessKeyId(keyId)
			if err == nil {
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", keyId)
			}
		}
	})
}
