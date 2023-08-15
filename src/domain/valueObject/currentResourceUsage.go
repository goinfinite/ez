package valueObject

type CurrentResourceUsage struct {
	CpuUsagePercent float64    `json:"cpuUsagePercent"`
	MemUsagePercent float64    `json:"memUsagePercent"`
	StorageUsage    []DiskInfo `json:"storageUsage"`
	NetUsage        NetUsage   `json:"netUsage"`
}

func NewCurrentResourceUsage(
	cpuUsagePercent float64,
	memUsagePercent float64,
	storageUsage []DiskInfo,
	netUsage NetUsage,
) CurrentResourceUsage {
	return CurrentResourceUsage{
		CpuUsagePercent: cpuUsagePercent,
		MemUsagePercent: memUsagePercent,
		StorageUsage:    storageUsage,
		NetUsage:        netUsage,
	}
}
