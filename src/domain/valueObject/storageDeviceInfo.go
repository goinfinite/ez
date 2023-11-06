package valueObject

type StorageDeviceInfo struct {
	DeviceName     DeviceName `json:"deviceName"`
	TotalBytes     Byte       `json:"totalBytes"`
	AvailableBytes Byte       `json:"availableBytes"`
	UsedBytes      Byte       `json:"usedBytes"`
}

func NewStorageDeviceInfo(
	deviceName DeviceName,
	total Byte,
	available Byte,
	used Byte,
) StorageDeviceInfo {
	return StorageDeviceInfo{
		DeviceName:     deviceName,
		TotalBytes:     total,
		AvailableBytes: available,
		UsedBytes:      used,
	}
}
