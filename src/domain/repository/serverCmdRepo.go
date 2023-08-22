package repository

type ServerCmdRepo interface {
	Reboot() error
	AddOneTimerSvc(name string, cmd string) error
	DeleteOneTimerSvc(name string) error
}
