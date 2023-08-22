package valueObject

import (
	"errors"
	"strings"
)

type ServerLogLevel string

const (
	debug ServerLogLevel = "debug"
	info  ServerLogLevel = "info"
	warn  ServerLogLevel = "warn"
	err   ServerLogLevel = "err"
	fatal ServerLogLevel = "fatal"
)

func NewServerLogLevel(value string) (ServerLogLevel, error) {
	value = strings.ToLower(value)
	sll := ServerLogLevel(value)
	if !sll.isValid() {
		return "", errors.New("InvalidServerLogLevel")
	}
	return sll, nil
}

func NewServerLogLevelPanic(value string) ServerLogLevel {
	sll, err := NewServerLogLevel(value)
	if err != nil {
		panic(err)
	}
	return sll
}

func (sll ServerLogLevel) isValid() bool {
	switch sll {
	case debug, info, warn, err, fatal:
		return true
	default:
		return false
	}
}

func (sll ServerLogLevel) String() string {
	return string(sll)
}
