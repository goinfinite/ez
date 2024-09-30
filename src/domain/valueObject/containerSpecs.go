package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/ez/src/domain/valueObject/helper"
)

type ContainerSpecs struct {
	Millicores              Millicores              `json:"millicores"`
	CpuCores                float64                 `json:"cpuCores"`
	MemoryBytes             Byte                    `json:"memoryBytes"`
	MemoryMebibytes         int64                   `json:"memoryMebibytes"`
	MemoryGibibytes         int64                   `json:"memoryGibibytes"`
	StoragePerformanceUnits StoragePerformanceUnits `json:"storagePerformanceUnits"`
}

func NewContainerSpecs(
	millicores Millicores, memoryBytes Byte, storagePerformanceUnits StoragePerformanceUnits,
) ContainerSpecs {
	cpuCores := millicores.ReadAsCores()
	memoryMebibytes := memoryBytes.ToMiB()
	memoryGibibytes := memoryBytes.ToGiB()

	return ContainerSpecs{
		Millicores:              millicores,
		CpuCores:                cpuCores,
		MemoryBytes:             memoryBytes,
		MemoryMebibytes:         memoryMebibytes,
		MemoryGibibytes:         memoryGibibytes,
		StoragePerformanceUnits: storagePerformanceUnits,
	}
}

// format: [millicores]:[memoryBytes]:[storagePerformanceUnits]
func NewContainerSpecsFromString(value string) (specs ContainerSpecs, err error) {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)

	specsRegex := `^(?:(?P<millicores>\d{1,19}))?:(?P<memoryBytes>\d{1,19}):(?P<storagePerformanceUnits>\d{1,19})$`
	specsParts := voHelper.FindNamedGroupsMatches(specsRegex, value)
	if len(specsParts) == 0 {
		return specs, errors.New("InvalidSpecsStructure")
	}

	specs = NewContainerSpecsWithDefaultValues()

	if specsParts["millicores"] != "" {
		specs.Millicores, err = NewMillicores(specsParts["millicores"])
		if err != nil {
			return specs, err
		}
	}

	if specsParts["memoryBytes"] != "" {
		specs.MemoryBytes, err = NewByte(specsParts["memoryBytes"])
		if err != nil {
			return specs, err
		}
	}

	if specsParts["storagePerformanceUnits"] != "" {
		specs.StoragePerformanceUnits, err = NewStoragePerformanceUnits(
			specsParts["storagePerformanceUnits"],
		)
		if err != nil {
			return specs, err
		}
	}

	return NewContainerSpecs(specs.Millicores, specs.MemoryBytes, specs.StoragePerformanceUnits), nil
}

func NewContainerSpecsWithDefaultValues() ContainerSpecs {
	millicores, _ := NewMillicores(500)
	memoryBytes, _ := NewByte(1073741824)
	storagePerformanceUnits, _ := NewStoragePerformanceUnits(1)

	return NewContainerSpecs(
		millicores, memoryBytes, storagePerformanceUnits,
	)
}

func (specs ContainerSpecs) String() string {
	return specs.Millicores.String() + ":" +
		specs.MemoryBytes.String() + ":" +
		specs.StoragePerformanceUnits.String()
}
