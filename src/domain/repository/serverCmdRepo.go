package repository

import "github.com/speedianet/sfm/src/domain/valueObject"

type ServerCmdRepo interface {
	Reboot() error
	AddOneTimerSvc(name string, cmd string) error
	DeleteOneTimerSvc(name string) error
	AddServerLog(
		level valueObject.ServerLogLevel,
		operation valueObject.ServerLogOperation,
		payload valueObject.ServerLogPayload,
	)
}
