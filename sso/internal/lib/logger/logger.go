package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func SetupLogger(env string) zerolog.Logger {
	var log zerolog.Logger

	switch env {
	case envLocal:
		log = newLogger(zerolog.DebugLevel)
	case envDev:
		log = newLogger(zerolog.DebugLevel)
	case envProd:
		log = newLogger(zerolog.InfoLevel)
	}

	return log
}

func newLogger(level zerolog.Level) zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stderr, NoColor: false, TimeFormat: time.RFC822,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
	}).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()
}
