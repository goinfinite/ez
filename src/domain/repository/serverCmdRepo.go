package repository

import "github.com/speedianet/control/src/domain/valueObject"

type ServerCmdRepo interface {
	Reboot() error
	AddSvc(name valueObject.SvcName, cmd valueObject.SvcCmd) error
	AddOneTimerSvc(name valueObject.SvcName, cmd valueObject.SvcCmd) error
	DeleteOneTimerSvc(name valueObject.SvcName) error
	AddServerLog(
		level valueObject.ServerLogLevel,
		operation valueObject.ServerLogOperation,
		payload valueObject.ServerLogPayload,
	)
	SendServerMessage(message string)
}
