package valueObject

type HostResourceUsage struct {
	CpuPercent    float64             `json:"cpuPercent"`
	MemoryPercent float64             `json:"memoryPercent"`
	StorageUsage  []StorageDeviceInfo `json:"storageUsage"`
	NetUsage      NetUsage            `json:"netUsage"`
}

func NewHostResourceUsage(
	cpuPercent float64,
	memoryPercent float64,
	storageUsage []StorageDeviceInfo,
	netUsage NetUsage,
) HostResourceUsage {
	return HostResourceUsage{
		CpuPercent:    cpuPercent,
		MemoryPercent: memoryPercent,
		StorageUsage:  storageUsage,
		NetUsage:      netUsage,
	}
}
