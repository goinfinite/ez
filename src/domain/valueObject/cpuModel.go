package valueObject

import "errors"

type CpuModelName string

func NewCpuModelName(value string) (CpuModelName, error) {
	modelName := CpuModelName(value)
	if !modelName.isValid() {
		return "", errors.New("InvalidCpuModelName")
	}
	return modelName, nil
}

func NewCpuModelNamePanic(value string) CpuModelName {
	modelName, err := NewCpuModelName(value)
	if err != nil {
		panic(err)
	}
	return modelName
}

func (modelName CpuModelName) isValid() bool {
	isTooShort := len(string(modelName)) < 2
	isTooLong := len(string(modelName)) > 100
	return !isTooShort && !isTooLong
}

func (modelName CpuModelName) String() string {
	return string(modelName)
}
