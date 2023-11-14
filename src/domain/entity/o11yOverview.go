package entity

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type O11yOverview struct {
	Hostname          valueObject.Fqdn              `json:"hostname"`
	UptimeSecs        uint64                        `json:"uptimeSecs"`
	PublicIpAddress   valueObject.IpAddress         `json:"publicIp"`
	HardwareSpecs     valueObject.HardwareSpecs     `json:"specs"`
	HostResourceUsage valueObject.HostResourceUsage `json:"resourceUsage"`
}

func NewO11yOverview(
	hostname valueObject.Fqdn,
	uptime uint64,
	publicIpAddress valueObject.IpAddress,
	hardwareSpecs valueObject.HardwareSpecs,
	currentResourceUsage valueObject.HostResourceUsage,
) O11yOverview {
	return O11yOverview{
		Hostname:          hostname,
		UptimeSecs:        uptime,
		PublicIpAddress:   publicIpAddress,
		HardwareSpecs:     hardwareSpecs,
		HostResourceUsage: currentResourceUsage,
	}
}
