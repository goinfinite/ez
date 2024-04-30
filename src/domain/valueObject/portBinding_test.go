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
			"1618/tcp",
			"unknown:12345",
			"unknown:12345:12345",
			"unknown:12345:12345/tcp",
			"unknown:12345:12345/tcp:40000",
			"dns-alt-name:53/udp",
			"dns alt name:53/udp",
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
			"uknown",
		}

		for _, path := range invalidPortBindings {
			_, err := NewPortBindingFromString(path)
			if err == nil {
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", path)
			}
		}
	})

	t.Run("GetPortBindingWithServiceName", func(t *testing.T) {
		portBindings, _ := NewPortBindingFromString("ssh")
		publicPort := portBindings[0].GetPublicPort()
		if publicPort.String() != "22" {
			t.Errorf(
				"GotWrongPublicPort: %s, Expected: %s", publicPort.String(), "22",
			)
		}
	})

	t.Run("GetPortBindingWithPublicPort", func(t *testing.T) {
		portBindings, _ := NewPortBindingFromString("22")
		serviceNameStr := portBindings[0].ServiceName.String()

		if serviceNameStr != "ssh" {
			t.Errorf(
				"GotWrongServiceName: %s, Expected: %s", serviceNameStr, "ssh",
			)
		}
	})

	t.Run("CheckIfPublicPortTakesPrecedenceOverServiceName", func(t *testing.T) {
		portBindings, _ := NewPortBindingFromString("mongodb:1000:1000/tcp")
		publicPortStr := portBindings[0].GetPublicPort().String()
		expectedPublicPortStr := "1000"
		if publicPortStr != expectedPublicPortStr {
			t.Errorf(
				"GotWrongPublicPort: %s, Expected: %s",
				publicPortStr, expectedPublicPortStr,
			)
		}
	})

	t.Run("CheckIfPublicPortTakesPrecedenceOverUnknownServiceName", func(t *testing.T) {
		portBindings, _ := NewPortBindingFromString("unknown-service:22")
		publicPortStr := portBindings[0].GetPublicPort().String()
		expectedPublicPortStr := "22"
		if publicPortStr != expectedPublicPortStr {
			t.Errorf(
				"GotWrongPublicPort: %s, Expected: %s",
				publicPortStr, expectedPublicPortStr,
			)
		}
	})
}
