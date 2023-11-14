package valueObject

type StorageUnitInfo struct {
	DeviceName     DeviceName     `json:"deviceName"`
	MountPoint     UnixFilePath   `json:"mountPoint"`
	FileSystem     UnixFileSystem `json:"fileSystem"`
	TotalBytes     Byte           `json:"totalBytes"`
	AvailableBytes Byte           `json:"availableBytes"`
	UsedBytes      Byte           `json:"usedBytes"`
	UsedPercent    float64        `json:"usedPercent"`
	ReadBytes      Byte           `json:"readBytes"`
	ReadOpsCount   uint64         `json:"readOpsCount"`
	WriteBytes     Byte           `json:"writeBytes"`
	WriteOpsCount  uint64         `json:"writeOpsCount"`
}

func NewStorageUnitInfo(
	deviceName DeviceName,
	mountPoint UnixFilePath,
	fileSystem UnixFileSystem,
	totalBytes Byte,
	availableBytes Byte,
	usedBytes Byte,
	usedPercent float64,
) StorageUnitInfo {
	return StorageUnitInfo{
		DeviceName:     deviceName,
		MountPoint:     mountPoint,
		FileSystem:     fileSystem,
		TotalBytes:     totalBytes,
		AvailableBytes: availableBytes,
		UsedBytes:      usedBytes,
		UsedPercent:    usedPercent,
		ReadBytes:      0,
		ReadOpsCount:   0,
		WriteBytes:     0,
		WriteOpsCount:  0,
	}
}
