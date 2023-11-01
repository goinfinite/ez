package valueObject

type ContainerResourceUsage struct {
	AccountId          AccountId   `json:"-"`
	ContainerId        ContainerId `json:"-"`
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

func NewContainerResourceUsage(
	accountId AccountId,
	containerId ContainerId,
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
) ContainerResourceUsage {
	return ContainerResourceUsage{
		AccountId:          accountId,
		ContainerId:        containerId,
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
