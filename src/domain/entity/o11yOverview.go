package entity

import (
	"github.com/goinfinite/ez/src/domain/valueObject"
)

type O11yOverview struct {
	Hostname          valueObject.Fqdn              `json:"hostname"`
	UptimeSecs        uint64                        `json:"uptimeSecs"`
	UptimeRelative    valueObject.RelativeTime      `json:"uptimeRelative"`
	PrivateIpAddress  valueObject.IpAddress         `json:"privateIp"`
	PublicIpAddress   valueObject.IpAddress         `json:"publicIp"`
	HardwareSpecs     valueObject.HardwareSpecs     `json:"specs"`
	HostResourceUsage valueObject.HostResourceUsage `json:"resourceUsage"`
}

func NewO11yOverview(
	hostname valueObject.Fqdn,
	uptimeSecs uint64,
	uptimeRelative valueObject.RelativeTime,
	privateIpAddress valueObject.IpAddress,
	publicIpAddress valueObject.IpAddress,
	hardwareSpecs valueObject.HardwareSpecs,
	currentResourceUsage valueObject.HostResourceUsage,
) O11yOverview {
	return O11yOverview{
		Hostname:          hostname,
		UptimeSecs:        uptimeSecs,
		UptimeRelative:    uptimeRelative,
		PrivateIpAddress:  privateIpAddress,
		PublicIpAddress:   publicIpAddress,
		HardwareSpecs:     hardwareSpecs,
		HostResourceUsage: currentResourceUsage,
	}
}
