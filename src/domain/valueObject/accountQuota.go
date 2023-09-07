package valueObject

type AccountQuota struct {
	CpuCores    CpuCoresCount
	MemoryBytes Byte
	DiskBytes   Byte
	Inodes      InodesCount
}

func NewAccountQuota(
	cpuCores CpuCoresCount,
	memoryBytes Byte,
	diskBytes Byte,
	inodes InodesCount,
) AccountQuota {
	return AccountQuota{
		CpuCores:    cpuCores,
		MemoryBytes: memoryBytes,
		DiskBytes:   diskBytes,
		Inodes:      inodes,
	}
}

func NewAccountQuotaWithDefaultValues() AccountQuota {
	return AccountQuota{
		CpuCores:    NewCpuCoresCountPanic(0.5),
		MemoryBytes: NewBytePanic(1073741824),
		DiskBytes:   NewBytePanic(5368709120),
		Inodes:      NewInodesCountPanic(500000),
	}
}
