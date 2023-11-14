package valueObject

type NetInterfaceInfo struct {
	DeviceName      DeviceName `json:"deviceName"`
	RecvBytes       Byte       `json:"recvBytes"`
	RecvPackets     uint64     `json:"recvPackets"`
	RecvDropPackets uint64     `json:"recvDropPackets"`
	RecvErrs        uint64     `json:"recvErrs"`
	SentBytes       Byte       `json:"sentBytes"`
	SentPackets     uint64     `json:"sentPackets"`
	SentDropPackets uint64     `json:"sentDropPackets"`
	SentErrs        uint64     `json:"sentErrs"`
}

func NewNetInterfaceInfo(
	deviceName DeviceName,
	recvBytes Byte,
	recvPackets uint64,
	recvDropPackets uint64,
	recvErrs uint64,
	sentBytes Byte,
	sentPackets uint64,
	sentDropPackets uint64,
	sentErrs uint64,
) NetInterfaceInfo {
	return NetInterfaceInfo{
		DeviceName:      deviceName,
		RecvBytes:       recvBytes,
		RecvPackets:     recvPackets,
		RecvDropPackets: recvDropPackets,
		RecvErrs:        recvErrs,
		SentBytes:       sentBytes,
		SentPackets:     sentPackets,
		SentDropPackets: sentDropPackets,
		SentErrs:        sentErrs,
	}
}
