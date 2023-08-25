package dto

import "github.com/speedianet/sfm/src/domain/valueObject"

type AddAccount struct {
	Username    valueObject.Username `json:"username"`
	Password    valueObject.Password `json:"password"`
	CpuCores    float64              `json:"cpuCores"`
	MemoryBytes valueObject.Byte     `json:"memoryBytes"`
	DiskBytes   valueObject.Byte     `json:"diskBytes"`
	Inodes      uint64               `json:"inodes"`
}

func NewAddAccount(
	username valueObject.Username,
	password valueObject.Password,
	cpuCores float64,
	memoryBytes valueObject.Byte,
	diskBytes valueObject.Byte,
	inodes uint64,
) AddAccount {
	return AddAccount{
		Username:    username,
		Password:    password,
		CpuCores:    cpuCores,
		MemoryBytes: memoryBytes,
		DiskBytes:   diskBytes,
		Inodes:      inodes,
	}
}
