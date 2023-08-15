package valueObject

type DiskInfo struct {
	Name           DiskName `json:"name"`
	TotalBytes     Byte     `json:"totalBytes"`
	AvailableBytes Byte     `json:"availableBytes"`
	UsedBytes      Byte     `json:"usedBytes"`
}

func NewDiskInfo(
	name DiskName,
	total Byte,
	available Byte,
	used Byte,
) DiskInfo {
	return DiskInfo{
		Name:           name,
		TotalBytes:     total,
		AvailableBytes: available,
		UsedBytes:      used,
	}
}
