package valueObject

type HostResourceUsage struct {
	CpuUsagePercent float64    `json:"cpuUsagePercent"`
	MemUsagePercent float64    `json:"memUsagePercent"`
	StorageUsage    []DiskInfo `json:"storageUsage"`
	NetUsage        NetUsage   `json:"netUsage"`
}

func NewHostResourceUsage(
	cpuUsagePercent float64,
	memUsagePercent float64,
	storageUsage []DiskInfo,
	netUsage NetUsage,
) HostResourceUsage {
	return HostResourceUsage{
		CpuUsagePercent: cpuUsagePercent,
		MemUsagePercent: memUsagePercent,
		StorageUsage:    storageUsage,
		NetUsage:        netUsage,
	}
}
