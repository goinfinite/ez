package valueObject

import (
	"testing"
)

func TestNewContainerImgAddress(t *testing.T) {
	t.Run("ValidContainerImgAddress", func(t *testing.T) {
		validContainerImgAddresses := []string{
			"nginx",
			"docker.io/nginx",
			"docker.io/nginx:latest",
			"docker.io/nginx:1.19.6",
			"node:18.17-1-alpine3.17",
			"docker.io/node:18.17-1-alpine3.17",
			"speedianet/os",
			"docker.io/speedianet/os",
			"http://docker.io/speedianet/os",
			"https://docker.io/speedianet/os",
			"speedianet/os:latest",
			"speedianet/os:0.0.1-alpha",
		}

		for _, path := range validContainerImgAddresses {
			_, err := NewContainerImgAddress(path)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), path)
			}
		}
	})

	t.Run("InvalidContainerImgAddress", func(t *testing.T) {
		invalidContainerImgAddresses := []string{
			"",
			"UNION SELECT * FROM USERS",
			"/path\n/path",
			"?param=value",
			"https://www.google.com",
			"/path/'; DROP TABLE users; --",
		}

		for _, path := range invalidContainerImgAddresses {
			_, err := NewContainerImgAddress(path)
			if err == nil {
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", path)
			}
		}
	})
}
