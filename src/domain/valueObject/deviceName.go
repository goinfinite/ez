package valueObject

import "errors"

type DeviceName string

func NewDeviceName(value string) (DeviceName, error) {
	deviceName := DeviceName(value)
	if !deviceName.isValid() {
		return "", errors.New("InvalidDeviceName")
	}
	return deviceName, nil
}

func NewDeviceNamePanic(value string) DeviceName {
	deviceName := DeviceName(value)
	if !deviceName.isValid() {
		panic("InvalidDeviceName")
	}
	return deviceName
}

func (deviceName DeviceName) isValid() bool {
	isTooShort := len(string(deviceName)) < 3
	isTooLong := len(string(deviceName)) > 64
	return !isTooShort && !isTooLong
}

func (deviceName DeviceName) String() string {
	return string(deviceName)
}
