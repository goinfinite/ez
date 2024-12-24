package valueObject

import (
	"testing"
)

func TestNewObjectStorageProviderSecretAccessKey(t *testing.T) {
	t.Run("ValidObjectStorageProviderSecretAccessKey", func(t *testing.T) {
		validObjectStorageProviderSecretAccessKeys := []string{
			"UKW0qMAv98AzFeBV96B9S5f8W9uHADSjY6DuAuIS",
			"dAtuFYH019L5qHufLrFqbwiEfqRpVwjO14DnMvZe",
			"q5VY543kymgWvEj27FKIfcGHc5DROltVuFZZGyXk",
			"B9JE59gxErcQcW1S9dHzD1+KqtFvQSKgpEs7Ki69",
			"SsB5J6sMYUJaqGXqeYptHXJScf0mMzjgR01gjoT1LJ8=",
		}

		for _, secretKey := range validObjectStorageProviderSecretAccessKeys {
			_, err := NewObjectStorageProviderSecretAccessKey(secretKey)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), secretKey)
			}
		}
	})

	t.Run("InvalidObjectStorageProviderSecretAccessKey", func(t *testing.T) {
		invalidObjectStorageProviderSecretAccessKeys := []string{
			"",
			"0GIV@XLGvTs#0L8t=zbD",
			"0GIV<alert>XLGvTs60L8tGzbD",
			"UNION SELECT * FROM USERS",
			"/secretKey\n/secretKey",
			"?param=value",
			"/secretKey/'; DROP TABLE users; --",
		}

		for _, secretKey := range invalidObjectStorageProviderSecretAccessKeys {
			_, err := NewObjectStorageProviderSecretAccessKey(secretKey)
			if err == nil {
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", secretKey)
			}
		}
	})
}
