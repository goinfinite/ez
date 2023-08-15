package valueObject

type NetUsage struct {
	ReceivedBytes Byte `json:"receivedBytes"`
	SentBytes     Byte `json:"sentBytes"`
}

func NewNetUsage(receivedBytes Byte, sentBytes Byte) NetUsage {
	return NetUsage{ReceivedBytes: receivedBytes, SentBytes: sentBytes}
}
