package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/speedianet/control/src/domain/valueObject/helper"
)

type PortBinding struct {
	ServiceName   ServiceName     `json:"serviceName"`
	PublicPort    NetworkPort     `json:"publicPort"`
	ContainerPort NetworkPort     `json:"containerPort"`
	Protocol      NetworkProtocol `json:"protocol"`
	PrivatePort   *NetworkPort    `json:"privatePort"`
}

type serviceBindingInfo struct {
	AltNames           []string
	PortBindings       []string
	PublicPortInterval string
}

var httpPublicPortInterval = "80"
var httpsPublicPortInterval = "443"
var databasePublicPortInterval = "30000-39999"

var ServiceNameInfoMap = map[string]serviceBindingInfo{
	"ftp": {
		AltNames:           []string{},
		PortBindings:       []string{"21"},
		PublicPortInterval: "21000-21999",
	},
	"ssh": {
		AltNames:           []string{"sftp"},
		PortBindings:       []string{"22"},
		PublicPortInterval: "22000-22999",
	},
	"telnet": {
		AltNames:     []string{},
		PortBindings: []string{"23"},
	},
	"dns": {
		AltNames:     []string{},
		PortBindings: []string{"53", "53/udp"},
	},
	"smtp": {
		AltNames:     []string{},
		PortBindings: []string{"25", "465", "587", "2525"},
	},
	"whois": {
		AltNames:     []string{},
		PortBindings: []string{"43"},
	},
	"http": {
		AltNames: []string{
			"http", "nginx", "caddy", "apache", "httpd", "php",
		},
		PortBindings: []string{"80", "8080"},
	},
	"kerberos": {
		AltNames:     []string{},
		PortBindings: []string{"88"},
	},
	"pop3": {
		AltNames:     []string{},
		PortBindings: []string{"110"},
	},
	"ntp": {
		AltNames:     []string{},
		PortBindings: []string{"123/udp"},
	},
	"imap": {
		AltNames:     []string{},
		PortBindings: []string{"143"},
	},
	"ldap": {
		AltNames:     []string{},
		PortBindings: []string{"389"},
	},
	"https": {
		AltNames: []string{
			"wss", "grpcs", "nginx", "caddy", "apache", "httpd", "php",
		},
		PortBindings: []string{"443", "8443"},
	},
	"syslog": {
		AltNames:     []string{},
		PortBindings: []string{"514/udp"},
	},
	"spamassasin": {
		AltNames:     []string{"spam-assassin"},
		PortBindings: []string{"783"},
	},
	"dot": {
		AltNames:     []string{"dns-over-tls"},
		PortBindings: []string{"853"},
	},
	"openvpn": {
		AltNames:     []string{},
		PortBindings: []string{"1194"},
	},
	"mssql": {
		AltNames:           []string{"ms-sql", "sqlserver", "sql-server"},
		PortBindings:       []string{"1433"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"oracledb": {
		AltNames:           []string{"oracle-db", "oracle"},
		PortBindings:       []string{"1521", "2483", "2484"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"mqtt": {
		AltNames:           []string{},
		PortBindings:       []string{"1883"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"nfs": {
		AltNames:     []string{},
		PortBindings: []string{"2049"},
	},
	"ghost": {
		AltNames:           []string{},
		PortBindings:       []string{"2368"},
		PublicPortInterval: httpPublicPortInterval,
	},
	"node": {
		AltNames: []string{
			"nodejs", "ruby-on-rails", "rails", "ruby",
		},
		PortBindings:       []string{"3000"},
		PublicPortInterval: httpPublicPortInterval,
	},
	"mysql": {
		AltNames:           []string{"mariadb", "percona"},
		PortBindings:       []string{"3306"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"nsq": {
		AltNames:           []string{},
		PortBindings:       []string{"4150", "4151", "4160", "4161"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"nats": {
		AltNames:           []string{},
		PortBindings:       []string{"4222", "6222", "8222"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"flask": {
		AltNames:           []string{},
		PortBindings:       []string{"5000"},
		PublicPortInterval: httpPublicPortInterval,
	},
	"kibana": {
		AltNames:           []string{},
		PortBindings:       []string{"5601"},
		PublicPortInterval: httpPublicPortInterval,
	},
	"postgresql": {
		AltNames:           []string{"postgres"},
		PortBindings:       []string{"5432"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"rabbitmq": {
		AltNames:           []string{},
		PortBindings:       []string{"5672"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"couchdb": {
		AltNames:           []string{"couch"},
		PortBindings:       []string{"5984"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"cassandra": {
		AltNames:           []string{},
		PortBindings:       []string{"9042"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"kafka": {
		AltNames:           []string{},
		PortBindings:       []string{"9092"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"redis": {
		AltNames:           []string{},
		PortBindings:       []string{"6379"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"neo4j": {
		AltNames:           []string{},
		PortBindings:       []string{"7474"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"django": {
		AltNames:           []string{},
		PortBindings:       []string{"8000"},
		PublicPortInterval: httpPublicPortInterval,
	},
	"kong": {
		AltNames:           []string{},
		PortBindings:       []string{"8001", "8444"},
		PublicPortInterval: httpPublicPortInterval,
	},
	"solr": {
		AltNames:           []string{},
		PortBindings:       []string{"8983"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"sonarqube": {
		AltNames:     []string{"sonar"},
		PortBindings: []string{"9000"},
	},
	"elasticsearch": {
		AltNames:           []string{"elastic", "elk"},
		PortBindings:       []string{"9200"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"memcached": {
		AltNames:           []string{"memcache"},
		PortBindings:       []string{"11211"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"mongodb": {
		AltNames:           []string{"mongo"},
		PortBindings:       []string{"27017"},
		PublicPortInterval: databasePublicPortInterval,
	},
	"wireguard": {
		AltNames:     []string{},
		PortBindings: []string{"51820"},
	},
}

func NewPortBinding(
	serviceName ServiceName,
	publicPort NetworkPort,
	containerPort NetworkPort,
	protocol NetworkProtocol,
	privatePort *NetworkPort,
) PortBinding {
	return PortBinding{
		ServiceName:   serviceName,
		PublicPort:    publicPort,
		ContainerPort: containerPort,
		Protocol:      protocol,
		PrivatePort:   privatePort,
	}
}

func NewPortBindingsByServiceName(serviceName ServiceName) ([]PortBinding, error) {
	portBindings := []PortBinding{}
	serviceNameStr := serviceName.String()

	serviceInfo, exists := ServiceNameInfoMap[serviceNameStr]
	if !exists {
		serviceExists := false
		for standardServiceName, bindingInfo := range ServiceNameInfoMap {
			for _, altName := range bindingInfo.AltNames {
				if serviceNameStr != altName {
					continue
				}
				serviceNameStr = standardServiceName
				serviceInfo = bindingInfo
				serviceExists = true
				break
			}
		}

		if !serviceExists {
			return portBindings, errors.New("UnknownServiceName")
		}
	}

	serviceName, err := NewServiceName(serviceNameStr)
	if err != nil {
		return portBindings, err
	}

	for _, portBindingStr := range serviceInfo.PortBindings {
		portBindingParts := strings.Split(portBindingStr, "/")
		if len(portBindingParts) == 0 {
			continue
		}

		publicPort, err := NewNetworkPort(portBindingParts[0])
		if err != nil {
			continue
		}

		containerPort := publicPort

		protocol := GuessNetworkProtocolByPort(publicPort)
		if len(portBindingParts) > 1 && protocol.String() == "tcp" {
			protocol, err = NewNetworkProtocol(portBindingParts[1])
			if err != nil {
				continue
			}
		}

		var privatePortPtr *NetworkPort

		portBinding := NewPortBinding(
			serviceName,
			publicPort,
			containerPort,
			protocol,
			privatePortPtr,
		)
		portBindings = append(portBindings, portBinding)
	}

	return portBindings, nil
}

// format: serviceName[:publicPort][:containerPort][/protocol][:privatePort]
func NewPortBindingFromString(value string) ([]PortBinding, error) {
	portBindings := []PortBinding{}

	value = strings.TrimSpace(value)
	value = strings.ToLower(value)
	portBindingRegex := `^(?P<serviceName>[a-z][\w\.\_\-]{0,128})(?::(?P<publicPort>\d{1,5}))?(?::(?P<containerPort>\d{1,5}))?(?:\/(?P<protocol>\w{1,5}))?(?::(?P<privatePort>\d{1,5}))?$`
	portBindingParts := voHelper.FindNamedGroupsMatches(portBindingRegex, string(value))

	serviceName, err := NewServiceName(portBindingParts["serviceName"])
	if err != nil {
		return portBindings, err
	}

	isServiceUnmapped := false
	servicePortBindings, err := NewPortBindingsByServiceName(serviceName)
	if err != nil {
		isServiceUnmapped = true
	}

	onlyServiceNameSent := portBindingParts["publicPort"] == "" &&
		portBindingParts["containerPort"] == ""
	if onlyServiceNameSent {
		if isServiceUnmapped {
			return portBindings, errors.New("UnknownServiceName")
		}
		return servicePortBindings, nil
	}

	rawPublicPortStr := portBindingParts["publicPort"]
	if rawPublicPortStr == "" {
		return portBindings, errors.New("UnknownPublicPort")
	}

	publicPort, err := NewNetworkPort(rawPublicPortStr)
	if err != nil {
		return portBindings, err
	}

	rawContainerPortStr := portBindingParts["containerPort"]
	if rawContainerPortStr == "" {
		if rawPublicPortStr == "0" {
			return portBindings, errors.New("UnknownContainerPort")
		}
		rawContainerPortStr = rawPublicPortStr
	}

	containerPort, err := NewNetworkPort(rawContainerPortStr)
	if err != nil {
		return portBindings, err
	}

	protocol := GuessNetworkProtocolByPort(publicPort)
	if portBindingParts["protocol"] != "" && protocol.String() == "tcp" {
		protocol, err = NewNetworkProtocol(portBindingParts["protocol"])
		if err != nil {
			return portBindings, err
		}
	}

	var privatePortPtr *NetworkPort
	if portBindingParts["privatePort"] != "" {
		privatePort, err := NewNetworkPort(portBindingParts["privatePort"])
		if err != nil {
			return portBindings, err
		}
		privatePortPtr = &privatePort
	}

	return []PortBinding{
		NewPortBinding(
			serviceName,
			publicPort,
			containerPort,
			protocol,
			privatePortPtr,
		),
	}, nil
}

func (portBinding PortBinding) GetPublicPort() NetworkPort {
	return portBinding.PublicPort
}

func (portBinding PortBinding) GetContainerPort() NetworkPort {
	return portBinding.ContainerPort
}

func (portBinding PortBinding) GetProtocol() NetworkProtocol {
	return portBinding.Protocol
}

func (portBinding PortBinding) String() string {
	portBindingStr := portBinding.ServiceName.String()

	if portBinding.PublicPort.String() != "" {
		portBindingStr += ":" + portBinding.PublicPort.String()
	}

	if portBinding.ContainerPort.String() != "" {
		portBindingStr += ":" + portBinding.ContainerPort.String()
	}

	if portBinding.Protocol.String() != "" {
		portBindingStr += "/" + portBinding.Protocol.String()
	}

	if portBinding.PrivatePort != nil {
		portBindingStr += ":" + portBinding.PrivatePort.String()
	}

	return portBindingStr
}
