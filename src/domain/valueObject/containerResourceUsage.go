package valueObject

type ContainerResourceUsage struct {
	CurrentCpuPercent   float64 `json:"currentCpuPercent"`
	AverageCpuPercent   float64 `json:"avgCpuPercent"`
	MemoryBytes         uint64  `json:"memoryBytes"`
	MemoryPercent       float64 `json:"memoryPercent"`
	StorageInputBytes   uint64  `json:"storageInputBytes"`
	StorageOutputBytes  uint64  `json:"storageOutputBytes"`
	StorageSpaceBytes   uint64  `json:"storageSpaceBytes"`
	StorageSpacePercent float64 `json:"storageSpacePercent"`
	StorageInodesCount  uint64  `json:"storageInodesCount"`
	NetInputBytes       uint64  `json:"netInputBytes"`
	NetOutputBytes      uint64  `json:"netOutputBytes"`
}

func NewContainerResourceUsage(
	currentCpuPercent float64,
	averageCpuPercent float64,
	memoryBytes uint64,
	memoryPercent float64,
	storageInputBytes uint64,
	storageOutputBytes uint64,
	storageSpaceBytes uint64,
	storageSpacePercent float64,
	storageInodesCount uint64,
	netInputBytes uint64,
	netOutputBytes uint64,
) ContainerResourceUsage {
	return ContainerResourceUsage{
		CurrentCpuPercent:   currentCpuPercent,
		AverageCpuPercent:   averageCpuPercent,
		MemoryBytes:         memoryBytes,
		MemoryPercent:       memoryPercent,
		StorageInputBytes:   storageInputBytes,
		StorageOutputBytes:  storageOutputBytes,
		StorageSpaceBytes:   storageSpaceBytes,
		StorageSpacePercent: storageSpacePercent,
		StorageInodesCount:  storageInodesCount,
		NetInputBytes:       netInputBytes,
		NetOutputBytes:      netOutputBytes,
	}
}
