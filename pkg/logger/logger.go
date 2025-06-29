package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New() zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	// Set log level based on environment
	level := zerolog.InfoLevel
	if os.Getenv("APP_ENV") == "development" {
		level = zerolog.DebugLevel
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	zerolog.SetGlobalLevel(level)

	return log.With().
		Str("service", "marketplace-api").
		Logger()
}
