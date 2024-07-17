package log

import (
	"os"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func init() {
	debug := os.Getenv("DEBUG") == "1"
	verbose := os.Getenv("DEBUG") == "2"
	level := zerolog.InfoLevel
	if debug {
		level = zerolog.DebugLevel
	}
	if verbose {
		level = zerolog.TraceLevel
	}
	console := zerolog.ConsoleWriter{Out: os.Stdout}
	logger = zerolog.New(console).With().Timestamp().Logger().Level(level)
	if verbose {
		logger = logger.With().Caller().Logger()
	}
}

func Debug() *zerolog.Event {
	return logger.Debug()
}

func Info() *zerolog.Event {
	return logger.Info()
}

func Warn() *zerolog.Event {
	return logger.Warn()
}

func Error() *zerolog.Event {
	return logger.Error()
}

func Fatal() *zerolog.Event {
	return logger.Fatal()
}

func Panic() *zerolog.Event {
	return logger.Panic()
}
