package valueObject

import (
	"testing"
)

func TestNewUnixHostname(t *testing.T) {
	t.Run("ValidUnixHostname", func(t *testing.T) {
		validUnixHostnames := []string{
			"localhost",
			"example.com",
			"sub.domain.com",
			"123-abc.com",
			"my-hostname",
			"hostname123",
			"host-name-123",
			"xn--d1acj3b",
			"xn--bcher-kva.example",
			"example.co.uk",
		}

		for _, hostname := range validUnixHostnames {
			_, err := NewUnixHostname(hostname)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), hostname)
			}
		}
	})

	t.Run("InvalidUnixHostname", func(t *testing.T) {
		invalidUnixHostnames := []string{
			"",
			"UNION SELECT * FROM USERS",
			"/path\n/path",
			"?param=value",
			"/path/'; DROP TABLE users; --",
		}

		for _, hostname := range invalidUnixHostnames {
			_, err := NewUnixHostname(hostname)
			if err == nil {
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", hostname)
			}
		}
	})
}
