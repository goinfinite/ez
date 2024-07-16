package repository

import "github.com/speedianet/control/src/domain/valueObject"

type ServerCmdRepo interface {
	Reboot() error
	AddSvc(name valueObject.ServiceName, cmd valueObject.SvcCmd) error
	AddOneTimerSvc(name valueObject.ServiceName, cmd valueObject.SvcCmd) error
	DeleteOneTimerSvc(name valueObject.ServiceName) error
	SendServerMessage(message string)
}
