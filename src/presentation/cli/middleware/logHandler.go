package cliMiddleware

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/rs/zerolog"
	slogZerolog "github.com/samber/slog-zerolog/v2"
	"golang.org/x/term"
)

type LogHandler struct {
}

func (LogHandler) ReadLevel() string {
	return os.Getenv("LOG_LEVEL")
}

func (LogHandler) SetLevel(logLevel string) {
	os.Setenv("LOG_LEVEL", logLevel)
}

func (logHandler LogHandler) Init() {
	var logWriter io.Writer = os.Stdout
	logLevel := zerolog.WarnLevel

	configuredLevel := logHandler.ReadLevel()
	if configuredLevel != "" {
		switch configuredLevel {
		case "DEBUG", "debug":
			stdoutFileDescriptor := int(os.Stdout.Fd())
			isInteractiveSession := term.IsTerminal(stdoutFileDescriptor)
			if isInteractiveSession {
				logWriter = zerolog.ConsoleWriter{
					Out: os.Stderr, TimeFormat: time.RFC3339,
				}
			}
			logLevel = zerolog.DebugLevel
		case "INFO", "info":
			logLevel = zerolog.InfoLevel
		case "WARN", "WARNING", "warn", "warning":
			logLevel = zerolog.WarnLevel
		case "ERROR", "error":
			logLevel = zerolog.ErrorLevel
		case "FATAL", "fatal":
			logLevel = zerolog.FatalLevel
		case "PANIC", "panic":
			logLevel = zerolog.PanicLevel
		}
	}
	zerologLogger := zerolog.New(logWriter).Level(logLevel)

	zerologHandler := slogZerolog.Option{Logger: &zerologLogger}.NewZerologHandler()
	slog.SetDefault(slog.New(zerologHandler))
}
