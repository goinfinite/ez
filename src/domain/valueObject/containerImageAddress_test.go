package valueObject

import (
	"testing"
)

func TestNewContainerImageAddress(t *testing.T) {
	t.Run("ValidContainerImageAddress", func(t *testing.T) {
		validContainerImageAddresses := []string{
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
			"rocket.chat",
			"docker.io/rocket.chat",
			"docker.io/rocket.chat:latest",
			"docker.io/rocket.chat:3.18.0",
			"https://docker.io/rocket.chat:3.18.0-rc1",
			"localhost/backup_image:latest",
			"localhost/1000/535a6943b88a:1724884472",
		}

		for _, path := range validContainerImageAddresses {
			_, err := NewContainerImageAddress(path)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), path)
			}
		}
	})

	t.Run("InvalidContainerImageAddress", func(t *testing.T) {
		invalidContainerImageAddresses := []string{
			"",
			"UNION SELECT * FROM USERS",
			"/path\n/path",
			"?param=value",
			"/path/'; DROP TABLE users; --",
		}

		for _, path := range invalidContainerImageAddresses {
			_, err := NewContainerImageAddress(path)
			if err == nil {
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", path)
			}
		}
	})
}
