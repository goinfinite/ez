package valueObject

type StorageDeviceInfo struct {
	DeviceName     DeviceName `json:"deviceName"`
	TotalBytes     Byte       `json:"totalBytes"`
	AvailableBytes Byte       `json:"availableBytes"`
	UsedBytes      Byte       `json:"usedBytes"`
	UsedPercent    float64    `json:"usedPercent"`
}

func NewStorageDeviceInfo(
	deviceName DeviceName,
	total Byte,
	available Byte,
	usedBytes Byte,
	usedPercent float64,
) StorageDeviceInfo {
	return StorageDeviceInfo{
		DeviceName:     deviceName,
		TotalBytes:     total,
		AvailableBytes: available,
		UsedBytes:      usedBytes,
		UsedPercent:    usedPercent,
	}
}
