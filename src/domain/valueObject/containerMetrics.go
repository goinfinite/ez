package valueObject

type ContainerMetrics struct {
	CurrentCpuPercent  float64     `json:"currentCpuPercent"`
	AverageCpuPercent  float64     `json:"avgCpuPercent"`
	MemoryBytes        Byte        `json:"memoryBytes"`
	MemoryPercent      float64     `json:"memoryPercent"`
	StorageInputBytes  Byte        `json:"storageInputBytes"`
	StorageOutputBytes Byte        `json:"storageOutputBytes"`
	StorageSpaceBytes  Byte        `json:"storageSpaceBytes"`
	StorageInodesCount InodesCount `json:"storageInodesCount"`
	NetInputBytes      Byte        `json:"netInputBytes"`
	NetOutputBytes     Byte        `json:"netOutputBytes"`
}

func NewContainerMetrics(
	currentCpuPercent float64,
	averageCpuPercent float64,
	memoryBytes Byte,
	memoryPercent float64,
	storageInputBytes Byte,
	storageOutputBytes Byte,
	storageSpaceBytes Byte,
	storageInodesCount InodesCount,
	netInputBytes Byte,
	netOutputBytes Byte,
) ContainerMetrics {
	return ContainerMetrics{
		CurrentCpuPercent:  currentCpuPercent,
		AverageCpuPercent:  averageCpuPercent,
		MemoryBytes:        memoryBytes,
		MemoryPercent:      memoryPercent,
		StorageInputBytes:  storageInputBytes,
		StorageOutputBytes: storageOutputBytes,
		StorageSpaceBytes:  storageSpaceBytes,
		StorageInodesCount: storageInodesCount,
		NetInputBytes:      netInputBytes,
		NetOutputBytes:     netOutputBytes,
	}
}

func NewContainerMetricsWithBlankValues() ContainerMetrics {
	return ContainerMetrics{
		CurrentCpuPercent:  0,
		AverageCpuPercent:  0,
		MemoryBytes:        0,
		MemoryPercent:      0,
		StorageInputBytes:  0,
		StorageOutputBytes: 0,
		StorageSpaceBytes:  0,
		StorageInodesCount: 0,
		NetInputBytes:      0,
		NetOutputBytes:     0,
	}
}
