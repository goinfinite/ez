package valueObject

type ContainerSpecs struct {
	CpuCores     float64 `json:"cpuCores"`
	MemoryBytes  Byte    `json:"memoryBytes"`
	StorageBytes Byte    `json:"storageBytes"`
}

func NewContainerSpecs(
	cpuCores float64,
	memoryBytes Byte,
	storageBytes Byte,
) ContainerSpecs {
	return ContainerSpecs{
		CpuCores:     cpuCores,
		MemoryBytes:  memoryBytes,
		StorageBytes: storageBytes,
	}
}
