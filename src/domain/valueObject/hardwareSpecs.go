package valueObject

import (
	"fmt"
	"strings"
)

type HardwareSpecs struct {
	CpuModelName     CpuModelName `json:"cpuModelName"`
	CpuCoresCount    float64      `json:"cpuCoresCount"`
	CpuFrequency     float64      `json:"cpuFrequency"`
	MemoryTotalBytes Byte         `json:"memoryTotalBytes"`
}

func NewHardwareSpecs(
	cpuModelName CpuModelName,
	cpuCoresCount float64,
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

func (vo HardwareSpecs) String() string {
	cpuModelNameStr := vo.CpuModelName.String()
	cpuModelNameParts := strings.Split(cpuModelNameStr, " ")
	if len(cpuModelNameParts) > 4 {
		cpuModelNameParts = cpuModelNameParts[:4]
	}
	cpuModelNameStr = strings.Join(cpuModelNameParts, " ")

	return fmt.Sprintf(
		"%s (%.0fc@%.1f GHz) ‖ %s RAM",
		cpuModelNameStr, vo.CpuCoresCount,
		vo.CpuFrequency, vo.MemoryTotalBytes.StringWithSuffix(),
	)
}
