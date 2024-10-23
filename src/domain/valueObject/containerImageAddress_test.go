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
			"goinfinite/os",
			"docker.io/goinfinite/os",
			"http://docker.io/goinfinite/os",
			"https://docker.io/goinfinite/os",
			"goinfinite/os:latest",
			"goinfinite/os:0.0.1-alpha",
			"rocket.chat",
			"docker.io/rocket.chat",
			"docker.io/rocket.chat:latest",
			"docker.io/rocket.chat:3.18.0",
			"https://docker.io/rocket.chat:3.18.0-rc1",
			"localhost/backup_image:latest",
			"localhost/1000/535a6943b88a:1724884472",
		}

		for _, addr := range validContainerImageAddresses {
			_, err := NewContainerImageAddress(addr)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), addr)
			}
		}
	})

	t.Run("InvalidContainerImageAddress", func(t *testing.T) {
		invalidContainerImageAddresses := []string{
			"",
			"UNION SELECT * FROM USERS",
			"/addr\n/addr",
			"?param=value",
			"/addr/'; DROP TABLE users; --",
		}

		for _, addr := range invalidContainerImageAddresses {
			_, err := NewContainerImageAddress(addr)
			if err == nil {
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", addr)
			}
		}
	})
}
