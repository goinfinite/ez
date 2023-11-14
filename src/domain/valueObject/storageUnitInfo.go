package valueObject

type StorageUnitInfo struct {
	DeviceName        DeviceName     `json:"deviceName"`
	MountPoint        UnixFilePath   `json:"mountPoint"`
	FileSystem        UnixFileSystem `json:"fileSystem"`
	TotalBytes        Byte           `json:"totalBytes"`
	FreeBytes         Byte           `json:"freeBytes"`
	UsedBytes         Byte           `json:"usedBytes"`
	UsedPercent       float64        `json:"usedPercent"`
	TotalInodes       InodesCount    `json:"totalInodes"`
	FreeInodes        InodesCount    `json:"freeInodes"`
	UsedInodes        InodesCount    `json:"usedInodes"`
	UsedInodesPercent float64        `json:"usedInodesPercent"`
	ReadBytes         Byte           `json:"readBytes"`
	ReadOpsCount      uint64         `json:"readOpsCount"`
	WriteBytes        Byte           `json:"writeBytes"`
	WriteOpsCount     uint64         `json:"writeOpsCount"`
}

func NewStorageUnitInfo(
	deviceName DeviceName,
	mountPoint UnixFilePath,
	fileSystem UnixFileSystem,
	totalBytes Byte,
	freeBytes Byte,
	usedBytes Byte,
	usedPercent float64,
	totalInodes InodesCount,
	freeInodes InodesCount,
	usedInodes InodesCount,
	usedInodesPercent float64,
	readBytes Byte,
	readOpsCount uint64,
	writeBytes Byte,
	writeOpsCount uint64,
) StorageUnitInfo {
	return StorageUnitInfo{
		DeviceName:        deviceName,
		MountPoint:        mountPoint,
		FileSystem:        fileSystem,
		TotalBytes:        totalBytes,
		FreeBytes:         freeBytes,
		UsedBytes:         usedBytes,
		UsedPercent:       usedPercent,
		TotalInodes:       totalInodes,
		FreeInodes:        freeInodes,
		UsedInodes:        usedInodes,
		UsedInodesPercent: usedInodesPercent,
		ReadBytes:         readBytes,
		ReadOpsCount:      readOpsCount,
		WriteBytes:        writeBytes,
		WriteOpsCount:     writeOpsCount,
	}
}
