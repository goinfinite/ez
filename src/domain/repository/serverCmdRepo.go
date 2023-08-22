package repository

import "github.com/speedianet/sfm/src/domain/entity"

type ServerCmdRepo interface {
	Reboot() error
	AddOneTimerSvc(name string, cmd string) error
	DeleteOneTimerSvc(name string) error
	AddServerLog(log entity.ServerLog)
}
