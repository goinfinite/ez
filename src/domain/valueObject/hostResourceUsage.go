package valueObject

type HostResourceUsage struct {
	CpuPercent    float64            `json:"cpuPercent"`
	MemoryPercent float64            `json:"memoryPercent"`
	StorageInfo   []StorageUnitInfo  `json:"storageInfo"`
	NetInfo       []NetInterfaceInfo `json:"netInfo"`
}

func NewHostResourceUsage(
	cpuPercent float64,
	memoryPercent float64,
	storageInfo []StorageUnitInfo,
	netInfo []NetInterfaceInfo,
) HostResourceUsage {
	return HostResourceUsage{
		CpuPercent:    cpuPercent,
		MemoryPercent: memoryPercent,
		StorageInfo:   storageInfo,
		NetInfo:       netInfo,
	}
}
