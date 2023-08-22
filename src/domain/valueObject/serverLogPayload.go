package valueObject

import "errors"

type ServerLogPayload string

func NewServerLogPayload(value string) (ServerLogPayload, error) {
	logPayload := ServerLogPayload(value)
	if !logPayload.isValid() {
		return "", errors.New("InvalidServerLogPayload")
	}
	return logPayload, nil
}

func NewServerLogPayloadPanic(value string) ServerLogPayload {
	logPayload, err := NewServerLogPayload(value)
	if err != nil {
		panic(err)
	}
	return logPayload
}

func (logPayload ServerLogPayload) isValid() bool {
	isTooShort := len(string(logPayload)) < 2
	isTooLong := len(string(logPayload)) > 4096
	return !isTooShort && !isTooLong
}

func (logPayload ServerLogPayload) String() string {
	return string(logPayload)
}
