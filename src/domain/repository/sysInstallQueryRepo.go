package repository

type SysInstallQueryRepo interface {
	IsInstalled() bool
	IsDataDiskMounted() bool
}
