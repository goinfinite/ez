package repository

type SysInstallCmdRepo interface {
	Install() error
	DisableDefaultSoftwares() error
	AddDataDisk() error
}
