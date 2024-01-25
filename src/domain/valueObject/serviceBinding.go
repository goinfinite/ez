package valueObject

import (
	"errors"
	"strings"
)

type ServiceBinding string

type serviceBindingInternal struct {
	ServiceNames []string
	PortBindings []string
}

var KnownServiceBindings = []serviceBindingInternal{
	{
		ServiceNames: []string{"ftp"},
		PortBindings: []string{"21"},
	},
	{
		ServiceNames: []string{"ssh", "sftp"},
		PortBindings: []string{"22"},
	},
	{
		ServiceNames: []string{"telnet"},
		PortBindings: []string{"23"},
	},
	{
		ServiceNames: []string{"dns"},
		PortBindings: []string{"53", "53/udp"},
	},
	{
		ServiceNames: []string{"smtp"},
		PortBindings: []string{"25", "465", "587"},
	},
	{
		ServiceNames: []string{"whois"},
		PortBindings: []string{"43"},
	},
	{
		ServiceNames: []string{"http"},
		PortBindings: []string{"80", "8080", "8081"},
	},
	{
		ServiceNames: []string{"kerberos"},
		PortBindings: []string{"88"},
	},
	{
		ServiceNames: []string{"pop3"},
		PortBindings: []string{"110"},
	},
	{
		ServiceNames: []string{"ntp"},
		PortBindings: []string{"123/udp"},
	},
	{
		ServiceNames: []string{"imap"},
		PortBindings: []string{"143"},
	},
	{
		ServiceNames: []string{"ldap"},
		PortBindings: []string{"389"},
	},
	{
		ServiceNames: []string{"https"},
		PortBindings: []string{"443", "8443"},
	},
	{
		ServiceNames: []string{"syslog"},
		PortBindings: []string{"514/udp"},
	},
	{
		ServiceNames: []string{"spamassassin", "spam-assassin"},
		PortBindings: []string{"783"},
	},
	{
		ServiceNames: []string{"dot", "dns-over-tls"},
		PortBindings: []string{"853"},
	},
	{
		ServiceNames: []string{"rsync"},
		PortBindings: []string{"873"},
	},
	{
		ServiceNames: []string{"openvpn"},
		PortBindings: []string{"1194"},
	},
	{
		ServiceNames: []string{"oracledb", "oracle-db", "oracle"},
		PortBindings: []string{"1521", "2483", "2484"},
	},
	{
		ServiceNames: []string{"mqtt"},
		PortBindings: []string{"1883"},
	},
	{
		ServiceNames: []string{"nfs"},
		PortBindings: []string{"2049"},
	},
	{
		ServiceNames: []string{"ghost"},
		PortBindings: []string{"2368"},
	},
	{
		ServiceNames: []string{"node", "nodejs", "ruby-on-rails", "rails", "ruby"},
		PortBindings: []string{"3000"},
	},
	{
		ServiceNames: []string{"mysql"},
		PortBindings: []string{"3306"},
	},
	{
		ServiceNames: []string{"nats"},
		PortBindings: []string{"4222", "6222", "8222"},
	},
	{
		ServiceNames: []string{"flask"},
		PortBindings: []string{"5000"},
	},
	{
		ServiceNames: []string{"kibana"},
		PortBindings: []string{"5601"},
	},
	{
		ServiceNames: []string{"postgres", "postgresql"},
		PortBindings: []string{"5432"},
	},
	{
		ServiceNames: []string{"rabbitmq"},
		PortBindings: []string{"5672"},
	},
	{
		ServiceNames: []string{"couchdb"},
		PortBindings: []string{"5984"},
	},
	{
		ServiceNames: []string{"cassandra"},
		PortBindings: []string{"9042"},
	},
	{
		ServiceNames: []string{"redis"},
		PortBindings: []string{"6379"},
	},
	{
		ServiceNames: []string{"neo4j"},
		PortBindings: []string{"7474"},
	},
	{
		ServiceNames: []string{"django"},
		PortBindings: []string{"8000"},
	},
	{
		ServiceNames: []string{"kong"},
		PortBindings: []string{"8001", "8444"},
	},
	{
		ServiceNames: []string{"solr"},
		PortBindings: []string{"8983"},
	},
	{
		ServiceNames: []string{"sonarqube"},
		PortBindings: []string{"9000"},
	},
	{
		ServiceNames: []string{"elastic", "elasticsearch", "elk"},
		PortBindings: []string{"9200"},
	},
	{
		ServiceNames: []string{"memcached", "memcache"},
		PortBindings: []string{"11211"},
	},
	{
		ServiceNames: []string{"mongodb", "mongo"},
		PortBindings: []string{"27017"},
	},
	{
		ServiceNames: []string{"wireguard"},
		PortBindings: []string{"51820"},
	},
}

func NewServiceBinding(value string) (ServiceBinding, error) {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)

	for _, bindingInfo := range KnownServiceBindings {
		standardName := bindingInfo.ServiceNames[0]
		if standardName == value {
			return ServiceBinding(standardName), nil
		}

		for _, serviceName := range bindingInfo.ServiceNames {
			if serviceName != value {
				continue
			}
			return ServiceBinding(serviceName), nil
		}
	}

	return "", errors.New("UnknownServiceBinding")
}

func NewServiceBindingPanic(value string) ServiceBinding {
	serviceBinding, err := NewServiceBinding(value)
	if err != nil {
		panic(err)
	}
	return serviceBinding
}

func NewServiceBindingByPort(port NetworkPort) (ServiceBinding, error) {
	for _, bindingInfo := range KnownServiceBindings {
		standardName := bindingInfo.ServiceNames[0]
		standardPort := bindingInfo.PortBindings[0]
		if standardPort == port.String() {
			return ServiceBinding(standardName), nil
		}

		for _, portBinding := range bindingInfo.PortBindings {
			portStr := strings.Split(portBinding, "/")[0]
			if portStr != port.String() {
				continue
			}
			return ServiceBinding(standardName), nil
		}
	}

	return "", errors.New("UnknownServiceBinding")
}

func (serviceBinding ServiceBinding) String() string {
	return string(serviceBinding)
}

func (serviceBinding ServiceBinding) GetAsPortBindings() ([]PortBinding, error) {
	portBindings := []PortBinding{}
	serviceBindingStr := serviceBinding.String()

	for _, bindingInfo := range KnownServiceBindings {
		standardName := bindingInfo.ServiceNames[0]
		if standardName != serviceBindingStr {
			continue
		}

		for _, portBinding := range bindingInfo.PortBindings {
			portBinding, err := NewPortBindingFromString(portBinding)
			if err != nil {
				continue
			}
			portBindings = append(portBindings, portBinding)
		}
	}

	if len(portBindings) == 0 {
		return portBindings, errors.New("UnknownServiceBinding")
	}

	return portBindings, nil
}
