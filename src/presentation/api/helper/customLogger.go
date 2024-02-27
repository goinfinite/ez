package apiHelper

import (
	"log"
	"strings"
)

type CustomLogger struct{}

func (*CustomLogger) Write(rawMessage []byte) (int, error) {
	messageStr := strings.TrimSpace(string(rawMessage))

	shouldLog := true
	if strings.HasSuffix(messageStr, "tls: unknown certificate") {
		shouldLog = false
	}

	messageLen := len(rawMessage)
	if !shouldLog {
		return messageLen, nil
	}

	return messageLen, log.Output(2, messageStr)
}

func NewCustomLogger() *log.Logger {
	return log.New(&CustomLogger{}, "", 0)
}
