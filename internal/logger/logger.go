package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/saleh-ghazimoradi/Cartopher/config"
)

func NewLogger(cfg *config.Config) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	if cfg.Server.GinMode != "release" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	return log.Logger
}
