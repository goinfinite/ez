package valueObject

import (
	"testing"
)

func TestNewPortBinding(t *testing.T) {
	t.Run("ValidPortBindingFromString", func(t *testing.T) {
		validPortBindings := []string{
			"22",
			"ssh",
			"ssh:22",
			"ssh:22:22",
			"ssh:22:22/tcp",
			"ssh:0:22/tcp",
			"ssh:22:22/tcp:40000",
			"53/udp",
		}

		for _, path := range validPortBindings {
			_, err := NewPortBindingFromString(path)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), path)
			}
		}
	})

	t.Run("InvalidPortBindingFromString", func(t *testing.T) {
		invalidPortBindings := []string{
			"",
			"UNION SELECT * FROM USERS",
			"/path\n/path",
			"?param=value",
			"https://www.google.com",
			"/path/'; DROP TABLE users; --",
		}

		for _, path := range invalidPortBindings {
			_, err := NewPortBindingFromString(path)
			if err == nil {
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", path)
			}
		}
	})
}
