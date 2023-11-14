package valueObject

type HostResourceUsage struct {
	CpuPercent    float64             `json:"cpuPercent"`
	MemoryPercent float64             `json:"memoryPercent"`
	StorageInfo   []StorageDeviceInfo `json:"storageInfo"`
	NetInfo       []NetInterfaceInfo  `json:"netInfo"`
}

func NewHostResourceUsage(
	cpuPercent float64,
	memoryPercent float64,
	storageInfo []StorageDeviceInfo,
	netInfo []NetInterfaceInfo,
) HostResourceUsage {
	return HostResourceUsage{
		CpuPercent:    cpuPercent,
		MemoryPercent: memoryPercent,
		StorageInfo:   storageInfo,
		NetInfo:       netInfo,
	}
}
