package valueObject

import (
	"testing"
)

func TestNewServiceBinding(t *testing.T) {
	t.Run("ValidServiceBinding", func(t *testing.T) {
		validServiceBindings := []string{
			"ssh",
			"dns",
			"spam-assassin",
			"dns-over-tls",
			"mysql",
			"wireguard",
		}

		for _, binding := range validServiceBindings {
			_, err := NewServiceBinding(binding)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), binding)
			}
		}
	})

	t.Run("ValidServiceBindingByPort", func(t *testing.T) {
		validServiceBindings := []NetworkPort{
			22,
			53,
			783,
			853,
			3306,
			51820,
		}

		for _, binding := range validServiceBindings {
			_, err := NewServiceBindingByPort(binding)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s [%s]", err.Error(), binding)
			}
		}
	})

	t.Run("InvalidServiceBinding", func(t *testing.T) {
		invalidServiceBindings := []string{
			"55/tcp",
			"UNION SELECT * FROM USERS",
			"/bindings\n/bindings",
			"?param=value",
			"https://www.google.com",
			"/bindings/'; DROP TABLE users; --",
		}

		for _, binding := range invalidServiceBindings {
			_, err := NewServiceBinding(binding)
			if err == nil {
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", binding)
			}
		}
	})

	t.Run("GetServiceAsPortBinding", func(t *testing.T) {
		publicPort, _ := NewNetworkPort(21)
		containerPort := publicPort
		protocol, _ := NewNetworkProtocol("tcp")

		portBinding := NewPortBinding(publicPort, containerPort, protocol, nil)
		portBindings := []PortBinding{portBinding}

		serviceBinding, _ := NewServiceBinding("ftp")
		serviceBindingPortBindings, err := serviceBinding.GetAsPortBindings()
		if err != nil {
			t.Errorf("ExpectingNoErrorButGot: %s", err.Error())
			return
		}

		if len(serviceBindingPortBindings) != len(portBindings) {
			t.Errorf(
				"ExpectingSameLengthButGot: %d %d",
				len(serviceBindingPortBindings),
				len(portBindings),
			)
			return
		}
	})
}
