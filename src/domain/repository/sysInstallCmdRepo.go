package repository

type SysInstallCmdRepo interface {
	Install() error
	AddDataDisk() error
}
