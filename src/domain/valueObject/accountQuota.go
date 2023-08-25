package valueObject

type AccountQuota struct {
	CpuCores    float64
	MemoryBytes Byte
	DiskBytes   Byte
	Inodes      uint64
}

func NewAccountQuota(
	cpuCores float64,
	memoryBytes Byte,
	diskBytes Byte,
	inodes uint64,
) AccountQuota {
	return AccountQuota{
		CpuCores:    cpuCores,
		MemoryBytes: memoryBytes,
		DiskBytes:   diskBytes,
		Inodes:      inodes,
	}
}
