package valueObject

type HardwareSpecs struct {
	CpuModelName     CpuModelName  `json:"cpuModelName"`
	CpuCoresCount    CpuCoresCount `json:"cpuCoresCount"`
	CpuFrequency     float64       `json:"cpuFrequency"`
	MemoryTotalBytes Byte          `json:"memoryTotalBytes"`
}

func NewHardwareSpecs(
	cpuModelName CpuModelName,
	cpuCoresCount CpuCoresCount,
	cpuFrequency float64,
	memoryTotalBytes Byte,
) HardwareSpecs {
	return HardwareSpecs{
		CpuModelName:     cpuModelName,
		CpuCoresCount:    cpuCoresCount,
		CpuFrequency:     cpuFrequency,
		MemoryTotalBytes: memoryTotalBytes,
	}
}
