package valueObject

import (
	"errors"
	"strings"
)

type ContainerSpecs struct {
	Millicores              Millicores              `json:"millicores"`
	MemoryBytes             Byte                    `json:"memoryBytes"`
	StoragePerformanceUnits StoragePerformanceUnits `json:"storagePerformanceUnits"`
}

func NewContainerSpecs(
	millicores Millicores, memoryBytes Byte, storagePerformanceUnits StoragePerformanceUnits,
) ContainerSpecs {
	return ContainerSpecs{
		Millicores:              millicores,
		MemoryBytes:             memoryBytes,
		StoragePerformanceUnits: storagePerformanceUnits,
	}
}

func NewContainerSpecsFromString(value string) (specs ContainerSpecs, err error) {
	if value == "" {
		return specs, errors.New("InvalidContainerSpecs")
	}

	if !strings.Contains(value, ":") {
		return specs, errors.New("InvalidContainerSpecs")
	}

	specParts := strings.Split(value, ":")
	if len(specParts) != 2 {
		return specs, errors.New("InvalidContainerSpecs")
	}

	millicores, err := NewMillicores(specParts[0])
	if err != nil {
		return specs, err
	}

	memory, err := NewByte(specParts[1])
	if err != nil {
		return specs, err
	}

	storagePerformanceUnits, _ := NewStoragePerformanceUnits(1)
	if len(specParts) == 3 {
		storagePerformanceUnits, err = NewStoragePerformanceUnits(specParts[2])
		if err != nil {
			return specs, err
		}
	}

	return NewContainerSpecs(millicores, memory, storagePerformanceUnits), nil
}

func (specs ContainerSpecs) String() string {
	return specs.Millicores.String() + ":" +
		specs.MemoryBytes.String() + ":" +
		specs.StoragePerformanceUnits.String()
}
