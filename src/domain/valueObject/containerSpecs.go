package valueObject

import (
	"errors"
	"strconv"
	"strings"
)

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

func NewContainerSpecsFromString(value string) (ContainerSpecs, error) {
	if value == "" {
		return ContainerSpecs{}, errors.New("InvalidContainerSpecs")
	}

	if !strings.Contains(value, ":") {
		return ContainerSpecs{}, errors.New("InvalidContainerSpecs")
	}

	specParts := strings.Split(value, ":")
	if len(specParts) != 2 {
		return ContainerSpecs{}, errors.New("InvalidContainerSpecs")
	}

	cpuCores, err := NewCpuCoresCount(specParts[0])
	if err != nil {
		return ContainerSpecs{}, err
	}

	memory, err := strconv.ParseUint(specParts[1], 10, 64)
	if err != nil {
		return ContainerSpecs{}, errors.New("InvalidMemoryLimit")
	}

	return NewContainerSpecs(
		cpuCores,
		Byte(int64(memory)),
	), nil
}

func (specs ContainerSpecs) GetCpuCores() CpuCoresCount {
	return specs.CpuCores
}

func (specs ContainerSpecs) GetMemoryBytes() Byte {
	return specs.MemoryBytes
}

func (specs ContainerSpecs) String() string {
	return specs.CpuCores.String() + ":" + specs.MemoryBytes.String()
}
