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
	ServiceNames       []string
	PortBindings       []string
	PublicPortInterval string
}

var httpPublicPortInterval = "80"
var httpsPublicPortInterval = "443"
var databasePublicPortInterval = "30000-39999"

var KnownServiceBindings = []serviceBindingInfo{
	{
		ServiceNames:       []string{"ftp"},
		PortBindings:       []string{"21"},
		PublicPortInterval: "21000-21999",
	},
	{
		ServiceNames:       []string{"ssh", "sftp"},
		PortBindings:       []string{"22"},
		PublicPortInterval: "22000-22999",
	},
	{
		ServiceNames: []string{"telnet"},
		PortBindings: []string{"23"},
	},
	{
		ServiceNames: []string{"smtp"},
		PortBindings: []string{"25"},
	},
	{
		ServiceNames: []string{"whois"},
		PortBindings: []string{"43"},
	},
	{
		ServiceNames: []string{"dns"},
		PortBindings: []string{"53", "53/udp"},
	},
	{
		ServiceNames: []string{
			"http", "ws", "grpc", "php", "nginx", "apache", "caddy", "traefik",
		},
		PortBindings: []string{"80"},
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
			"https", "wss", "grpcs", "php",
		},
		PortBindings: []string{"443"},
	},
	{
		ServiceNames: []string{"smtps"},
		PortBindings: []string{"465"},
	},
	{
		ServiceNames: []string{"syslog"},
		PortBindings: []string{"514/udp"},
	},
	{
		ServiceNames: []string{"smtp-tls"},
		PortBindings: []string{"587"},
	},
	{
		ServiceNames: []string{"spamassasin", "spam-assassin"},
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
		ServiceNames: []string{
			"mssql", "ms-sql", "sqlserver", "sql-server",
		},
		PortBindings:       []string{"1433"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames: []string{
			"oracledb", "oracle-db", "oracle",
		},
		PortBindings:       []string{"1521", "2483", "2484"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames: []string{
			"sos", "speediaos", "speedia-os", "speedia",
		},
		PortBindings:       []string{"1618"},
		PublicPortInterval: httpsPublicPortInterval,
	},
	{
		ServiceNames:       []string{"mqtt"},
		PortBindings:       []string{"1883"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames: []string{"nfs"},
		PortBindings: []string{"2049"},
	},
	{
		ServiceNames:       []string{"orientdb"},
		PortBindings:       []string{"2424", "2480"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"ghost"},
		PortBindings:       []string{"2368"},
		PublicPortInterval: httpPublicPortInterval,
	},
	{
		ServiceNames: []string{"smtp-tls-alt"},
		PortBindings: []string{"2525"},
	},
	{
		ServiceNames: []string{"custom", "node", "nodejs", "ruby-on-rails", "rails", "ruby", "aerospike"},
		PortBindings: []string{
			"3000", "3001", "3002", "3003", "3004", "3005", "3006", "3007", "3008", "3009"},
		PublicPortInterval: httpPublicPortInterval,
	},
	{
		ServiceNames:       []string{"mysql", "mariadb", "percona"},
		PortBindings:       []string{"3306"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"nsq"},
		PortBindings:       []string{"4150", "4151", "4160", "4161"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"cratedb", "crate"},
		PortBindings:       []string{"4200", "4300"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"nats"},
		PortBindings:       []string{"4222", "6222", "8222"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames: []string{
			"flask", "distribution",
		},
		PortBindings:       []string{"5000"},
		PublicPortInterval: httpPublicPortInterval,
	},
	{
		ServiceNames:       []string{"kibana"},
		PortBindings:       []string{"5601"},
		PublicPortInterval: httpPublicPortInterval,
	},
	{
		ServiceNames:       []string{"postgresql", "postgres", "cratedb", "crate"},
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
		ServiceNames:       []string{"flink"},
		PortBindings:       []string{"6123", "8081"},
		PublicPortInterval: httpPublicPortInterval,
	},
	{
		ServiceNames:       []string{"redis"},
		PortBindings:       []string{"6379"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames: []string{"storm"},
		PortBindings: []string{"6627"},
	},
	{
		ServiceNames:       []string{"meilisearch", "meili"},
		PortBindings:       []string{"7700"},
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
		PortBindings:       []string{"8000", "8001", "8002"},
		PublicPortInterval: httpPublicPortInterval,
	},
	{
		ServiceNames:       []string{"mattermost"},
		PortBindings:       []string{"8065"},
		PublicPortInterval: httpPublicPortInterval,
	},
	{
		ServiceNames:       []string{"odoo"},
		PortBindings:       []string{"8069", "8072"},
		PublicPortInterval: httpPublicPortInterval,
	},
	{
		ServiceNames: []string{
			"http", "php",
		},
		PortBindings: []string{"8080"},
	},
	{
		ServiceNames:       []string{"influxdb", "influx"},
		PortBindings:       []string{"8086"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames: []string{"couchbase"},
		PortBindings: []string{
			"8091",
			"8092",
			"8093",
			"8094",
			"8095",
			"8096",
			"8097",
			"9123",
			"11207",
			"11210",
			"11280",
			"18091",
			"18092",
			"18093",
			"18094",
			"18095",
			"18096",
			"18097",
		},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"clickhouse"},
		PortBindings:       []string{"8123", "9000", "9009"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames: []string{
			"https",
			"wss",
			"grpcs",
			"php",
		},
		PortBindings: []string{"8443"},
	},
	{
		ServiceNames: []string{"kong-secure"},
		PortBindings: []string{
			"8444", "8445",
		},
		PublicPortInterval: httpsPublicPortInterval,
	},
	{
		ServiceNames:       []string{"arangodb", "arango"},
		PortBindings:       []string{"8529"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"chronograf"},
		PortBindings:       []string{"8888"},
		PublicPortInterval: httpPublicPortInterval,
	},
	{
		ServiceNames:       []string{"solr"},
		PortBindings:       []string{"8983"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames: []string{"sonarqube", "sonar", "sentry"},
		PortBindings: []string{"9000"},
	},
	{
		ServiceNames:       []string{"cassandra"},
		PortBindings:       []string{"9042"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"kafka", "kapacitor"},
		PortBindings:       []string{"9092"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"elasticsearch", "elastic", "elk"},
		PortBindings:       []string{"9200"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames:       []string{"logstash"},
		PortBindings:       []string{"9600"},
		PublicPortInterval: httpPublicPortInterval,
	},
	{
		ServiceNames: []string{"teamspeak", "ts3"},
		PortBindings: []string{"9987/udp", "10011", "30033"},
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
		ServiceNames:       []string{"rethinkdb", "rethink"},
		PortBindings:       []string{"28015", "29015", "8080"},
		PublicPortInterval: databasePublicPortInterval,
	},
	{
		ServiceNames: []string{"jenkins"},
		PortBindings: []string{"50000", "8080"},
	},
	{
		ServiceNames: []string{"wireguard"},
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

func findKnownServiceBindingByName(
	receivedServiceName ServiceName,
) (serviceBindingInfo, error) {
	var serviceBinding serviceBindingInfo
	receivedServiceNameStr := receivedServiceName.String()

	desiredServiceBindingIndex := -1

	for bindingIndex, bindingInfo := range KnownServiceBindings {
		standardName := bindingInfo.ServiceNames[0]
		if receivedServiceNameStr == standardName {
			desiredServiceBindingIndex = bindingIndex
			break
		}

		for _, altName := range bindingInfo.ServiceNames {
			if receivedServiceNameStr != altName {
				continue
			}
			desiredServiceBindingIndex = bindingIndex
			break
		}
	}

	if desiredServiceBindingIndex == -1 {
		return serviceBinding, errors.New("UnknownServiceName")
	}

	return KnownServiceBindings[desiredServiceBindingIndex], nil
}

func knownServiceBindingsPortBindingFactory(
	serviceName ServiceName,
	portBindingStr string,
) (PortBinding, error) {
	var portBinding PortBinding

	portBindingParts := strings.Split(portBindingStr, "/")
	if len(portBindingParts) == 0 {
		return portBinding, errors.New("InvalidPortBindingStructure")
	}

	publicPort, err := NewNetworkPort(portBindingParts[0])
	if err != nil {
		return portBinding, err
	}

	containerPort := publicPort

	likelyProtocol := GuessNetworkProtocolByPort(publicPort)
	isLikelyProtocolGeneric := likelyProtocol.String() == "tcp"
	protocol := likelyProtocol
	if len(portBindingParts) > 1 && isLikelyProtocolGeneric {
		protocol, err = NewNetworkProtocol(portBindingParts[1])
		if err != nil {
			return portBinding, err
		}
	}

	var privatePortPtr *NetworkPort

	return NewPortBinding(
		serviceName,
		publicPort,
		containerPort,
		protocol,
		privatePortPtr,
	), nil
}

func NewPortBindingsByServiceName(
	receivedServiceName ServiceName,
) ([]PortBinding, error) {
	portBindings := []PortBinding{}

	desiredServiceBinding, err := findKnownServiceBindingByName(receivedServiceName)
	if err != nil {
		return portBindings, err
	}

	for _, portBindingStr := range desiredServiceBinding.PortBindings {
		portBinding, err := knownServiceBindingsPortBindingFactory(
			receivedServiceName,
			portBindingStr,
		)
		if err != nil {
			return portBindings, err
		}

		portBindings = append(portBindings, portBinding)
	}

	return portBindings, nil
}

func NewPortBindingByPort(port NetworkPort) (PortBinding, error) {
	portStr := port.String()
	portBinding := PortBinding{}

	for _, serviceBinding := range KnownServiceBindings {
		for _, portBindingStr := range serviceBinding.PortBindings {
			if !strings.Contains(portBindingStr, portStr) {
				continue
			}

			serviceName, err := NewServiceName(serviceBinding.ServiceNames[0])
			if err != nil {
				continue
			}

			portBinding, err = knownServiceBindingsPortBindingFactory(
				serviceName,
				portBindingStr,
			)
			if err != nil {
				continue
			}

			return portBinding, nil
		}
	}

	return portBinding, errors.New("UnknownPort")
}

// format: [serviceName][:publicPort][:containerPort][/protocol][:privatePort]
func NewPortBindingFromString(value string) ([]PortBinding, error) {
	portBindings := []PortBinding{}

	value = strings.TrimSpace(value)
	value = strings.ToLower(value)
	portBindingRegex := `^(?:(?P<serviceName>[a-z][\w\.\_\ \-]{0,128}[a-z0-9]))?(?::?(?P<publicPort>\d{1,5}))?(?::(?P<containerPort>\d{1,5}))?(?:\/(?P<protocol>\w{1,5}))?(?::(?P<privatePort>\d{1,5}))?$`
	portBindingParts := voHelper.FindNamedGroupsMatches(portBindingRegex, string(value))

	serviceNameSent := portBindingParts["serviceName"] != ""
	publicPortSent := portBindingParts["publicPort"] != ""
	protocolSent := portBindingParts["protocol"] != ""
	nothingSent := !serviceNameSent && !publicPortSent
	if nothingSent {
		return portBindings, errors.New("ServiceNameOrPortRequired")
	}

	var err error
	serviceName, _ := NewServiceName("unmapped")
	if serviceNameSent {
		serviceName, err = NewServiceName(portBindingParts["serviceName"])
		if err != nil {
			serviceName, _ = NewServiceName("unknown")
		}
	}

	if serviceNameSent && !publicPortSent {
		servicePortBindings, err := NewPortBindingsByServiceName(serviceName)
		if err != nil {
			return portBindings, err
		}

		return servicePortBindings, nil
	}

	publicPort, err := NewNetworkPort(portBindingParts["publicPort"])
	if err != nil {
		return portBindings, err
	}

	isKnownPublicPort := true
	likelyPortBinding, err := NewPortBindingByPort(publicPort)
	if err != nil {
		isKnownPublicPort = false
	}

	if isKnownPublicPort && !serviceNameSent {
		if !protocolSent {
			return []PortBinding{likelyPortBinding}, nil
		}
		serviceName = likelyPortBinding.ServiceName
	}

	likelyProtocol := GuessNetworkProtocolByPort(publicPort)
	isLikelyProtocolGeneric := likelyProtocol.String() == "tcp"
	protocol := likelyProtocol
	if protocolSent && isLikelyProtocolGeneric {
		protocol, err = NewNetworkProtocol(portBindingParts["protocol"])
		if err != nil {
			return portBindings, err
		}
	}

	rawContainerPortStr := portBindingParts["containerPort"]
	if rawContainerPortStr == "" {
		if publicPort.Get() == 0 {
			return portBindings, errors.New("UnknownContainerPort")
		}
		rawContainerPortStr = publicPort.String()
	}

	containerPort, err := NewNetworkPort(rawContainerPortStr)
	if err != nil {
		return portBindings, err
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

func (portBinding PortBinding) GetPublicPort() NetworkPort {
	return portBinding.PublicPort
}

func (portBinding PortBinding) GetContainerPort() NetworkPort {
	return portBinding.ContainerPort
}

func (portBinding PortBinding) GetProtocol() NetworkProtocol {
	return portBinding.Protocol
}

func (portBinding PortBinding) GetPublicPortInterval() (NetworkPortInterval, error) {
	var portInterval NetworkPortInterval

	serviceInfo, err := findKnownServiceBindingByName(portBinding.ServiceName)
	if err != nil {
		return portInterval, err
	}

	if serviceInfo.PublicPortInterval == "" {
		return portInterval, errors.New("UnknownPublicPortInterval")
	}

	intervalParts := strings.Split(serviceInfo.PublicPortInterval, "-")
	if len(intervalParts) <= 1 {
		preDefinedPublicPort, err := NewNetworkPort(serviceInfo.PublicPortInterval)
		if err != nil {
			return portInterval, err
		}

		return NewNetworkPortInterval(preDefinedPublicPort, nil)
	}

	minPublicPortStr := intervalParts[0]
	minPublicPort, err := NewNetworkPort(minPublicPortStr)
	if err != nil {
		return portInterval, err
	}

	maxPublicPortStr := intervalParts[1]
	maxPublicPort, err := NewNetworkPort(maxPublicPortStr)
	if err != nil {
		return portInterval, err
	}

	return NewNetworkPortInterval(minPublicPort, &maxPublicPort)
}
