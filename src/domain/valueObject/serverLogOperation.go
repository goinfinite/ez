package valueObject

import "errors"

type ServerLogOperation string

func NewServerLogOperation(value string) (ServerLogOperation, error) {
	logOperation := ServerLogOperation(value)
	if !logOperation.isValid() {
		return "", errors.New("InvalidServerLogOperation")
	}
	return logOperation, nil
}

func NewServerLogOperationPanic(value string) ServerLogOperation {
	logOperation, err := NewServerLogOperation(value)
	if err != nil {
		panic(err)
	}
	return logOperation
}

func (logOperation ServerLogOperation) isValid() bool {
	isTooShort := len(string(logOperation)) < 3
	isTooLong := len(string(logOperation)) > 128
	return !isTooShort && !isTooLong
}

func (logOperation ServerLogOperation) String() string {
	return string(logOperation)
}
