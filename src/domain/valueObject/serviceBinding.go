package valueObject

import (
	"errors"
	"strings"
)

type ServiceBinding string

type serviceBindingInternal struct {
	ServiceNames       []string
	PortBindings       []string
	PublicPortInterval []string
}

var httpPublicPortInterval = []string{"80", "443"}
var databasePublicPortInterval = []string{"30000-39999"}

var KnownServiceBindings = []serviceBindingInternal{
	{
		ServiceNames:       []string{"ftp"},
		PortBindings:       []string{"21"},
		PublicPortInterval: []string{"21000-21999"},
	},
	{
		ServiceNames:       []string{"ssh", "sftp"},
		PortBindings:       []string{"22"},
		PublicPortInterval: []string{"22000-22999"},
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
		ServiceNames: []string{
			"http", "nginx", "caddy", "apache", "httpd", "php",
		},
		PortBindings: []string{"80", "8080"},
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
		ServiceNames: []string{
			"https", "wss", "grpcs", "nginx", "caddy", "apache", "httpd", "php",
		},
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
		ServiceNames: []string{"openvpn"},
		PortBindings: []string{"1194"},
	},
	{
		ServiceNames:       []string{"ms-sql", "mssql", "sqlserver", "sql-server"},
		PortBindings:       []string{"1433"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"oracledb", "oracle-db", "oracle"},
		PortBindings:       []string{"1521", "2483", "2484"},
		PublicPortInterval: databasePublicPortInterval,
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
		ServiceNames:       []string{"ghost"},
		PortBindings:       []string{"2368"},
		PublicPortInterval: httpPublicPortInterval,
	},
	{
		ServiceNames: []string{
			"node", "nodejs", "ruby-on-rails", "rails", "ruby",
		},
		PortBindings:       []string{"3000"},
		PublicPortInterval: httpPublicPortInterval,
	},
	{
		ServiceNames:       []string{"mysql"},
		PortBindings:       []string{"3306"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"nsq"},
		PortBindings:       []string{"4150", "4151", "4160", "4161"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"nats"},
		PortBindings:       []string{"4222", "6222", "8222"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"flask"},
		PortBindings:       []string{"5000"},
		PublicPortInterval: httpPublicPortInterval,
	},
	{
		ServiceNames:       []string{"kibana"},
		PortBindings:       []string{"5601"},
		PublicPortInterval: httpPublicPortInterval,
	},
	{
		ServiceNames:       []string{"postgres", "postgresql"},
		PortBindings:       []string{"5432"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"rabbitmq"},
		PortBindings:       []string{"5672"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"couchdb"},
		PortBindings:       []string{"5984"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"cassandra"},
		PortBindings:       []string{"9042"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"kafka"},
		PortBindings:       []string{"9092"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"redis"},
		PortBindings:       []string{"6379"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"neo4j"},
		PortBindings:       []string{"7474"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"django"},
		PortBindings:       []string{"8000"},
		PublicPortInterval: httpPublicPortInterval,
	},
	{
		ServiceNames:       []string{"kong"},
		PortBindings:       []string{"8001", "8444"},
		PublicPortInterval: httpPublicPortInterval,
	},
	{
		ServiceNames:       []string{"solr"},
		PortBindings:       []string{"8983"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames: []string{"sonarqube"},
		PortBindings: []string{"9000"},
	},
	{
		ServiceNames:       []string{"elastic", "elasticsearch", "elk"},
		PortBindings:       []string{"9200"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"memcached", "memcache"},
		PortBindings:       []string{"11211"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"mongodb", "mongo"},
		PortBindings:       []string{"27017"},
		PublicPortInterval: databasePublicPortInterval,
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
