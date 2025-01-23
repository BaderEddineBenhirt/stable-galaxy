package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	*zerolog.Logger
}

func NewLogger(level string, isJSON bool) *Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	var l zerolog.Level
	switch level {
	case "debug":
		l = zerolog.DebugLevel
	case "info":
		l = zerolog.InfoLevel
	case "warn":
		l = zerolog.WarnLevel
	case "error":
		l = zerolog.ErrorLevel
	default:
		l = zerolog.InfoLevel
	}

	var logger zerolog.Logger
	if isJSON {
		logger = zerolog.New(os.Stdout).
			Level(l).
			With().
			Timestamp().
			Logger()
	} else {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		logger = zerolog.New(output).
			Level(l).
			With().
			Timestamp().
			Logger()
	}

	return &Logger{&logger}
}
