package valueObject

type HostResourceUsage struct {
	CpuPercent          float64            `json:"cpuPercent"`
	CpuPercentStr       string             `json:"cpuPercentStr"`
	MemoryPercent       float64            `json:"memoryPercent"`
	MemoryPercentStr    string             `json:"memoryPercentStr"`
	UserDataStorageInfo StorageUnitInfo    `json:"userDataStorageInfo"`
	StorageInfo         []StorageUnitInfo  `json:"storageInfo"`
	NetInfo             []NetInterfaceInfo `json:"netInfo"`
	NetInfoAggregated   NetInterfaceInfo   `json:"netInfoAggregated"`
}

func NewHostResourceUsage(
	cpuPercent float64,
	cpuPercentStr string,
	memoryPercent float64,
	memoryPercentStr string,
	userDataStorageInfo StorageUnitInfo,
	storageInfo []StorageUnitInfo,
	netInfo []NetInterfaceInfo,
	netInfoAggregated NetInterfaceInfo,
) HostResourceUsage {
	return HostResourceUsage{
		CpuPercent:          cpuPercent,
		CpuPercentStr:       cpuPercentStr,
		MemoryPercent:       memoryPercent,
		MemoryPercentStr:    memoryPercentStr,
		UserDataStorageInfo: userDataStorageInfo,
		StorageInfo:         storageInfo,
		NetInfo:             netInfo,
		NetInfoAggregated:   netInfoAggregated,
	}
}
