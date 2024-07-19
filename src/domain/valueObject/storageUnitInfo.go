package valueObject

type StorageUnitInfo struct {
	DeviceName        DeviceName     `json:"deviceName"`
	MountPoint        UnixFilePath   `json:"mountPoint"`
	FileSystem        UnixFileSystem `json:"fileSystem"`
	TotalBytes        Byte           `json:"totalBytes"`
	FreeBytes         Byte           `json:"freeBytes"`
	UsedBytes         Byte           `json:"usedBytes"`
	UsedPercent       float64        `json:"usedPercent"`
	TotalInodes       uint64         `json:"totalInodes"`
	FreeInodes        uint64         `json:"freeInodes"`
	UsedInodes        uint64         `json:"usedInodes"`
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
	totalBytes, freeBytes, usedBytes Byte,
	usedPercent float64,
	totalInodes, freeInodes, usedInodes uint64,
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
