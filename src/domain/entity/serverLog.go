package entity

import (
	"github.com/speedianet/control/src/domain/valueObject"
)

type ServerLog struct {
	Timestamp valueObject.UnixTime           `json:"timestamp"`
	Level     valueObject.ServerLogLevel     `json:"level"`
	Operation valueObject.ServerLogOperation `json:"operation"`
	Payload   valueObject.ServerLogPayload   `json:"payload"`
}

func NewServerLog(
	level valueObject.ServerLogLevel,
	operation valueObject.ServerLogOperation,
	payload valueObject.ServerLogPayload,
) ServerLog {
	return ServerLog{
		Timestamp: valueObject.NewUnixTimeNow(),
		Level:     level,
		Operation: operation,
		Payload:   payload,
	}
}
