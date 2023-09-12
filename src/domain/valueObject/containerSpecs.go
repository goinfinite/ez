package valueObject

type ContainerSpecs struct {
	CpuCores    CpuCoresCount `json:"cpuCores"`
	MemoryBytes Byte          `json:"memoryBytes"`
}

func NewContainerSpecs(cpuCores CpuCoresCount, memoryBytes Byte) ContainerSpecs {
	return ContainerSpecs{
		CpuCores:    cpuCores,
		MemoryBytes: memoryBytes,
	}
}

func (specs ContainerSpecs) GetCpuCores() CpuCoresCount {
	return specs.CpuCores
}

func (specs ContainerSpecs) GetMemoryBytes() Byte {
	return specs.MemoryBytes
}
