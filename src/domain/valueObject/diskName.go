package valueObject

import "errors"

type DiskName string

func NewDiskName(value string) (DiskName, error) {
	diskName := DiskName(value)
	if !diskName.isValid() {
		return "", errors.New("InvalidDiskName")
	}
	return diskName, nil
}

func NewDiskNamePanic(value string) DiskName {
	diskName := DiskName(value)
	if !diskName.isValid() {
		panic("InvalidDiskName")
	}
	return diskName
}

func (diskName DiskName) isValid() bool {
	isTooShort := len(string(diskName)) < 3
	isTooLong := len(string(diskName)) > 64
	return !isTooShort && !isTooLong
}

func (diskName DiskName) String() string {
	return string(diskName)
}
