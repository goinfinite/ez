package valueObject

import "errors"

type SvcName string

func NewSvcName(value string) (SvcName, error) {
	name := SvcName(value)
	if !name.isValid() {
		return "", errors.New("InvalidSvcName")
	}
	return name, nil
}

func NewSvcNamePanic(value string) SvcName {
	name, err := NewSvcName(value)
	if err != nil {
		panic(err)
	}
	return name
}

func (name SvcName) isValid() bool {
	isTooShort := len(string(name)) < 3
	isTooLong := len(string(name)) > 128
	return !isTooShort && !isTooLong
}

func (name SvcName) String() string {
	return string(name)
}
